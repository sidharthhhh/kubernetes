# AWS EKS Setup (UI Method)

Welcome to the cloud phase of our Kubernetes journey! Setting up a cluster via the UI is a great way to visually understand the moving parts before we move to Infrastructure as Code (Terraform).

Here is the simplified, beginner-friendly process to spin up an EKS cluster via the AWS Console.

## Prerequisites
1. An **AWS Account**.
2. Completed the **[AWS CLI Setup Guide](./aws-cli-setup.md)**.

---

## Step 1: Create the Cluster IAM Role
Your cluster needs permission to manage AWS resources (like Load Balancers).

1. Go to the **IAM** dashboard in the AWS Console.
2. Click **Roles** -> **Create role**.
3. Choose **AWS service** -> **EKS** -> **EKS - Cluster**.
4. Click **Next**, name it `my-eks-cluster-role`, and click **Create role**.

---

## Step 2: Create the VPC (Network)
Your cluster needs a secure, isolated network.

1. Go to the **VPC** dashboard.
2. Click **Create VPC** and select **VPC and more**.
3. Name it `eks-vpc`.
4. Choose **2 Availability Zones**, **2 Public Subnets**, and **2 Private Subnets**.
5. Set NAT gateways to **1 AZ**.
6. Click **Create VPC**.

---

## Step 3: Create the EKS Cluster (Control Plane)
1. Go to the **EKS** dashboard and click **Add cluster** -> **Create**.
2. Name it `my-prod-cluster`.
3. Select the `my-eks-cluster-role` from Step 1. Click Next.
4. Select your `eks-vpc-vpc` from Step 2.
5. Click **Next** on the remaining screens, then **Create**.

*Wait ~10 minutes for this to finish.*

---

## Step 4: Create the Node Group IAM Role
Worker nodes need permissions to run your pods and pull images.

1. Go back to **IAM** -> **Roles** -> **Create role**.
2. Choose **AWS service** -> **EC2**.
3. Attach these 3 policies:
   - `AmazonEKSWorkerNodePolicy`
   - `AmazonEC2ContainerRegistryReadOnly`
   - `AmazonEKS_CNI_Policy`
4. Click **Next**, name it `my-eks-node-role`, and click **Create**.

---

## Step 5: Add Compute (Nodes)
1. Go to your EKS Cluster, wait until it is **Active**.
2. Click the **Compute** tab -> **Add node group**.
3. Name it `my-node-group` and select `my-eks-node-role`.
4. Instance type: `t3.medium`.
5. Under Subnets, **uncheck** the Public subnets (we want nodes in Private subnets for security).
6. Click **Create**.

*Wait ~5 minutes for the nodes to spin up.*

---

## Step 6: Connect from your Terminal
Run this command to configure `kubectl`:
```bash
aws eks update-kubeconfig --region us-east-1 --name my-prod-cluster
```
*(Replace `us-east-1` with your actual region)*

Test it out!
```bash
kubectl get nodes
```

---

## 🧹 Crucial Cleanup
When you are done, **delete your cluster** so you don't get billed!
1. Go to EKS -> Compute tab -> select your Node Group -> **Delete**.
2. After the Node Group is gone, click **Delete cluster**.
3. Go to VPC dashboard and delete `eks-vpc`.
