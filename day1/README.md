# Day 1: The Pod & The Container - Lab Guide

This guide covers the commands used to set up the cluster, build the application, and deploy the first Pod.

## 1. Cluster Setup (If not already running)
Create a k3d cluster with 1 server and 2 agents, exposing ports 80 and 443.
```powershell
k3d cluster create k8s-class --servers 1 --agents 2 -p "80:80@loadbalancer" -p "443:443@loadbalancer"
```

## 2. Application Build
Move into the app directory and build the Docker image.
```powershell
cd day1/app
docker build -t k8s-day1:v1 .
```

## 3. Image Import
Since k3d runs in containers, we must import our local image into the cluster nodes so they can "pull" it.
```powershell
k3d image import k8s-day1:v1 -c k8s-class
```

## 4. Deployment
Deploy the Pod manifest using `kubectl`.
```powershell
cd ..  # Go back to day1 root
kubectl apply -f pod.yaml
```

## 5. Inspection & Debugging
Verify the Pod is running and check its details.

**Check Status:**
```powershell
kubectl get pods
```

**See System Pods (CoreDNS, Traefik, etc.):**
```powershell
kubectl get pods -A
```

**See IP and Node Assignment:**
```powershell
kubectl get pods -o wide
```

**Check Application Logs:**
```powershell
kubectl logs go-web-app
```

**Deep Dive (Events, IP, Node info):**
```powershell
kubectl describe pod go-web-app
```

## 6. Testing & Debugging the Application

### Method 1: Port Forwarding (Test API Locally)
Forward the pod's port to your local machine to test the API:
```powershell
kubectl port-forward pod/go-web-app 8080:8080
```
Then test in another terminal or browser:
```powershell
curl http://localhost:8080
# Or open http://localhost:8080 in browser
```

### Method 2: Shell Access (Debug Inside Container)
If your image includes a shell (like alpine), you can exec into the pod:
```powershell
# Interactive shell
kubectl exec -it go-web-app -- /bin/sh

# Run single commands
kubectl exec go-web-app -- ls -la /
kubectl exec go-web-app -- ps aux
kubectl exec go-web-app -- curl localhost:8080
```

**Why might exec fail?**
If you get `exec: "/bin/sh": no such file or directory`, your Docker image uses a minimal base like `scratch` that has no shell. See section below for solution.

### Understanding Docker Base Images

| Base Image | Size | Has Shell? | Use Case |
|------------|------|------------|----------|
| `scratch` | ~0 MB | ❌ No | Production (max security, minimal attack surface) |
| `distroless` | ~2 MB | ❌ No | Production (slightly more features) |
| `alpine` | ~5 MB | ✅ Yes (`/bin/sh`) | **Learning/Debugging** (small but functional) |
| `ubuntu/debian` | ~50+ MB | ✅ Yes (`/bin/bash`) | Development (full tooling) |

### Switching from scratch to alpine (for debugging)

Our Dockerfile v1 used `scratch` (no shell). To enable debugging, we switched to `alpine`:

**Edit `app/Dockerfile`:**
```dockerfile
# Change line 8 from:
FROM scratch

# To:
FROM alpine:latest  # Now includes /bin/sh for debugging
```

**Rebuild and redeploy:**
```powershell
# Build new version
cd app
docker build -t k8s-day1:v2 .

# Import to k3d
cd ..
k3d image import k8s-day1:v2 -c k8s-class

# Update pod.yaml to use v2, then redeploy
kubectl delete pod go-web-app
kubectl apply -f pod.yaml
```

💡 **Tip:** Use `alpine` for learning Kubernetes, switch to `scratch` for production deployments.

## 7. Cleanup (Optional)
To delete the pod:
```powershell
kubectl delete pod go-web-app
```
