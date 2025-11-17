# 🚀 Plan de Migración: Plot AI (PHP → Go + React)

## 📋 Resumen

Migrar la aplicación **Plot AI** (PHP) a la misma tecnología que TaskBoard:
- **Backend**: Go (Gin framework)
- **Frontend**: React + TypeScript
- **Base de datos**: PostgreSQL (compartida con TaskBoard)
- **Autenticación**: JWT (compartida con TaskBoard)

## 🎯 Funcionalidades Actuales (PHP)

1. **Autenticación**: Sesiones PHP (`$_SESSION`)
2. **Chat con Gemini AI**: Integración con Google Gemini API
3. **Gestión de archivos**: 
   - `manual_entrenamiento.txt` - Manual del asistente
   - `contactos_plotcenter.txt` - Base de contactos
   - `soportes_publicitarios.txt` - Información de soportes
4. **Interfaz**: PHP con Tailwind CSS (CDN)

## 🏗️ Arquitectura Propuesta

```
┌─────────────────────────────────────┐
│         Frontend (React)            │
│  - TaskBoard                        │
│  - Plot AI                          │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Backend Go (Gin)               │
│  ┌──────────────┬─────────────────┐ │
│  │ TaskBoard API│  Plot AI API    │ │
│  └──────────────┴─────────────────┘ │
│  ┌─────────────────────────────────┐ │
│  │  Shared: JWT Auth, PostgreSQL   │ │
│  └─────────────────────────────────┘ │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  External Services                  │
│  - Google Gemini API                │
│  - PostgreSQL                       │
│  - Redis                            │
└─────────────────────────────────────┘
```

## 📁 Estructura de Carpetas

```
task-board/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go (actualizado)
│   ├── internal/
│   │   ├── handler/
│   │   │   ├── plot_ai_handler.go (nuevo)
│   │   │   └── ...
│   │   ├── service/
│   │   │   ├── plot_ai_service.go (nuevo)
│   │   │   └── gemini_service.go (nuevo)
│   │   └── domain/
│   │       └── plot_ai.go (nuevo)
│   └── ...
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── PlotAI/
│   │   │   │   ├── Chat.tsx (nuevo)
│   │   │   │   ├── Message.tsx (nuevo)
│   │   │   │   └── ChatInput.tsx (nuevo)
│   │   │   └── ...
│   │   └── services/
│   │       └── plotAI.ts (nuevo)
│   └── ...
└── data/
    ├── manual_entrenamiento.txt (mover)
    ├── contactos_plotcenter.txt (mover)
    └── soportes_publicitarios.txt (mover)
```

## 🔧 Pasos de Migración

### Fase 1: Backend Go

1. **Crear modelos de dominio**
   - `ChatMessage` - Mensajes del chat
   - `ChatHistory` - Historial de conversaciones
   - `PlotAIConfig` - Configuración del asistente

2. **Servicio de Gemini**
   - Integración con Google Gemini API
   - Manejo de prompts del sistema
   - Gestión de historial de conversación

3. **Handlers**
   - `POST /api/v1/plot-ai/chat` - Enviar mensaje
   - `GET /api/v1/plot-ai/history` - Obtener historial
   - `GET /api/v1/plot-ai/config` - Obtener configuración

4. **Base de datos**
   - Tabla `chat_messages` - Almacenar conversaciones
   - Tabla `plot_ai_config` - Configuración del sistema

### Fase 2: Frontend React

1. **Componentes**
   - `Chat.tsx` - Interfaz principal del chat
   - `Message.tsx` - Componente de mensaje
   - `ChatInput.tsx` - Input para enviar mensajes

2. **Servicios**
   - `plotAI.ts` - API client para Plot AI

3. **Integración**
   - Agregar ruta `/plot-ai` en React Router
   - Compartir autenticación con TaskBoard

### Fase 3: Integración

1. **Docker Compose**
   - Agregar Plot AI al mismo stack
   - Compartir base de datos

2. **Nginx**
   - Configurar rutas para Plot AI
   - SSL para todos los servicios

## 🔐 Autenticación Compartida

Ambas aplicaciones usarán:
- **JWT tokens** del mismo sistema
- **Misma base de datos de usuarios**
- **Mismo middleware de autenticación**

## 📊 Base de Datos

### Nuevas Tablas

```sql
-- Historial de conversaciones
CREATE TABLE chat_messages (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    role VARCHAR(20) NOT NULL, -- 'user' o 'assistant'
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Configuración de Plot AI
CREATE TABLE plot_ai_config (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## 🔑 Variables de Entorno

```env
# Gemini API
GEMINI_API_KEY=tu-api-key
GEMINI_MODEL=gemini-pro
GEMINI_TEMPERATURE=0.7
GEMINI_MAX_TOKENS=2000

# Plot AI Config
PLOT_AI_MANUAL_PATH=/data/manual_entrenamiento.txt
PLOT_AI_CONTACTOS_PATH=/data/contactos_plotcenter.txt
PLOT_AI_SOPORTES_PATH=/data/soportes_publicitarios.txt
```

## 🚀 Ventajas de la Migración

1. **Unificación**: Misma tecnología que TaskBoard
2. **Rendimiento**: Go es más rápido que PHP
3. **Escalabilidad**: Mejor manejo de concurrencia
4. **Mantenimiento**: Un solo stack tecnológico
5. **Autenticación**: Sistema unificado de usuarios
6. **Deployment**: Mismo proceso de despliegue

## 📝 Próximos Pasos

1. ✅ Crear estructura de backend Go
2. ⏳ Implementar servicio de Gemini
3. ⏳ Crear handlers de API
4. ⏳ Migrar frontend a React
5. ⏳ Integrar con TaskBoard
6. ⏳ Testing y deployment

