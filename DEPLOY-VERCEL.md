# 🚀 Despliegue: Frontend en Vercel + Backend en VPS

Esta guía te ayudará a desplegar el **frontend de TaskBoard en Vercel** y el **backend en tu VPS** junto con n8n.

## 🎯 Arquitectura

```
┌─────────────────┐
│   Vercel CDN    │  ← Frontend (React) - Global CDN
│  (Frontend)     │
└────────┬────────┘
         │ HTTPS
         │
┌────────▼─────────────────────────┐
│         Tu VPS                   │
│  ┌──────────────────────────┐   │
│  │   Nginx (Reverse Proxy)  │   │
│  └──────┬───────────┬───────┘   │
│         │           │            │
│  ┌──────▼──┐  ┌────▼─────────┐  │
│  │   n8n   │  │ TaskBoard    │  │
│  │ :5678   │  │ Backend :8081│  │
│  └─────────┘  └────┬─────────┘  │
│                    │             │
│           ┌────────┴────────┐   │
│           │  PostgreSQL     │   │
│           │  Redis          │   │
│           └─────────────────┘   │
└─────────────────────────────────┘
```

## ✅ Ventajas de esta Configuración

- ⚡ **CDN Global**: Frontend servido desde múltiples ubicaciones
- 🔒 **SSL Automático**: Vercel maneja certificados SSL
- 💰 **Gratis**: Plan gratuito generoso de Vercel
- 🚀 **Despliegue Automático**: Desde Git push
- 📉 **Menos Carga en VPS**: Solo backend + DB
- 🔄 **Escalabilidad**: Vercel escala automáticamente

## 📋 Requisitos Previos

