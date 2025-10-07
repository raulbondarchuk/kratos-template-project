// ui.go
package swagger

import "html/template"

var uiTpl = template.Must(template.New("ui").Parse(`<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css"/>
  <style>
    :root{
      --hdr-bg:#ffffff; --hdr-fg:#0f172a; --border:#e5e7eb;
      --green:#10b981; --green-700:#047857; --gray:#374151;
      --danger:#ef4444;
      --input-bg:#ffffff; --input-fg:#0f172a; --ph:#9ca3af;
      --ring:rgba(16,185,129,.18);
      --h:44px;

      --chip-bg:#ecfdf5; --chip-br:#a7f3d0; --chip-fg:#065f46;

      --link-bg:#f1f5ff; --link-br:#dbe4ff; --link-fg:#1e3a8a; --link-bg-h:#e6edff;
    }
    html,body{margin:0}
    .app-header{
      position:sticky; top:0; z-index:1000;
      display:flex; align-items:center; gap:12px;
      padding:12px 16px; background:var(--hdr-bg); color:var(--hdr-fg);
      border-bottom:1px solid var(--border); box-shadow:0 1px 6px rgba(0,0,0,.04);
      flex-wrap:wrap;
      font-family:ui-sans-serif,system-ui,-apple-system,Segoe UI,Roboto,Ubuntu;
    }
    .brand{display:flex; align-items:center; gap:10px}
    .brand img{height:28px}
    .controls{margin-left:auto; display:flex; align-items:center; gap:10px; justify-content:flex-end; flex:1 1 auto; min-width:320px; flex-wrap:wrap;}

    .scheme-chip{
      display:inline-flex; align-items:center; justify-content:center;
      height:var(--h); padding:0 18px;
      border-radius:14px; border:1px solid var(--chip-br);
      background:var(--chip-bg); color:var(--chip-fg);
      font-weight:700; letter-spacing:.02em; text-transform:uppercase;
      box-shadow:0 10px 24px rgba(16,185,129,.12), inset 0 0 0 1px rgba(16,185,129,.08);
    }

    .link-pill{
      display:inline-flex; align-items:center;
      height:var(--h); padding:0 14px; border-radius:12px;
      background:var(--link-bg); color:var(--link-fg); text-decoration:none;
      font-weight:700; border:1px solid var(--link-br);
      max-width:48vw; overflow:hidden; text-overflow:ellipsis; white-space:nowrap;
      box-shadow:0 8px 18px rgba(30,58,138,.08);
      transition:background .15s ease, transform .05s ease, border-color .15s ease;
    }
    .link-pill:hover{ background:var(--link-bg-h); border-color:#cdd8ff }
    .link-pill:active{ transform:translateY(1px) }

    .token-field{
      position:relative; display:flex; align-items:center; gap:8px; height:var(--h); padding:0 12px 0 10px;
      border:1px solid var(--border); border-radius:14px; background:var(--input-bg);
      box-shadow:0 1px 2px rgba(0,0,0,.04); transition:border .2s, box-shadow .2s;
      flex:0 1 min(560px, 46vw); min-width:260px;
    }
    .token-field:focus-within{border-color:var(--green); box-shadow:0 0 0 3px var(--ring), 0 4px 14px rgba(0,0,0,.06);}
    .token-ic{width:18px;height:18px;stroke:#9ca3af;flex:0 0 auto}
    .token-input{flex:1 1 auto; min-width:120px; height:100%; border:0; outline:none; background:transparent; color:var(--input-fg); font:inherit; padding:0 38px 0 2px;}
    .token-input::placeholder{color:var(--ph)}
    .copy-ic{
      position:absolute; right:6px; top:50%; transform:translateY(-50%);
      display:inline-flex; align-items:center; justify-content:center;
      width:32px; height:32px; border-radius:9px; border:0; background:transparent; cursor:pointer;
      color:#6b7280; transition:background .15s ease, color .15s ease, opacity .15s ease;
    }
    .copy-ic:hover{ background:#f3f4f6; color:#374151 }
    .copy-ic[disabled]{ opacity:.45; cursor:not-allowed }
    .copy-ic svg{width:18px;height:18px;stroke:currentColor;fill:none}

    .btn{display:inline-flex; align-items:center; gap:8px; height:var(--h); padding:0 14px; border:1px solid transparent; border-radius:12px; font-weight:800; cursor:pointer; user-select:none; white-space:nowrap; transition:transform .08s ease, filter .15s ease, box-shadow .15s ease, background .2s ease, color .2s ease; flex:0 0 auto;}
    .btn:active{transform:translateY(1px)}
    .btn svg{width:18px; height:18px; stroke:currentColor; fill:none}
    .btn-save{ color:#fff; background:linear-gradient(135deg, #10b981, #0ea371); box-shadow:0 6px 16px rgba(16,185,129,.22); }
    .btn-save:hover{filter:brightness(1.03)}
    .btn-clear{ color:#374151; background:#fff; border-color:var(--border); }
    .btn-clear:hover{background:#f9fafb; box-shadow:0 4px 12px rgba(0,0,0,.06)}
    .btn-logout{ color:#fff; background:linear-gradient(135deg, #ef4444, #dc2626); box-shadow:0 6px 16px rgba(239,68,68,.22); }
    .btn-logout:hover{filter:brightness(1.03)}

    @media (max-width: 900px){ .controls{flex:1 1 100%; justify-content:flex-end} }

    .toast{position:fixed; right:16px; top:72px; z-index:1100; padding:10px 14px; border-radius:10px; background:#ecfdf5; color:#065f46; border:1px solid #a7f3d0; box-shadow:0 8px 24px rgba(0,0,0,.12); opacity:0; transform:translateY(-6px); transition:opacity .18s, transform .18s; font-weight:700; pointer-events:none;}
    .toast.show{opacity:1; transform:translateY(0)}
    .content{padding-top:6px}
    form.logout{margin:0}
  </style>
</head>
<body>
  <header class="app-header">
    <div class="brand">
      <img src="{{.Base}}/docs/logo.png" alt="logo" onerror="this.style.display='none'">
    </div>

    <div class="controls">
      <div id="schemeChip" class="scheme-chip">HTTPS</div>

      <a id="baseLink" class="link-pill" href="#" title="Click to copy" rel="noopener">
        <span id="baseText">loading…</span>
      </a>

      <div class="token-field">
        <svg class="token-ic" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 9h14M5 15h14M10 3L8 21M16 3l-2 18"/>
        </svg>
        <input id="authToken" class="token-input" type="text" placeholder="Paste access token..."/>
        <button id="copyTokenBtn" class="copy-ic" type="button" title="Copy token" aria-label="Copy token" onclick="copyToken()" disabled>
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <rect x="4.5" y="4.5" width="11" height="13" rx="2" ry="2" stroke-width="2"/>
            <rect x="8.5" y="6.5" width="11" height="13" rx="2" ry="2" stroke-width="2"/>
          </svg>
        </button>
      </div>

      <button class="btn btn-save" onclick="saveToken()" type="button">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
        <span>Save</span>
      </button>

      <button class="btn btn-clear" onclick="clearToken()" type="button">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
        <span>Clear</span>
      </button>

      <form class="logout" method="post" action="{{.Base}}/docs/logout">
        <button class="btn btn-logout" type="submit" title="Salir">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4-4-4"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12H9"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5H7a2 2 0 00-2 2v10a2 2 0 002 2h6"/>
          </svg>
          <span>Logout</span>
        </button>
      </form>
    </div>
  </header>

  <div class="toast" id="toast">Saved</div>
  <div class="content"><div id="swagger-ui"></div></div>

  <script src="{{.Base}}/docs/bootstrap.js" defer></script>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>
    const STATIC_BASE   = '{{.Base}}';
    const PROJECT_PREF  = '{{.DefaultProj}}';
    const FIXED_SCHEME  = '{{.FixedScheme}}';
    const SPEC_URL      = STATIC_BASE + '/docs/openapi.yaml';

    const K='swagger_access_token', K_OVR='swagger_user_override';
    function toast(msg){ const t=document.getElementById('toast'); t.textContent=msg||'Saved'; t.classList.add('show'); setTimeout(()=>t.classList.remove('show'),1300); }
    function saveToken(){ const el=document.getElementById('authToken'); const v=(el?.value||'').trim(); if(v){ localStorage.setItem(K,v); localStorage.setItem(K_OVR,'1'); toast('Token saved'); } else { localStorage.removeItem(K); localStorage.setItem(K_OVR,'1'); toast('Token cleared'); } syncCopyBtn(); }
    function clearToken(){ const el=document.getElementById('authToken'); if(el) el.value=''; localStorage.removeItem(K); localStorage.setItem(K_OVR,'1'); toast('Token cleared'); syncCopyBtn(); }
    async function copyToken(){ const v=(document.getElementById('authToken')?.value||'').trim(); if(!v){ toast('Nothing to copy'); return; } try{ if(navigator.clipboard?.writeText){ await navigator.clipboard.writeText(v); } else { const ta=document.createElement('textarea'); ta.value=v; document.body.appendChild(ta); ta.select(); document.execCommand('copy'); document.body.removeChild(ta); } toast('Copied'); }catch(_){ toast('Copy failed'); } }
    function syncCopyBtn(){ const el=document.getElementById('authToken'); const btn=document.getElementById('copyTokenBtn'); if(btn){ btn.disabled = !(el && el.value.trim().length); } }

    function readScheme(){ const s=(FIXED_SCHEME||'').toLowerCase(); if(s==='http'||s==='https') return s; return (location.protocol||'http:').replace(':',''); }
    function effectivePrefix(){ return readScheme()==='https' ? (PROJECT_PREF || '') : ''; }

    function buildAbsoluteUrl(pathOrUrl){
      const scheme = readScheme();
      let u; try { u = new URL(pathOrUrl, window.location.origin); } catch(_) { return pathOrUrl; }
      let path = u.pathname + u.search + u.hash;
      const pref = effectivePrefix();
      if (pref && pref !== '/' && path.startsWith('/') && !path.startsWith(pref + '/') && path !== pref){
        path = pref + path;
      }
      return scheme + '://' + window.location.host + path;
    }

    function paintHeader(){
      const scheme = readScheme();
      const pref   = effectivePrefix();
      const base   = scheme + '://' + window.location.host + (pref||'');

      const chip = document.getElementById('schemeChip');
      if (chip) chip.textContent = scheme.toUpperCase();

      const link = document.getElementById('baseLink');
      const txt  = document.getElementById('baseText');
      if (txt)  txt.textContent = base;
      if (link){
        link.href = base + '/';
        link.onclick = async (e)=>{
          if (e.metaKey || e.ctrlKey) return;
          e.preventDefault();
          try{ await navigator.clipboard.writeText(base); toast('URL copied'); }catch(_){ toast('Copy failed'); }
        };
      }
    }

    document.addEventListener('DOMContentLoaded', ()=>{
      const input = document.getElementById('authToken');
      if(input){
        input.value = localStorage.getItem(K) || input.value || '';
        input.addEventListener('input', ()=>{ localStorage.setItem(K, input.value.trim()); localStorage.setItem(K_OVR,'1'); syncCopyBtn(); });
        syncCopyBtn();
      }
      paintHeader();
    });

    window.ui = SwaggerUIBundle({
      url: SPEC_URL,
      dom_id:'#swagger-ui',
      presets:[SwaggerUIBundle.presets.apis],
      persistAuthorization:true,
      requestInterceptor:(req)=>{
        const isSpec = /\/docs\/openapi\.yaml(?:\?|$)/.test(req.url||'');
        if (!isSpec) req.url = buildAbsoluteUrl(req.url || '/');
        const token = localStorage.getItem(K) || '';
        if (token){ req.headers['Authorization'] = token; }
        return req;
      }
    });
  </script>
</body>
</html>`))
