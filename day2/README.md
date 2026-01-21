# Day 2: Controllers & Scalability - Lab Guide

This lab covers **Deployments**, **ReplicaSets**, and **Health Probes** for production-ready applications.

## 📚 Learning Objectives

By the end of this lab, you will understand:
- How Deployments manage ReplicaSets and Pods
- Label selectors and how they connect components
- Rolling updates and rollback strategies
- Liveness, Readiness, and Startup probes
- How to debug common deployment issues

## 🎯 Prerequisites

- Completed Day 1 (Pod basics)
- k3d cluster running with `k8s-day1:v2` image imported
- `kubectl` configured

## 📋 Lab Exercises

### 1. Deploy the Application (3 Replicas)

Create a Deployment that manages 3 replicas of your Go web app.

```powershell
# Deploy the application
kubectl apply -f deployment.yaml

# Watch pods being created
kubectl get pods -w
# Press Ctrl+C to stop watching

# Check deployment status
kubectl get deployments
kubectl get replicasets
kubectl get pods
```

**Expected Output:**
- 1 Deployment
- 1 ReplicaSet (managed by Deployment)
- 3 Pods (managed by ReplicaSet)

**Observe:**
```powershell
# See the hierarchy
kubectl get all -l app=go-web

# See detailed deployment info
kubectl describe deployment go-web-app
```

### 2. Expose the Deployment with a Service

Services provide stable networking to your pods.

```powershell
# Create the service
kubectl apply -f service.yaml

# Verify the service
kubectl get svc go-web-app-service
kubectl describe svc go-web-app-service

# Check endpoints (should show 3 pod IPs)
kubectl get endpoints go-web-app-service
```

**Test the service:**
```powershell
# Port forward to the service (not individual pods)
kubectl port-forward svc/go-web-app-service 8080:80

# In another terminal, test it
curl http://localhost:8080
```

### 3. Scaling the Deployment

Scale up and down to see ReplicaSets in action.

```powershell
# Scale to 5 replicas
kubectl scale deployment go-web-app --replicas=5

# Watch it happen
kubectl get pods -w

# Scale back to 3
kubectl scale deployment go-web-app --replicas=3
```

**Key Observation:** The ReplicaSet controller automatically creates/deletes pods to match desired state.

### 4. Rolling Update (Zero Downtime)

Update your application without downtime.

```powershell
# Check current rollout status
kubectl rollout status deployment/go-web-app

# Apply the v2 deployment (has VERSION=v2.0.0 env var)
kubectl apply -f deployment-v2.yaml

# Watch the rolling update in real-time
kubectl rollout status deployment/go-web-app

# See both old and new ReplicaSets
kubectl get replicasets
```

**What's happening:**
1. New ReplicaSet created with updated template
2. New pods created gradually (maxSurge: 1)
3. Old pods terminated gradually (maxUnavailable: 1)
4. Old ReplicaSet scaled to 0 (but kept for rollback)

### 5. Pause and Resume a Rollout

Useful for gradual rollouts or canary deployments.

```powershell
# Make a change (just to trigger rollout)
kubectl set image deployment/go-web-app web-server=k8s-day1:v2

# Immediately pause it
kubectl rollout pause deployment/go-web-app

# Check status (should show paused)
kubectl rollout status deployment/go-web-app

# Resume when ready
kubectl rollout resume deployment/go-web-app
```

### 6. Rollback to Previous Version

Oh no! The new version has bugs. Roll back immediately.

```powershell
# Check rollout history
kubectl rollout history deployment/go-web-app

# Rollback to previous version
kubectl rollout undo deployment/go-web-app

# Check status
kubectl rollout status deployment/go-web-app

# Rollback to specific revision
kubectl rollout undo deployment/go-web-app --to-revision=1
```

### 7. Debug Scenario: Broken Readiness Probe

Deploy the intentionally broken configuration.

```powershell
# Apply the broken deployment
kubectl apply -f deployment-broken-probe.yaml

# Watch the pods
kubectl get pods

# Notice: Pods are Running but NOT Ready (0/1)
```

**Debug steps:**
```powershell
# Check pod events
kubectl describe pod <pod-name>
# Look for: "Readiness probe failed: HTTP probe failed"

# Check service endpoints (should be empty!)
kubectl get endpoints go-web-app-service
# No traffic will be routed to unhealthy pods

# Check the probe definition
kubectl get deployment go-web-app -o yaml | grep -A 5 readinessProbe
```

**Fix it:**
```powershell
# Reapply the correct deployment
kubectl apply -f deployment.yaml

# Watch pods become Ready
kubectl get pods -w
```

**Key Learning:** Readiness probes protect your service from routing traffic to broken pods.

## 🔍 Understanding Health Probes

### Readiness Probe
- **Purpose:** Is the pod ready to receive traffic?
- **Action on failure:** Remove from Service endpoints (no traffic)
- **Use case:** App is starting, warming up cache, or temporarily unavailable

### Liveness Probe
- **Purpose:** Is the pod alive or stuck/deadlocked?
- **Action on failure:** Restart the container
- **Use case:** App is frozen, deadlocked, or corrupted

### Startup Probe (Advanced)
- **Purpose:** Has the app finished starting? (for slow-starting apps)
- **Action on failure:** Kill container if startup takes too long
- **Use case:** Legacy apps with long initialization

## 🧪 Additional Experiments

### Experiment 1: What if I delete a Pod?
```powershell
kubectl delete pod <pod-name>
# Watch it get recreated immediately by ReplicaSet!
```

### Experiment 2: What if I delete the ReplicaSet?
```powershell
kubectl delete replicaset <replicaset-name>
# Watch it get recreated immediately by Deployment!
```

### Experiment 3: Test Auto-Healing
```powershell
# Exec into a pod and kill the server process
kubectl exec -it <pod-name> -- sh
kill 1  # Kill the main process

# Liveness probe will fail → Pod restarts
kubectl get pods -w
```

## 📊 Key kubectl Commands Summary

```powershell
# View all resources with label
kubectl get all -l app=go-web

# Check rollout status
kubectl rollout status deployment/go-web-app

# Rollout history
kubectl rollout history deployment/go-web-app

# Undo rollout
kubectl rollout undo deployment/go-web-app

# Pause/Resume rollout
kubectl rollout pause deployment/go-web-app
kubectl rollout resume deployment/go-web-app

# Scale
kubectl scale deployment/go-web-app --replicas=5

# Update image
kubectl set image deployment/go-web-app web-server=new-image:tag
```

## 🎓 Interview Questions (Practice)

1. **Q:** What happens if I delete a ReplicaSet that's managed by a Deployment?  
   **A:** The Deployment controller immediately recreates it to maintain desired state.

2. **Q:** What's the difference between Liveness and Readiness probes?  
   **A:** Liveness restarts the pod; Readiness removes it from Service endpoints (no traffic).

3. **Q:** Why would a Pod show `Running` but `0/1` in READY column?  
   **A:** The container is running, but the Readiness probe is failing.

4. **Q:** How does a rolling update achieve zero downtime?  
   **A:** It gradually creates new pods (maxSurge) while terminating old ones (maxUnavailable), ensuring minimum replicas are always running.

## 🧹 Cleanup

```powershell
# Delete all Day 2 resources
kubectl delete -f deployment.yaml
kubectl delete -f service.yaml

# Or delete by label
kubectl delete all -l app=go-web
```

## 🚀 Next Steps

- **Day 3:** Networking I - Services & Service Discovery (how pods communicate)
- **Challenge:** Try setting up a multi-tier app (frontend → backend → database) with proper Services

Good luck! 🎉
