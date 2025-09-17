# Makefile — Guía rápida

> **Entorno:** Este proyecto usa un **Makefile para Windows** que ejecuta **PowerShell**:  
> `SHELL := powershell.exe` con `-NoProfile -ExecutionPolicy Bypass`.  
> Ejecuta los comandos desde la carpeta `service/`.

## Requisitos

- **Go** (versión reciente).
- **PowerShell** en Windows.
- Git instalado y configurado.
- Acceso a `origin` (remote por defecto).
- 
---

## Ayuda integrada

Muestra ayuda en **español** (por defecto):

```bash
make help
```

En **inglés**:

```bash
make help LANG=en
```

---

## Preparación del entorno

Instala las herramientas necesarias y actualiza dependencias:

```bash
make init
```

Esto instalará: `buf`, `protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-go-http`, `protoc-gen-openapi`, `grpc-gateway`, `protoc-gen-openapiv2`, `wire`, `kratos`, y ejecutará `go mod tidy`.

> Si hay errores de imports de Protobuf, prueba:
> ```bash
> buf --version
> buf dep update
> ```

---

## Generación de código

Genera código a partir de los protos y corre **wire**:

```bash
make gen
```

Hace:
- `buf dep update`
- `buf generate --template buf.gen.yaml`
- `wire` dentro de `cmd/service` (genera `wire_gen.go`)

---

## Compilar

Compila binario (usa script PowerShell):

```bash
make build
```

Salida: `bin/service.linux`

---

## Ejecutar

Con **Kratos** (hot reload) **(Recomendado)**:

```bash
make krun
```

Con `go run`:

```bash
make run
```

---

## Tidy y limpieza

Ordenar módulos:

```bash
make tidy
```

Limpiar binarios y archivos generados por wire:

```bash
make clean
```

---

## Commits con versionado automático

### Resumen del flujo

El objetivo `commit` llama a `scripts/ps/git-release.ps1`, que:

1. Lee **solo** la versión base desde `configs/config.yaml`, campo `app.version` (p. ej., `v1`, `v2`).  
   > Tú **no** pasas versiones manualmente.
2. Calcula el siguiente **patch** como `vX.N` consultando **origin** (remoto) para evitar colisiones entre personas (usa `git ls-remote --tags`).  
   - Ej.: si `app.version = v1`, saldrán `v1.1`, `v1.2`, …
   - Si cambias a `v2`, el siguiente será `v2.1`, etc.
3. **Política de ramas**:
   - **Prohibido** en `main` y `master`.
   - **Prohibido** en estado **detached HEAD**.
4. Flujo de push:
   - Empuja primero **la rama actual** (crea upstream si no existe).
   - Empuja **el tag** con **reintento automático** si el nombre ya existe (elige `vX.(N+1)` y reintenta).
5. Si no hay cambios (`git status` vacío), sale sin error.

El **mensaje del tag** es **igual que el mensaje de commit** (título + descripción).

### Uso (sintaxis completa, corta `t=` / `d=`)

```bash
make commit t="Título del commit" d="Descripción del commit"
```

Alias (Hagan lo mismo):

```bash
make release t="Título" d="Descripción"
```

```bash
make commit t="Título" d="Descripción"
```

Compatibilidad: si alguna vez usas `TITLE=` / `DESC=`, también funciona (pero se recomienda `t=` y `d=`).

> ⚠️ PowerShell: si tu texto incluye el signo `$`, escápalo con acento grave:
> ```powershell
> make commit t="feat: ingest LTA" d="support `\$MSG:11"
> ```

---

## Ejemplos rápidos

```bash
# 1) Preparar herramientas y dependencias
make init

# 2) Generar código (protos + wire)
make gen

# 3) Compilar
make build

# 4) Ejecutar
make run
# o
make krun

# 5) Commit + tag + push automático (vX.N desde config.yaml)
make commit t="fix: reconexión mqtt" d="backoff exponencial y logs"

# 6) Limpiar (Si es necesario)
make clean
```

---

## Personalización

- Cambiar nombre de app/binario:
  ```bash
  make build APP_NAME_TO_BUILD=otro CMD_DIR=cmd/otro
  ```
- Cambiar rutas (si mueves scripts o config):
  - `RELEASE_SCRIPT := ./scripts/ps/git-release.ps1`
  - `CONFIG_PATH := ./configs/config.yaml`
- Cambiar idioma de la ayuda:
  ```bash
  make help LANG=en
  ```

---

## Solución de problemas

- **“Script not found”** al hacer commit:  
  Asegúrate de que `RELEASE_SCRIPT` apunta a `./scripts/ps/git-release.ps1` **desde la carpeta `service/`**.
- **Error por rama bloqueada**:  
  Cambia a una rama de feature (`git checkout -b feature/...`). `main`/`master` no están permitidas.
- **Colisión de tag al empujar**:  
  El script reintenta automáticamente calculando un nuevo patch. Si ves errores repetidos, ejecuta:
  ```powershell
  git fetch --tags
  ```
  y vuelve a intentarlo.
- **Caracteres especiales en PowerShell**:  
  Escapa `$` con acento grave \` como se mostró arriba.

---

## Notas

- Este Makefile está diseñado para **Windows/PowerShell**.
- Si usas otro entorno, adapta el `SHELL` o ejecuta los scripts directamente con PowerShell.  
- El control de versiones **siempre** se deriva de `configs/config.yaml` → `app.version`.



---
NOTAS: 

(Ctrl+Shift+P) -> "Go: Restart Language Server" para wire 