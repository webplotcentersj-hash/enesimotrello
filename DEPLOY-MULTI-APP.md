# 🚀 Guía de Despliegue: n8n + TaskBoard en el mismo VPS

Esta guía te ayudará a desplegar **n8n** y **TaskBoard** en el mismo VPS usando Nginx como reverse proxy.

> 💡 **Recomendación**: Si quieres mejor rendimiento y escalabilidad, considera desplegar el frontend en [Vercel](https://vercel.com). Ver [DEPLOY-VERCEL.md](./DEPLOY-VERCEL.md) para la guía completa.

## 📋 Requisitos Previos

- VPS con Ubuntu 24.04 (o similar)
- Acceso SSH root
- Docker y Docker Compose instalados
- Dominios configurados apuntando a tu VPS:
  - `n8n.tudominio.com` → n8n
  - `taskboard.tudominio.com` → TaskBoard Frontend
  - `api.taskboard.tudominio.com` → TaskBoard Backend

## 📊 Recursos del VPS

Tu VPS actual tiene:
- ✅ 50 GB de disco (6 GB usados - 12%)
- ✅ 23% de memoria usada
- ✅ 1% de CPU
- ✅ n8n ya instalado

**Recursos estimados necesarios:**
- n8n: ~500 MB RAM, 2 GB disco
- TaskBoard: ~1 GB RAM, 5 GB disco
- **Total disponible: ~44 GB disco, ~77% RAM libre** ✅

## 🔧 Paso 1: Preparar el VPS

### 1.1 Conectarse al VPS

```bash
ssh root@93.127.211.98
```

### 1.2 Verificar Docker

```bash
docker --version
docker compose version
```

Si no está instalado:
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
```

### 1.3 Crear directorio para TaskBoard

```bash
mkdir -p /opt/apps
cd /opt/apps
```

## 📥 Paso 2: Clonar TaskBoard

```bash
cd /opt/apps
git clone https://github.com/lunareclipsemontaigne667/task-board.git
cd task-board
```

## ⚙️ Paso 3: Configurar Variables de Entorno

```bash
cp env.production.example .env
nano .env
```

Configura las siguientes variables:

```env
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=TU_CONTRASEÑA_FUERTE_AQUI
DB_NAME=taskboard

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# JWT
JWT_SECRET=TU_SECRET_JWT_MUY_LARGO_MINIMO_32_CARACTERES
JWT_EXPIRY=24h

# Server
HOST=0.0.0.0
PORT=8080

# CORS - Actualizar después de configurar Nginx
CORS_ORIGIN=https://taskboard.tudominio.com

# Frontend
REACT_APP_API_URL=https://api.taskboard.tudominio.com/api/v1
```

**Generar contraseñas seguras:**
```bash
# Generar DB_PASSWORD
openssl rand -base64 32

# Generar JWT_SECRET
openssl rand -base64 48
```

## 🐳 Paso 4: Desplegar TaskBoard

```bash
cd /opt/apps/task-board
docker compose -f docker-compose.multi-app.yml up -d --build
```

Verificar que todo esté corriendo:
```bash
docker compose -f docker-compose.multi-app.yml ps
docker compose -f docker-compose.multi-app.yml logs -f
```

## 🔒 Paso 5: Configurar Nginx

### 5.1 Hacer el script ejecutable

```bash
chmod +x scripts/setup-nginx-multi-app.sh
```

### 5.2 Ejecutar el script

```bash
sudo ./scripts/setup-nginx-multi-app.sh
```

El script te pedirá:
- Dominio de n8n (ej: `n8n.tudominio.com`)
- Dominio del frontend de TaskBoard (ej: `taskboard.tudominio.com`)
- Dominio del API de TaskBoard (ej: `api.taskboard.tudominio.com`)
- Email para certificados SSL
- Puerto de n8n (por defecto: 5678)

### 5.3 Verificar configuración de Nginx

```bash
sudo nginx -t
sudo systemctl status nginx
```

## 🔄 Paso 6: Actualizar Variables de Entorno

Después de configurar los dominios, actualiza el `.env`:

```bash
nano /opt/apps/task-board/.env
```

Actualiza:
```env
CORS_ORIGIN=https://taskboard.tudominio.com
REACT_APP_API_URL=https://api.taskboard.tudominio.com/api/v1
```

Reinicia TaskBoard:
```bash
cd /opt/apps/task-board
docker compose -f docker-compose.multi-app.yml restart
```

## ✅ Paso 7: Verificar Todo

### Verificar servicios Docker

```bash
docker ps
```

Deberías ver:
- `taskboard-postgres`
- `taskboard-redis`
- `taskboard-backend`
- `taskboard-frontend`
- `n8n` (si está en Docker)

### Verificar Nginx

```bash
sudo systemctl status nginx
sudo nginx -t
```

### Probar las aplicaciones

```bash
# Probar n8n
curl -I https://n8n.tudominio.com

# Probar TaskBoard Frontend
curl -I https://taskboard.tudominio.com

# Probar TaskBoard Backend
curl -I https://api.taskboard.tudominio.com/api/v1/health
```

## 📊 Estructura de Puertos

| Aplicación | Puerto Interno | Puerto Externo | Acceso |
|------------|---------------|----------------|--------|
| n8n | 5678 | 127.0.0.1:5678 | Solo localhost |
| TaskBoard Backend | 8080 | 127.0.0.1:8081 | Solo localhost |
| TaskBoard Frontend | 80 | 127.0.0.1:3001 | Solo localhost |
| Nginx | 80, 443 | 0.0.0.0:80, 443 | Público |

**Todas las aplicaciones solo son accesibles a través de Nginx con SSL.**

## 🛠️ Comandos Útiles

### TaskBoard

```bash
# Ver logs
cd /opt/apps/task-board
docker compose -f docker-compose.multi-app.yml logs -f

# Reiniciar
docker compose -f docker-compose.multi-app.yml restart

# Detener
docker compose -f docker-compose.multi-app.yml down

# Iniciar
docker compose -f docker-compose.multi-app.yml up -d

# Ver estado
docker compose -f docker-compose.multi-app.yml ps
```

### Nginx

```bash
# Recargar configuración
sudo nginx -t && sudo systemctl reload nginx

# Ver logs
sudo tail -f /var/log/nginx/error.log
sudo tail -f /var/log/nginx/access.log

# Ver configuración
sudo cat /etc/nginx/sites-available/n8n
sudo cat /etc/nginx/sites-available/taskboard
```

### SSL Certificates

```bash
# Renovar certificados manualmente
sudo certbot renew

# Ver certificados
sudo certbot certificates
```

## 🔐 Seguridad

### Firewall (UFW)

```bash
# Permitir SSH
sudo ufw allow 22/tcp

# Permitir HTTP/HTTPS
sudo ufw allow 'Nginx Full'

# Activar firewall
sudo ufw enable
sudo ufw status
```

### Backups

```bash
# Backup de base de datos TaskBoard
docker exec taskboard-postgres pg_dump -U postgres taskboard > backup_$(date +%Y%m%d).sql

# Backup de volúmenes Docker
docker run --rm -v task-board_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup_$(date +%Y%m%d).tar.gz /data
```

## 🐛 Solución de Problemas

### TaskBoard no inicia

```bash
# Ver logs
docker compose -f docker-compose.multi-app.yml logs backend
docker compose -f docker-compose.multi-app.yml logs frontend

# Verificar variables de entorno
docker compose -f docker-compose.multi-app.yml config
```

### Nginx no funciona

```bash
# Verificar configuración
sudo nginx -t

# Ver logs de error
sudo tail -f /var/log/nginx/error.log

# Verificar que los servicios estén corriendo
curl http://127.0.0.1:8081/api/v1/health
curl http://127.0.0.1:3001
```

### Certificados SSL no se generan

1. Verifica que los dominios apunten a tu IP:
   ```bash
   curl -s ifconfig.me  # Tu IP pública
   dig n8n.tudominio.com
   dig taskboard.tudominio.com
   ```

2. Verifica que los puertos 80 y 443 estén abiertos:
   ```bash
   sudo ufw status
   ```

3. Intenta generar certificados manualmente:
   ```bash
   sudo certbot --nginx -d n8n.tudominio.com -d taskboard.tudominio.com -d api.taskboard.tudominio.com
   ```

## 📈 Monitoreo

### Ver uso de recursos

```bash
# Uso de CPU y memoria
htop

# Uso de disco
df -h

# Uso de Docker
docker stats
```

### Logs centralizados

Considera usar:
- **Dozzle**: `docker run -d --name dozzle -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock amir20/dozzle`

## 🎉 ¡Listo!

Ahora tienes:
- ✅ n8n corriendo en `https://n8n.tudominio.com`
- ✅ TaskBoard Frontend en `https://taskboard.tudominio.com`
- ✅ TaskBoard Backend en `https://api.taskboard.tudominio.com`
- ✅ SSL automático con Let's Encrypt
- ✅ Todo en el mismo VPS

## 📚 Recursos Adicionales

- [Documentación de Nginx](https://nginx.org/en/docs/)
- [Documentación de Certbot](https://certbot.eff.org/)
- [Documentación de Docker Compose](https://docs.docker.com/compose/)

