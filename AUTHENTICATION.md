# 🔐 Guía de Autenticación - TaskBoard

TaskBoard soporta dos tipos de autenticación:
1. **Autenticación JWT** - Para usuarios registrados
2. **Autenticación Anónima** - Para usuarios no registrados (modo invitado)

## 📋 Tipos de Autenticación

### 1. Autenticación JWT (Usuarios Registrados)

Los usuarios pueden registrarse e iniciar sesión para obtener un token JWT que les permite acceder a todas las funcionalidades.

**Endpoints:**
- `POST /api/v1/auth/register` - Registrar nuevo usuario
- `POST /api/v1/auth/login` - Iniciar sesión
- `GET /api/v1/auth/profile` - Obtener perfil del usuario (requiere autenticación)
- `PUT /api/v1/auth/profile` - Actualizar perfil (requiere autenticación)

**Flujo:**
1. Usuario se registra o inicia sesión
2. Backend genera un token JWT
3. Frontend almacena el token en `localStorage`
4. Todas las peticiones incluyen el token en el header `Authorization: Bearer <token>`

### 2. Autenticación Anónima (Modo Invitado)

Los usuarios no registrados pueden usar la aplicación de forma anónima. El sistema genera un UUID único que se almacena en `localStorage` y se envía en cada petición.

**Header requerido:**
- `X-Anonymous-User-Id: <uuid>`

**Flujo:**
1. Frontend genera un UUID único al cargar
2. UUID se almacena en `localStorage` como `anonymous_user_id`
3. Todas las peticiones incluyen el UUID en el header `X-Anonymous-User-Id`
4. Backend crea o recupera un usuario anónimo basado en el UUID

## 🔧 Configuración

### Backend

#### Variables de Entorno

```env
# JWT Configuration
JWT_SECRET=tu-secret-jwt-muy-largo-minimo-32-caracteres
JWT_EXPIRY=24h

# CORS Configuration
CORS_ORIGIN=https://taskboard.vercel.app,https://taskboard.tudominio.com
```

**Importante:**
- `JWT_SECRET` debe ser una cadena segura de al menos 32 caracteres
- `CORS_ORIGIN` debe incluir todos los dominios del frontend (separados por comas)
- Si usas `*` en `CORS_ORIGIN`, no podrás usar cookies/credenciales

#### Generar JWT Secret

```bash
# Generar un secret seguro
openssl rand -base64 48
```

### Frontend

#### Variables de Entorno

```env
REACT_APP_API_URL=https://api.taskboard.tudominio.com/api/v1
```

#### Almacenamiento

**JWT Token:**
- Clave: `taskboard_auth_token`
- Ubicación: `localStorage`
- Se incluye automáticamente en todas las peticiones si existe

**Usuario:**
- Clave: `taskboard_user`
- Ubicación: `localStorage`
- Contiene información del usuario autenticado

**Usuario Anónimo:**
- Clave: `anonymous_user_id`
- Ubicación: `localStorage`
- Se genera automáticamente si no existe

## 📝 Uso en el Frontend

### Servicio de Autenticación

```typescript
import authService from './services/auth';

// Registrar usuario
const response = await authService.register({
  email: 'user@example.com',
  username: 'johndoe',
  password: 'password123',
  first_name: 'John',
  last_name: 'Doe'
});

// Iniciar sesión
const response = await authService.login({
  email: 'user@example.com',
  password: 'password123'
});

// Verificar si está autenticado
if (authService.isAuthenticated()) {
  // Usuario autenticado
}

// Obtener usuario actual
const user = authService.getUser();

// Cerrar sesión
authService.logout();

// Obtener perfil
const profile = await authService.getProfile();

// Actualizar perfil
const updated = await authService.updateProfile({
  first_name: 'Jane',
  last_name: 'Smith'
});
```

### Interceptor de Axios

El interceptor de axios maneja automáticamente:
- Agregar el token JWT si el usuario está autenticado
- Agregar el UUID anónimo si no está autenticado
- Manejar errores 401 (token inválido/expirado) y cerrar sesión automáticamente

