function previewCard() {
    var username = document.getElementById('username-input').value.trim();
    if (!username) return;
    var origin = window.location.origin;
    document.getElementById('card-error').style.display = 'none';
    document.getElementById('card-full').src = '/card/' + username + '.svg?t=' + Date.now();
    document.getElementById('snippet-full').textContent = '![GitHub RPG](' + origin + '/card/' + username + '.svg)';
    document.getElementById('profile-link').href = '/u/' + username;
    document.getElementById('preview-section').style.display = 'block';
}

function copySnippet(snippetId, btnId) {
    var text = document.getElementById(snippetId).textContent;
    navigator.clipboard.writeText(text).then(function() {
        var btn = document.getElementById(btnId);
        btn.textContent = 'COPIED!';
        btn.style.borderColor = 'var(--green)';
        btn.style.color = 'var(--green)';
        setTimeout(function() {
            btn.innerHTML = '<i aria-hidden="true" data-lucide="copy" style="width:11px;height:11px;stroke:var(--cyan);stroke-width:2;"></i>COPY';
            btn.style.borderColor = '';
            btn.style.color = '';
            lucide.createIcons();
        }, 1500);
    });
}

document.getElementById('preview-card-btn').addEventListener('click', previewCard);
document.getElementById('copy-full-btn').addEventListener('click', function() {
    copySnippet('snippet-full', 'copy-full-btn');
});

document.getElementById('username-input').addEventListener('keydown', function(e) {
    if (e.key === 'Enter') previewCard();
});

var cardImg = document.getElementById('card-full');
cardImg.addEventListener('error', function() {
    document.getElementById('preview-section').style.display = 'none';
    document.getElementById('card-error').style.display = 'flex';
});
cardImg.addEventListener('load', function() {
    document.getElementById('card-error').style.display = 'none';
});
