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

## 🧹 Cleanup

```powershell
kubectl delete -f day3/manifests/
```
