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

    .recommendations {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 16px;
      width: 100%;
      max-width: 720px;
      margin-top: 32px;
      opacity: 0;
      transform: translateY(12px);
      transition: opacity 0.5s ease, transform 0.5s ease;
      perspective: 1000px;
    }

    .recommendations.visible {
      opacity: 1;
      transform: translateY(0);
    }

    .rec-card {
      position: relative;
      background: rgba(24, 24, 27, 0.5);
      border: 1px solid rgba(255, 255, 255, 0.06);
      border-radius: 16px;
      overflow: hidden;
      text-align: left;
      cursor: pointer;
      backdrop-filter: blur(12px);
      -webkit-backdrop-filter: blur(12px);
      transition: box-shadow 0.3s ease, border-color 0.3s ease, background 0.3s ease;
      opacity: 0;
      transform: translateY(30px) scale(0.9) rotateX(10deg);
      will-change: transform, opacity;
    }

    .rec-card.show {
      animation: card-spring 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
    }

    @keyframes card-spring {
      0% {
        opacity: 0;
        transform: translateY(30px) scale(0.9) rotateX(10deg);
      }
      60% {
        opacity: 1;
        transform: translateY(-4px) scale(1.02) rotateX(0deg);
      }
      80% {
        transform: translateY(2px) scale(0.99);
      }
      100% {
        opacity: 1;
        transform: translateY(0) scale(1) rotateX(0deg);
      }
    }

    .rec-card:hover {
      background: rgba(24, 24, 27, 0.8);
      border-color: rgba(124, 58, 237, 0.4);
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.4), 0 0 24px rgba(124, 58, 237, 0.15);
      z-index: 2;
    }

    .rec-card:active {
      transform: translateY(0) scale(0.96) !important;
    }

    /* Thumbnail area */
    .rec-thumb {
      position: relative;
      width: 100%;
      height: 130px;
      border-radius: 12px 12px 0 0;
      overflow: hidden;
      background: rgba(255, 255, 255, 0.03);
    }

    .rec-thumb-img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      display: none;
      transition: transform 0.4s ease;
    }

    .rec-thumb-img.loaded {
      display: block;
    }

    .rec-card:hover .rec-thumb-img.loaded {
      transform: scale(1.08);
    }

    .rec-thumb-fallback {
      position: absolute;
      inset: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: opacity 0.3s ease;
    }

    .rec-thumb-img.loaded ~ .rec-thumb-fallback {
      opacity: 0;
    }

    .rec-thumb-emoji {
      font-size: 40px;
      filter: grayscale(0.3);
      transition: transform 0.3s ease;
    }

    .rec-card:hover .rec-thumb-emoji {
      transform: scale(1.2) rotate(10deg);
    }

    /* Category-specific fallback gradients */
    .rec-card[data-category="novel"] .rec-thumb-fallback { background: linear-gradient(135deg, rgba(236, 72, 153, 0.15), rgba(236, 72, 153, 0.05)); }
    .rec-card[data-category="video"] .rec-thumb-fallback { background: linear-gradient(135deg, rgba(59, 130, 246, 0.15), rgba(59, 130, 246, 0.05)); }
    .rec-card[data-category="blog"] .rec-thumb-fallback { background: linear-gradient(135deg, rgba(16, 185, 129, 0.15), rgba(16, 185, 129, 0.05)); }
    .rec-card[data-category="news"] .rec-thumb-fallback { background: linear-gradient(135deg, rgba(245, 158, 11, 0.15), rgba(245, 158, 11, 0.05)); }
    .rec-card[data-category="lifestyle"] .rec-thumb-fallback { background: linear-gradient(135deg, rgba(244, 114, 182, 0.15), rgba(244, 114, 182, 0.05)); }

    /* Glint/shine sweep */
    .rec-shine {
      position: absolute;
      top: 0;
      left: -100%;
      width: 60%;
      height: 100%;
      background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.04), transparent);
      transform: skewX(-20deg);
      transition: left 0.6s ease;
      pointer-events: none;
    }

    .rec-card.show:nth-child(1) .rec-shine {
      animation: shine-sweep 0.8s 0.35s ease-out forwards;
    }
    .rec-card.show:nth-child(2) .rec-shine {
      animation: shine-sweep 0.8s 0.5s ease-out forwards;
    }
    .rec-card.show:nth-child(3) .rec-shine {
      animation: shine-sweep 0.8s 0.65s ease-out forwards;
    }

    @keyframes shine-sweep {
      0% { left: -100%; }
      100% { left: 150%; }
    }

    .rec-body {
      padding: 14px 16px 16px;
    }

    .rec-badge {
      display: inline-block;
      padding: 4px 10px;
      background: rgba(124, 58, 237, 0.12);
      border: 1px solid rgba(124, 58, 237, 0.25);
      border-radius: 99px;
      color: #c4b5fd;
      font-size: 11px;
      font-weight: 600;
      letter-spacing: 0.04em;
      text-transform: uppercase;
      margin-bottom: 10px;
    }

    .rec-title {
      font-size: 14px;
      font-weight: 600;
      color: #fafafa;
      line-height: 1.4;
      margin-bottom: 4px;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .rec-category {
      font-size: 11px;
      color: #71717a;
      text-transform: capitalize;
      display: flex;
      align-items: center;
      gap: 4px;
    }

    .rec-category::before {
      content: '';
      display: inline-block;
      width: 6px;
      height: 6px;
      border-radius: 50%;
      background: currentColor;
    }

    .rec-card[data-category="novel"] .rec-badge { background: rgba(236, 72, 153, 0.12); border-color: rgba(236, 72, 153, 0.25); color: #f9a8d4; }
    .rec-card[data-category="novel"] .rec-category::before { background: #ec4899; }
    .rec-card[data-category="video"] .rec-badge { background: rgba(59, 130, 246, 0.12); border-color: rgba(59, 130, 246, 0.25); color: #93c5fd; }
    .rec-card[data-category="video"] .rec-category::before { background: #3b82f6; }
    .rec-card[data-category="blog"] .rec-badge { background: rgba(16, 185, 129, 0.12); border-color: rgba(16, 185, 129, 0.25); color: #6ee7b7; }
    .rec-card[data-category="blog"] .rec-category::before { background: #10b981; }
    .rec-card[data-category="news"] .rec-badge { background: rgba(245, 158, 11, 0.12); border-color: rgba(245, 158, 11, 0.25); color: #fcd34d; }
    .rec-card[data-category="news"] .rec-category::before { background: #f59e0b; }
    .rec-card[data-category="lifestyle"] .rec-badge { background: rgba(244, 114, 182, 0.12); border-color: rgba(244, 114, 182, 0.25); color: #f9a8d4; }
    .rec-card[data-category="lifestyle"] .rec-category::before { background: #f472b6; }

    /* Floating particles canvas */
    #particles {
      position: fixed;
      inset: 0;
      z-index: 101;
      pointer-events: none;
      opacity: 0;
      transition: opacity 0.6s ease;
    }
    #particles.active {
      opacity: 1;
    }

    /* Pick burst effect */
    .rec-card.picked {
      animation: card-pick-burst 0.45s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
      border-color: rgba(124, 58, 237, 0.8) !important;
      box-shadow: 0 0 48px rgba(124, 58, 237, 0.5), 0 0 96px rgba(56, 189, 248, 0.3) !important;
    }

    @keyframes card-pick-burst {
      0% { transform: scale(1); filter: brightness(1); }
      40% { transform: scale(1.08); filter: brightness(1.4); }
      100% { transform: scale(1.04); filter: brightness(1.2); }
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
      .recommendations {
        grid-template-columns: 1fr;
        gap: 12px;
        max-width: 100%;
      }
      .rec-thumb {
        height: 160px;
      }
    }

    /* ── User Area ───────────────────────────── */
    .user-area {
      position: absolute;
      top: 16px;
      right: 20px;
      z-index: 2;
    }

    .user-pill {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      padding: 6px 14px;
      background: rgba(124, 58, 237, 0.1);
      border: 1px solid rgba(124, 58, 237, 0.2);
      border-radius: 99px;
      color: #c4b5fd;
      font-size: 13px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s ease;
    }

    .user-pill:hover {
      background: rgba(124, 58, 237, 0.18);
      border-color: rgba(124, 58, 237, 0.35);
    }

    .user-pill .user-avatar {
      width: 24px;
      height: 24px;
      border-radius: 50%;
      background: linear-gradient(135deg, #7c3aed, #38bdf8);
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 11px;
      font-weight: 700;
      color: #fff;
    }

    .user-pill .logout-hint {
      font-size: 11px;
      color: #71717a;
      margin-left: 4px;
      opacity: 0;
      transition: opacity 0.2s ease;
    }

    .user-pill:hover .logout-hint {
      opacity: 1;
    }

    .login-btn {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 6px 14px;
      background: rgba(255, 255, 255, 0.06);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 99px;
      color: #a1a1aa;
      font-size: 13px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s ease;
    }

    .login-btn:hover {
      background: rgba(255, 255, 255, 0.1);
      border-color: rgba(255, 255, 255, 0.2);
      color: #fafafa;
    }

    /* ── Auth Modal ─────────────────────────── */
    .auth-modal-overlay {
      position: fixed;
      inset: 0;
      background: rgba(9, 9, 11, 0.7);
      backdrop-filter: blur(8px);
      -webkit-backdrop-filter: blur(8px);
      z-index: 200;
      display: flex;
      align-items: center;
      justify-content: center;
      opacity: 0;
      pointer-events: none;
      transition: opacity 0.3s ease;
    }

    .auth-modal-overlay.active {
      opacity: 1;
      pointer-events: all;
    }

    .auth-modal {
      background: rgba(24, 24, 27, 0.9);
      border: 1px solid rgba(255, 255, 255, 0.08);
      border-radius: 20px;
      padding: 36px 32px 28px;
      width: 100%;
      max-width: 380px;
      backdrop-filter: blur(20px);
      -webkit-backdrop-filter: blur(20px);
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.5), 0 0 40px rgba(124, 58, 237, 0.1);
      transform: translateY(12px) scale(0.96);
      transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
    }

    .auth-modal-overlay.active .auth-modal {
      transform: translateY(0) scale(1);
    }

    .auth-tabs {
      display: flex;
      gap: 0;
      margin-bottom: 24px;
      background: rgba(255, 255, 255, 0.03);
      border-radius: 10px;
      padding: 3px;
    }

    .auth-tab {
      flex: 1;
      text-align: center;
      padding: 8px 0;
      border-radius: 8px;
      font-size: 14px;
      font-weight: 500;
      color: #71717a;
      cursor: pointer;
      transition: all 0.2s ease;
      border: none;
      background: transparent;
    }

    .auth-tab.active {
      background: rgba(124, 58, 237, 0.15);
      color: #c4b5fd;
    }

    .auth-field {
      margin-bottom: 16px;
    }

    .auth-field label {
      display: block;
      font-size: 13px;
      color: #a1a1aa;
      margin-bottom: 6px;
      font-weight: 500;
    }

    .auth-field input {
      width: 100%;
      padding: 10px 14px;
      background: rgba(255, 255, 255, 0.04);
      border: 1px solid rgba(255, 255, 255, 0.08);
      border-radius: 10px;
      color: #fafafa;
      font-size: 14px;
      outline: none;
      transition: border-color 0.2s ease, background 0.2s ease;
      box-sizing: border-box;
    }

    .auth-field input:focus {
      border-color: rgba(124, 58, 237, 0.4);
      background: rgba(255, 255, 255, 0.06);
    }

    .auth-field input::placeholder {
      color: #52525b;
    }

    .auth-submit {
      width: 100%;
      padding: 11px 0;
      background: #fafafa;
      color: #09090b;
      border: none;
      border-radius: 10px;
      font-size: 14px;
      font-weight: 600;
      cursor: pointer;
      margin-top: 8px;
      transition: background 0.2s ease, transform 0.2s ease;
    }

    .auth-submit:hover {
      background: #e4e4e7;
      transform: translateY(-1px);
    }

    .auth-submit:active {
      transform: translateY(0) scale(0.98);
    }

    .auth-error {
      color: #f87171;
      font-size: 13px;
      margin-top: 12px;
      text-align: center;
      min-height: 20px;
    }

    .auth-close {
      position: absolute;
      top: 12px;
      right: 16px;
      background: none;
      border: none;
      color: #71717a;
      font-size: 20px;
      cursor: pointer;
      padding: 4px;
      line-height: 1;
      transition: color 0.2s ease;
    }

    .auth-close:hover {
      color: #fafafa;
    }

    .auth-modal {
      position: relative;
    }

    /* ── User dropdown ──────────────────────── */
    .user-dropdown {
      position: absolute;
      top: calc(100% + 6px);
      right: 0;
      background: rgba(24, 24, 27, 0.95);
      border: 1px solid rgba(255, 255, 255, 0.08);
      border-radius: 12px;
      padding: 8px 0;
      min-width: 160px;
      backdrop-filter: blur(16px);
      -webkit-backdrop-filter: blur(16px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
      opacity: 0;
      pointer-events: none;
      transform: translateY(-4px);
      transition: all 0.2s ease;
    }

    .user-dropdown.open {
      opacity: 1;
      pointer-events: all;
      transform: translateY(0);
    }

    .user-dropdown-item {
      padding: 8px 16px;
      font-size: 13px;
      color: #a1a1aa;
      cursor: pointer;
      transition: all 0.15s ease;
      display: block;
      width: 100%;
      text-align: left;
      background: none;
      border: none;
    }

    .user-dropdown-item:hover {
      background: rgba(255, 255, 255, 0.04);
      color: #fafafa;
    }

    .user-dropdown-item.danger {
      color: #f87171;
    }

    .user-dropdown-item.danger:hover {
      background: rgba(248, 113, 113, 0.08);
    }

    .user-dropdown-divider {
      height: 1px;
      background: rgba(255, 255, 255, 0.06);
      margin: 4px 0;
    }

    .user-dropdown .history-count {
      font-size: 11px;
      color: #52525b;
      padding: 4px 16px 8px;
    }
  </style>
</head>
<body>
  <div class="grid-bg"></div>
  
  <div class="container">
    <div class="card">
      <div class="badge">Chaos Engine</div>
      <div class="user-area" id="userArea">
        <button class="login-btn" id="loginBtn" onclick="showAuthModal('login')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
          Sign In
        </button>
        <div class="user-pill" id="userPill" style="display:none" onclick="toggleUserDropdown(event)">
          <div class="user-avatar" id="userAvatar">?</div>
          <span id="userName">user</span>
          <span class="logout-hint">▼</span>
        </div>
        <div class="user-dropdown" id="userDropdown">
          <div class="history-count" id="historyCount">0 clicks recorded</div>
          <div class="user-dropdown-divider"></div>
          <button class="user-dropdown-item danger" onclick="handleLogout()">Sign Out</button>
        </div>
      </div>

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

    <div id="recommendations" class="recommendations">
      <div id="recCard0" class="rec-card" onclick="pickTarget(0)">
        <div class="rec-thumb">
          <img id="recThumb0" class="rec-thumb-img" src="" alt="" />
          <div class="rec-thumb-fallback" id="recFallback0">
            <span class="rec-thumb-emoji">🎲</span>
          </div>
        </div>
        <div class="rec-body">
          <span class="rec-badge" id="recBadge0">SOURCE</span>
          <div class="rec-title" id="recTitle0">Loading...</div>
          <div class="rec-category" id="recCat0">category</div>
        </div>
        <div class="rec-shine"></div>
      </div>
      <div id="recCard1" class="rec-card" onclick="pickTarget(1)">
        <div class="rec-thumb">
          <img id="recThumb1" class="rec-thumb-img" src="" alt="" />
          <div class="rec-thumb-fallback" id="recFallback1">
            <span class="rec-thumb-emoji">🎲</span>
          </div>
        </div>
        <div class="rec-body">
          <span class="rec-badge" id="recBadge1">SOURCE</span>
          <div class="rec-title" id="recTitle1">Loading...</div>
          <div class="rec-category" id="recCat1">category</div>
        </div>
        <div class="rec-shine"></div>
      </div>
      <div id="recCard2" class="rec-card" onclick="pickTarget(2)">
        <div class="rec-thumb">
          <img id="recThumb2" class="rec-thumb-img" src="" alt="" />
          <div class="rec-thumb-fallback" id="recFallback2">
            <span class="rec-thumb-emoji">🎲</span>
          </div>
        </div>
        <div class="rec-body">
          <span class="rec-badge" id="recBadge2">SOURCE</span>
          <div class="rec-title" id="recTitle2">Loading...</div>
          <div class="rec-category" id="recCat2">category</div>
        </div>
        <div class="rec-shine"></div>
      </div>
    </div>
    <div class="overlay-actions">
      <button id="retry" class="retry-btn" type="button">Retry another jump</button>
    </div>
  </div>


  <div id="authModalOverlay" class="auth-modal-overlay" onclick="closeAuthModal(event)">
    <div class="auth-modal" onclick="event.stopPropagation()">
      <button class="auth-close" onclick="closeAuthModal()">×</button>
      <div class="auth-tabs">
        <button class="auth-tab active" id="tabLogin" onclick="switchAuthTab('login')">Sign In</button>
        <button class="auth-tab" id="tabRegister" onclick="switchAuthTab('register')">Register</button>
      </div>
      <form id="authForm" onsubmit="handleAuthSubmit(event)">
        <div class="auth-field" id="fieldUsername" style="display:none">
          <label for="authUsername">Username</label>
          <input id="authUsername" type="text" placeholder="Pick a username" autocomplete="username" />
        </div>
        <div class="auth-field">
          <label for="authEmail">Email</label>
          <input id="authEmail" type="email" placeholder="you@example.com" autocomplete="email" />
        </div>
        <div class="auth-field">
          <label for="authPassword">Password</label>
          <input id="authPassword" type="password" placeholder="At least 6 characters" autocomplete="current-password" />
        </div>
        <div class="auth-error" id="authError"></div>
        <button type="submit" class="auth-submit" id="authSubmit">Sign In</button>
      </form>
    </div>
  </div>

  <script>
    const overlay = document.getElementById('overlay');
    const overlayTitle = document.getElementById('overlayTitle');
    const overlayDesc = document.getElementById('overlayDesc');
    const recommendations = document.getElementById('recommendations');
    const trigger = document.getElementById('trigger');
    const retry = document.getElementById('retry');
    const loader = document.querySelector('.loader');

    let isLoading = false;
    let audioCtx;
    let jumpController;
    let pendingJumpTimer;
    let targets = [];
    // ── Auth state ───────────────────────────────
    let authToken = localStorage.getItem('wuck_token') || '';
    let currentUser = null;

    function updateAuthUI() {
      const loginBtn = document.getElementById('loginBtn');
      const userPill = document.getElementById('userPill');
      const userName = document.getElementById('userName');
      const userAvatar = document.getElementById('userAvatar');
      const historyCount = document.getElementById('historyCount');

      if (currentUser) {
        loginBtn.style.display = 'none';
        userPill.style.display = 'inline-flex';
        userName.textContent = currentUser.username;
        userAvatar.textContent = currentUser.username.charAt(0).toUpperCase();
        const clicks = currentUser.clickHistory ? currentUser.clickHistory.length : 0;
        historyCount.textContent = clicks + ' click' + (clicks !== 1 ? 's' : '') + ' recorded';
      } else {
        loginBtn.style.display = 'inline-flex';
        userPill.style.display = 'none';
      }
    }

    async function checkAuth() {
      if (!authToken) return;
      try {
        const res = await fetch('/api/auth/me', {
          headers: { 'Authorization': 'Bearer ' + authToken }
        });
        if (res.ok) {
          const payload = await res.json();
          currentUser = payload.data;
        } else {
          authToken = '';
          localStorage.removeItem('wuck_token');
        }
      } catch (e) {
        // offline, keep token and try later
      }
      updateAuthUI();
    }

    function showAuthModal(mode) {
      document.getElementById('authModalOverlay').classList.add('active');
      document.getElementById('authError').textContent = '';
      document.getElementById('authEmail').value = '';
      document.getElementById('authPassword').value = '';
      document.getElementById('authUsername').value = '';
      switchAuthTab(mode);
      if (mode === 'login') {
        document.getElementById('authEmail').focus();
      } else {
        document.getElementById('authUsername').focus();
      }
    }

    function closeAuthModal(e) {
      if (e && e.target !== document.getElementById('authModalOverlay')) return;
      document.getElementById('authModalOverlay').classList.remove('active');
    }

    function switchAuthTab(mode) {
      const tabLogin = document.getElementById('tabLogin');
      const tabRegister = document.getElementById('tabRegister');
      const fieldUsername = document.getElementById('fieldUsername');
      const authSubmit = document.getElementById('authSubmit');
      const authError = document.getElementById('authError');

      if (mode === 'register') {
        tabLogin.classList.remove('active');
        tabRegister.classList.add('active');
        fieldUsername.style.display = 'block';
        authSubmit.textContent = 'Create Account';
      } else {
        tabLogin.classList.add('active');
        tabRegister.classList.remove('active');
        fieldUsername.style.display = 'none';
        authSubmit.textContent = 'Sign In';
      }
      authError.textContent = '';
      document.getElementById('authForm').dataset.mode = mode;
    }

    async function handleAuthSubmit(e) {
      e.preventDefault();
      const mode = document.getElementById('authForm').dataset.mode || 'login';
      const authError = document.getElementById('authError');
      authError.textContent = '';

      const email = document.getElementById('authEmail').value.trim();
      const password = document.getElementById('authPassword').value;

      if (!email || !password) {
        authError.textContent = 'Please fill in all fields';
        return;
      }

      if (mode === 'register' && password.length < 6) {
        authError.textContent = 'Password must be at least 6 characters';
        return;
      }

      let body;
      if (mode === 'register') {
        const username = document.getElementById('authUsername').value.trim();
        if (!username) {
          authError.textContent = 'Please enter a username';
          return;
        }
        body = JSON.stringify({ username, email, password });
      } else {
        body = JSON.stringify({ email, password });
      }

      try {
        const endpoint = mode === 'register' ? '/api/auth/register' : '/api/auth/login';
        const res = await fetch(endpoint, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: body
        });
        const payload = await res.json();

        if (!res.ok) {
          authError.textContent = payload.msg || 'Something went wrong';
          return;
        }

        authToken = payload.data.token;
        currentUser = payload.data.user;
        localStorage.setItem('wuck_token', authToken);
        updateAuthUI();
        closeAuthModal();
      } catch (err) {
        authError.textContent = 'Network error, please try again';
      }
    }

    function handleLogout() {
      authToken = '';
      currentUser = null;
      localStorage.removeItem('wuck_token');
      updateAuthUI();
      document.getElementById('userDropdown').classList.remove('open');
    }

    function toggleUserDropdown(e) {
      e.stopPropagation();
      document.getElementById('userDropdown').classList.toggle('open');
    }

    document.addEventListener('click', function() {
      document.getElementById('userDropdown').classList.remove('open');
    });

    async function recordClick(item) {
      if (!authToken || !currentUser) return;
      try {
        await fetch('/api/auth/click', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + authToken
          },
          body: JSON.stringify({
            url: item.url,
            title: item.title,
            source: item.source,
            category: item.category
          })
        });
        // Update local click count
        if (currentUser.clickHistory) {
          currentUser.clickHistory.push({ url: item.url, title: item.title, source: item.source, category: item.category });
        }
        updateAuthUI();
      } catch (e) {}
    }

    let particlesCanvas;
    let particlesCtx;
    let particlesActive = false;
    let particlesRaf;

    // Category emoji map
    const categoryEmoji = {
      novel: '📖',
      video: '🎬',
      blog: '📝',
      news: '📰',
      lifestyle: '🌟',
    };

    // ── Particles ──────────────────────────────────
    function initParticles() {
      particlesCanvas = document.createElement('canvas');
      particlesCanvas.id = 'particles';
      document.body.appendChild(particlesCanvas);
      particlesCtx = particlesCanvas.getContext('2d');
      resizeParticles();
      window.addEventListener('resize', resizeParticles);
    }

    function resizeParticles() {
      if (!particlesCanvas) return;
      particlesCanvas.width = window.innerWidth;
      particlesCanvas.height = window.innerHeight;
    }

    let particlePool = [];
    const MAX_PARTICLES = 60;

    function spawnParticles(x, y, color) {
      const now = performance.now();
      for (let i = 0; i < 18; i++) {
        const angle = (Math.PI * 2 * i) / 18 + Math.random() * 0.5;
        const speed = 60 + Math.random() * 140;
        particlePool.push({
          x, y,
          vx: Math.cos(angle) * speed,
          vy: Math.sin(angle) * speed,
          life: 0.5 + Math.random() * 0.7,
          born: now,
          size: 2 + Math.random() * 3,
          color: color || 'hsl(' + (260 + Math.random() * 40) + ', 80%, ' + (60 + Math.random() * 30) + '%)',
        });
      }
      if (particlePool.length > MAX_PARTICLES) {
        particlePool = particlePool.slice(-MAX_PARTICLES);
      }
      if (!particlesActive) {
        particlesActive = true;
        particlesCanvas.classList.add('active');
        animateParticles();
      }
    }

    function animateParticles() {
      if (!particlesCtx) return;
      const now = performance.now();
      particlesCtx.clearRect(0, 0, particlesCanvas.width, particlesCanvas.height);

      particlePool = particlePool.filter(p => {
        const age = (now - p.born) / 1000;
        if (age > p.life) return false;
        const alpha = 1 - age / p.life;
        const x = p.x + p.vx * age;
        const y = p.y + p.vy * age + 100 * age * age; // gravity
        particlesCtx.beginPath();
        particlesCtx.arc(x, y, p.size * alpha, 0, Math.PI * 2);
        particlesCtx.fillStyle = p.color;
        particlesCtx.globalAlpha = alpha * 0.7;
        particlesCtx.fill();
        return true;
      });

      particlesCtx.globalAlpha = 1;

      if (particlePool.length === 0) {
        particlesActive = false;
        particlesCanvas.classList.remove('active');
        return;
      }
      particlesRaf = requestAnimationFrame(animateParticles);
    }

    function stopParticles() {
      if (particlesRaf) {
        cancelAnimationFrame(particlesRaf);
        particlesRaf = null;
      }
      particlePool = [];
      particlesActive = false;
      if (particlesCanvas) particlesCanvas.classList.remove('active');
    }

    // ── 3D Tilt ────────────────────────────────────
    function bindTilt(card) {
      card.addEventListener('mousemove', (e) => {
        if (card.classList.contains('picked')) return;
        const rect = card.getBoundingClientRect();
        const cx = rect.left + rect.width / 2;
        const cy = rect.top + rect.height / 2;
        const dx = e.clientX - cx;
        const dy = e.clientY - cy;
        const rx = (dy / rect.height) * -10;
        const ry = (dx / rect.width) * 10;
        card.style.transform = 'perspective(600px) rotateX(' + rx + 'deg) rotateY(' + ry + 'deg) translateY(-4px)';
      });
      card.addEventListener('mouseleave', () => {
        if (card.classList.contains('picked')) return;
        card.style.transform = '';
      });
    }

    // ── Audio ──────────────────────────────────────
    function playTeleportSound() {
      if (!window.AudioContext && !window.webkitAudioContext) return;
      if (!audioCtx) {
        audioCtx = new (window.AudioContext || window.webkitAudioContext)();
      }
      if (audioCtx.state === 'suspended') audioCtx.resume();

      const osc = audioCtx.createOscillator();
      const gain = audioCtx.createGain();
      osc.type = 'sine';
      osc.frequency.setValueAtTime(150, audioCtx.currentTime);
      osc.frequency.exponentialRampToValueAtTime(800, audioCtx.currentTime + 0.5);
      gain.gain.setValueAtTime(0, audioCtx.currentTime);
      gain.gain.linearRampToValueAtTime(0.1, audioCtx.currentTime + 0.1);
      gain.gain.exponentialRampToValueAtTime(0.01, audioCtx.currentTime + 0.6);
      osc.connect(gain);
      gain.connect(audioCtx.destination);
      osc.start();
      osc.stop(audioCtx.currentTime + 0.6);
    }

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

    function playPickSound() {
      if (!audioCtx) return;
      const osc = audioCtx.createOscillator();
      const gain = audioCtx.createGain();
      osc.type = 'sine';
      osc.frequency.setValueAtTime(500, audioCtx.currentTime);
      osc.frequency.exponentialRampToValueAtTime(1200, audioCtx.currentTime + 0.2);
      osc.frequency.exponentialRampToValueAtTime(900, audioCtx.currentTime + 0.35);
      gain.gain.setValueAtTime(0, audioCtx.currentTime);
      gain.gain.linearRampToValueAtTime(0.12, audioCtx.currentTime + 0.05);
      gain.gain.linearRampToValueAtTime(0, audioCtx.currentTime + 0.35);
      osc.connect(gain);
      gain.connect(audioCtx.destination);
      osc.start();
      osc.stop(audioCtx.currentTime + 0.35);
    }

    // ── Card state ─────────────────────────────────
    function clearJumpTimers() {
      if (pendingJumpTimer) {
        window.clearTimeout(pendingJumpTimer);
        pendingJumpTimer = null;
      }
    }

    function resetRecommendationCards() {
      for (let i = 0; i < 3; i++) {
        const card = document.getElementById('recCard' + i);
        const badge = document.getElementById('recBadge' + i);
        const title = document.getElementById('recTitle' + i);
        const cat = document.getElementById('recCat' + i);
        const thumbImg = document.getElementById('recThumb' + i);
        const fallback = document.getElementById('recFallback' + i);
        const emoji = fallback ? fallback.querySelector('.rec-thumb-emoji') : null;

        card.classList.remove('show', 'picked');
        card.removeAttribute('data-category');
        card.style.transform = '';
        card.style.borderColor = '';
        card.style.boxShadow = '';
        card.style.filter = '';
        badge.textContent = 'SOURCE';
        title.textContent = 'Loading...';
        cat.textContent = 'category';
        if (thumbImg) {
          thumbImg.classList.remove('loaded');
          thumbImg.src = '';
        }
        if (emoji) emoji.textContent = '🎲';
      }
      recommendations.classList.remove('visible');
      stopParticles();
    }

    function showRecommendationCards(items) {
      targets = items;

      // Hide loader
      loader.style.display = 'none';

      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        const card = document.getElementById('recCard' + i);
        const badge = document.getElementById('recBadge' + i);
        const title = document.getElementById('recTitle' + i);
        const cat = document.getElementById('recCat' + i);
        const thumbImg = document.getElementById('recThumb' + i);
        const fallback = document.getElementById('recFallback' + i);
        const emojiEl = fallback ? fallback.querySelector('.rec-thumb-emoji') : null;

        card.setAttribute('data-category', item.category);
        badge.textContent = item.source;
        title.textContent = item.title;
        cat.textContent = item.category;

        // Set category emoji
        if (emojiEl) {
          emojiEl.textContent = categoryEmoji[item.category] || '🎲';
        }

        // Load thumbnail
        if (thumbImg && item.thumbnailUrl) {
          const img = new Image();
          img.onload = function () {
            thumbImg.src = item.thumbnailUrl;
            thumbImg.classList.add('loaded');
          };
          img.onerror = function () {
            // Keep fallback visible
          };
          img.src = item.thumbnailUrl;
        }
      }

      // Show grid
      recommendations.classList.add('visible');

      // Staggered spring entrance
      for (let i = 0; i < items.length; i++) {
        const card = document.getElementById('recCard' + i);
        setTimeout(() => {
          card.classList.add('show');
          bindTilt(card);
        }, 100 + i * 130);
      }
    }

    function resetOverlayState() {
      clearJumpTimers();
      if (jumpController) {
        jumpController.abort();
        jumpController = null;
      }

      isLoading = false;
      targets = [];
      overlay.classList.remove('active');
      overlayTitle.textContent = 'Summoning Chaos';
      overlayDesc.textContent = 'Rolling the dice across the internet, please wait...';
      loader.style.display = '';
      resetRecommendationCards();
    }

    // ── Pick & navigate ────────────────────────────
    function pickTarget(index) {
      if (!targets[index]) return;
      const item = targets[index];

      // Record click if logged in
      recordClick(item);

      // Fade out unselected cards
      for (let i = 0; i < 3; i++) {
        if (i !== index) {
          const other = document.getElementById('recCard' + i);
          other.style.transition = 'opacity 0.3s ease, transform 0.3s ease';
          other.style.opacity = '0.3';
          other.style.transform = 'scale(0.92)';
          other.style.pointerEvents = 'none';
        }
      }

      const card = document.getElementById('recCard' + index);
      const rect = card.getBoundingClientRect();
      const cx = rect.left + rect.width / 2;
      const cy = rect.top + rect.height / 2;

      playPickSound();
      card.classList.add('picked');
      card.style.pointerEvents = 'none';

      // Spawn particles from card center
      spawnParticles(cx, cy, '#7c3aed');

      // Overlay flash
      const flash = document.createElement('div');
      flash.style.cssText = 'position:fixed;inset:0;background:rgba(124,58,237,0.06);z-index:102;pointer-events:none;';
      document.body.appendChild(flash);
      setTimeout(() => flash.remove(), 300);

      setTimeout(() => {
        window.location.href = item.url;
      }, 400);
    }

    // ── Main jump flow ─────────────────────────────
    async function startJump(forceRetry) {
      if (isLoading && !forceRetry) return;

      clearJumpTimers();
      if (jumpController) {
        jumpController.abort();
      }

      jumpController = new AbortController();
      isLoading = true;

      resetRecommendationCards();
      loader.style.display = '';
      overlayTitle.textContent = 'Navigating the Web...';
      overlayDesc.textContent = 'RNG is picking targets from live trending sources';
      overlay.classList.add('active');

      playTeleportSound();

      try {
        const headers = { 'Accept': 'application/json' };
        if (authToken) {
          headers['Authorization'] = 'Bearer ' + authToken;
        }
        const res = await fetch('/api/biu/three', {
          headers: headers,
          signal: jumpController.signal,
        });
        const payload = await res.json();

        if (!res.ok || !payload.data || !Array.isArray(payload.data) || payload.data.length === 0) {
          throw new Error(payload.msg || 'API returned invalid data');
        }

        pendingJumpTimer = window.setTimeout(() => {
          playSuccessSound();
          overlayTitle.textContent = 'Pick Your Journey';
          overlayDesc.textContent = 'The web has spoken — choose your next destination';
          showRecommendationCards(payload.data);
          isLoading = false;
        }, 400);

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

    // ── Init ───────────────────────────────────────
    initParticles();
    checkAuth();
    window.addEventListener('pageshow', function() {
      resetOverlayState();
      checkAuth();
    });
    window.addEventListener('pagehide', clearJumpTimers);
    trigger.addEventListener('click', startJump);
    retry.addEventListener('click', () => startJump(true));
  </script>
</body>
</html>`