- Cuenta en [Vercel](https://vercel.com) (gratis)
- Repositorio Git (GitHub, GitLab, Bitbucket)
- VPS con Docker instalado
- Dominios configurados:
  - `api.taskboard.tudominio.com` → Backend API
  - `n8n.tudominio.com` → n8n
  - (Opcional) Dominio personalizado para Vercel

## 🔧 Parte 1: Desplegar Backend en VPS

### 1.1 Preparar el VPS

```bash
ssh root@93.127.211.98
cd /opt/apps
git clone https://github.com/lunareclipsemontaigne667/task-board.git
cd task-board
```

### 1.2 Configurar Variables de Entorno

```bash
cp env.production.example .env
nano .env
```

Configura:

```env
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=TU_CONTRASEÑA_FUERTE
DB_NAME=taskboard

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# JWT - Genera uno seguro
JWT_SECRET=TU_SECRET_JWT_MUY_LARGO
JWT_EXPIRY=24h

# CORS - Actualizar después de obtener dominio de Vercel
CORS_ORIGIN=https://taskboard.vercel.app

# Backend
HOST=0.0.0.0
PORT=8080
```

**Generar contraseñas:**
```bash
openssl rand -base64 32  # Para DB_PASSWORD
openssl rand -base64 48  # Para JWT_SECRET
```

### 1.3 Desplegar Backend

```bash
docker compose -f docker-compose.backend-only.yml up -d --build
```

Verificar:
```bash
docker compose -f docker-compose.backend-only.yml ps
docker compose -f docker-compose.backend-only.yml logs -f backend
```

### 1.4 Configurar Nginx

```bash
chmod +x scripts/setup-nginx-vercel.sh
sudo ./scripts/setup-nginx-vercel.sh
```

El script te pedirá:
- Dominio de n8n
- Dominio del API (ej: `api.taskboard.tudominio.com`)
- Dominio de Vercel (ej: `taskboard.vercel.app` o tu dominio personalizado)
- Email para SSL

### 1.5 Actualizar CORS y JWT

Después de obtener el dominio de Vercel, actualiza `.env`:

```bash
nano .env
```

Configura:
```env
# CORS - Dominio de Vercel (puedes usar múltiples separados por comas)
CORS_ORIGIN=https://taskboard.vercel.app

# JWT Secret - Genera uno seguro
JWT_SECRET=tu-secret-jwt-muy-largo-minimo-32-caracteres
JWT_EXPIRY=24h
```

**Generar JWT Secret:**
```bash
openssl rand -base64 48
```

Reinicia el backend:
```bash
docker compose -f docker-compose.backend-only.yml restart backend
```

## 🎨 Parte 2: Desplegar Frontend en Vercel

### 2.1 Preparar el Repositorio

Asegúrate de que el repositorio tenga:
- ✅ `frontend/vercel.json` (ya creado)
- ✅ `frontend/package.json`
- ✅ Código del frontend en `frontend/`

### 2.2 Conectar con Vercel

1. Ve a [vercel.com](https://vercel.com) e inicia sesión
2. Click en **"Add New Project"**
3. Importa tu repositorio de GitHub/GitLab/Bitbucket
4. Selecciona el repositorio `task-board`

### 2.3 Configurar el Proyecto

**Configuración del Proyecto:**

- **Framework Preset**: Create React App
- **Root Directory**: `frontend`
- **Build Command**: `npm run build` (automático)
- **Output Directory**: `build` (automático)
- **Install Command**: `npm install` (automático)

### 2.4 Variables de Entorno

En la configuración del proyecto, agrega estas **Environment Variables**:

```
REACT_APP_API_URL=https://api.taskboard.tudominio.com/api/v1
```

**Importante**: 
- Reemplaza `api.taskboard.tudominio.com` con tu dominio real del backend
- Asegúrate de usar `https://` (no `http://`)

### 2.5 Desplegar

1. Click en **"Deploy"**
2. Vercel construirá y desplegará automáticamente
3. Obtendrás una URL como: `https://task-board-xxxxx.vercel.app`

### 2.6 Dominio Personalizado (Opcional)

Si quieres usar tu propio dominio:

1. Ve a **Settings** → **Domains**
2. Agrega tu dominio: `taskboard.tudominio.com`
3. Configura los registros DNS según las instrucciones de Vercel
4. Espera a que se propague (puede tardar unos minutos)

**Registros DNS necesarios:**
```
Tipo: CNAME
Nombre: taskboard (o @)
Valor: cname.vercel-dns.com
```

## 🔄 Parte 3: Actualizar Configuraciones

### 3.1 Actualizar CORS en Backend

Después de obtener el dominio de Vercel, actualiza el `.env` en el VPS:

```bash
ssh root@93.127.211.98
cd /opt/apps/task-board
nano .env
```

Actualiza:
```env
CORS_ORIGIN=https://taskboard.vercel.app
# O si usas dominio personalizado:
CORS_ORIGIN=https://taskboard.tudominio.com
```

Reinicia:
```bash
docker compose -f docker-compose.backend-only.yml restart backend
```

### 3.2 Verificar Conexión

Abre el frontend en Vercel y verifica que:
- ✅ La aplicación carga correctamente
- ✅ Puedes hacer login/registro
- ✅ Las peticiones al API funcionan (abre DevTools → Network)

## 🧪 Verificar Todo

### Backend API

```bash
# Health check
curl https://api.taskboard.tudominio.com/api/v1/health

# Debería responder: {"status":"ok"}
```

### Frontend

1. Abre `https://taskboard.vercel.app` (o tu dominio)
2. Abre DevTools (F12) → Console
3. No debería haber errores de CORS
4. Prueba crear un tablero o tarea

### n8n

```bash
curl -I https://n8n.tudominio.com
```

## 🔐 Autenticación

TaskBoard soporta dos tipos de autenticación:

1. **JWT (Usuarios Registrados)**: Los usuarios pueden registrarse e iniciar sesión
2. **Anónima (Modo Invitado)**: Los usuarios pueden usar la app sin registrarse

**Endpoints de autenticación:**
- `POST /api/v1/auth/register` - Registrar usuario
- `POST /api/v1/auth/login` - Iniciar sesión
- `GET /api/v1/auth/profile` - Obtener perfil (requiere autenticación)
- `PUT /api/v1/auth/profile` - Actualizar perfil (requiere autenticación)

Ver [AUTHENTICATION.md](./AUTHENTICATION.md) para más detalles.

## 🔐 Configuración de Seguridad

### Backend CORS

El backend debe aceptar peticiones desde Vercel. Verifica en `backend/internal/middleware/cors.go` que esté configurado correctamente.

### Variables de Entorno Sensibles

**Nunca** subas el archivo `.env` al repositorio. Está en `.gitignore`, pero verifica:

```bash
cat .gitignore | grep .env
```

## 📊 Comandos Útiles

### Vercel CLI (Opcional)

Instalar Vercel CLI:
```bash
npm i -g vercel
```

Comandos útiles:
```bash
# Login
vercel login

# Desplegar
cd frontend
vercel

# Ver logs
vercel logs

# Listar proyectos
vercel ls
```

### Backend en VPS

```bash
# Ver logs
docker compose -f docker-compose.backend-only.yml logs -f backend

# Reiniciar
docker compose -f docker-compose.backend-only.yml restart

# Detener
docker compose -f docker-compose.backend-only.yml down

# Iniciar
docker compose -f docker-compose.backend-only.yml up -d
```

## 🐛 Solución de Problemas

### Error de CORS

**Síntoma**: Error en consola del navegador sobre CORS

**Solución**:
1. Verifica que `CORS_ORIGIN` en `.env` del backend coincida exactamente con el dominio de Vercel
2. Incluye el protocolo: `https://taskboard.vercel.app` (no solo `taskboard.vercel.app`)
3. Reinicia el backend después de cambiar `.env`

### Frontend no se conecta al API

**Síntoma**: Errores 404 o "Network Error"

**Solución**:
1. Verifica que `REACT_APP_API_URL` en Vercel sea correcta
2. Verifica que el backend esté corriendo: `docker ps`
3. Prueba el API directamente: `curl https://api.taskboard.tudominio.com/api/v1/health`

### Build falla en Vercel

**Síntoma**: Error durante el build

**Solución**:
1. Verifica que `Root Directory` esté configurado como `frontend`
2. Revisa los logs de build en Vercel
3. Prueba localmente: `cd frontend && npm run build`

### SSL no funciona

**Síntoma**: Certificado SSL inválido

**Solución**:
1. Verifica que el dominio apunte correctamente
2. Espera unos minutos para que se propague
3. Revisa la configuración DNS

### Error de Autenticación

**Síntoma**: Error 401 o "Invalid token"

**Solución**:
1. Verifica que `JWT_SECRET` esté configurado correctamente
2. Verifica que el token no haya expirado
3. El usuario debe iniciar sesión nuevamente

## 📈 Monitoreo

### Vercel Analytics

Vercel ofrece analytics gratuitos:
1. Ve a tu proyecto en Vercel
2. Click en **Analytics**
3. Habilita Vercel Analytics

### Backend Logs

```bash
# Ver logs en tiempo real
docker compose -f docker-compose.backend-only.yml logs -f backend

# Ver últimas 100 líneas
docker compose -f docker-compose.backend-only.yml logs --tail=100 backend
```

## 🎉 ¡Listo!

Ahora tienes:
- ✅ Frontend en Vercel con CDN global
- ✅ Backend en tu VPS
- ✅ n8n en tu VPS
- ✅ SSL automático en todos los servicios
- ✅ Despliegue automático desde Git
- ✅ Autenticación JWT y anónima funcionando

## 📚 Recursos

- Guía completa: `DEPLOY-VERCEL.md`
- Autenticación: `AUTHENTICATION.md`
- Configuración de Vercel: `frontend/vercel.json`
- Docker Compose backend: `docker-compose.backend-only.yml`
- [Documentación de Vercel](https://vercel.com/docs)
- [Vercel CLI](https://vercel.com/docs/cli)
- [Configuración de CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
