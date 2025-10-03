package swagger

import (
	"html/template"
	stdhttp "net/http"
)

func serveLoginPage(w stdhttp.ResponseWriter, _ *stdhttp.Request, errMsg string) {
	setOnlyContentType(w, "text/html; charset=utf-8")
	setNoCache(w)
	w.WriteHeader(stdhttp.StatusOK)
	_ = loginTpl.Execute(w, struct{ Error string }{Error: errMsg})
}

var loginTpl = template.Must(template.New("login").Parse(`<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>Auth | API Docs</title>
  <style>
    :root{
      --bg:#f3f4f6; --card:#ffffff;
      --text:#111827; --border:#d1d5db;
      --green:#10b981; --green-h:#059669;
      --red:#ef4444; --muted:#6b7280;
    }
    *{box-sizing:border-box}
    html,body{margin:0;height:100%;font-family:ui-sans-serif,system-ui,-apple-system,Segoe UI,Roboto,Ubuntu}
    body{background:var(--bg);display:flex;align-items:center;justify-content:center;min-height:100vh;padding:20px}
    .card{width:100%;max-width:420px;background:var(--card);border:1px solid var(--border);border-radius:20px;box-shadow:0 12px 40px rgba(0,0,0,.1);padding:36px 32px;animation:fadeIn .6s ease}
    @keyframes fadeIn{from{opacity:0;transform:translateY(12px)}to{opacity:1;transform:none}}
    .logo{text-align:center;margin-bottom:20px}
    .logo img{height:50px}
    .title{text-align:center;font-weight:800;font-size:24px;margin-bottom:6px;color:var(--text)}
    .subtitle{text-align:center;color:var(--muted);font-size:14px;margin-bottom:20px}
    .err{color:#fff;background:var(--red);font-weight:600;font-size:14px;margin:0 0 16px;padding:10px 12px;border-radius:10px;display:{{if .Error}}block{{else}}none{{end}};text-align:center}
    form{display:grid;gap:16px}
    label{font-weight:600;font-size:13px;margin-bottom:4px;display:block;color:var(--text)}
    .input-wrap{position:relative}
    .input-wrap svg.icon{position:absolute;left:12px;top:50%;transform:translateY(-50%);width:20px;height:20px;stroke:#9ca3af;pointer-events:none}
    input{width:100%;padding:10px 48px 10px 44px;border:1px solid var(--border);border-radius:12px;background:#fff;color:var(--text);outline:none;transition:border .2s,box-shadow .2s}
    input:focus{border-color:var(--green);box-shadow:0 0 0 3px rgba(16,185,129,.2)}
    .toggle-eye{position:absolute;right:6px;top:50%;transform:translateY(-50%);background:none;border:none;cursor:pointer;padding:6px;width:36px;height:36px;border-radius:8px;display:flex;align-items:center;justify-content:center}
    .toggle-eye:focus-visible{outline:2px solid rgba(16,185,129,.6);outline-offset:2px}
    .toggle-eye svg{width:22px;height:22px;stroke:#6b7280;fill:none}
    .eye-core,.eye-slash{transition:opacity .18s ease, transform .18s ease; transform-origin:center}
    .toggle-eye[data-state="hidden"] .eye-core{opacity:1;transform:scale(1)}
    .toggle-eye[data-state="hidden"] .eye-slash{opacity:1;transform:scale(1)}
    .toggle-eye[data-state="shown"]  .eye-core{opacity:1;transform:scale(1)}
    .toggle-eye[data-state="shown"]  .eye-slash{opacity:0;transform:scale(.94)}

    /* Button + Loader */
    .btn{
      position:relative;display:inline-flex;align-items:center;justify-content:center;gap:8px;
      border:0;border-radius:12px;padding:12px 14px;font-weight:700;cursor:pointer;
      color:#fff;background:var(--green);margin-top:8px;transition:background .2s,transform .1s,opacity .2s; min-height:44px;
    }
    .btn:hover{background:var(--green-h)}
    .btn:active{transform:scale(.98)}
    .btn[disabled]{opacity:.7;cursor:not-allowed}
    .btn .spinner{
      width:16px;height:16px;border:2px solid rgba(255,255,255,.5);border-top-color:#fff;border-radius:50%;
      animation:spin .9s linear infinite;opacity:0;transform:scale(.8);
      transition:opacity .15s ease; margin-right:6px;
    }
    .btn.loading .spinner{opacity:1}
    @keyframes spin{to{transform:rotate(360deg)}}
  </style>
</head>
<body>
  <div class="card">
    <div class="logo">
      <img src="/swagger/logo.png" alt="logo" onerror="this.style.display='none'">
    </div>
    <div class="title">Autenticación requerida</div>
    <div class="subtitle">Introduce tus credenciales para acceder a la documentación técnica del servicio.</div>
    <div class="err">{{.Error}}</div>

    <form id="loginForm" method="post" action="/swagger/login">
      <div>
        <label>Usuario</label>
        <div class="input-wrap">
          <svg class="icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.121 17.804A13.937 13.937 0 0112 15c2.5 0 4.847.655 6.879 1.804M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
          </svg>
          <input type="text" name="username" autocomplete="username" required />
        </div>
      </div>

      <div>
        <label>Contraseña</label>
        <div class="input-wrap">
          <svg class="icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10V8a4 4 0 118 0v2"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 10h12v8a2 2 0 01-2 2H8a2 2 0 01-2-2v-8z"/>
          </svg>
          <input id="passwordInput" type="password" name="password" autocomplete="current-password" required />
          <button
            type="button"
            id="toggleBtn"
            class="toggle-eye"
            data-state="hidden"
            aria-label="Show password"
            aria-pressed="false"
            onclick="togglePassword()">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <g class="eye-core">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5s8.268 2.943 9.542 7c-1.274 4.057-5.065 7-9.542 7S3.732 16.057 2.458 12z"/>
                <circle cx="12" cy="12" r="3" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </g>
              <path class="eye-slash" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3l18 18"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Button + Loader -->
      <button id="loginBtn" class="btn" type="submit" aria-live="polite">
        <span class="spinner" aria-hidden="true"></span>
        <span class="btn-label">Acceder</span>
      </button>
    </form>
  </div>

  <script>
    // toggle password visibility
    function togglePassword(){
      const input = document.getElementById('passwordInput');
      const btn   = document.getElementById('toggleBtn');
      const isShown = btn.getAttribute('data-state') === 'shown';
      if (isShown){
        input.type = 'password';
        btn.setAttribute('data-state','hidden');
        btn.setAttribute('aria-label','Show password');
        btn.setAttribute('aria-pressed','false');
      } else {
        input.type = 'text';
        btn.setAttribute('data-state','shown');
        btn.setAttribute('aria-label','Hide password');
        btn.setAttribute('aria-pressed','true');
      }
    }

    // block the button and show the loader during submission
    (function(){
      const form = document.getElementById('loginForm');
      const btn  = document.getElementById('loginBtn');
      if(!form || !btn) return;

      form.addEventListener('submit', function(){
        // protect against double submission
        if (btn.disabled) return;

        btn.disabled = true;
        btn.classList.add('loading');
        const label = btn.querySelector('.btn-label');
        if (label) label.textContent = 'Accediendo…';

        // freeze the fields, so the user cannot change them during submission
        form.querySelectorAll('input').forEach(function(i){ i.setAttribute('readonly','true'); });
      });
    })();
  </script>
</body>
</html>`))
