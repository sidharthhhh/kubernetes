# Day 3: Networking I - Services & Service Discovery

## 📚 Learning Objectives

By the end of this lab, you will understand:
- How **Kubernetes Services** provide stable IPs and DNS names for Pods.
- The difference between **ClusterIP**, **NodePort**, and **LoadBalancer**.
- How **CoreDNS** enables service discovery (e.g., accessing Redis via `redis-service`).
- How to connect microservices (Go App -> Redis) using internal DNS.

## 🎯 Prerequisites

- Completed Day 2
- k3d cluster running
- Docker installed (to build the app image)

## 📋 Lab Exercises

### 1. Build and Import the Application Image

Since we are adding a Redis client to our code, we need to build a new version of our Go app.

```powershell
# Navigate to the day3/app directory
cd day3/app

# Build the Docker image
docker build -t k8s-day3-app:v1 .

# Import the image into your k3d cluster (replace 'my-cluster' with your cluster name)
k3d image import k8s-day3-app:v1 -c my-cluster
```

### 2. Deploy Redis (The Backend)

We need a database. We'll deploy Redis and expose it internally using a Service.

```powershell
# Navigate back to the day3 root
cd ..

# Apply the Redis deployment and service
kubectl apply -f manifests/redis.yaml

# Verify Redis is running
kubectl get pods -l app=redis
kubectl get svc redis-service
```

**Key Concept:** The Service `redis-service` creates a stable DNS name. Any pod in the cluster can now reach Redis at `redis-service:6379`.

### 3. Deploy the Go Web App

Now deploy the Go app. Look at `manifests/app.yaml` to see how we pass the Redis hostname:

```yaml
env:
- name: REDIS_HOST
  value: "redis-service"
```

Deploy it:

```powershell
kubectl apply -f manifests/app.yaml

# Wait for pods to be ready
kubectl get pods -l app=go-web-app -w
```

### 4. Test Service Discovery

Now verify that the Go app can talk to Redis.

```powershell
# Port forward to the Go app service
kubectl port-forward svc/go-web-app-service 8080:80

# In another terminal:
# 1. Check health
curl http://localhost:8080
# Output: Hello from Go! Running on Pod: ...

# 2. Increment the counter (Connecting to Redis)
curl http://localhost:8080/incr
# Output: Hits: 1
curl http://localhost:8080/incr
# Output: Hits: 2
```

If you see the "Hits" increasing, **Service Discovery is working!** The Go app resolved `redis-service` to the Redis Pod IP and connected successfully.

### 5. Inspecting DNS (Debugging)

How does this actually work? Let's look inside a Pod.

```powershell
# Exec into one of your Go app pods
kubectl exec -it <go-app-pod-name> -- sh

# Inside the pod, check DNS resolution
nslookup redis-service

# You should see something like:
# Server:    10.43.0.10
# Address:   10.43.0.10:53
# Name:      redis-service.default.svc.cluster.local
# Address:   10.43.x.x
```

## 🔍 Understanding Services

- **ClusterIP (Default):** Exposes the Service on a cluster-internal IP. Reachable only from within the cluster.
- **NodePort:** Exposes the Service on each Node's IP at a static port. Reachable from outside the cluster.
- **LoadBalancer:** Exposes the Service externally using a cloud provider's load balancer.

In this lab, we used **ClusterIP** for Redis (database should be internal) and the Go App (accessed via port-forward for security/learning).

## 🧹 Cleanup

```powershell
kubectl delete -f manifests/
```

## 🚀 Next Steps

- **Day 4:** Storage & State - Persistence with Volumes (PVs and PVCs)
