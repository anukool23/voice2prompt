// Frontend logic for Voice2Prompt. Talks to the Go backend via window.go.main.App.*
// and listens for live events via window.runtime.

const App = () => window.go.main.App;
let settings = null;
let dict = {};

function toast(msg) {
  const t = document.getElementById('toast');
  t.textContent = msg;
  t.classList.add('show');
  setTimeout(() => t.classList.remove('show'), 1600);
}

async function loadSettings() {
  settings = await App().GetSettings();
  document.getElementById('trigger').value = settings.trigger || 'chord';
  document.getElementById('hotkey').value = settings.hotkey;
  onTriggerChange();
  document.getElementById('language').value = settings.language;
  document.getElementById('whisperModel').value = settings.whisperModel;
  document.getElementById('cleanupEnabled').checked = settings.cleanupEnabled;
  document.getElementById('commandsEnabled').checked = settings.commandsEnabled;
  document.getElementById('llmModel').value = settings.llmModel;
  dict = settings.dictionary || {};
  renderDict();

  if (!settings.onboardingComplete) {
    document.getElementById('onboard').classList.add('show');
  }
}

function renderDict() {
  const list = document.getElementById('dictList');
  list.innerHTML = '';
  const keys = Object.keys(dict);
  if (keys.length === 0) {
    list.innerHTML = '<div class="empty">No entries yet.</div>';
    return;
  }
  for (const k of keys) {
    const row = document.createElement('div');
    row.className = 'dict-row';
    row.innerHTML = `<input type="text" value="${escapeHtml(k)}" disabled />
      <input type="text" value="${escapeHtml(dict[k])}" disabled />
      <button title="Remove">✕</button>`;
    row.querySelector('button').onclick = () => { delete dict[k]; renderDict(); };
    list.appendChild(row);
  }
}

function addDict() {
  const k = document.getElementById('dictKey').value.trim();
  const v = document.getElementById('dictVal').value.trim();
  if (!k || !v) return;
  dict[k] = v;
  document.getElementById('dictKey').value = '';
  document.getElementById('dictVal').value = '';
  renderDict();
}

async function save() {
  const updated = {
    trigger: document.getElementById('trigger').value,
    hotkey: document.getElementById('hotkey').value,
    language: document.getElementById('language').value,
    whisperModel: document.getElementById('whisperModel').value,
    cleanupEnabled: document.getElementById('cleanupEnabled').checked,
    commandsEnabled: document.getElementById('commandsEnabled').checked,
    llmModel: document.getElementById('llmModel').value,
    dictionary: dict,
    onboardingComplete: settings ? settings.onboardingComplete : true,
  };
  try {
    await App().SaveSettings(updated);
    settings = updated;
    toast('Settings saved');
    refreshStatus();
  } catch (e) {
    toast('Save failed: ' + e);
  }
}

async function refreshStatus() {
  const running = await App().EngineRunning();
  setEngineUI(running);
  const ax = await App().AccessibilityTrusted();
  const pill = document.getElementById('axPill');
  pill.textContent = ax ? 'Granted' : 'Not granted';
  pill.className = 'pill ' + (ax ? 'ok' : 'warn');
  document.getElementById('axActionRow').style.display = ax ? 'none' : 'flex';

  const mic = await App().MicStatus();
  const micPill = document.getElementById('micPill');
  const micOk = mic === 'authorized';
  micPill.textContent = micOk ? 'Granted' : (mic === 'denied' ? 'Denied' : 'Not granted');
  micPill.className = 'pill ' + (micOk ? 'ok' : 'warn');
  document.getElementById('micActionRow').style.display = micOk ? 'none' : 'flex';

  // Input Monitoring — only relevant for the Fn trigger.
  const useFn = document.getElementById('trigger').value === 'fn';
  const imRow = document.getElementById('imRow');
  if (useFn) {
    const im = await App().InputMonitoringStatus();
    const imOk = im === 'authorized';
    imRow.style.display = 'flex';
    const imPill = document.getElementById('imPill');
    imPill.textContent = imOk ? 'Granted' : (im === 'denied' ? 'Denied' : 'Not granted');
    imPill.className = 'pill ' + (imOk ? 'ok' : 'warn');
    document.getElementById('imActionRow').style.display = imOk ? 'none' : 'flex';
  } else {
    imRow.style.display = 'none';
    document.getElementById('imActionRow').style.display = 'none';
  }

  document.getElementById('launchAtLogin').checked = await App().LaunchAtLoginEnabled();
}

