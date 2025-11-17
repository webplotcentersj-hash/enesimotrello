#!/bin/bash

# Quick deployment script for Kubernetes

set -e

echo "🚀 Deploying TaskBoard to Kubernetes"
echo "====================================="
echo ""

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found. Please install kubectl first."
    exit 1
fi

# Check if secrets exist
echo "🔍 Checking for secrets..."
if ! kubectl get secret taskboard-secrets -n taskboard &> /dev/null; then
    echo "⚠️  Secrets not found. Please create secrets first:"
    echo ""
    echo "kubectl create secret generic taskboard-secrets -n taskboard \\"
    echo "  --from-literal=db-password=YOUR_STRONG_DB_PASSWORD \\"
    echo "  --from-literal=jwt-secret=YOUR_STRONG_JWT_SECRET"
    echo ""
    read -p "Press enter after creating secrets, or Ctrl+C to cancel..."
fi

# Apply manifests in order
echo ""
echo "📦 Creating namespace..."
kubectl apply -f namespace.yaml

echo "📦 Creating configmap..."
kubectl apply -f configmap.yaml

echo "📦 Creating persistent volume claims..."
kubectl apply -f postgres-pvc.yaml
kubectl apply -f redis-pvc.yaml

echo "📦 Deploying PostgreSQL..."
kubectl apply -f postgres-deployment.yaml

echo "📦 Deploying Redis..."
kubectl apply -f redis-deployment.yaml

echo "⏳ Waiting for databases to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres -n taskboard --timeout=120s
kubectl wait --for=condition=ready pod -l app=redis -n taskboard --timeout=120s

echo "📦 Deploying backend..."
kubectl apply -f backend-deployment.yaml

echo "📦 Deploying frontend..."
kubectl apply -f frontend-deployment.yaml

echo "⏳ Waiting for applications to be ready..."
kubectl wait --for=condition=ready pod -l app=taskboard-backend -n taskboard --timeout=120s
kubectl wait --for=condition=ready pod -l app=taskboard-frontend -n taskboard --timeout=120s

echo "📦 Deploying ingress..."
kubectl apply -f ingress.yaml

echo ""
echo "====================================="
echo "✅ Deployment Complete!"
echo "====================================="
echo ""
echo "📊 Check status:"
echo "  kubectl get pods -n taskboard"
echo "  kubectl get services -n taskboard"
echo "  kubectl get ingress -n taskboard"
echo ""
echo "📝 View logs:"
echo "  kubectl logs -f deployment/taskboard-backend -n taskboard"
echo "  kubectl logs -f deployment/taskboard-frontend -n taskboard"
echo ""

