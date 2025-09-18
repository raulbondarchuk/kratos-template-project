# 🚀 Kratos Template Project (PowerShell + Make)

Este repositorio es un **template** para proyectos basados en [Go Kratos](https://go-kratos.dev/), con soporte para:

- Generación de código con **buf**  
- Inyección de dependencias con **wire**
- Flujo de trabajo completo de **módulos** (proto, feature, repo, biz, service)  
- **Documentación OpenAPI** autogenerada
  - [Swagger](https://swagger.io/)
  - [Scalar](https://scalar.com/)
- Scripts **PowerShell** para Windows, integrados en `make`
- Soporte de dos tipos de base de datos.
  - **MySQL** (GORM)
    - Ensure - Comprobar si existe un esquema y autogeneración utilizando scripts `.sql`
    - Migrate - Migraciones de tablas y configuración de campos.
    - Seed - Autogeneración de valores por defecto utilizando scripts `.sql`
  - **Postgres** (GORM)
    - Ensure - Comprobar si existe un esquema y autogeneración utilizando scripts `.sql`
    - Migrate - Migraciones de tablas y configuración de campos.
    - Seed - Autogeneración de valores por defecto utilizando scripts `.sql`

---

## 📦 Requisitos

- [Go](https://go.dev/) ≥ 1.25 
- [buf](https://buf.build) con extención recomendada `Buf`
  - En caso de utilizar [Cursor](https://cursor.com/dashboard) hay que instalar la extención utilizando [VSIX](https://www.vsixhub.com/vsix/155966/) 
- [Protocol Buffers](https://protobuf.dev/) - Compilador de protobuf (`protoc`)
  - Windows: `choco install protoc`
- [wire](https://github.com/google/wire)
- PowerShell (Windows o [pwsh cross-platform](https://learn.microsoft.com/en-us/powershell/))  

---

## ⚡ Configuraciines `.env`

```sh
# Database (Puede ser MySQL o Postgres)
DB_DRIVER= mysql # mysql | postgres
DB_USER=root
DB_PASSWORD=passw@rd
DB_HOST=127.0.0.1
DB_PORT=3307
DB_SCHEMA=kratos-template # Nombre de la esquema

# Configuraciones adicionales para Postgres
# DB_SSLMODE=disable # disable|require|verify-ca|verify-full
# DB_TS=UTC # Europe/Madrid

# MQTT Broker
MQTT_USERNAME="usernamemqtt"                
MQTT_PASSWORD="passw@rdmqtt"
```
---

## ⚡ Flujo recomendado

Lo mejor para ejecutar el servicio por primera vez es usar `make all`.

Comando `make all` ejecuta otros `make` en orden correcto.

```sh
make all        # init -> gen -> wire -> run
# Orden de make
make init       # instala herramientas necesarias
make gen        # genera código protobuf (buf generate)
make wire       # genera inyección con wire
make run        # ejecuta con kratos run (hot reload)
```

## 📚 Flujo recomendado módulos

Al iniciar el proyecto y asegurarnos de que todas las configuraciones sean correctas, podemos crear y eliminar módulos automáticamente en nuestro proyecto.

> **Nota:** Es importante saber que, en caso de **no usar base de datos**, será necesario **comentar `data.ProviderSet`** dentro de **`cmd/service/wire.go`** para evitar **errores de wire**.

Crear un módulo completo (proto + feature + repo + biz + service + wire + generación .proto y docs):

```sh
make module name="foo"
```

Eliminar un módulo:

```sh
make module-delete name="foo"               # todas las versiones
make module-delete name="foo" version="v2"  # sólo v2
```

## 📝 Commit + versionado automático


```sh
make commit t="Titulo" d="Descripcion"
```

#### 🔧 Comandos principales 

```sh
make help      # Ayuda interactiva (colores y ejemplos)
make init      # Instalar/actualizar herramientas + go mod tidy
make deps      # Actualizar buf.lock si cambió buf.yaml
make gen       # Generar código protobuf (usa buf.gen.yaml)
make wire      # Generar inyección (wire) en cmd/service
make build     # Compilar binario en bin/
make run       # Ejecutar con kratos run
make gorun     # Ejecutar con go run directamente
make clean     # Limpiar binarios, wire_gen.go, archivos .pb.go
make docs      # Regenerar documentación (docs/ y docs/openapi)
```
#### 🔧 Comandos de módulos

```sh
make module-proto name="foo"    # Generar sólo .proto
make module-feature name="foo"  # Generar sólo feature
make module-repo name="foo"     # Generar sólo repo
make module-biz name="foo"      # Generar sólo biz
make module-service name="foo"  # Generar sólo service
make module-wire name="foo"     # Generar sólo wire
```