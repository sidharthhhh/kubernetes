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

## 6. Cleanup (Optional)
To delete the pod:
```powershell
kubectl delete pod go-web-app
```
