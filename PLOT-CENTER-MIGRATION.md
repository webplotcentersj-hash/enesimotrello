# 🏭 Migración Completa: Sistema Plot Center (PHP → Go + React)

## 📋 Resumen del Sistema

El sistema actual es una **plataforma completa de gestión de taller gráfico** con los siguientes módulos:

### Módulos Principales

1. **Gestión de Órdenes de Trabajo (Kanban Board)**
   - Sistema de estados: Diseño Gráfico, En Espera, Taller Gráfico, etc.
   - Prioridades: Alta, Normal, Baja
   - Complejidad: Simple, Media, Compleja
   - Sectores y materiales asociados
   - Archivos adjuntos (imágenes, documentos)
   - Historial de movimientos
   - Tracking de usuario trabajando

2. **Sistema de Chat General**
   - Salas de chat (públicas/privadas)
   - Mensajes en tiempo real
   - Notificaciones

3. **Plot AI (Asistente Virtual)**
   - Integración con Google Gemini
   - Chat con contexto del sistema
   - Acceso a contactos y soportes

4. **Sistema de Usuarios y Autenticación**
   - Login/Logout
   - Roles: administración, mostrador, operario, etc.
   - Permisos por rol

5. **Sistema de Notificaciones y Alertas**
   - Alertas de órdenes estancadas (>3 días)
   - Notificaciones de usuarios
   - Sistema de alertas enviadas

6. **Sistema de Estadísticas**
   - Dashboard con métricas
   - Análisis de productividad

7. **Gestión de Archivos**
   - Subida de archivos adjuntos
   - Optimización de imágenes
   - Almacenamiento en carpeta `uploads/`

## 🗄️ Estructura de Base de Datos

### Tablas Principales

```sql
-- Usuarios
usuarios (id, nombre, password_hash, rol, last_seen)

-- Órdenes de Trabajo
ordenes_trabajo (
    id, numero_op, cliente, descripcion,
    fecha_creacion, fecha_entrega, estado,
    prioridad, operario_asignado, complejidad,
    sector, hora_estimada_entrega,
    id_usuario_creador,
    usuario_trabajando_id,
    usuario_trabajando_nombre,
    timestamp_inicio_trabajo
)

-- Materiales
materiales (id, descripcion)
orden_materiales (id_orden, id_material, cantidad)

-- Sectores
sectores (id, nombre, color, orden_visualizacion)
orden_sectores (id_orden, id_sector)

-- Archivos
archivos_adjuntos (id, id_orden, nombre_archivo, nombre_original)

-- Historial
historial_movimientos (
    id, id_orden, estado_anterior, estado_nuevo,
    usuario, timestamp, observaciones
)

-- Chat
chat_rooms (id, nombre, tipo, created_at)
chat_messages (
    id, room_id, id_usuario, nombre_usuario,
    mensaje, timestamp, message_type
)

-- Notificaciones
user_notifications (
    id, user_id, message, type, is_read, created_at
)
alertas_enviadas (id, id_orden, tipo_alerta)
```

## 🏗️ Arquitectura Propuesta (Go + React)

```
┌─────────────────────────────────────────┐
│      Frontend React (Vercel)            │
│  ┌──────────────┬────────────────────┐  │
│  │  TaskBoard   │  Plot Center       │  │
│  │  (Kanban)    │  (Gestión Taller)  │  │
│  └──────────────┴────────────────────┘  │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│      Backend Go (Gin) - VPS             │
│  ┌────────────────────────────────────┐ │
│  │  TaskBoard API                     │ │
│  │  Plot Center API                   │ │
│  │  - Orders Management               │ │
│  │  - Chat System                     │ │
│  │  - Plot AI                         │ │
│  │  - Notifications                   │ │
│  │  - File Management                 │ │
│  └────────────────────────────────────┘ │
│  ┌────────────────────────────────────┐ │
│  │  Shared: JWT Auth, PostgreSQL      │ │
│  └────────────────────────────────────┘ │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│  External Services                       │
│  - Google Gemini API                    │
│  - PostgreSQL (compartida)              │
│  - Redis (compartido)                   │
│  - File Storage (uploads/)              │
└─────────────────────────────────────────┘
```

## 📁 Estructura de Carpetas Propuesta

