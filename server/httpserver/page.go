package httpserver

const homePage = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>FlipOff · Random Teleporter</title>
  <style>
    :root {
      color-scheme: dark;
      --bg: #09090b;
      --card-bg: rgba(24, 24, 27, 0.6);
      --card-border: rgba(255, 255, 255, 0.08);
      --text-main: #fafafa;
      --text-muted: #a1a1aa;
      --accent-glow: rgba(124, 58, 237, 0.4);
      --btn-bg: #fafafa;
      --btn-text: #09090b;
      --btn-hover: #e4e4e7;
    }

    * {
      box-sizing: border-box;
      margin: 0;
      padding: 0;
    }

    body {
      min-height: 100vh;
      font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
      color: var(--text-main);
      background-color: var(--bg);
      background-image: 
        radial-gradient(circle at 15% 50%, rgba(124, 58, 237, 0.08), transparent 25%),
        radial-gradient(circle at 85% 30%, rgba(56, 189, 248, 0.08), transparent 25%);
      display: flex;
      align-items: center;
      justify-content: center;
      overflow: hidden;
      padding: 24px;
    }

    /* 背景网格效果 */
    .grid-bg {
      position: absolute;
      inset: 0;
      background-size: 40px 40px;
      background-image: linear-gradient(to right, rgba(255, 255, 255, 0.02) 1px, transparent 1px),
                        linear-gradient(to bottom, rgba(255, 255, 255, 0.02) 1px, transparent 1px);
      mask-image: radial-gradient(circle at center, black, transparent 80%);
      -webkit-mask-image: radial-gradient(circle at center, black, transparent 80%);
      z-index: -1;
    }

    .container {
      position: relative;
      width: 100%;
      max-width: 600px;
      z-index: 1;
    }

    .card {
      background: var(--card-bg);
      border: 1px solid var(--card-border);
      border-radius: 24px;
      padding: 48px 40px;
      backdrop-filter: blur(20px);
      -webkit-backdrop-filter: blur(20px);
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4), inset 0 1px 0 rgba(255, 255, 255, 0.05);
      text-align: center;
      transition: transform 0.3s ease, box-shadow 0.3s ease;
    }

    .card:hover {
      transform: translateY(-2px);
      box-shadow: 0 30px 60px rgba(0, 0, 0, 0.5), 0 0 40px var(--accent-glow), inset 0 1px 0 rgba(255, 255, 255, 0.05);
    }

    .badge {
      display: inline-block;
      padding: 6px 12px;
      background: rgba(124, 58, 237, 0.15);
      border: 1px solid rgba(124, 58, 237, 0.3);
      border-radius: 99px;
      color: #c4b5fd;
      font-size: 13px;
      font-weight: 500;
      letter-spacing: 0.05em;
      margin-bottom: 24px;
      text-transform: uppercase;
    }

    h1 {
      font-size: 42px;
      font-weight: 800;
      letter-spacing: -0.03em;
      line-height: 1.1;
      margin-bottom: 20px;
      background: linear-gradient(135deg, #fff 0%, #a1a1aa 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    p.desc {
      font-size: 16px;
      color: var(--text-muted);
      line-height: 1.6;
      margin-bottom: 40px;
      max-width: 480px;
      margin-inline: auto;
    }

    .btn-container {
      position: relative;
      display: inline-block;
    }

    .btn-glow {
      position: absolute;
      inset: -4px;
      background: linear-gradient(90deg, #7c3aed, #38bdf8, #7c3aed);
      background-size: 200% auto;
      border-radius: 18px;
      filter: blur(12px);
      opacity: 0;
      transition: opacity 0.3s ease;
      animation: gradient-shift 3s linear infinite;
      z-index: -1;
    }

    .btn-container:hover .btn-glow {
      opacity: 0.8;
    }

    .trigger-btn {
      position: relative;
      appearance: none;
      background: var(--btn-bg);
      color: var(--btn-text);
      border: none;
      border-radius: 14px;
      padding: 0 40px;
      height: 56px;
      font-size: 17px;
      font-weight: 600;
      cursor: pointer;
      display: inline-flex;
      align-items: center;
      gap: 10px;
      transition: all 0.2s ease;
      box-shadow: 0 4px 14px rgba(255, 255, 255, 0.1);
    }

    .trigger-btn:hover {
      background: var(--btn-hover);
      transform: scale(1.02);
    }

    .trigger-btn:active {
      transform: scale(0.98);
    }

    .trigger-btn svg {
      width: 20px;
      height: 20px;
      transition: transform 0.3s ease;
    }

    .trigger-btn:hover svg {
      transform: translateX(4px);
    }

    .features {
      display: flex;
      justify-content: center;
      gap: 24px;
      margin-top: 36px;
      border-top: 1px solid var(--card-border);
      padding-top: 24px;
    }

    .feature-item {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 13px;
      color: var(--text-muted);
    }

    .feature-item svg {
      width: 14px;
      height: 14px;
      color: #7c3aed;
    }

    /* 加载遮罩层 */
    .overlay {
      position: fixed;
      inset: 0;
      background: rgba(9, 9, 11, 0.85);
      backdrop-filter: blur(12px);
      -webkit-backdrop-filter: blur(12px);
      z-index: 100;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      opacity: 0;
      pointer-events: none;
      transition: opacity 0.4s ease;
    }

    .overlay.active {
      opacity: 1;
      pointer-events: all;
    }

    .loader {
      position: relative;
      width: 80px;
      height: 80px;
      margin-bottom: 32px;
    }

    .loader-ring {
      position: absolute;
      inset: 0;
      border-radius: 50%;
      border: 3px solid transparent;
      border-top-color: #7c3aed;
      border-right-color: #38bdf8;
      animation: spin 1s cubic-bezier(0.68, -0.55, 0.265, 1.55) infinite;
    }

    .loader-ring:nth-child(2) {
      inset: 10px;
      border-top-color: transparent;
      border-right-color: transparent;
      border-bottom-color: #38bdf8;
      border-left-color: #7c3aed;
      animation-direction: reverse;
      animation-duration: 1.5s;
    }

    .overlay-title {
      font-size: 24px;
      font-weight: 700;
      margin-bottom: 12px;
      background: linear-gradient(90deg, #fff, #a1a1aa);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }

    .overlay-desc {
      font-size: 15px;
      color: var(--text-muted);
      max-width: 400px;
      text-align: center;
      line-height: 1.5;
    }

    .target-info {
      margin-top: 24px;
      padding: 16px 24px;
      background: rgba(255, 255, 255, 0.03);
      border: 1px solid rgba(255, 255, 255, 0.08);
      border-radius: 12px;
      opacity: 0;
      transform: translateY(10px);
      transition: all 0.4s ease;
    }

    .target-info.visible {
      opacity: 1;
      transform: translateY(0);
    }

    .target-tag {
      font-size: 12px;
      color: #38bdf8;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      margin-bottom: 4px;
      display: block;
    }

    .target-title {
      font-size: 18px;
      color: #fff;
      font-weight: 600;
    }

    .overlay-actions {
      margin-top: 18px;
      display: flex;
      justify-content: center;
    }

    .retry-btn {
      appearance: none;
      border: 1px solid rgba(255, 255, 255, 0.12);
      background: rgba(255, 255, 255, 0.06);
      color: var(--text-main);
      border-radius: 999px;
      padding: 10px 16px;
      font-size: 13px;
      font-weight: 600;
      cursor: pointer;
      transition: background 0.2s ease, transform 0.2s ease, border-color 0.2s ease;
    }

    .retry-btn:hover {
      background: rgba(255, 255, 255, 0.1);
      border-color: rgba(255, 255, 255, 0.24);
      transform: translateY(-1px);
    }

    .retry-btn:active {
      transform: translateY(0);
    }

    @keyframes spin {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
    }

    @keyframes gradient-shift {
      0% { background-position: 0% 50%; }
      100% { background-position: 200% 50%; }
    }

    @media (max-width: 640px) {
      .card {
        padding: 40px 24px;
      }
      h1 {
        font-size: 32px;
      }
      .features {
        flex-direction: column;
        gap: 12px;
        align-items: center;
      }
    }
  </style>
</head>
<body>
  <div class="grid-bg"></div>
  
  <div class="container">
    <div class="card">
      <div class="badge">Chaos Engine</div>
      <h1>Wuck start your journal</h1>
      <p class="desc">
        bored?  click the button below and find something amuse you
		Simple. Weird. Fun.
		Colorless and transparent
      </p>
      
      <div class="btn-container">
        <div class="btn-glow"></div>
        <button id="trigger" class="trigger-btn">
          <span>Start Random Jump</span>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="5" y1="12" x2="19" y2="12"></line>
            <polyline points="12 5 19 12 12 19"></polyline>
          </svg>
        </button>
      </div>

      <div class="features">
        <div class="feature-item">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          Live Trending Fetch
        </div>
        <div class="feature-item">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>
          Multi-source Aggregation
        </div>
        <div class="feature-item">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>
          Seamless Redirection
        </div>
      </div>
    </div>
  </div>

  <div id="overlay" class="overlay">
    <div class="loader">
      <div class="loader-ring"></div>
      <div class="loader-ring"></div>
    </div>
    <div class="overlay-title" id="overlayTitle">Summoning Chaos</div>
    <div class="overlay-desc" id="overlayDesc">Rolling the dice across the internet, please wait...</div>
    
    <div id="targetInfo" class="target-info">
      <span class="target-tag" id="targetTag">SOURCE</span>
      <div class="target-title" id="targetTitle">Loading...</div>
    </div>
    <div class="overlay-actions">
      <button id="retry" class="retry-btn" type="button">Retry another jump</button>
    </div>
  </div>

  <script>
    const overlay = document.getElementById('overlay');
    const overlayTitle = document.getElementById('overlayTitle');
    const overlayDesc = document.getElementById('overlayDesc');
    const targetInfo = document.getElementById('targetInfo');
    const targetTag = document.getElementById('targetTag');
    const targetTitle = document.getElementById('targetTitle');
    const trigger = document.getElementById('trigger');
    const retry = document.getElementById('retry');

    let isLoading = false;
    let audioCtx;
    let jumpController;
    let pendingJumpTimer;
    let pendingRedirectTimer;

    function playTeleportSound() {
      if (!window.AudioContext && !window.webkitAudioContext) return;
      if (!audioCtx) {
        audioCtx = new (window.AudioContext || window.webkitAudioContext)();
      }
      if (audioCtx.state === 'suspended') audioCtx.resume();

      const osc = audioCtx.createOscillator();
      const gain = audioCtx.createGain();
      
      osc.type = 'sine';
      // 频率从低到高滑音
      osc.frequency.setValueAtTime(150, audioCtx.currentTime);
      osc.frequency.exponentialRampToValueAtTime(800, audioCtx.currentTime + 0.5);
      
      // 音量淡入淡出
      gain.gain.setValueAtTime(0, audioCtx.currentTime);
      gain.gain.linearRampToValueAtTime(0.1, audioCtx.currentTime + 0.1);
      gain.gain.exponentialRampToValueAtTime(0.01, audioCtx.currentTime + 0.6);
      
      osc.connect(gain);
      gain.connect(audioCtx.destination);
      
      osc.start();
      osc.stop(audioCtx.currentTime + 0.6);
    }

    // 成功锁定目标的提示音
    function playSuccessSound() {
      if (!audioCtx) return;
      const osc = audioCtx.createOscillator();
      const gain = audioCtx.createGain();
      
      osc.type = 'triangle';
      osc.frequency.setValueAtTime(600, audioCtx.currentTime);
      osc.frequency.setValueAtTime(800, audioCtx.currentTime + 0.1);
      
      gain.gain.setValueAtTime(0, audioCtx.currentTime);
      gain.gain.linearRampToValueAtTime(0.1, audioCtx.currentTime + 0.05);
      gain.gain.linearRampToValueAtTime(0, audioCtx.currentTime + 0.3);
      
      osc.connect(gain);
      gain.connect(audioCtx.destination);
      
      osc.start();
      osc.stop(audioCtx.currentTime + 0.3);
    }

    function clearJumpTimers() {
      if (pendingJumpTimer) {
        window.clearTimeout(pendingJumpTimer);
        pendingJumpTimer = null;
      }

      if (pendingRedirectTimer) {
        window.clearTimeout(pendingRedirectTimer);
        pendingRedirectTimer = null;
      }
    }

    function resetOverlayState() {
      clearJumpTimers();
      if (jumpController) {
        jumpController.abort();
        jumpController = null;
      }

      isLoading = false;
      overlay.classList.remove('active');
      targetInfo.classList.remove('visible');
      overlayTitle.textContent = 'Summoning Chaos';
      overlayDesc.textContent = 'Rolling the dice across the internet, please wait...';
      targetTag.textContent = 'SOURCE';
      targetTitle.textContent = 'Loading...';
    }

    async function startJump(forceRetry) {
      if (isLoading && !forceRetry) return;

      clearJumpTimers();
      if (jumpController) {
        jumpController.abort();
      }

      jumpController = new AbortController();
      isLoading = true;

      targetInfo.classList.remove('visible');
      overlayTitle.textContent = 'Navigating the Web...';
      overlayDesc.textContent = 'RNG is picking a target from live trending sources';
      overlay.classList.add('active');

      playTeleportSound();

      try {
        const res = await fetch('/api/biu', {
          headers: { 'Accept': 'application/json' },
          signal: jumpController.signal,
        });
        const payload = await res.json();

        if (!res.ok || !payload.data || !payload.data.url) {
          throw new Error(payload.msg || 'API returned an invalid target');
        }

        pendingJumpTimer = window.setTimeout(() => {
          playSuccessSound();
          overlayTitle.textContent = 'Target Locked';
          overlayDesc.textContent = 'Portal is open, initiating jump sequence';

          targetTag.textContent = payload.data.source + ' · ' + payload.data.category;
          targetTitle.textContent = payload.data.title;
          targetInfo.classList.add('visible');

          pendingRedirectTimer = window.setTimeout(() => {
            window.location.href = payload.data.url;
          }, 700);
        }, 250);

      } catch (err) {
        if (err && err.name === 'AbortError') {
          return;
        }

        setTimeout(() => {
          overlayTitle.textContent = 'Jump Failed';
          overlayDesc.textContent = err.message || 'Network error or service unavailable, please try again';

          isLoading = false;
        }, 200);
      } finally {
        jumpController = null;
      }
    }

    window.addEventListener('pageshow', resetOverlayState);
    window.addEventListener('pagehide', clearJumpTimers);
    trigger.addEventListener('click', startJump);
    retry.addEventListener('click', () => startJump(true));
  </script>
</body>
</html>`
