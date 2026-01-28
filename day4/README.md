# Day 4: Storage & State - Persistent Data in Kubernetes

Welcome to Day 4! Today we tackle one of the most important concepts in Kubernetes: **State**. 

By default, Pods are ephemeral. If a Pod crashes or is deleted, any data written to its local filesystem is **LOST**. This is fine for stateless applications (like our Go API from Day 1-3), but terrible for databases like MySQL.

To save data, we use **PersistentVolumes (PV)** and **PersistentVolumeClaims (PVC)**.

## 📚 Concepts

### 1. PersistentVolume (PV)
Think of this as the actual hard drive or cloud storage (e.g., AWS EBS, Google Disk, or a local folder on the node). It is a piece of storage in the cluster.

### 2. PersistentVolumeClaim (PVC)
Think of this as a **"Ticket"** or **"Request"** for storage. A developer creates a PVC saying "I need 1GB of storage". Kubernetes looks for a PV that matches this request and "binds" them together.

### 3. StorageClass
This is the "magician" that automatically creates PVs for you (Dynamic Provisioning). In `k3s`, the default StorageClass uses the host's disk (Local Path).

---

## 🛠️ Lab: Deploying MySQL with Persistence

We will deploy MySQL, but instead of letting it write to the container's temporary storage, we will mount a PVC.

### Step 1: Create the Storage "Ticket" (PVC)

Create `mysql-pvc.yaml`:

```yaml
# mysql-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc # The name we will reference in the Deployment
spec:
  accessModes:
    - ReadWriteOnce # RWO: Can be mounted by one node as read-write
  resources:
    requests:
      storage: 1Gi # We request 1 Gigabyte of space
```

Apply it:
```bash
kubectl apply -f mysql-pvc.yaml
kubectl get pvc
# Status should be 'Bound' (if strictly dynamic) or 'Pending' (until a Pod uses it, depends on StorageClass)
```

### Step 2: Deploy MySQL using the PVC

Create `mysql-deployment.yaml`:

```yaml
# mysql-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
spec:
  selector:
    matchLabels:
      app: mysql
  strategy:
    type: Recreate # IMPORTANT: For RWO volumes, we can't have 2 pods running at once during update
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - image: mysql:5.7
        name: mysql
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: "password" # (In Day 5 we will secure this!)
        ports:
        - containerPort: 3306
          name: mysql
        volumeMounts:
        - name: mysql-persistent-storage # Must match volume name below
          mountPath: /var/lib/mysql # Where MySQL stores data inside container
      volumes:
      - name: mysql-persistent-storage
        persistentVolumeClaim:
          claimName: mysql-pvc # Must match the PVC name we created
```

Apply it:
```bash
kubectl apply -f mysql-deployment.yaml
kubectl get pods
```

### Step 3: Verify Persistence (The "Kill" Test)

1.  **Exec into the MySQL Pod and create data:**
    ```bash
    # Get pod name
    kubectl get pods
    
    # Enter pod
    kubectl exec -it <mysql-pod-name> -- mysql -uroot -ppassword
    
    # Inside MySQL:
    CREATE DATABASE day4k8s;
    USE day4k8s;
    CREATE TABLE important_data (message VARCHAR(255));
    INSERT INTO important_data VALUES ('Data persists even if Pod dies!');
    SELECT * FROM important_data;
    EXIT;
    ```

2.  **Delete the Pod:**
    ```bash
    kubectl delete pod <mysql-pod-name>
    ```
    *Kubernetes will notice the Pod is gone and create a NEW one automatically.*

3.  **Verify Data is still there:**
    ```bash
    # Wait for new pod to be Running
    kubectl get pods
    
    # Check data in NEW pod
    kubectl exec -it <new-pod-name> -- mysql -uroot -ppassword -e "SELECT * FROM day4k8s.important_data;"
    ```
    
    ✅ If you see the message, congrats! You have persistent storage.

---

## 🐛 Debugging Common Issues

### 1. PVC is "Pending" forever
*   **Cause:** No PV available that matches the claim, or no StorageClass defined.
*   **Fix:** Check events: `kubectl describe pvc mysql-pvc`.

### 2. Pod is "Pending"
*   **Cause:** The PVC is not bound yet, or the node has no disk space.
*   **Fix:** `kubectl describe pod <pod-name>`.

### 3. "Multi-Attach" Error
*   **Cause:** You have `ReadWriteOnce` volume and tried to run 2 Replicas of MySQL.
*   **Fix:** Use `strategy: Recreate` in Deployment or ensure only 1 replica exists.

---

## 🧹 Cleanup

When you are done with the practical, it's good practice to clean up your resources.

1.  **Delete the MySQL App (Pod & Deployment):**
    ```bash
    kubectl delete -f mysql-deployment.yaml
    ```

2.  **Delete the Storage (PVC & PV):**
    *   *Warning: This will delete the persistence "ticket". In most clouds (and local-path), this will effectively WIPE the data from the disk.*
    ```bash
    kubectl delete -f mysql-pvc.yaml
    ```

---

## ⚡ Command Cheat Sheet

### Lifecycle
```bash
kubectl apply -f mysql-pvc.yaml         # 1. Create Storage Request
kubectl apply -f mysql-deployment.yaml  # 2. Create App
kubectl delete -f mysql-deployment.yaml # 3. Delete App
kubectl delete -f mysql-pvc.yaml        # 4. Delete Storage
```

### Verification
```bash
kubectl get pvc              # Check storage status (Pending -> Bound)
kubectl get pods             # Check app status
kubectl get sc               # List StorageClasses
kubectl get pv               # List physical volumes (cluster-wide)
```

### Persistence Test (The "Kill" Logic)
```bash
# 1. Enter Pod
kubectl exec -it <pod-name> -- mysql -uroot -ppassword

# 2. Database Commands
# CREATE DATABASE test;
# CREATE TABLE data (id INT);
# EXIT;

# 3. Kill Pod
kubectl delete pod <pod-name>

# 4. Verify in New Pod
kubectl exec -it <new-pod-name> -- mysql -uroot -ppassword -e "SHOW DATABASES;"
```

### 🐛 Debugging
```bash
# Why is my PVC Pending?
kubectl describe pvc mysql-pvc

# Why is my Pod Pending?
kubectl describe pod <pod-name>

# Check Storage Provisioner logs (Advanced)
kubectl logs -n kube-system -l app=local-path-provisioner
```