```
task-board/
├── backend/
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── order.go          (nuevo)
│   │   │   ├── material.go       (nuevo)
│   │   │   ├── sector.go         (nuevo)
│   │   │   ├── chat_room.go      (nuevo)
│   │   │   ├── notification.go   (nuevo)
│   │   │   ├── file.go           (nuevo)
│   │   │   └── plot_ai.go        (ya creado)
│   │   ├── handler/
│   │   │   ├── order_handler.go      (nuevo)
│   │   │   ├── material_handler.go   (nuevo)
│   │   │   ├── chat_handler.go       (nuevo)
│   │   │   ├── notification_handler.go (nuevo)
│   │   │   ├── file_handler.go       (nuevo)
│   │   │   └── plot_ai_handler.go    (ya creado)
│   │   ├── service/
│   │   │   ├── order_service.go      (nuevo)
│   │   │   ├── material_service.go   (nuevo)
│   │   │   ├── chat_service.go       (nuevo)
│   │   │   ├── notification_service.go (nuevo)
│   │   │   ├── file_service.go       (nuevo)
│   │   │   ├── plot_ai_service.go    (ya creado)
│   │   │   └── gemini_service.go     (ya creado)
│   │   └── repository/
│   │       ├── order_repository.go      (nuevo)
│   │       ├── material_repository.go   (nuevo)
│   │       ├── chat_repository.go       (nuevo)
│   │       ├── notification_repository.go (nuevo)
│   │       └── plot_ai_repository.go    (ya creado)
│   └── pkg/
│       └── storage/ (nuevo - gestión de archivos)
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── TaskBoard/ (ya existe)
│   │   │   ├── PlotCenter/
│   │   │   │   ├── Orders/
│   │   │   │   │   ├── KanbanBoard.tsx (nuevo)
│   │   │   │   │   ├── OrderCard.tsx (nuevo)
│   │   │   │   │   ├── OrderDetail.tsx (nuevo)
│   │   │   │   │   └── OrderForm.tsx (nuevo)
│   │   │   │   ├── Chat/
│   │   │   │   │   ├── ChatRoom.tsx (nuevo)
│   │   │   │   │   ├── ChatMessage.tsx (nuevo)
│   │   │   │   │   └── ChatInput.tsx (nuevo)
│   │   │   │   ├── PlotAI/ (ya iniciado)
│   │   │   │   ├── Dashboard/
│   │   │   │   │   └── StatsDashboard.tsx (nuevo)
│   │   │   │   └── Notifications/
│   │   │   │       └── NotificationCenter.tsx (nuevo)
│   │   │   └── ...
│   │   └── services/
│   │       ├── orders.ts (nuevo)
│   │       ├── chat.ts (nuevo)
│   │       ├── plotAI.ts (ya iniciado)
│   │       └── ...
│   └── ...
│
└── data/
    ├── manual_entrenamiento.txt
    ├── contactos_plotcenter.txt
    └── soportes_publicitarios.txt
```

## 🔧 Plan de Migración por Fases

### Fase 1: Base de Datos y Autenticación ✅
- [x] Analizar estructura de BD actual
- [ ] Crear migraciones de base de datos
- [ ] Migrar sistema de usuarios y roles
- [ ] Implementar JWT (ya hecho para TaskBoard)

### Fase 2: Gestión de Órdenes de Trabajo
- [ ] Crear modelos de dominio (Order, Material, Sector)
- [ ] Implementar repositorios
- [ ] Crear servicios de negocio
- [ ] Implementar handlers REST
- [ ] Crear componente Kanban Board en React
- [ ] Implementar drag & drop
- [ ] Sistema de filtros y búsqueda

### Fase 3: Sistema de Chat
- [ ] Crear modelos (ChatRoom, ChatMessage)
- [ ] Implementar WebSocket para tiempo real
- [ ] Crear componentes de chat en React
- [ ] Integrar con sistema de notificaciones

### Fase 4: Plot AI
- [x] Servicio Gemini (ya creado)
- [x] Handler básico (ya creado)
- [ ] Completar integración con contexto del sistema
- [ ] Frontend React (ya iniciado)

### Fase 5: Gestión de Archivos
- [ ] Sistema de almacenamiento de archivos
- [ ] API de subida/descarga
- [ ] Optimización de imágenes
- [ ] Componentes de gestión de archivos

### Fase 6: Notificaciones y Alertas
- [ ] Sistema de notificaciones en tiempo real
- [ ] Agentes de monitoreo (órdenes estancadas)
- [ ] Dashboard de alertas

### Fase 7: Estadísticas y Dashboard
- [ ] Endpoints de estadísticas
- [ ] Componentes de visualización
- [ ] Gráficos y métricas

### Fase 8: Integración y Testing
- [ ] Integrar todos los módulos
- [ ] Testing end-to-end
- [ ] Migración de datos existentes
- [ ] Deployment

