# Kubernetes Deployment for TaskBoard

This directory contains Kubernetes manifests for deploying TaskBoard to a Kubernetes cluster.

## Prerequisites

- Kubernetes cluster (v1.20+)
- `kubectl` configured to access your cluster
- Ingress controller (nginx, traefik, etc.)
- cert-manager (for SSL certificates)

## Quick Start

### 1. Create Secrets

```bash
# Create secrets for sensitive data
kubectl create secret generic taskboard-secrets -n taskboard \
  --from-literal=db-password=YOUR_STRONG_DB_PASSWORD \
  --from-literal=jwt-secret=YOUR_STRONG_JWT_SECRET
```

### 2. Update ConfigMap

Edit `configmap.yaml` and update:
- `CORS_ORIGIN`: Your frontend domain
- `REACT_APP_API_URL`: Your backend API URL

### 3. Build and Push Docker Images

```bash
# Build backend image
cd backend
docker build -t your-registry/taskboard-backend:latest .
docker push your-registry/taskboard-backend:latest

# Build frontend image
cd ../frontend
docker build -f Dockerfile.prod -t your-registry/taskboard-frontend:latest .
docker push your-registry/taskboard-frontend:latest
```

### 4. Update Deployment Images

Edit the following files and update the image references:
- `backend-deployment.yaml`: Update `image:` field
- `frontend-deployment.yaml`: Update `image:` field

### 5. Deploy to Kubernetes

```bash
# Create namespace
kubectl apply -f namespace.yaml

# Create secrets (if not done in step 1)
kubectl apply -f secrets.yaml

# Apply configurations
kubectl apply -f configmap.yaml

# Create persistent volumes
kubectl apply -f postgres-pvc.yaml
kubectl apply -f redis-pvc.yaml

# Deploy databases
kubectl apply -f postgres-deployment.yaml
kubectl apply -f redis-deployment.yaml

# Wait for databases to be ready
kubectl wait --for=condition=ready pod -l app=postgres -n taskboard --timeout=120s
kubectl wait --for=condition=ready pod -l app=redis -n taskboard --timeout=120s

# Deploy applications
kubectl apply -f backend-deployment.yaml
kubectl apply -f frontend-deployment.yaml

# Deploy ingress (update domains first)
kubectl apply -f ingress.yaml
```

### 6. Verify Deployment

```bash
# Check pod status
kubectl get pods -n taskboard

# Check services
kubectl get services -n taskboard

# Check ingress
kubectl get ingress -n taskboard

# View logs
kubectl logs -f deployment/taskboard-backend -n taskboard
kubectl logs -f deployment/taskboard-frontend -n taskboard
```

## Configuration

### Scaling

Scale backend replicas:
```bash
kubectl scale deployment/taskboard-backend --replicas=3 -n taskboard
```

Scale frontend replicas:
```bash
kubectl scale deployment/taskboard-frontend --replicas=3 -n taskboard
```

### Rolling Updates

Update backend image:
```bash
kubectl set image deployment/taskboard-backend backend=your-registry/taskboard-backend:v2 -n taskboard
```

Update frontend image:
```bash
kubectl set image deployment/taskboard-frontend frontend=your-registry/taskboard-frontend:v2 -n taskboard
```

### Resource Limits

Adjust resource limits in:
- `backend-deployment.yaml`
- `frontend-deployment.yaml`
- `postgres-deployment.yaml`
- `redis-deployment.yaml`

## Ingress Setup

### Install Nginx Ingress Controller

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml
```

### Install cert-manager (for SSL)

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
```

### Create ClusterIssuer for Let's Encrypt

```bash
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOF
```

## Monitoring

### View Logs

```bash
# Backend logs
kubectl logs -f deployment/taskboard-backend -n taskboard

# Frontend logs
kubectl logs -f deployment/taskboard-frontend -n taskboard

# Database logs
kubectl logs -f deployment/postgres -n taskboard
```

### Port Forwarding (for testing)

```bash
# Forward backend
kubectl port-forward -n taskboard svc/taskboard-backend 8080:8080

# Forward frontend
kubectl port-forward -n taskboard svc/taskboard-frontend 8080:80
```

### Execute Commands in Pods

```bash
# Access PostgreSQL
kubectl exec -it deployment/postgres -n taskboard -- psql -U postgres taskboard

# Access Redis
kubectl exec -it deployment/redis -n taskboard -- redis-cli

# Access backend shell
kubectl exec -it deployment/taskboard-backend -n taskboard -- sh
```

## Backup and Restore

### Database Backup

```bash
# Backup PostgreSQL
kubectl exec -n taskboard deployment/postgres -- pg_dump -U postgres taskboard > backup.sql

# Restore
kubectl exec -i -n taskboard deployment/postgres -- psql -U postgres taskboard < backup.sql
```

## Troubleshooting

### Pods not starting

```bash
# Check pod events
kubectl describe pod <pod-name> -n taskboard

# Check logs
kubectl logs <pod-name> -n taskboard
```

### Database connection issues

```bash
# Test database connectivity from backend pod
kubectl exec -it deployment/taskboard-backend -n taskboard -- sh
# Inside the pod:
nc -zv postgres 5432
```

### Ingress not working

```bash
# Check ingress status
kubectl describe ingress taskboard-ingress -n taskboard

# Check ingress controller logs
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller
```

## Cleanup

To remove all TaskBoard resources:

```bash
kubectl delete namespace taskboard
```

Or delete individual resources:

```bash
kubectl delete -f .
```

## Production Recommendations

1. **Use external managed databases**: AWS RDS, GCP Cloud SQL, Azure Database
2. **Use external Redis**: AWS ElastiCache, GCP Memorystore
3. **Set up monitoring**: Prometheus + Grafana
4. **Configure autoscaling**: HorizontalPodAutoscaler
5. **Set up backup automation**: CronJobs for database backups
6. **Use secrets management**: HashiCorp Vault, AWS Secrets Manager
7. **Implement network policies**: Restrict pod-to-pod communication
8. **Set resource limits**: Prevent resource exhaustion
9. **Use affinity rules**: Distribute pods across nodes
10. **Enable logging**: ELK stack, Loki, or cloud-native logging

## Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [Ingress Nginx Documentation](https://kubernetes.github.io/ingress-nginx/)
- [cert-manager Documentation](https://cert-manager.io/docs/)