async function requestMic() {
  await App().RequestMicrophone();
  toast('Approve the microphone prompt, then Start dictation.');
  setTimeout(refreshStatus, 1500);
}

function onTriggerChange() {
  const fn = document.getElementById('trigger').value === 'fn';
  document.getElementById('hotkeyRow').style.display = fn ? 'none' : 'flex';
  document.getElementById('fnHint').style.display = fn ? 'block' : 'none';
  refreshStatus();
}

async function requestIM() {
  await App().RequestInputMonitoring();
  toast('Enable Voice2Prompt under Input Monitoring, then relaunch.');
  setTimeout(refreshStatus, 1500);
}

async function toggleLaunchAtLogin() {
  const on = document.getElementById('launchAtLogin').checked;
  try {
    await App().SetLaunchAtLogin(on);
    toast(on ? 'Will launch at login' : 'Launch at login off');
  } catch (e) {
    toast('Failed: ' + e);
    document.getElementById('launchAtLogin').checked = !on;
  }
}

function setEngineUI(running) {
  document.getElementById('engineDot').className = 'dot' + (running ? ' live' : '');
  document.getElementById('engineLabel').textContent = running ? 'Listening' : 'Stopped';
  const btn = document.getElementById('engineBtn');
  btn.textContent = running ? 'Stop dictation' : 'Start dictation';
  btn.className = running ? 'stop' : '';
  const cp = document.getElementById('cleanupPill');
  if (!running) { cp.textContent = '—'; cp.className = 'pill off'; }
}

async function toggleEngine() {
  const running = await App().EngineRunning();
  const btn = document.getElementById('engineBtn');
  btn.disabled = true;
  try {
    if (running) {
      await App().StopEngine();
    } else {
      btn.textContent = 'Starting…';
      await App().StartEngine();
    }
  } catch (e) {
    toast('Engine error: ' + e);
  } finally {
    btn.disabled = false;
    refreshStatus();
  }
}

async function requestAX() {
  await App().RequestAccessibility();
  toast('Check System Settings → Accessibility, then relaunch.');
}

async function finishOnboarding() {
  await App().CompleteOnboarding();
  if (settings) settings.onboardingComplete = true;
  document.getElementById('onboard').classList.remove('show');
}

function addUtterance(u) {
  const feed = document.getElementById('feed');
  if (feed.querySelector('.empty')) feed.innerHTML = '';
  const div = document.createElement('div');
  div.className = 'utt';
  if (u.error) {
    div.innerHTML = `<div class="meta" style="color:var(--err)">⚠ ${escapeHtml(u.error)}</div>`;
  } else if (!u.raw && !u.command) {
    div.innerHTML = `<div class="meta">… no speech detected (check the microphone)</div>`;
  } else if (u.command) {
    div.innerHTML =
      `<div class="clean">⌘ ${escapeHtml(u.command)}` +
      (u.cleaned && u.cleaned !== u.raw ? ` → ${escapeHtml(u.cleaned)}` : '') + `</div>` +
      `<div class="meta">heard “${escapeHtml(u.raw)}” · ${u.totalMS}ms · ${escapeHtml(u.app || '?')}</div>`;
  } else {
    const showRaw = u.cleaned && u.raw && u.cleaned !== u.raw;
    div.innerHTML =
      (showRaw ? `<div class="raw">${escapeHtml(u.raw)}</div>` : '') +
      `<div class="clean">${escapeHtml(u.cleaned || u.raw)}</div>` +
      `<div class="meta">${u.audioSecs.toFixed(1)}s · infer ${u.inferMS}ms` +
      (u.cleanMS ? ` · clean ${u.cleanMS}ms` : '') +
      ` · ${u.totalMS}ms · ${escapeHtml(u.app || '?')} · ${u.method}</div>`;
  }
  feed.insertBefore(div, feed.firstChild);
  // Reflect cleanup activity in the status pill.
  const cp = document.getElementById('cleanupPill');
  if (u.cleanMS) { cp.textContent = 'Active'; cp.className = 'pill ok'; }
}

function escapeHtml(s) {
  return String(s).replace(/[&<>"']/g, c =>
    ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));
}

// Wait for the Wails runtime to be ready.
window.addEventListener('DOMContentLoaded', async () => {
  await loadSettings();
  await refreshStatus();
  const hist = await App().History();
  hist.slice().reverse().forEach(addUtterance);

  window.runtime.EventsOn('utterance', addUtterance);
  window.runtime.EventsOn('engine:state', setEngineUI);
});
