(function () {
    var canvas = document.getElementById('game-canvas');
    if (!canvas) return;

    var W = 480, H = 360;
    var MAX_WAVE = 5;
    var ctx = canvas.getContext('2d');
    ctx.imageSmoothingEnabled = false;

    var str = parseInt(canvas.dataset.str, 10) || 0;
    var dex = parseInt(canvas.dataset.dex, 10) || 0;
    var level = parseInt(canvas.dataset.level, 10) || 1;
    var accent = canvas.dataset.accent || '#FFD700';

    var overlay = document.getElementById('game-overlay');
    var overlayTitle = document.getElementById('overlay-title');
    var overlaySub = document.getElementById('overlay-sub');
    var retryBtn = document.getElementById('retry-btn');
    var menuBtn = document.getElementById('menu-btn');
    var menu = document.getElementById('arena-menu');
    var classPreview = document.getElementById('class-preview');
    var className = document.getElementById('class-name');
    var prevClassBtn = document.getElementById('prev-class');
    var nextClassBtn = document.getElementById('next-class');
    var weaponLabel = document.getElementById('weapon-label');
    var startBtn = document.getElementById('start-btn');
    var instructionsBtn = document.getElementById('instructions-btn');
    var instructions = document.getElementById('arena-instructions');

    function img(src) {
        var i = new Image();
        i.src = src;
        return i;
    }
    var classWeapons = {
        guardian: { kind: 'melee', label: 'BLADE ARC', color: '#00add8', cd: 0.42, mult: 1.6, range: 62, arc: 1.0, knock: 18, width: 4 },
        knight: { kind: 'melee', label: 'LONGSWORD', color: '#9b72cf', cd: 0.52, mult: 2.1, range: 76, arc: 0.55, knock: 14, width: 3 },
        warlord: { kind: 'melee', label: 'WAR CLEAVER', color: '#f34b7d', cd: 0.78, mult: 3.2, range: 68, arc: 1.5, knock: 34, width: 7, shake: 4 },
        berserker: { kind: 'spin', label: 'WHIRLWIND', color: '#e05d44', cd: 0.5, mult: 1.2, range: 60, knock: 16, width: 5 },
        sage: { kind: 'orb', label: 'ARCANE ORB', color: '#4b8bbe', cd: 0.55, mult: 1.5, speed: 200, blast: 52, radius: 7 },
        battlemage: { kind: 'orb', label: 'EMBER BOLT', color: '#c07d28', cd: 0.32, mult: 0.85, speed: 300, blast: 32, radius: 4 },
        rogue: { kind: 'spread', label: 'FAN OF KNIVES', color: '#e8c94a', cd: 0.32, mult: 0.55, count: 5, spread: 0.2, speed: 360 },
        paladin: { kind: 'pierce', label: 'HOLY LANCE', color: '#3178c6', cd: 0.45, mult: 1.4, speed: 420, width: 5 },
        wanderer: { kind: 'pierce', label: 'PIERCING ARROW', color: '#6e7681', cd: 0.3, mult: 0.85, speed: 540, width: 2 }
    };
    var classOrder = ['guardian', 'berserker', 'paladin', 'rogue', 'sage', 'knight', 'battlemage', 'warlord', 'wanderer'];

    function spritePath(cls) {
        return '/static/assets/sprites/' + cls + '.png';
    }

    var selectedClass = canvas.dataset.classname || 'wanderer';
    var playerImg = img(spritePath(selectedClass));
    var playerImgB = img(spritePath(selectedClass).replace('.png', '_b.png'));
    var sprites = {
        slime: img('/static/assets/sprites/slime.png'),
        imp: img('/static/assets/sprites/imp.png'),
        wraith: img('/static/assets/sprites/wraith.png'),
        heart: img('/static/assets/sprites/heart.png'),
        star: img('/static/assets/sprites/star.png'),
        bolt: img('/static/assets/sprites/bolt.png'),
        boss: img('/static/assets/sprites/boss.png')
    };

    var SIZE = 32;
    var wpnColor = accent;
    var keys = {};
    var mouse = { x: W / 2, y: 0 };
    var mode = 'menu'; // menu | play | ending | win | over
    var state;

    function setClass(cls) {
        selectedClass = cls;
        wpnColor = classWeapons[cls].color;
        playerImg = img(spritePath(cls));
        playerImgB = img(spritePath(cls).replace('.png', '_b.png'));
        classPreview.src = spritePath(cls);
        className.textContent = cls.toUpperCase();
        className.style.color = classWeapons[cls].color;
        weaponLabel.textContent = classWeapons[cls].label;
    }

    function cycleClass(step) {
        var i = classOrder.indexOf(selectedClass);
        setClass(classOrder[(i + step + classOrder.length) % classOrder.length]);
    }

    var enemyTypes = {
        slime: { hp: 22, hpW: 6, speed: 60, speedW: 4, maxSpeed: 140, touch: 13, color: '#5cb85c' },
        imp: { hp: 17, hpW: 4, speed: 70, speedW: 3, maxSpeed: 120, touch: 9, color: '#e05d44' },
        wraith: { hp: 14, hpW: 3, speed: 40, speedW: 2, maxSpeed: 80, touch: 16, color: '#9b72cf' },
        boss: { hp: 400, hpW: 0, speed: 55, speedW: 0, maxSpeed: 55, touch: 25, color: '#8e2a3c' }
    };

    // Damage swings ~10-50 with POWER; without this, strong characters one-shot everything.
    function hpScale() {
        return 0.8 + state.player.dmg / 30;
    }

    function newState() {
        return {
            player: {
                x: W / 2 - SIZE / 2, y: H / 2 - SIZE / 2,
                maxHp: 80 + level * 2, hp: 80 + level * 2,
                speed: 120 + dex * 0.8,
                dmg: 10 + str * 0.4,
                fireCd: 0, iframes: 0, recoil: 0
            },
            enemies: [], shots: [], enemyShots: [], pickups: [],
            particles: [], texts: [],
            buffs: { star: 0, bolt: 0 },
            spawnQueue: [], spawnTimer: 0,
            wave: 0, kills: 0, banner: 0, shake: 0,
            moving: false, animTime: 0,
            pending: null, slash: null
        };
    }

    function endGame(nextMode, title, sub, delay) {
        mode = 'ending';
        state.pending = { mode: nextMode, title: title, sub: sub, t: delay };
    }

    function maxAlive() {
        return 4 + state.wave;
    }

    // Farthest of three edge candidates, so spawns miss a cornered player.
    function spawnPoint(size) {
        var p = state.player;
        var pick = { x: 0, y: 0 }, best = -1;
        for (var i = 0; i < 3; i++) {
            var x = Math.random() * (W - size);
            var y = Math.random() * (H - size);
            switch (Math.floor(Math.random() * 4)) {
                case 0: y = -size; break;
                case 1: y = H; break;
                case 2: x = -size; break;
                default: x = W;
            }
            var dx = x - p.x, dy = y - p.y;
            if (dx * dx + dy * dy > best) {
                best = dx * dx + dy * dy;
                pick.x = x;
                pick.y = y;
            }
        }
        return pick;
    }

    function spawnEnemy(type) {
        var t = enemyTypes[type];
        var size = type === 'boss' ? 64 : SIZE;
        var at = spawnPoint(size);
        var x = at.x;
        var y = at.y;
        var step = state.wave - 1;
        var hp = Math.round((t.hp + step * t.hpW) * hpScale());
        if (type === 'boss') {
            x = W / 2 - size / 2;
            y = -size;
            hp = Math.round((400 + level * 5) * hpScale());
        }
        state.enemies.push({
            type: type, x: x, y: y, size: size,
            hp: hp, maxHp: hp,
            speed: Math.min(t.speed + step * t.speedW, t.maxSpeed),
            flash: 0, timer: 1 + Math.random() * 1.5, fade: 1,
            phase: 'chase', phaseT: 2, dashX: 0, dashY: 0,
            lunge: 0, wind: 0, burst: 0, orbit: Math.random() < 0.5 ? 1 : -1,
            shotT: 2.5
        });
    }

    function spawnWave() {
        state.wave++;
        state.banner = 1.5;
        var w = state.wave;
        var q = [];
        var i;
        if (w >= MAX_WAVE) {
            state.spawnQueue = ['wraith', 'wraith', 'wraith', 'imp', 'imp', 'imp', 'imp',
                'slime', 'slime', 'slime', 'slime', 'slime', 'boss'];
            state.spawnTimer = 0.3;
            return;
        }
        for (i = 0; i < 2 + w * 2; i++) q.push('slime');
        if (w >= 2) for (i = 0; i < w; i++) q.push('imp');
        if (w >= 3) for (i = 0; i < w - 1; i++) q.push('wraith');
        for (i = q.length - 1; i > 0; i--) {
            var j = Math.floor(Math.random() * (i + 1));
            var tmp = q[i];
            q[i] = q[j];
            q[j] = tmp;
        }
        state.spawnQueue = q;
        state.spawnTimer = 0.3;
    }

    function burst(x, y, color, n, speed) {
        for (var i = 0; i < n; i++) {
            var a = Math.random() * Math.PI * 2;
            var v = speed * (0.4 + Math.random() * 0.6);
            state.particles.push({
                x: x, y: y, vx: Math.cos(a) * v, vy: Math.sin(a) * v,
                life: 0.35 + Math.random() * 0.3, color: color, size: 2 + Math.random() * 2
            });
        }
    }

    function floatText(x, y, s, color) {
        state.texts.push({ x: x, y: y, s: s, color: color, life: 0.7 });
    }

    function aimDir() {
        var p = state.player;
        var dx = mouse.x - (p.x + SIZE / 2);
        var dy = mouse.y - (p.y + SIZE / 2);
        var d = Math.sqrt(dx * dx + dy * dy) || 1;
        return { x: dx / d, y: dy / d };
    }

    function dmgFor(mult) {
        return Math.round(state.player.dmg * mult * (state.buffs.star > 0 ? 2 : 1));
    }

    function hitEnemy(e, dmg, kx, ky) {
        e.hp -= dmg;
        e.flash = 0.12;
        floatText(e.x + e.size / 2, e.y, '' + dmg, state.buffs.star > 0 ? '#FFD700' : '#ffffff');
        if (e.type !== 'boss') {
            e.x += kx;
            e.y += ky;
        }
    }

    function shoot() {
        var p = state.player;
        if (p.fireCd > 0) return;
        var wpn = classWeapons[selectedClass];
        p.fireCd = wpn.cd * (state.buffs.bolt > 0 ? 0.5 : 1);
        p.recoil = 4;
        var dir = aimDir();
        var cx = p.x + SIZE / 2, cy = p.y + SIZE / 2;
        var ang = Math.atan2(dir.y, dir.x);
        var i, e, ex, ey, d;

        if (wpn.kind === 'melee' || wpn.kind === 'spin') {
            var full = wpn.kind === 'spin';
            state.slash = {
                x: cx, y: cy, ang: ang, t: 0.15, full: full,
                r: wpn.range - 8, arc: wpn.arc, color: wpn.color, width: wpn.width
            };
            for (i = 0; i < state.enemies.length; i++) {
                e = state.enemies[i];
                ex = e.x + e.size / 2 - cx;
                ey = e.y + e.size / 2 - cy;
                d = Math.sqrt(ex * ex + ey * ey);
                if (d > wpn.range + e.size / 4) continue;
                if (!full) {
                    var diff = Math.abs(Math.atan2(ey, ex) - ang);
                    if (diff > Math.PI) diff = Math.PI * 2 - diff;
                    if (diff > wpn.arc) continue;
                }
                d = d || 1;
                hitEnemy(e, dmgFor(wpn.mult), ex / d * wpn.knock, ey / d * wpn.knock);
            }
            if (wpn.shake) state.shake = Math.max(state.shake, wpn.shake);
            burst(cx + dir.x * 30, cy + dir.y * 30, wpn.color, 4, 60);
        } else if (wpn.kind === 'orb') {
            state.shots.push({
                kind: 'orb', mult: wpn.mult, color: wpn.color,
                radius: wpn.radius, blast: wpn.blast,
                x: cx, y: cy,
                vx: dir.x * wpn.speed, vy: dir.y * wpn.speed
            });
        } else if (wpn.kind === 'spread') {
            var half = (wpn.count - 1) / 2;
            for (i = 0; i < wpn.count; i++) {
                var a = ang + (i - half) * wpn.spread;
                state.shots.push({
                    kind: 'shot', mult: wpn.mult, color: wpn.color,
                    x: cx, y: cy,
                    vx: Math.cos(a) * wpn.speed, vy: Math.sin(a) * wpn.speed
                });
            }
        } else if (wpn.kind === 'pierce') {
            state.shots.push({
                kind: 'pierce', mult: wpn.mult, color: wpn.color, width: wpn.width,
                x: cx, y: cy,
                vx: dir.x * wpn.speed, vy: dir.y * wpn.speed, hitList: []
            });
        }
        if (wpn.kind !== 'melee' && wpn.kind !== 'spin') {
            burst(cx + dir.x * 18, cy + dir.y * 18, wpn.color, 2, 40);
        }
    }

    function dropPickup(x, y) {
        if (Math.random() > 0.2) return;
        var types = ['heart', 'star', 'bolt'];
        state.pickups.push({ type: types[Math.floor(Math.random() * types.length)], x: x, y: y, ttl: 10 });
    }

    function hurtPlayer(dmg) {
        var p = state.player;
        if (mode !== 'play' || p.iframes > 0) return;
        p.hp -= dmg;
        p.iframes = 0.8;
        state.shake = 6;
        burst(p.x + SIZE / 2, p.y + SIZE / 2, '#e05d44', 8, 90);
        if (p.hp <= 0) {
            p.hp = 0;
            state.shake = 10;
            burst(p.x + SIZE / 2, p.y + SIZE / 2, accent, 20, 140);
            endGame('over', 'GAME OVER', 'WAVE ' + state.wave + ' · ' + state.kills + ' KILLS', 1.3);
        }
    }

    function killEnemy(i) {
        var e = state.enemies[i];
        var t = enemyTypes[e.type];
        burst(e.x + e.size / 2, e.y + e.size / 2, t.color, e.type === 'boss' ? 30 : 10, e.type === 'boss' ? 160 : 110);
        dropPickup(e.x, e.y);
        state.enemies.splice(i, 1);
        state.kills++;
        if (e.type === 'boss') {
            state.shake = 10;
            floatText(W / 2, H / 2 - 30, 'THE WARDEN FALLS', '#FFD700');
        }
        if (state.enemies.length === 0 && state.spawnQueue.length === 0 && state.wave >= MAX_WAVE) {
            endGame('win', 'VICTORY!', 'ARENA CLEARED · ' + state.kills + ' KILLS', 1.8);
        }
    }

    function showOverlay(title, sub) {
        overlayTitle.textContent = title;
        overlaySub.textContent = sub;
        overlay.style.display = 'flex';
    }

    function hideOverlay() {
        overlay.style.display = 'none';
    }

    function restart() {
        state = newState();
        mode = 'play';
        hideOverlay();
        menu.style.display = 'none';
        retryBtn.textContent = 'TRY AGAIN';
        spawnWave();
    }

    function toMenu() {
        state = newState();
        mode = 'menu';
        hideOverlay();
        menu.style.display = 'flex';
        retryBtn.textContent = 'TRY AGAIN';
    }

    function pauseGame() {
        if (mode !== 'play') return;
        mode = 'pause';
        retryBtn.textContent = 'RESUME';
        showOverlay('PAUSED', 'WAVE ' + state.wave + '/' + MAX_WAVE + ' · ' + state.kills + ' KILLS');
    }

    function resumeGame() {
        mode = 'play';
        retryBtn.textContent = 'TRY AGAIN';
        hideOverlay();
    }

    function updateEnemy(e, dt) {
        var p = state.player;
        var t = enemyTypes[e.type];
        var dx = p.x - e.x, dy = p.y - e.y;
        var d = Math.sqrt(dx * dx + dy * dy) || 1;

        if (e.type === 'slime') {
            if (e.lunge > 0) {
                e.lunge -= dt;
                e.x += e.dashX * dt;
                e.y += e.dashY * dt;
            } else if (e.wind > 0) {
                e.wind -= dt;
                e.flash = 0.05;
                if (e.wind <= 0) {
                    e.lunge = 0.32;
                    e.dashX = dx / d * e.speed * 3.4;
                    e.dashY = dy / d * e.speed * 3.4;
                }
            } else {
                e.x += dx / d * e.speed * dt;
                e.y += dy / d * e.speed * dt;
                e.timer -= dt;
                if (e.timer <= 0) {
                    e.timer = 2.5 + Math.random() * 2;
                    if (d < 170 && state.wave >= 2) e.wind = 0.35;
                }
            }
        } else if (e.type === 'imp') {
            var want = 130;
            var dir = d > want ? 1 : -0.7;
            e.x += (dx / d * dir - dy / d * e.orbit * 0.8) * e.speed * dt;
            e.y += (dy / d * dir + dx / d * e.orbit * 0.8) * e.speed * dt;
            e.timer -= dt;
            if (e.timer <= 0 && e.x > 0 && e.x < W && e.y > 0 && e.y < H) {
                e.burst = state.wave >= 3 ? 3 : 2;
                e.timer = 2.8 - state.wave * 0.1 + Math.random();
            }
            if (e.burst > 0) {
                e.shotT -= dt;
                if (e.shotT <= 0) {
                    e.burst--;
                    e.shotT = 0.16;
                    state.enemyShots.push({
                        x: e.x + SIZE / 2, y: e.y + SIZE / 2,
                        vx: dx / d * 170, vy: dy / d * 170
                    });
                }
            } else {
                e.shotT = 0;
            }
        } else if (e.type === 'wraith') {
            if (e.lunge > 0) {
                e.lunge -= dt;
                e.x += e.dashX * dt;
                e.y += e.dashY * dt;
            } else {
                e.x += dx / d * e.speed * dt;
                e.y += dy / d * e.speed * dt;
            }
            e.timer -= dt;
            if (e.timer <= 0) {
                e.timer = 2.8 - state.wave * 0.2 + Math.random() * 1.2;
                burst(e.x + SIZE / 2, e.y + SIZE / 2, t.color, 6, 70);
                var a = Math.random() * Math.PI * 2;
                e.x = Math.max(0, Math.min(W - SIZE, p.x + Math.cos(a) * 80));
                e.y = Math.max(0, Math.min(H - SIZE, p.y + Math.sin(a) * 80));
                e.fade = 0;
                burst(e.x + SIZE / 2, e.y + SIZE / 2, t.color, 6, 70);
                e.lunge = 0.4;
                e.dashX = (p.x - e.x) / 80 * 150;
                e.dashY = (p.y - e.y) / 80 * 150;
            }
            e.fade = Math.min(1, e.fade + dt * 4);
        } else if (e.type === 'boss') {
            updateBoss(e, dt, dx, dy, d);
        }

        var reach = (SIZE + e.size) / 2 * 0.7;
        if (Math.abs(p.x + SIZE / 2 - (e.x + e.size / 2)) < reach &&
            Math.abs(p.y + SIZE / 2 - (e.y + e.size / 2)) < reach) {
            hurtPlayer(t.touch);
        }
    }

    function updateBoss(e, dt, dx, dy, d) {
        var i, a;
        e.phaseT -= dt;
        if (e.phase === 'chase') {
            e.x += dx / d * e.speed * dt;
            e.y += dy / d * e.speed * dt;
            e.shotT -= dt;
            if (e.shotT <= 0) {
                e.shotT = 2.2 + Math.random();
                for (i = -1; i <= 1; i++) {
                    a = Math.atan2(dy, dx) + i * 0.28;
                    state.enemyShots.push({
                        x: e.x + e.size / 2, y: e.y + e.size / 2,
                        vx: Math.cos(a) * 190, vy: Math.sin(a) * 190
                    });
                }
            }
            if (e.phaseT <= 0) {
                var r = Math.random();
                if (r < 0.55) {
                    e.phase = 'telegraph';
                    e.phaseT = 0.45;
                    e.burst = Math.random() < 0.25 ? 2 : 1;
                } else if (r < 0.82) {
                    var spin = Math.random() * Math.PI * 2;
                    for (i = 0; i < 12; i++) {
                        a = spin + i / 12 * Math.PI * 2;
                        state.enemyShots.push({
                            x: e.x + e.size / 2, y: e.y + e.size / 2,
                            vx: Math.cos(a) * 150, vy: Math.sin(a) * 150
                        });
                    }
                    burst(e.x + e.size / 2, e.y + e.size / 2, '#e05d44', 10, 100);
                    e.phase = 'chase';
                    e.phaseT = 1.2;
                } else if (state.enemies.length <= maxAlive() - 3) {
                    spawnEnemy('slime');
                    spawnEnemy('slime');
                    spawnEnemy(Math.random() < 0.5 ? 'imp' : 'wraith');
                    floatText(e.x + e.size / 2, e.y - 8, 'SUMMON', '#9b72cf');
                    e.phase = 'chase';
                    e.phaseT = 1.8;
                } else {
                    e.phase = 'chase';
                    e.phaseT = 1;
                }
            }
        } else if (e.phase === 'telegraph') {
            e.flash = 0.05;
            if (e.phaseT <= 0) {
                e.phase = 'dash';
                e.phaseT = 0.5;
                e.dashX = dx / d * 420;
                e.dashY = dy / d * 420;
                floatText(e.x + e.size / 2, e.y - 8, 'CHARGE!', '#FFD700');
            }
        } else if (e.phase === 'dash') {
            e.x += e.dashX * dt;
            e.y += e.dashY * dt;
            e.x = Math.max(0, Math.min(W - e.size, e.x));
            e.y = Math.max(0, Math.min(H - e.size, e.y));
            if (e.phaseT <= 0) {
                e.burst--;
                if (e.burst > 0) {
                    e.phase = 'telegraph';
                    e.phaseT = 0.35;
                } else {
                    e.phase = 'chase';
                    e.phaseT = 1.4;
                }
            }
        }
    }

    function update(dt) {
        if (mode !== 'play' && mode !== 'ending') return;
        var p = state.player;
        var i, e, s;

        if (mode === 'play') {
            var mx = 0, my = 0;
            if (keys.ArrowLeft || keys.a) mx -= 1;
            if (keys.ArrowRight || keys.d) mx += 1;
            if (keys.ArrowUp || keys.w) my -= 1;
            if (keys.ArrowDown || keys.s) my += 1;
            state.moving = mx !== 0 || my !== 0;
            if (state.moving) {
                state.animTime += dt;
                var len = Math.sqrt(mx * mx + my * my);
                var spd = p.speed + (state.buffs.bolt > 0 ? 40 : 0);
                p.x += mx / len * spd * dt;
                p.y += my / len * spd * dt;
            }
            p.x = Math.max(0, Math.min(W - SIZE, p.x));
            p.y = Math.max(0, Math.min(H - SIZE, p.y));

            if (state.spawnQueue.length > 0) {
                state.spawnTimer -= dt;
                if (state.spawnTimer <= 0) {
                    if (state.enemies.length < maxAlive()) {
                        spawnEnemy(state.spawnQueue.pop());
                        state.spawnTimer = 0.35 + Math.random() * 0.6;
                    } else {
                        state.spawnTimer = 0.3;
                    }
                }
            }

        }

        p.recoil = Math.max(0, p.recoil - dt * 24);
        p.fireCd = Math.max(0, p.fireCd - dt);
        p.iframes = Math.max(0, p.iframes - dt);
        state.buffs.star = Math.max(0, state.buffs.star - dt);
        state.buffs.bolt = Math.max(0, state.buffs.bolt - dt);

        for (i = state.shots.length - 1; i >= 0; i--) {
            s = state.shots[i];
            s.x += s.vx * dt;
            s.y += s.vy * dt;
            if (s.x < -10 || s.x > W + 10 || s.y < -10 || s.y > H + 10) {
                state.shots.splice(i, 1);
                continue;
            }
            if (s.kind === 'pierce') {
                for (var j0 = 0; j0 < state.enemies.length; j0++) {
                    e = state.enemies[j0];
                    var pr = e.size / 2;
                    var px0 = e.x + pr - s.x, py0 = e.y + pr - s.y;
                    if (px0 * px0 + py0 * py0 < pr * pr && s.hitList.indexOf(e) === -1) {
                        s.hitList.push(e);
                        hitEnemy(e, dmgFor(s.mult), 0, 0);
                        burst(s.x, s.y, s.color, 3, 60);
                    }
                }
                continue;
            }
            for (var j = 0; j < state.enemies.length; j++) {
                e = state.enemies[j];
                var er = e.size / 2;
                var ex = e.x + er - s.x, ey = e.y + er - s.y;
                if (ex * ex + ey * ey < er * er) {
                    if (s.kind === 'orb') {
                        burst(s.x, s.y, s.color, 14, 130);
                        state.shake = Math.max(state.shake, 3);
                        for (var k = 0; k < state.enemies.length; k++) {
                            var o = state.enemies[k];
                            var ox = o.x + o.size / 2 - s.x, oy = o.y + o.size / 2 - s.y;
                            if (ox * ox + oy * oy < s.blast * s.blast) {
                                hitEnemy(o, dmgFor(s.mult), 0, 0);
                            }
                        }
                    } else {
                        hitEnemy(e, dmgFor(s.mult), 0, 0);
                        burst(s.x, s.y, s.color, 4, 70);
                    }
                    state.shots.splice(i, 1);
                    break;
                }
            }
        }

        for (i = state.enemyShots.length - 1; i >= 0; i--) {
            s = state.enemyShots[i];
            s.x += s.vx * dt;
            s.y += s.vy * dt;
            if (s.x < -10 || s.x > W + 10 || s.y < -10 || s.y > H + 10) {
                state.enemyShots.splice(i, 1);
                continue;
            }
            var px = p.x + SIZE / 2 - s.x, py = p.y + SIZE / 2 - s.y;
            if (px * px + py * py < 14 * 14) {
                state.enemyShots.splice(i, 1);
                hurtPlayer(10);
            }
        }

        for (i = state.enemies.length - 1; i >= 0; i--) {
            e = state.enemies[i];
            e.flash = Math.max(0, e.flash - dt);
            if (e.hp <= 0) {
                killEnemy(i);
                continue;
            }
            updateEnemy(e, dt);
        }

        for (i = state.pickups.length - 1; i >= 0; i--) {
            var pk = state.pickups[i];
            pk.ttl -= dt;
            if (pk.ttl <= 0) {
                state.pickups.splice(i, 1);
                continue;
            }
            var cx = p.x + SIZE / 2 - (pk.x + 12), cy = p.y + SIZE / 2 - (pk.y + 12);
            if (cx * cx + cy * cy < 24 * 24) {
                if (pk.type === 'heart') {
                    p.hp = Math.min(p.maxHp, p.hp + 25);
                    floatText(p.x + SIZE / 2, p.y, '+25', '#5cb85c');
                } else if (pk.type === 'star') {
                    state.buffs.star = 8;
                    floatText(p.x + SIZE / 2, p.y, '2X DMG', '#FFD700');
                } else {
                    state.buffs.bolt = 8;
                    floatText(p.x + SIZE / 2, p.y, 'RAPID', '#00add8');
                }
                burst(pk.x + 12, pk.y + 12, '#FFD700', 8, 80);
                state.pickups.splice(i, 1);
            }
        }

        for (i = state.particles.length - 1; i >= 0; i--) {
            var pt = state.particles[i];
            pt.life -= dt;
            if (pt.life <= 0) {
                state.particles.splice(i, 1);
                continue;
            }
            pt.x += pt.vx * dt;
            pt.y += pt.vy * dt;
            pt.vx *= 0.92;
            pt.vy *= 0.92;
        }

        for (i = state.texts.length - 1; i >= 0; i--) {
            var tx = state.texts[i];
            tx.life -= dt;
            tx.y -= 24 * dt;
            if (tx.life <= 0) state.texts.splice(i, 1);
        }

        if (state.slash) {
            state.slash.t -= dt;
            if (state.slash.t <= 0) state.slash = null;
        }

        state.shake = Math.max(0, state.shake - dt * 14);
        state.banner = Math.max(0, state.banner - dt);
        if (mode === 'play' && state.enemies.length === 0 && state.spawnQueue.length === 0 && state.banner <= 0 && state.wave < MAX_WAVE) spawnWave();

        if (state.pending) {
            state.pending.t -= dt;
            if (state.pending.t <= 0) {
                mode = state.pending.mode;
                showOverlay(state.pending.title, state.pending.sub);
                state.pending = null;
            }
        }
    }

    function drawBar(x, y, w, h, ratio, color) {
        ctx.fillStyle = '#180830';
        ctx.fillRect(x, y, w, h);
        ctx.fillStyle = color;
        ctx.fillRect(x, y, Math.max(0, w * ratio), h);
        ctx.strokeStyle = '#3A1A7A';
        ctx.strokeRect(x + 0.5, y + 0.5, w - 1, h - 1);
    }

    function render() {
        ctx.setTransform(1, 0, 0, 1, 0, 0);
        ctx.fillStyle = '#050010';
        ctx.fillRect(0, 0, W, H);

        if (state.shake > 0) {
            ctx.translate((Math.random() - 0.5) * state.shake, (Math.random() - 0.5) * state.shake);
        }

        ctx.strokeStyle = '#3A1A7A';
        ctx.strokeRect(1.5, 1.5, W - 3, H - 3);

        var p = state.player;
        var i;

        for (i = 0; i < state.pickups.length; i++) {
            var pk = state.pickups[i];
            if (pk.ttl < 2 && Math.floor(pk.ttl * 8) % 2 === 0) continue;
            var bob = Math.sin(Date.now() / 250 + pk.x) * 3;
            ctx.drawImage(sprites[pk.type], pk.x, pk.y + bob, 24, 24);
        }

        for (i = 0; i < state.enemies.length; i++) {
            var e = state.enemies[i];
            ctx.globalAlpha = (e.type === 'wraith' ? 0.85 * e.fade : 1) * (e.flash > 0 ? 0.5 : 1);
            ctx.drawImage(sprites[e.type], e.x, e.y, e.size, e.size);
            ctx.globalAlpha = 1;
            if (e.type !== 'boss' && e.hp < e.maxHp) drawBar(e.x, e.y - 6, e.size, 3, e.hp / e.maxHp, enemyTypes[e.type].color);
        }

        if ((mode === 'play' || mode === 'ending') && p.hp > 0) {
            if (p.iframes > 0 && Math.floor(p.iframes * 10) % 2 === 0) ctx.globalAlpha = 0.4;
            var frame = state.moving && Math.floor(state.animTime * 8) % 2 === 1 ? playerImgB : playerImg;
            var bob = state.moving ? 0 : Math.sin(Date.now() / 350) * 1.5;
            var aim0 = aimDir();
            ctx.drawImage(frame, p.x - aim0.x * p.recoil, p.y + bob - aim0.y * p.recoil, SIZE, SIZE);
            ctx.globalAlpha = 1;
        }

        if (state.slash) {
            ctx.strokeStyle = state.slash.color;
            ctx.lineWidth = state.slash.width;
            ctx.globalAlpha = Math.min(1, state.slash.t * 8);
            ctx.beginPath();
            if (state.slash.full) {
                ctx.arc(state.slash.x, state.slash.y, state.slash.r, 0, Math.PI * 2);
            } else {
                ctx.arc(state.slash.x, state.slash.y, state.slash.r, state.slash.ang - state.slash.arc, state.slash.ang + state.slash.arc);
            }
            ctx.stroke();
            ctx.lineWidth = 1;
            ctx.globalAlpha = 1;
        }

        if (mode === 'play') {
            var dir = aimDir();
            var ax = p.x + SIZE / 2 + dir.x * 30;
            var ay = p.y + SIZE / 2 + dir.y * 30;
            var ang = Math.atan2(dir.y, dir.x);
            ctx.fillStyle = wpnColor;
            ctx.beginPath();
            ctx.moveTo(ax + Math.cos(ang) * 7, ay + Math.sin(ang) * 7);
            ctx.lineTo(ax + Math.cos(ang + 2.5) * 5, ay + Math.sin(ang + 2.5) * 5);
            ctx.lineTo(ax + Math.cos(ang - 2.5) * 5, ay + Math.sin(ang - 2.5) * 5);
            ctx.fill();
        }

        for (i = 0; i < state.shots.length; i++) {
            var sh = state.shots[i];
            ctx.fillStyle = sh.color;
            if (sh.kind === 'orb') {
                var pulse = sh.radius + Math.sin(Date.now() / 60) * 1.5;
                ctx.beginPath();
                ctx.arc(sh.x, sh.y, pulse, 0, Math.PI * 2);
                ctx.fill();
                ctx.strokeStyle = '#ffffff';
                ctx.globalAlpha = 0.5;
                ctx.beginPath();
                ctx.arc(sh.x, sh.y, pulse + 2, 0, Math.PI * 2);
                ctx.stroke();
                ctx.globalAlpha = 1;
            } else if (sh.kind === 'pierce') {
                var sl = Math.sqrt(sh.vx * sh.vx + sh.vy * sh.vy) || 1;
                ctx.strokeStyle = sh.color;
                ctx.lineWidth = sh.width;
                ctx.beginPath();
                ctx.moveTo(sh.x - sh.vx / sl * 14, sh.y - sh.vy / sl * 14);
                ctx.lineTo(sh.x + sh.vx / sl * 5, sh.y + sh.vy / sl * 5);
                ctx.stroke();
                ctx.lineWidth = 1;
            } else {
                var kl = Math.sqrt(sh.vx * sh.vx + sh.vy * sh.vy) || 1;
                ctx.strokeStyle = sh.color;
                ctx.lineWidth = 2;
                ctx.beginPath();
                ctx.moveTo(sh.x - sh.vx / kl * 6, sh.y - sh.vy / kl * 6);
                ctx.lineTo(sh.x + sh.vx / kl * 3, sh.y + sh.vy / kl * 3);
                ctx.stroke();
                ctx.lineWidth = 1;
            }
        }
        ctx.fillStyle = '#e05d44';
        for (i = 0; i < state.enemyShots.length; i++) {
            ctx.fillRect(state.enemyShots[i].x - 3, state.enemyShots[i].y - 3, 6, 6);
        }

        for (i = 0; i < state.particles.length; i++) {
            var pt = state.particles[i];
            ctx.globalAlpha = Math.min(1, pt.life * 3);
            ctx.fillStyle = pt.color;
            ctx.fillRect(pt.x - pt.size / 2, pt.y - pt.size / 2, pt.size, pt.size);
        }
        ctx.globalAlpha = 1;

        ctx.font = 'bold 9px monospace';
        for (i = 0; i < state.texts.length; i++) {
            var tx = state.texts[i];
            ctx.globalAlpha = Math.min(1, tx.life * 2);
            ctx.fillStyle = tx.color;
            ctx.textAlign = 'center';
            ctx.fillText(tx.s, tx.x, tx.y);
        }
        ctx.globalAlpha = 1;
        ctx.textAlign = 'left';

        drawBar(10, 10, 120, 8, p.hp / p.maxHp, accent);
        ctx.fillStyle = '#AA88DD';
        ctx.font = '8px monospace';
        ctx.fillText('HP ' + Math.ceil(p.hp) + '/' + p.maxHp, 10, 28);
        ctx.textAlign = 'right';
        ctx.fillText('WAVE ' + state.wave + '/' + MAX_WAVE + '  KILLS ' + state.kills, W - 10, 18);
        ctx.textAlign = 'left';

        var bossE = null;
        for (i = 0; i < state.enemies.length; i++) {
            if (state.enemies[i].type === 'boss') bossE = state.enemies[i];
        }
        if (bossE) {
            ctx.fillStyle = '#FFD700';
            ctx.textAlign = 'center';
            ctx.fillText('THE WARDEN', W / 2, 26);
            ctx.textAlign = 'left';
            drawBar(W / 2 - 100, 30, 200, 8, bossE.hp / bossE.maxHp, '#e05d44');
        }

        var bx = 140;
        if (state.buffs.star > 0) {
            ctx.drawImage(sprites.star, bx, 6, 14, 14);
            ctx.fillStyle = '#FFD700';
            ctx.fillText(Math.ceil(state.buffs.star), bx + 16, 17);
            bx += 34;
        }
        if (state.buffs.bolt > 0) {
            ctx.drawImage(sprites.bolt, bx, 6, 14, 14);
            ctx.fillStyle = '#00add8';
            ctx.fillText(Math.ceil(state.buffs.bolt), bx + 16, 17);
        }

        if (state.banner > 0 && mode === 'play') {
            ctx.fillStyle = accent;
            ctx.font = 'bold 16px monospace';
            ctx.textAlign = 'center';
            ctx.fillText(state.wave >= MAX_WAVE ? 'FINAL WAVE' : 'WAVE ' + state.wave, W / 2, H / 2 - 40);
            ctx.textAlign = 'left';
        }

    }

    function canvasPos(evt) {
        var r = canvas.getBoundingClientRect();
        return {
            x: (evt.clientX - r.left) * (W / r.width),
            y: (evt.clientY - r.top) * (H / r.height)
        };
    }

    canvas.addEventListener('mousemove', function (e) {
        var pos = canvasPos(e);
        mouse.x = pos.x;
        mouse.y = pos.y;
    });
    canvas.addEventListener('mousedown', function (e) {
        if (e.button !== 0) return;
        e.preventDefault();
        if (mode === 'play') shoot();
    });
    canvas.addEventListener('contextmenu', function (e) { e.preventDefault(); });

    retryBtn.addEventListener('click', function () {
        if (mode === 'pause') resumeGame();
        else restart();
    });
    menuBtn.addEventListener('click', toMenu);
    startBtn.addEventListener('click', restart);
    prevClassBtn.addEventListener('click', function () { cycleClass(-1); });
    nextClassBtn.addEventListener('click', function () { cycleClass(1); });
    instructionsBtn.addEventListener('click', function () {
        instructions.style.display = instructions.style.display === 'none' ? 'block' : 'none';
    });

    var gameKeys = ['ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown', 'w', 'a', 's', 'd', 'r', 'Escape'];

    function onKeyDown(e) {
        if (!canvas.isConnected) {
            document.removeEventListener('keydown', onKeyDown);
            document.removeEventListener('keyup', onKeyUp);
            return;
        }
        var k = e.key.length === 1 ? e.key.toLowerCase() : e.key;
        if (gameKeys.indexOf(k) === -1) return;
        e.preventDefault();
        if (mode === 'menu') {
            if (k === 'ArrowLeft' || k === 'a') cycleClass(-1);
            if (k === 'ArrowRight' || k === 'd') cycleClass(1);
            return;
        }
        if (k === 'Escape') {
            if (mode === 'play') pauseGame();
            else if (mode === 'pause') resumeGame();
            return;
        }
        keys[k] = true;
        if (k === 'r' && (mode === 'over' || mode === 'win')) restart();
    }

    function onKeyUp(e) {
        var k = e.key.length === 1 ? e.key.toLowerCase() : e.key;
        keys[k] = false;
    }

    document.addEventListener('keydown', onKeyDown);
    document.addEventListener('keyup', onKeyUp);

    var last = 0;
    function loop(ts) {
        if (!canvas.isConnected) return;
        var dt = Math.min((ts - last) / 1000, 0.05);
        last = ts;
        update(dt);
        render();
        requestAnimationFrame(loop);
    }

    state = newState();
    setClass(selectedClass);
    requestAnimationFrame(function (ts) {
        last = ts;
        requestAnimationFrame(loop);
    });
})();
