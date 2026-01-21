# Day 3: Networking I - Services & Service Discovery

## 📚 Learning Objectives

By the end of this lab, you will understand:
- **Service Discovery:** How Frontend finds Backend, and Backend finds Redis using DNS.
- **Service Chaining:** Connecting multiple microservices (`Frontend` -> `Backend` -> `Redis`).
- **ClusterIP:** Why internal services used for backend are secure by default.

## 🎯 Architecture

```mermaid
graph LR
    User[User (Browser)] -- "Port Forward (8080)" --> Frontend[Frontend (Go)]
    Frontend -- "DNS: todo-api-service" --> Backend[Todo API (Go)]
    Backend -- "DNS: redis-service" --> Redis[(Redis)]
    Backend -- "DNS: audit-service" --> Audit[Audit Service (Go)]
```

## 📋 Lab Exercises

### 1. Build & Import Images

We need to build two images: one for the Frontend, one for the Backend.

```powershell
# 1. Build Backend
cd day3/todo-api
docker build -t k8s-day3-todo-api:v1 .
k3d image import k8s-day3-todo-api:v1 -c k8s-class

# 2. Build Frontend
cd ../app
docker build -t k8s-day3-app:v1 .
k3d image import k8s-day3-app:v1 -c k8s-class

# 3. Build Audit Service
cd ../audit-service
docker build -t k8s-day3-audit-service:v1 .
k3d image import k8s-day3-audit-service:v1 -c k8s-class
```

### 2. Deploy Services

```powershell
# Navigate to root
cd ../..

# Apply all manifests (Redis, Todo API, Frontend)
kubectl apply -f day3/manifests/
```

### 3. Verify Deployment

```powershell
kubectl get pods -w
# Wait until all 3 pods (redis, todo-api, go-web-app) are Running
```

### 4. Test the Application

Port-forward to the **Frontend** service only.

```powershell
kubectl port-forward svc/go-web-app-service 8080:80
```

1. Open your browser to [http://localhost:8080](http://localhost:8080).
2. You should see the Frontend UI.
3. Type a task (e.g., "Learn K8s") and click **Add**.
4. The page should reload and show the todo item.

### 5. What just happened?
1. Browser POSTs to Frontend (`/create`).
2. Frontend POSTs JSON to `http://todo-api-service/todos`.
3. Todo API writes data to `redis-service`.
4. Success!

## 📈 Scaling

You can scale the microservices (e.g., Todo API) horizontally to handle more load.

```powershell
# Scale to 2 replicas
kubectl scale deployment todo-api --replicas=2

# Verify horizontal scaling
kubectl get pods -w
```

## 🐞 Debugging Guide

If things don't work as expected, use these commands to inspect the cluster.

### 1. Check Service Status
```powershell
# See if pods are Running or crashing
kubectl get pods

# See if services are created and have ClusterIPs
kubectl get svc
```

### 2. Inspect Logs
View the standard output/error of your application containers.
```powershell
# Frontend Logs
kubectl logs -f -l app=go-web-app

# Backend Logs (Todo API)
kubectl logs -f -l app=todo-api

# Audit Service Logs
kubectl logs -f -l app=audit-service
```

### 3. Describe Resources
Detailed events and configuration. Useful for `ImagePullBackOff` or `CrashLoopBackOff`.
```powershell
# Describe a Pod
kubectl describe pod <pod-name>

# Describe a Deployment
kubectl describe deployment todo-api
```

### 4. Direct Access (Port Forwarding)
Bypass the frontend and query services directly from your local machine.

```powershell
# Forward Todo API to local 8081
kubectl port-forward svc/todo-api-service 8081:80

# Windows PowerShell: Test API directly
Invoke-WebRequest -Uri http://localhost:8081/todos -Method GET
```

### 5. Check Environment Variables
Verify if pods are getting correct service hostnames (`REDIS_HOST`, etc).

```powershell
kubectl exec deployment/todo-api -- env
```

## 🧹 Cleanup

```powershell
kubectl delete -f day3/manifests/
```
