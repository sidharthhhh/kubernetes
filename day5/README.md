# Day 5: Rust API Experiment

This lab builds a simple Rust API (using Actix-web and Tokio-Postgres), containerizes it, and runs it on Kubernetes with a Postgres database.

## Prerequisites
- Docker
- Kubectl
- Cargo (optional, for local development)

## 1. Build the Docker Image
First, build the Rust application image.

```powershell
cd day5
docker build -t rust-api:v1 .
```

> **For K3d Users:**
> k3d clusters cannot see your local Docker images by default. You must import the image:
> ```powershell
> k3d image import rust-api:v1 -c k8s-class
> ``` 
> *(Replace `k8s-class` with your cluster name if different)*

> **For K3s/VM Users:** 
> ```powershell
> docker save rust-api:v1 -o rust-api.tar
> sudo k3s ctr images import rust-api.tar
> ```

## 2. Deploy Infrastructure
Deploy Postgres and the Rust API to your cluster.

```powershell
# Deploy Postgres
kubectl apply -f k8s/postgres.yaml

# Wait for Postgres to be ready
kubectl rollout status deployment/postgres

# Deploy Rust API
kubectl apply -f k8s/app.yaml

# Wait for API to be ready
kubectl rollout status deployment/rust-api
```

## 3. Verify Application
Access the API. It is exposed via NodePort 30080.
If you are on localhost (Docker Desktop): `http://localhost:30080`
If you are on K3s/VM: `http://<node-ip>:30080`

```powershell
# Check health (should return "Rust API is running! 🦀")
curl http://localhost:30080/

# Check Database Connection (should return Postgres version)
curl http://localhost:30080/db
```

## 4. Infrastructure Debugging
If something goes wrong (e.g., `CrashLoopBackOff` or `Connection Refused`), use these commands:

### Check Pod Status
```powershell
kubectl get pods
```

### View Application Logs
If the Rust API fails to start:
```powershell
kubectl logs -l app=rust-api
```
*Look for "Connection refused" errors indicating Postgres isn't ready yet or DNS issues.*

### View Database Logs
```powershell
kubectl logs -l app=postgres
```

### Network Debugging
Spawn a debug pod to test connectivity from within the cluster:
```powershell
kubectl run -it --rm debug --image=curlimages/curl -- sh
# Inside the shell:
curl http://rust-api:8080/
curl http://rust-api:8080/db
```

### Verify Environment Variables
Ensure the Rust app has the correct DB credentials:
```powershell
kubectl describe pod -l app=rust-api
```
*Check the `Environment` section.*

## 5. Destroy (Cleanup)
To clean up all resources created in this lab:

```powershell
kubectl delete -f k8s/app.yaml
kubectl delete -f k8s/postgres.yaml
# Or simply
kubectl delete deployment rust-api postgres
kubectl delete service rust-api postgres
```

REMOVE the local docker image:
```powershell
docker rmi rust-api:v1
```

> **For K3d Users:**
> The image is cached inside ALL K3d cluster nodes. To free up space, you should remove it from every node.
> Run this PowerShell command to remove it from all agents and servers:
> ```powershell
> # Loop through all k3d nodes and remove the image
> docker ps --format "{{.Names}}" | Select-String "k3d-k8s-class-server", "k3d-k8s-class-agent" | ForEach-Object {
>     Write-Host "Cleaning $_..."
>     docker exec $_ crictl rmi "docker.io/library/rust-api:v1"
> }
> ```
> *Note: If you see "context deadline exceeded" or "no such image", it usually means the image is already gone or the node is busy. You can safely ignore it if `crictl images ls` confirms it's gone.*
