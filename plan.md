# 15-Day Kubernetes Crash Course: From Zero to Production Hero

**Target Audience:** Mid-level DevOps Engineer  
**Goal:** Production reliability, debugging mastery, and "real-world" operations.  
**Environment:** k3s on Windows (Production-like local setup), migrating to Cloud later.  
**App Stack:** Go-based Microservices (API, Worker), Redis, MySQL.

---

## **Phase 1: The Core Foundation (Days 1-5)**
*Focus: Internalizing the primitives and basic networking.*

### **Day 1: The Pod & The Container (The Real Basics)**
*Goal: Understand exactly what runs where and how to debug startup failures.*

*   **Concepts:**
    *   Pods vs Containers (Inter-process communication, shared network ns).
    *   Static Pods vs Managed Pods.
    *   `kubectl` deep dive: `run`, `get`, `describe`, `logs`, `exec`.
    *   YAML Anatomy: `apiVersion`, `kind`, `metadata`, `spec`, `status`.
*   **Lab:**
    1.  Deploy a simple Go web app Pod manually (imperative vs declarative).
    2.  Write the YAML by hand (no copy-paste).
    3.  Multi-container Pod: Main Go app + Sidecar (simple logger).
*   **Debug Scenario:**
    *   **Break:** Misconfigure the container command or image tag.
    *   **Fix:** Debug `ImagePullBackOff` and `CrashLoopBackOff`.
    *   **Command:** `kubectl describe pod <name>`, `kubectl logs <name> -c <container> --previous`.
*   **Pro Tip:** Always check `Events` in `describe` first.

### **Day 2: Controllers & Scalability (ReplicaSets & Deployments)**
*Goal: Managing state and updates without downtime.*

*   **Concepts:**
    *   ReplicaSets (Selectors & Labels).
    *   Deployments (Rolling Updates, Rollbacks, Strategies).
    *   Liveness vs Readiness vs Startup Probes (Crucial for Go apps).
*   **Lab:**
    1.  Create a Deployment for the Go API (3 replicas).
    2.  Perform a Rolling Update (change image version).
    3.  Pause and Resume a rollout.
    4.  Configure a Readiness probe that fails, observe traffic cut-off.
*   **Debug Scenario:**
    *   **Break:** Set a Readiness probe to a non-existent path.
    *   **Fix:** Observe Pods running but `READY 0/1`. Fix the probe.
*   **Interview Q:** "What happens if I delete a ReplicaSet managed by a Deployment?"

### **Day 3: Networking I - Services & Service Discovery**
*Goal: How Pods talk to each other reliably.*

*   **Concepts:**
    *   ClusterIP vs NodePort vs LoadBalancer vs ExternalName.
    *   CoreDNS & Service Discovery (DNS names: `svc.namespace.svc.cluster.local`).
    *   Endpoints and EndpointSlices.
*   **Lab:**
    1.  Deploy a Multi-Tier Architecture: Frontend (Go) -> Backend (Go/Todo API) -> Redis.
    2.  Implement Service Discovery: Frontend calls Backend via `http://todo-service`.
    3.  Backend connects to Redis via `http://redis-service` (Service Chaining).
    4.  Expose Frontend via Port Forwarding and test the full flow.
*   **Debug Scenario:**
    *   **Break:** Misconfigure the Backend URL in the Frontend (wrong DNS name).
    *   **Fix:** Use `nslookup` inside the pod to find the correct Service name.
*   **Pro Tip:** `kubectl port-forward` is your best friend for local debugging.

### **Day 4: Storage & State (PV, PVC, StorageClass)**
*Goal: Persistent data in a stateless world.*

*   **Concepts:**
    *   PV vs PVC vs StorageClass.
    *   Access Modes (RWO, RWX).
    *   Static vs Dynamic Provisioning (Local Path Provisioner in k3s).
*   **Lab:**
    1.  Deploy MySQL.
    2.  Create a PVC for MySQL data.
    3.  Kill the MySQL Pod, ensure data persists on restart.
*   **Debug Scenario:**
    *   **Break:** Request storage size larger than available or wrong StorageClass.
    *   **Fix:** Debug `Pending` PVCs.
*   **Interview Q:** "How do I expand a PVC without deleting it?"