```typescript
// El interceptor ya está configurado en api.ts
// No necesitas hacer nada adicional
```

## 🔒 Seguridad

### JWT Token

- **Expiración**: Configurable (por defecto 24 horas)
- **Almacenamiento**: `localStorage` (considera usar `httpOnly` cookies en producción)
- **Validación**: El backend valida el token en cada petición protegida

### CORS

- El backend valida el origen de las peticiones
- Solo los dominios en `CORS_ORIGIN` pueden hacer peticiones
- Las credenciales están habilitadas para dominios específicos

### Recomendaciones de Seguridad

1. **JWT Secret**: Usa un secret fuerte y único en producción
2. **HTTPS**: Siempre usa HTTPS en producción
3. **Token Expiry**: Considera reducir el tiempo de expiración del token
4. **Refresh Tokens**: Considera implementar refresh tokens para mejor seguridad
5. **Rate Limiting**: Implementa rate limiting en los endpoints de autenticación

## 🐛 Solución de Problemas

### Error: "Authorization header required"

**Causa**: El endpoint requiere autenticación pero no se envió el token.

**Solución**:
- Verifica que el usuario haya iniciado sesión
- Verifica que el token esté en `localStorage`
- Verifica que el interceptor esté agregando el header correctamente

### Error: "Invalid token"

**Causa**: El token JWT es inválido o ha expirado.

**Solución**:
- El usuario debe iniciar sesión nuevamente
- Verifica que `JWT_SECRET` sea el mismo en backend y frontend
- Verifica que el token no haya expirado

### Error de CORS

**Causa**: El dominio del frontend no está en `CORS_ORIGIN`.

**Solución**:
- Agrega el dominio del frontend a `CORS_ORIGIN` en el backend
- Reinicia el backend después de cambiar la variable
- Verifica que uses `https://` (no `http://`) en producción

### Token no se envía en las peticiones

**Causa**: El token no está en `localStorage` o el interceptor no está funcionando.

**Solución**:
- Verifica `localStorage.getItem('taskboard_auth_token')`
- Verifica que el interceptor esté configurado correctamente
- Verifica la consola del navegador para errores

## 📚 Ejemplos de Peticiones

### Registrar Usuario

```bash
curl -X POST https://api.taskboard.tudominio.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Iniciar Sesión

```bash
curl -X POST https://api.taskboard.tudominio.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Petición Autenticada

```bash
curl -X GET https://api.taskboard.tudominio.com/api/v1/auth/profile \
  -H "Authorization: Bearer <tu-token-jwt>"
```

### Petición Anónima

```bash
curl -X GET https://api.taskboard.tudominio.com/api/v1/boards \
  -H "X-Anonymous-User-Id: <uuid-generado>"
```

## 🔄 Migración de Anónimo a Autenticado

Si un usuario anónimo decide registrarse:

1. El usuario crea una cuenta o inicia sesión
2. El token JWT reemplaza el UUID anónimo
3. Los datos anónimos pueden migrarse al usuario autenticado (requiere implementación adicional)

## 📊 Flujo Completo

```
┌─────────────┐
│   Usuario   │
└──────┬──────┘
       │
       ├─→ Registro/Login
       │   ↓
       │   Token JWT almacenado
       │   ↓
       │   Peticiones con Authorization: Bearer <token>
       │
       └─→ Modo Anónimo
           ↓
           UUID generado y almacenado
           ↓
           Peticiones con X-Anonymous-User-Id: <uuid>
```

## 🎯 Mejores Prácticas

1. **Siempre usa HTTPS** en producción
2. **Valida el token** en el frontend antes de hacer peticiones
3. **Maneja errores 401** para cerrar sesión automáticamente
4. **No almacenes información sensible** en el token JWT
5. **Implementa refresh tokens** para mejor UX y seguridad
6. **Usa httpOnly cookies** en lugar de localStorage para mayor seguridad (requiere cambios adicionales)