## 🔐 Autenticación y Roles

### Roles del Sistema
- `administracion` / `administración` - Acceso completo
- `mostrador` - Crear/editar órdenes, ver chat
- `operario` - Ver órdenes asignadas, actualizar estados
- `diseñador` - Ver órdenes en diseño, actualizar

### Permisos por Módulo

| Módulo | Administración | Mostrador | Operario | Diseñador |
|--------|---------------|-----------|----------|-----------|
| Ver Órdenes | ✅ | ✅ | ✅ (asignadas) | ✅ (diseño) |
| Crear Órdenes | ✅ | ✅ | ❌ | ❌ |
| Editar Órdenes | ✅ | ✅ | ✅ (estado) | ✅ (diseño) |
| Eliminar Órdenes | ✅ | ✅ | ❌ | ❌ |
| Subir Archivos | ✅ | ✅ | ✅ | ✅ |
| Chat | ✅ | ✅ | ✅ | ✅ |
| Plot AI | ✅ | ✅ | ✅ | ✅ |
| Estadísticas | ✅ | ❌ | ❌ | ❌ |

## 📊 Endpoints API Propuestos

### Órdenes de Trabajo
```
GET    /api/v1/orders              - Listar órdenes (con filtros)
GET    /api/v1/orders/:id          - Obtener orden
POST   /api/v1/orders              - Crear orden
PUT    /api/v1/orders/:id          - Actualizar orden
DELETE /api/v1/orders/:id          - Eliminar orden
POST   /api/v1/orders/:id/move     - Mover orden (cambiar estado)
GET    /api/v1/orders/:id/history  - Historial de movimientos
```

### Materiales y Sectores
```
GET    /api/v1/materials           - Listar materiales
GET    /api/v1/sectors             - Listar sectores
```

### Chat
```
GET    /api/v1/chat/rooms          - Listar salas
GET    /api/v1/chat/rooms/:id/messages - Mensajes de sala
POST   /api/v1/chat/rooms/:id/messages - Enviar mensaje
WS     /api/v1/chat/ws             - WebSocket para tiempo real
```

### Archivos
```
POST   /api/v1/files/upload        - Subir archivo
GET    /api/v1/files/:id           - Descargar archivo
DELETE /api/v1/files/:id           - Eliminar archivo
```

### Notificaciones
```
GET    /api/v1/notifications       - Listar notificaciones
PUT    /api/v1/notifications/:id/read - Marcar como leída
```

## 🔄 Migración de Datos

### Estrategia
1. **Migración paralela**: Mantener PHP funcionando durante migración
2. **Sincronización**: Scripts para migrar datos existentes
3. **Validación**: Verificar integridad de datos migrados
4. **Cutover**: Cambio gradual por módulos

### Scripts de Migración
- Migrar usuarios
- Migrar órdenes de trabajo
- Migrar materiales y sectores
- Migrar archivos adjuntos
- Migrar historial
- Migrar mensajes de chat

## 🚀 Deployment

### Configuración Docker Compose
- Agregar volumen para `uploads/`
- Configurar variables de entorno
- Servir archivos estáticos

### Nginx
- Configurar rutas para archivos
- Proxy para API
- WebSocket para chat en tiempo real

## 📝 Variables de Entorno

```env
# Base de datos (compartida)
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=plotcenter

# Gemini AI
GEMINI_API_KEY=tu-api-key
GEMINI_MODEL=gemini-pro

# File Storage
UPLOAD_DIR=/data/uploads
MAX_FILE_SIZE=50MB

# JWT
JWT_SECRET=tu-secret
JWT_EXPIRY=24h

# CORS
CORS_ORIGIN=https://plotcenter.vercel.app
```

## ⏱️ Estimación de Tiempo

- **Fase 1**: 2-3 días
- **Fase 2**: 5-7 días (más compleja)
- **Fase 3**: 3-4 días
- **Fase 4**: 2-3 días
- **Fase 5**: 2-3 días
- **Fase 6**: 2-3 días
- **Fase 7**: 3-4 días
- **Fase 8**: 3-5 días

**Total estimado**: 22-32 días de desarrollo

## 🎯 Prioridades

1. **Alta**: Autenticación, Órdenes de Trabajo, Chat
2. **Media**: Plot AI, Archivos, Notificaciones
3. **Baja**: Estadísticas, Optimizaciones

## 📚 Próximos Pasos Inmediatos

1. Crear modelos de dominio completos
2. Crear migraciones de base de datos
3. Implementar repositorios básicos
4. Crear estructura de frontend React
5. Implementar Kanban Board básico