### **Day 5: Configuration & Secrets (Decoupling Config from Code)**
*Goal: 12-Factor App compliance.*

*   **Concepts:**
    *   ConfigMaps (Env vars, Mounted volumes).
    *   Secrets (Opaque, TLS, DockerRegistry).
    *   Immutable ConfigMaps.
*   **Lab:**
    1.  Move Go app DB credentials to a Secret.
    2.  Move app configuration (API keys, feature flags) to ConfigMap.
    3.  Update ConfigMap and verify hot-reload (if app supports it) or requires restart.
*   **Debug Scenario:**
    *   **Break:** Refer to a non-existent Secret key in Deployment.
    *   **Fix:** Debug `CreateContainerConfigError`.

---

## **Phase 2: Production Operations (Days 6-10)**
*Focus: Traffic management, Observability, and Advanced Scheduling.*

### **Day 6: Networking II - Ingress & TLS**
*Goal: Exposing services to the world securely.*

*   **Concepts:**
    *   Ingress Controllers (Traefik in k3s vs Nginx).
    *   Ingress Resources (Paths, Host-based routing).
    *   TLS Termination (Cert-Manager basics).
*   **Lab:**
    1.  Create an Ingress rule for `api.local` -> Go Service.
    2.  Map `api.local` in Windows `hosts` file to k3s IP.
    3.  Generate self-signed cert and secure it.
*   **Debug Scenario:**
    *   **Break:** Wrong Service Port in Ingress YAML.
    *   **Fix:** Debug `502 Bad Gateway` vs `404 Not Found`.
*   **Pro Tip:** Always check Ingress Controller logs for backend connectivity issues.

### **Day 7: Advanced Scheduling & Node Management**
*Goal: Controlling exactly where workloads run.*

*   **Concepts:**
    *   Taints & Tolerations (Node-centric).
    *   Node Affinity & Anti-Affinity (Pod-centric).
    *   `cordon` and `drain` (Maintenance mode).
*   **Lab:**
    1.  Taint a node `env=prod:NoSchedule`.
    2.  Try to deploy a normal Pod (it stays Pending).
    3.  Add Toleration to the Go API deployment.
    4.  Drain a node safely.
*   **Debug Scenario:**
    *   **Break:** Use `requiredDuringScheduling...` with no matching nodes.
    *   **Fix:** Debug `Pending` Pods with `SchedulingFailed`.

### **Day 8: Resource Management & QoS**
*Goal: Preventing the "Noisy Neighbor" problem.*

*   **Concepts:**
    *   Requests vs Limits (CPU/Memory).
    *   QoS Classes: Guaranteed, Burstable, BestEffort.
    *   LimitRanges & ResourceQuotas (Namespace level).
    *   OOMKilled.
*   **Lab:**
    1.  Set CPU/Memory requests/limits for Go API.
    2.  Stress test the Pod to trigger `OOMKilled`.
    3.  Set a ResourceQuota on a namespace and try to exceed it.
*   **Interview Q:** "Why is setting CPU limits sometimes controversial?" (Throttling).

### **Day 9: Observability I - Monitoring (Prometheus & Grafana)**
*Goal: Knowing *what* is happening.*

*   **Concepts:**
    *   Metrics Server (for HPA/`kubectl top`).
    *   Prometheus Architecture (Scraping, Exporters).
    *   Grafana Dashboards.
*   **Lab:**
    1.  Install Prometheus/Grafana stack (Helm chart).
    2.  View Cluster stats (CPU/Mem usage).
    3.  Run `kubectl top pod` and `kubectl top node`.
    4.  Instrument Go app to expose `/metrics`.
*   **Debug Scenario:**
    *   **Break:** Metrics Server down.
    *   **Fix:** Cannot run `kubectl top` or HPA fails.

### **Day 10: Observability II - Logging & Tracing**
*Goal: Knowing *why* it happened.*

*   **Concepts:**
    *   EFK/ELK Stack (Elasticsearch, Fluentd/Fluentbit, Kibana) or Loki/Promtail.
    *   Structured Logging (JSON).
    *   Distributed Tracing (Jaeger/OpenTelemetry concept).
*   **Lab:**
    1.  Deploy Loki & Promtail (lightweight for learning).
    2.  Query logs in Grafana using LogQL.
    3.  Search for "error" in logs across all pods.
