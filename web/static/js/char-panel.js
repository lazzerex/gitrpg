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
    document.querySelectorAll('[data-stat]').forEach(function(el) {
        counter(el, parseInt(el.dataset.stat, 10), 900);
    });
    var bar = document.getElementById('xp-bar');
    if (bar) { bar.style.width = '0%'; void bar.offsetWidth; requestAnimationFrame(function() { bar.style.width = bar.dataset.width + '%'; }); }
})();
