# Kubernetes 15-Day Crash Course

Welcome to the **"Zero to Production Hero"** Kubernetes challenge!

## 🎯 Goal
Master Kubernetes in 15 days using a hands-on, debug-first approach. By the end of this course, you will be able to:
- Deploy production-ready microservices (Go-based).
- Debug complex failures (CrashLoopBackOff, PVC issues, Ingress errors).
- Manage resources, security, and auto-scaling.
- Answer senior DevOps interview questions with confidence.

## 📚 The Plan
The complete day-by-day syllabus is available in **[plan.md](./plan.md)**.

## 🛠️ Prerequisites & Setup
We are using **k3s** on **Windows** (via WSL2 recommended) for a lightweight, production-grade local cluster.

### 1. Install WSL2
Ensure you have WSL2 enabled on your Windows machine.
```powershell
wsl --install
```

### 2. Install k3s (via Rancher Desktop or K3d)
For Windows users, the easiest way to get a full k3s experience is using **Rancher Desktop** or **k3d** (k3s in docker).

**Option A: Rancher Desktop (Recommended for GUI)**
1. Download [Rancher Desktop](https://rancherdesktop.io/).
2. Enable "Kubernetes" in settings.
3. Select "dockerd" as the container runtime.

**Option B: k3d (Recommended for CLI speed)**
1. Install Docker Desktop for Windows.
2. Install k3d: `winget install k3d`
3. Create a cluster:
   ```bash
   k3d cluster create my-cluster --servers 1 --agents 2 -p "80:80@loadbalancer" -p "443:443@loadbalancer"
   ```

### 3. Verify Installation
```bash
kubectl get nodes
kubectl version
```

## 🚀 Getting Started

### Course Progress
- ✅ **[Day 1: The Pod & The Container](./day1/README.md)** - Pods, containers, debugging basics
- ✅ **[Day 2: Controllers & Scalability](./day2/README.md)** - Deployments, ReplicaSets, health probes
- ✅ **[Day 3: Networking I - Services](./day3/README.md)** - Microservices, ClusterIP, Service Discovery, Scaling
- ✅ **[Day 4: Storage & State](./day4/README.md)** - PV, PVC, StorageClass, Persistent Databases
- ✅ **[Day 5: Rust API Experiment](./day5/README.md)** - Rust + Postgres, Multi-stage Docker Builds, k3d Image Importing

Open [plan.md](./plan.md) for the complete 15-day syllabus

Good luck!
