# TaskBoard Deployment Resources

This directory contains deployment configurations for various platforms.

## Directory Structure

```
deployment/
├── kubernetes/          # Kubernetes manifests
│   ├── README.md       # Detailed Kubernetes instructions
│   ├── deploy.sh       # Quick deployment script
│   └── *.yaml         # K8s manifests
└── README.md           # This file
```

## Quick Links

- **Kubernetes Deployment**: See [kubernetes/README.md](kubernetes/README.md)
- **Main Deployment Guide**: See [../DEPLOYMENT.md](../DEPLOYMENT.md)

## Available Deployment Options

### 1. Docker Compose (Easiest)

For VPS, dedicated servers, or local development:

```bash
# Production deployment
docker-compose -f docker-compose.prod.yml up -d --build
```

See [../DEPLOYMENT.md](../DEPLOYMENT.md) for detailed instructions.

### 2. Kubernetes

For scalable cloud deployments:

```bash
cd kubernetes
./deploy.sh
```

See [kubernetes/README.md](kubernetes/README.md) for detailed instructions.

### 3. Cloud Platforms

Specific instructions for:
- AWS (ECS, Lightsail, EC2)
- Google Cloud (Cloud Run, GKE)
- Azure (Container Instances, AKS)
- DigitalOcean (App Platform, Droplets)
- Heroku

See [../DEPLOYMENT.md](../DEPLOYMENT.md) for platform-specific guides.

## Prerequisites

Before deploying, ensure you have:

1. **Environment Variables**: Configure all required variables
2. **Database**: PostgreSQL (managed or self-hosted)
3. **Cache**: Redis (managed or self-hosted)
4. **Domain**: Optional but recommended for production
5. **SSL Certificate**: Required for production (Let's Encrypt recommended)

## Security Checklist

- [ ] Strong database password
- [ ] Strong JWT secret (min 32 characters)
- [ ] CORS configured for your domain
- [ ] HTTPS/SSL enabled
- [ ] Firewall configured
- [ ] Database not publicly accessible
- [ ] Regular backups configured

## Support

For deployment help:
1. Check the main [DEPLOYMENT.md](../DEPLOYMENT.md)
2. Review platform-specific README files
3. Check application logs
4. Open an issue on GitHub

## Contributing

To add new deployment configurations:
1. Create a new directory for the platform
2. Include detailed README with instructions
3. Provide example configuration files
4. Test the deployment thoroughly
5. Submit a pull request

