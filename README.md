# 🚀 Kratos-Template Project v3

![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)
![Kratos](https://img.shields.io/badge/Kratos-v2.8.4-green.svg)

Un proyecto template completo basado en **Kratos Framework** para desarrollar microservicios robustos en Go, con soporte para HTTP/gRPC, base de datos, MQTT, webhooks y mucho más.

## 📋 Tabla de Contenidos

- [✨ Características](#-características)
- [🏗️ Arquitectura](#️-arquitectura)
- [🚀 Inicio Rápido](#-inicio-rápido)
- [⚙️ Comandos Make](#️-comandos-make)
- [📁 Estructura del Proyecto](#-estructura-del-proyecto)
- [🔧 Configuración](#-configuración)
- [📦 Módulos y Workflow](#-módulos-y-workflow)
- [🧪 Testing](#-testing)
- [📚 Ejemplos](#-ejemplos)
- [🤝 Contribución](#-contribución)

## ✨ Características

### 🔥 Principales
- **🏛️ Arquitectura Hexagonal** - Clean Architecture con separación clara de responsabilidades
- **⚡ Kratos Framework** - Framework moderno para microservicios en Go
- **🌐 HTTP + gRPC** - Servidores duales con soporte completo
- **📊 Inyección de Dependencias** - Google Wire para DI automática
- **🔐 Autenticación** - PASETO tokens para seguridad robusta
- **📡 MQTT Support** - Cliente MQTT integrado con reconexión automática
- **🔗 Webhooks** - Sistema de webhooks configurable
- **📖 OpenAPI/Swagger** - Documentación automática de APIs
- **🗄️ Multi-Database** - Soporte para MySQL y PostgreSQL
- **🔄 Migraciones** - Sistema de migraciones automáticas
- **🌱 Seeding** - Datos iniciales automáticos
- **📝 Logging** - Logrus con formato estructurado

### 🛠️ Herramientas
- **📦 Buf** - Generación de código Protocol Buffers
- **🧪 Testing** - Tests unitarios automáticos
- **📊 Métricas** - Prometheus metrics integradas
- **🔧 Make** - Automatización completa del workflow
- **📋 PowerShell** - Scripts optimizados para Windows

## 🏗️ Arquitectura

┌──────────────────────────────────────────────────────────────┐
│                    🌐 Capa de Presentación                   
├──────────────────────────────────────────────────────────────┤
│  Servidor HTTP  │  Servidor gRPC  │  Cliente MQTT  │ Webhooks 
├──────────────────────────────────────────────────────────────┤
│                    🏢 Capa de Lógica de Negocio             
├──────────────────────────────────────────────────────────────┤
│  Servicios  │  Casos de Uso  │  Validadores  │  Autenticación 
├──────────────────────────────────────────────────────────────┤
│                    💾 Capa de Acceso a Datos                 
├──────────────────────────────────────────────────────────────┤
│  Repositorios  │  Base de Datos  │  Migraciones  │  Seeds     
└──────────────────────────────────────────────────────────────┘

| Capa | Componentes | Descripción |
|------|--------------|-------------|
| 🌐 **Capa de Presentación** | Servidor HTTP, Servidor gRPC, Cliente MQTT, Webhooks | Proporciona la interacción con el mundo exterior: APIs REST/gRPC, MQTT y Webhooks. |
| 🏢 **Capa de Lógica de Negocio** | Servicios, Casos de Uso, Validadores, Autenticación | Contiene la lógica principal del negocio, las reglas y la autenticación. |
| 💾 **Capa de Acceso a Datos** | Repositorios, Base de Datos, Migraciones, Seeds | Gestiona la persistencia de datos y las operaciones con la base de datos. |

### 📦 Componentes Principales

- **`cmd/main.go`** - Punto de entrada de la aplicación
- **`internal/feature/`** - Lógica de negocio y casos de uso
- **`internal/data/`** - Capa de acceso a datos
- **`internal/server/`** - Servidores HTTP/gRPC
- **`internal/middleware/`** - Middleware personalizado
- **`pkg/`** - Utilidades y librerías compartidas

## 🚀 Inicio Rápido

### 📋 Prerrequisitos

- **Go 1.25.0+**
- **Make** (para Windows: [GNU Make](https://www.gnu.org/software/make/))
- **PowerShell 5.1+**
- **Docker** (opcional, para bases de datos)

### 🛠️ Instalación

1. **Clonar el repositorio**
2. **Configurar el entorno**
3. **Inicializar el proyecto**
```bash
# Instalar herramientas y dependencias
make init

# Generar código protobuf
make gen

# Configurar inyección de dependencias
make wire

# Ejecutar la aplicación
make run
```

4. **Verificar la instalación**
```bash
# Verificar que el servidor HTTP está corriendo
curl http://localhost:8000/health

# Verificar que el servidor gRPC está corriendo
grpcurl -plaintext localhost:9000 list
```

## ⚙️ Comandos Make

El proyecto incluye un sistema completo de comandos Make optimizado para **Windows PowerShell**.

### 🎯 Comandos Principales

| Comando | Descripción | Ejemplo |
|---------|-------------|---------|
| `make help` | Mostrar ayuda completa | `make help` |
| `make init` | Instalar herramientas y dependencias | `make init` |
| `make all` | Ejecutar todo el flujo (init → gen → wire → run) | `make all` |
| `make deps` | Actualizar dependencias de protobuf | `make deps` |
| `make gen` | Generar código protobuf | `make gen` |
| `make wire` | Generar inyección de dependencias | `make wire` |
| `make run` | Ejecutar con Kratos (hot reload) | `make run` |
| `make gorun` | Ejecutar con go run | `make gorun` |
| `make docs` | Regenerar documentación | `make docs` |

### 🔄 Flujo Recomendado

```bash
# 1. Configuración inicial (solo una vez)
make init

# 2. Desarrollo diario
make all          # Ejecutar todo el flujo
# o paso a paso:
make deps         # Actualizar dependencias
make gen          # Generar código
make wire         # Configurar DI
make run          # Ejecutar aplicación
```

### 🏷️ Gestión de Versiones

```bash
# Commit con versión automática
make commit t="Nueva funcionalidad" d="Descripción detallada"
```

## 📁 Estructura del Proyecto

```
service/
├── 📁 cmd/ # Punto de entrada
│   ├── main.go             # Aplicación principal
│   ├── wire.go             # Configuración Wire
│   └── wire_gen.go         # Código generado por Wire
├── 📁 configs/             # Configuraciones
│   └── config.yaml         # Configuración principal
├── 📁 docs/                # Documentación
│   ├── openapi/            # OpenAPI/Swagger
│   └── logo.png            # Logo del proyecto
├── 📁 internal/            # Código interno
│   ├── 📁 conf/v1/          # Configuración protobuf (la estructure de config.yaml)
│   ├── 📁 data/             # Capa de datos
│   │   ├── adapters/        # Adaptadores de BD
│   │   ├── model/           # Modelos de datos
│   │   └── migrations/      # Migraciones
│   ├── 📁 feature/          # Lógica de negocio
│   ├── 📁 middleware/       # Middleware personalizado
│   ├── 📁 out/              # Salidas externas
│   │   ├── broker/          # MQTT broker
│   │   └── webhooks/        # Sistema de webhooks
│   └── 📁 server/           # Servidores
│       ├── grpc/            # Servidor gRPC
│       └── http/            # Servidor HTTP
├── 📁 pkg/                  # Utilidades compartidas
│   ├── generic/             # Utilidades genéricas
│   ├── logger/              # Sistema de logging
│   ├── mqtt/                # Cliente MQTT
│   └── utils/               # Utilidades generales
├── 📁 scripts/              # Scripts de automatización
│   ├── ps/                  # Scripts PowerShell
│   ├── mysql/               # Scripts MySQL
│   └── postgres/            # Scripts PostgreSQL
├── Makefile                 # Comandos de automatización
├── go.mod                   # Dependencias Go
├── buf.yaml                 # Configuración Buf
└── README.md                # Esta documentación
```

## 🔧 Configuración

### 📄 Archivo de Configuración Principal

El archivo `configs/config.yaml` contiene toda la configuración de la aplicación:

```yaml
app:
  mode: dev                    # Modo: dev, prod, test
  name: kratos-template        # Nombre de la aplicación
  version: v3                  # Versión de la aplicación

server:
  http:
    addr: 0.0.0.0:8000        # Dirección del servidor HTTP
    timeout: 1s               # Timeout de requests
  grpc:
    addr: 0.0.0.0:9000        # Dirección del servidor gRPC
    timeout: 1s               # Timeout de requests

data:
  database:
    active: false             # Activar base de datos
    migrations: false         # Aplicar migraciones
    seed: false              # Llenar con datos iniciales
  mqtt:
    active: false             # Activar cliente MQTT
    source: "tcp://10.70.20.40:1883"  # URI del broker MQTT
    client_id: "client_kratos_template"  # ID del cliente
    max_reconnect_interval: "60s"        # Intervalo de reconexión
    topics:                   # Topics para suscripción
      - "receiver/ltm/#"
      - "receiver/lta/#"
      - "receiver/ltc/#"
      - "receiver/scr/#"
    publish:                  # Topics para publicación
      topic1: "topic1/test"
      topic2: "topic2/test"

webhooks:
  webhook:
    url: http://localhost:3000  # URL base del webhook
    timeout: 5s                # Timeout del webhook
    routes:                    # Rutas configuradas
      route1: /v1/route1
      route2: /v2/route2
```

### 🗄️ Configuración de Base de Datos

#### MySQL
```yaml
data:
  database:
    active: true
    migrations: true
    seed: true
  mysql:
    dsn: "user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"
```

#### PostgreSQL
```yaml
data:
  database:
    active: true
    migrations: true
    seed: true
  postgres:
    dsn: "host=localhost user=user password=password dbname=database port=5432 sslmode=disable"
```

### 📡 Configuración MQTT

```yaml
data:
  mqtt:
    active: true
    source: "tcp://localhost:1883"
    client_id: "kratos_client"
    max_reconnect_interval: "30s"
    topics:
      - "sensors/#"
      - "devices/+/status"
    publish:
      notifications: "notifications/all"
      alerts: "alerts/critical"
```

## 📦 Módulos y Workflow

### 🏗️ Creación de Módulos

El sistema incluye un workflow completo para crear módulos nuevos:

#### 📝 Crear Módulo Completo

```bash
# Crear módulo sin endpoints (solo mock)
make module name="user"

# Crear módulo con operaciones específicas
make module name="user" ops="get,upsert,delete"

# Crear módulo con operaciones personalizadas
make module name="product" ops="get,create,update"
```

#### 🔧 Operaciones Disponibles

| Operación | Alias | Descripción | Archivos Generados |
|-----------|-------|-------------|-------------------|
| `get` | `find`, `list`, `read` | Obtener/listar recursos | `s_find.go`, `r_list.go` |
| `upsert` | `create`, `update` | Crear/actualizar recursos | `s_upsert.go`, `r_upsert.go` |
| `delete` | `del`, `remove` | Eliminar recursos | `s_delete_by_id.go`, `r_delete_by_id.go` |

#### 📁 Estructura de Módulo Generado

```
internal/
├── feature/user/v1/
│   ├── /biz              # Lógica de negocio
│   ├── /repo             # Implementación repo
│   ├── /service           # Handlers
└── api/user/v1/              # API protobuf
    └── user.proto            # Definición de la API
```

### 🗑️ Gestión de Módulos

```bash
# Eliminar módulo completo (todas las versiones)
make module-delete name="user"

# Eliminar versión específica
make module-delete name="user" version="v2"

```

### 🧪 Testing de Módulos

```bash
# Generar tests para la última versión
make tmodule name="user"

# Generar tests para versión específica
make tmodule name="user" version="v2"

# Regenerar tests (sobrescribir existentes)
make tmodule name="user" force=1

# Eliminar tests
make tmodule-delete name="user"
make tmodule-delete name="user" version="v2"
```

## 🧪 Testing

### 🔧 Configuración de Tests

```bash
# Ejecutar todos los tests
go test ./...

# Ejecutar tests con cobertura
go test -cover ./...

# Ejecutar tests específicos
go test ./internal/feature/user/...

# Ejecutar tests con verbose
go test -v ./...
```

### 📊 Generación de Tests Automáticos

```bash
# Generar tests para un módulo
make tmodule name="user"

# Los tests se generan en:
# internal/feature/user/v1/service/user_test.go

```

### 🎯 Ejemplos de Tests

```go
func TestUserService_FindUser(t *testing.T) {
    // Setup
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    // Test
    user, err := service.FindUser(context.Background(), "user-id")
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "user-id", user.ID)
}
```

## 📚 Ejemplos

### 🌐 Ejemplo: API HTTP

```bash
# Crear módulo de usuarios
make module name="user" ops="get,upsert"

# Ejecutar aplicación
make run

# Probar endpoints
curl -X GET http://localhost:8000/v1/users
curl -X POST http://localhost:8000/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Juan", "email": "juan@example.com"}'
```

### 🔗 Ejemplo: API gRPC

```bash
# Usar grpcurl para probar gRPC
grpcurl -plaintext localhost:9000 list
grpcurl -plaintext localhost:9000 user.v1.UserService/FindUser \
  -d '{"id": "user-123"}'
```

### 📡 Ejemplo: MQTT

```yaml
# configs/config.yaml
data:
  mqtt:
    active: true
    source: "tcp://localhost:1883"
    client_id: "kratos_client"
    topics:
      - "sensors/temperature"
      - "sensors/humidity"
    publish:
      alerts: "alerts/system"
```

```go
// En tu código
func (s *SensorService) PublishAlert(ctx context.Context, alert *Alert) error {
    return s.mqttClient.Publish(ctx, "alerts/system", alert)
}
```

### 🔗 Ejemplo: Webhooks

```yaml
# configs/config.yaml
webhooks:
  webhook:
    url: http://localhost:3000
    timeout: 5s
    routes:
      user_created: /v1/webhooks/user/created
      user_updated: /v1/webhooks/user/updated
```

```go
// En tu código
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    // Crear usuario
    err := s.repo.CreateUser(ctx, user)
    if err != nil {
        return err
    }
    
    // Enviar webhook
    return s.webhookClient.Send(ctx, "user_created", user)
}
```

### 🗄️ Ejemplo: Base de Datos

```go
// Modelo de datos
type User struct {
    ID        string    `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"not null" json:"name"`
    Email     string    `gorm:"uniqueIndex" json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Repositorio
func (r *userRepository) CreateUser(ctx context.Context, user *User) error {
    return r.DB(ctx).Create(user).Error
}

func (r *userRepository) FindUser(ctx context.Context, id string) (*User, error) {
    var user User
    err := r.db.DB(ctx).Where("id = ?", id).First(&user).Error
    return &user, err
}
```

## 🤝 Contribución

### 📋 Guías de Contribución

1. **Fork** el repositorio
2. **Crear** una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. **Commit** tus cambios (`git commit -m 'Agregar nueva funcionalidad'`)
4. **Push** a la rama (`git push origin feature/nueva-funcionalidad`)
5. **Abrir** un Pull Request

### 📝 Estándares de Código

- **Go**: Seguir las convenciones estándar de Go
- **Commits**: Usar mensajes descriptivos en español
- **Documentación**: Documentar funciones públicas
- **Tests**: Escribir tests para nueva funcionalidad

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.

## 🙏 Agradecimientos

- [Kratos Framework](https://github.com/go-kratos/kratos) - Framework principal
- [Google Wire](https://github.com/google/wire) - Inyección de dependencias
- [Buf](https://buf.build/) - Herramientas Protocol Buffers
- [GORM](https://gorm.io/) - ORM para Go

---

<div align="center">

**¡La documentación fue construida con AI!**

[📖 Documentación](docs/) • [🐛 Reportar Bug](issues/) • [💡 Solicitar Feature](issues/) • [💬 Discusiones](discussions/)

</div>