*   **Pro Tip:** Don't rely on `kubectl logs` for historical investigation; use a centralized system.

---

## **Phase 3: Production Hardening (Days 11-15)**
*Focus: Security, Automation, and Incident Response.*

### **Day 11: Auto-Scaling (HPA & VPA)**
*Goal: Handling traffic spikes automatically.*

*   **Concepts:**
    *   HPA (Horizontal Pod Autoscaler).
    *   VPA (Vertical Pod Autoscaler).
    *   Custom Metrics (Scaling on Queue depth vs CPU).
*   **Lab:**
    1.  Set up HPA for Go API (scale on CPU > 50%).
    2.  Generate load (use `hey` or `k6`).
    3.  Watch replicas increase and decrease.
*   **Debug Scenario:**
    *   **Break:** Remove resource requests from Deployment.
    *   **Fix:** HPA will not calculate metrics (`<unknown>`).

### **Day 12: Security & RBAC**
*Goal: Least privilege access.*

*   **Concepts:**
    *   Role, ClusterRole, RoleBinding, ClusterRoleBinding.
    *   ServiceAccounts.
    *   `kubectl auth can-i`.
*   **Lab:**
    1.  Create a "developer" user (Certificate based or rigid ServiceAccount for testing).
    2.  Grant "view" access only to a specific namespace.
    3.  Try to delete a pod as "developer" (Permission Denied).
    4.  Grant "edit" access and retry.
*   **Interview Q:** "Difference between RoleBinding and ClusterRoleBinding?"

### **Day 13: StatefulSets & Headless Services**
*Goal: Running databases or distributed systems.*

*   **Concepts:**
    *   StatefulSet vs Deployment (Identity, Ordering, Storage).
    *   Headless Service (No ClusterIP).
    *   Stable Network ID (`pod-0.svc`, `pod-1.svc`).
*   **Lab:**
    1.  Deploy a Redis Cluster or CockroachDB using StatefulSet.
    2.  Scale up and down; observe ordered termination.
    3.  Check DNS resolution of individual pods.

### **Day 14: Helm & GitOps (ArgoCD)**
*Goal: Managing "YAML Hell" and Deployment pipelines.*

*   **Concepts:**
    *   Helm Charts (Templates, Values).
    *   GitOps Principles.
    *   ArgoCD (Syncing Git to Cluster).
*   **Lab:**
    1.  Package the Go App into a Helm Chart.
    2.  Install ArgoCD.
    3.  Connect ArgoCD to a Git repo.
    4.  Push a change to Git -> Watch ArgoCD sync it to the cluster.

### **Day 15: The Final Exam - "Chaos Engineering"**
*Goal: Simulation of a bad day on-call.*

*   **Scenario:**
    1.  Deploy the full stack: Go API, Worker, Redis, MySQL.
    2.  **Break it:**
        *   Delete the ConfigMap used by the API.
        *   Taint the generic nodes so pods can't schedule.
        *   Fill up the disk (simulate PV full).
        *   Corrupt the DNS.
    3.  **Fix it:** Use all the tools learned (`events`, `logs`, `describe`, `top`, `nslookup`) to restore the service.

---

## **Kubectl Cheat Sheet for Mastery**

### **Inspection**
```bash
kubectl get pods -o wide --show-labels
kubectl describe pod <pod-name>
kubectl logs <pod-name> -c <container-name> --previous
kubectl get events --sort-by='.lastTimestamp'
```

### **Debugging**
```bash
kubectl run -it --rm debug --image=curlimages/curl -- sh  # Network debug shell
kubectl port-forward svc/<service-name> 8080:80           # Local access
kubectl auth can-i delete pods --as system:serviceaccount:default:my-sa
```

### **Quick Fixes**
```bash
kubectl rollout restart deployment/<name>
kubectl scale deployment/<name> --replicas=0  # Stop
kubectl scale deployment/<name> --replicas=3  # Start
kubectl delete pod <name> --grace-period=0 --force # Stuck Terminating (Careful!)
```

## **Next Steps**
1.  Setup your environment (k3s on Windows/WSL2).
2.  Start with Day 1.
3.  Document your failures! The learning happens when things break.
