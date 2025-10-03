package swagger

const uiHTML = `<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css"/>
  <style>
    :root{
      --hdr-bg:#ffffff; --hdr-fg:#0f172a; --border:#e5e7eb;
      --green:#10b981; --green-h:#059669; --gray:#374151;
      --danger:#ef4444; --danger-h:#dc2626;
      --input-bg:#ffffff; --input-fg:#0f172a; --ph:#9ca3af;
      --ring:rgba(16,185,129,.18);
      --h:44px;
    }
    html,body{margin:0}
    .app-header{
      position:sticky; top:0; z-index:1000;
      display:flex; align-items:center; gap:16px;
      padding:12px 16px; background:var(--hdr-bg); color:var(--hdr-fg);
      border-bottom:1px solid var(--border); box-shadow:0 1px 6px rgba(0,0,0,.04);
      flex-wrap:wrap;
    }
    .brand{display:flex; align-items:center; gap:10px; flex:0 0 auto;}
    .brand img{height:28px; width:auto; display:block}
    .brand .title{font-weight:800; letter-spacing:.2px}

    /* everything on the right */
    .controls{
      margin-left:auto;
      display:flex; align-items:center; gap:10px;
      justify-content:flex-end;
      flex:1 1 auto; min-width:320px; flex-wrap:wrap;
    }

    /* TOKEN FIELD */
    .token-field{
      display:flex; align-items:center; gap:8px;
      height:var(--h);
      padding:0 12px;
      border:1px solid var(--border); border-radius:14px;
      background:var(--input-bg);
      box-shadow:0 1px 2px rgba(0,0,0,.04);
      transition:border .2s, box-shadow .2s;
      flex:0 1 min(560px, 46vw);
      min-width:260px;
    }
    .token-ic{width:18px;height:18px;stroke:#9ca3af;flex:0 0 auto}
    .token-input{
      flex:1 1 auto; min-width:120px; height:100%;
      border:0; outline:none; background:transparent;
      color:var(--input-fg); font:inherit; padding:0 2px;
    }
    .token-input::placeholder{color:var(--ph)}
    .token-field:focus-within{
      border-color:var(--green);
      box-shadow:0 0 0 3px var(--ring), 0 4px 14px rgba(0,0,0,.06);
    }

    /* ICON BUTTON (copy, right) */
    .icon-btn{
      display:inline-flex; align-items:center; justify-content:center;
      width:32px; height:32px; border:0; background:transparent; cursor:pointer;
      border-radius:8px; transition:background .15s ease, opacity .15s ease, transform .05s ease;
      color:#6b7280; flex:0 0 auto;
    }
    .icon-btn:hover{ background:#f3f4f6 }
    .icon-btn:active{ transform:translateY(1px) }
    .icon-btn[disabled]{ opacity:.45; cursor:not-allowed }
    .icon-btn svg{ width:18px; height:18px; stroke:currentColor; fill:none }

    /* BUTTONS */
    .btn{
      display:inline-flex; align-items:center; gap:8px;
      height:var(--h); padding:0 14px;
      border:1px solid transparent; border-radius:12px;
      font-weight:800; cursor:pointer; user-select:none; white-space:nowrap;
      transition:transform .08s ease, filter .15s ease, box-shadow .15s ease, background .2s ease, color .2s ease;
      flex:0 0 auto;
    }
    .btn:active{transform:translateY(1px)}
    .btn svg{width:18px; height:18px; stroke:currentColor; fill:none}

    .btn-save{
      color:#fff; background:linear-gradient(135deg, #10b981, #0ea371);
      box-shadow:0 6px 16px rgba(16,185,129,.22);
    }
    .btn-save:hover{filter:brightness(1.03)}
    .btn-save:focus-visible{outline:2px solid rgba(16,185,129,.55); outline-offset:2px}

    .btn-clear{
      color:var(--gray); background:#fff; border-color:var(--border);
    }
    .btn-clear:hover{background:#f9fafb; box-shadow:0 4px 12px rgba(0,0,0,.06)}
    .btn-clear:focus-visible{outline:2px solid rgba(15,23,42,.18); outline-offset:2px}

    .btn-logout{
      color:#fff; background:linear-gradient(135deg, #ef4444, #dc2626);
      box-shadow:0 6px 16px rgba(239,68,68,.22);
    }
    .btn-logout:hover{filter:brightness(1.03)}
    .btn-logout:focus-visible{outline:2px solid rgba(239,68,68,.5); outline-offset:2px}

    @media (max-width: 900px){
      .controls{flex:1 1 100%; justify-content:flex-end}
    }

    .toast{
      position:fixed; right:16px; top:72px; z-index:1100;
      padding:10px 14px; border-radius:10px; background:#ecfdf5; color:#065f46;
      border:1px solid #a7f3d0; box-shadow:0 8px 24px rgba(0,0,0,.12);
      opacity:0; transform:translateY(-6px); transition:opacity .18s, transform .18s; font-weight:700; pointer-events:none;
    }
    .toast.show{opacity:1; transform:translateY(0)}
    .content{padding-top:6px}
    form.logout{margin:0}
  </style>
</head>
<body>
  <header class="app-header">
    <div class="brand">
      <img src="/swagger/logo.png" alt="logo" onerror="this.style.display='none'">
    </div>

    <div class="controls">
      <!-- TOKEN FIELD -->
      <div class="token-field">
        <!-- hash (#) as the token icon (left) -->
        <svg class="token-ic" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 9h14M5 15h14M10 3L8 21M16 3l-2 18"/>
        </svg>

        <!-- the input -->
        <input id="authToken" class="token-input" type="text" placeholder="Paste access token..."/>

        <!-- COPY BUTTON (right) -->
        <button id="copyTokenBtn" class="icon-btn" type="button" title="Copy token" aria-label="Copy token" onclick="copyToken()" disabled>
          <!-- normal «copy»: two overlapping squares -->
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <!-- bottom leaf -->
            <rect x="4.5" y="4.5" width="11" height="13" rx="2" ry="2" stroke-width="2"/>
            <!-- top leaf -->
            <rect x="8.5" y="6.5" width="11" height="13" rx="2" ry="2" stroke-width="2"/>
          </svg>
        </button>
      </div>

      <!-- BUTTONS -->
      <button class="btn btn-save" onclick="saveToken()" type="button">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
        </svg>
        <span>Save</span>
      </button>

      <button class="btn btn-clear" onclick="clearToken()" type="button">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
        </svg>
        <span>Clear</span>
      </button>

      <form class="logout" method="post" action="/swagger/logout">
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

  <div class="content">
    <div id="swagger-ui"></div>
  </div>

  <script src="/swagger/bootstrap.js" defer></script>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>
    const K = 'swagger_access_token';
    const K_OVR = 'swagger_user_override';

    function toast(msg){
      const t=document.getElementById('toast'); t.textContent=msg||'Saved';
      t.classList.add('show'); setTimeout(()=>t.classList.remove('show'),1300);
    }

    function saveToken(){
      const el=document.getElementById('authToken');
      const t=(el?.value||'').trim();
      if(t){ localStorage.setItem(K, t); localStorage.setItem(K_OVR,'1'); toast('Token saved'); }
      else { localStorage.removeItem(K); localStorage.setItem(K_OVR,'1'); toast('Token cleared'); }
      syncCopyBtn();
    }

    function clearToken(){
      const el=document.getElementById('authToken');
      if(el) el.value='';
      localStorage.removeItem(K);
      localStorage.setItem(K_OVR,'1');
      toast('Token cleared');
      syncCopyBtn();
    }

    async function copyToken(){
      const el=document.getElementById('authToken');
      const val=(el?.value||'').trim();
      if(!val){ toast('Nothing to copy'); return; }
      try{
        if(navigator.clipboard && navigator.clipboard.writeText){
          await navigator.clipboard.writeText(val);
        }else{
          // fallback
          const ta=document.createElement('textarea');
          ta.value=val; document.body.appendChild(ta); ta.select();
          document.execCommand('copy');
          document.body.removeChild(ta);
        }
        toast('Copied');
      }catch(e){
        toast('Copy failed');
      }
    }

    function syncCopyBtn(){
      const el=document.getElementById('authToken');
      const btn=document.getElementById('copyTokenBtn');
      if(btn){ btn.disabled = !(el && el.value.trim().length); }
    }

    document.addEventListener('DOMContentLoaded',()=>{
      const input = document.getElementById('authToken');
      if(input){
        input.value = localStorage.getItem(K) || input.value || '';
        input.addEventListener('input', ()=>{
          localStorage.setItem(K, input.value.trim());
          localStorage.setItem(K_OVR,'1');
          syncCopyBtn();
        });
        syncCopyBtn();
      }
    });

    // Swagger UI
    window.ui = SwaggerUIBundle({
      url:'/swagger/openapi.yaml',
      dom_id:'#swagger-ui',
      presets:[SwaggerUIBundle.presets.apis],
      persistAuthorization:true,
      requestInterceptor:(req)=>{
        const token=localStorage.getItem(K)||'';
        if(token){
          req.headers['Authorization'] = token; // add 'Bearer ' if needed
        }
        return req;
      }
    });
  </script>
</body>
</html>`
