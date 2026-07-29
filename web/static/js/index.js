var classDescs = {
    guardian:  'RELIABLE DEFENDER. GO ENGINEERS BUILD SCALABLE SYSTEMS WITH PRECISION AND DISCIPLINE.',
    berserker: 'UNSTOPPABLE FORCE. RUST WARRIORS CONQUER MEMORY AND PERFORMANCE FEARLESSLY.',
    paladin:   'CODE WITH HONOR. TYPESCRIPT PALADINS BRING ORDER TO DYNAMIC SYSTEMS.',
    rogue:     'MOVE FAST. JAVASCRIPT ROGUES SHIP FEATURES AT THE SPEED OF THOUGHT.',
    sage:      'DATA SORCERER. PYTHON SAGES WIELD ALGORITHMS AND MACHINE LEARNING AS SPELLS.',
    knight:    'ENTERPRISE WARRIOR. C# KNIGHTS UPHOLD THE REALM WITH STRUCTURED MIGHT.'
};

function selectClass(el) {
    document.querySelectorAll('.class-card').forEach(function(c) {
        c.style.borderColor = 'var(--border)';
        c.style.boxShadow = '';
    });
    el.style.borderColor = 'var(--gold)';
    el.style.boxShadow = '4px 4px 0 #000';
    var desc = document.getElementById('class-desc');
    document.getElementById('class-desc-text').textContent = classDescs[el.dataset.class] || '';
    desc.style.display = 'block';
}

document.querySelectorAll('.class-card').forEach(function(el) {
    el.addEventListener('click', function() { selectClass(el); });
});

(function () {
    function counter(el, target, ms) {
        var s = performance.now();
        function tick(now) {
            var p = Math.min((now - s) / ms, 1);
            el.textContent = Math.floor(p * target);
            if (p < 1) requestAnimationFrame(tick);
            else el.textContent = target;
        }
        requestAnimationFrame(tick);
    }

    document.querySelectorAll('[data-demo]').forEach(function(el) {
        counter(el, parseInt(el.dataset.demo, 10), 1400);
    });

    var bar = document.getElementById('demo-xp');
    if (bar) requestAnimationFrame(function() { requestAnimationFrame(function() { bar.style.width = '73%'; }); });
})();

function rollStats() {
    var els = document.querySelectorAll('[data-demo]');
    els.forEach(function(el) {
        var target = Math.floor(Math.random() * 60) + 30;
        var start = performance.now();
        function tick(now) {
            var p = Math.min((now - start) / 600, 1);
            el.textContent = Math.floor(p * target);
            if (p < 1) requestAnimationFrame(tick);
            else el.textContent = target;
        }
        requestAnimationFrame(tick);
    });
}

var rollStatsBtn = document.getElementById('roll-stats-btn');
if (rollStatsBtn) rollStatsBtn.addEventListener('click', rollStats);
