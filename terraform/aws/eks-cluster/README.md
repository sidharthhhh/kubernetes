# AWS EKS Setup via Terraform (Infrastructure as Code)

Welcome to the professional way to manage cloud infrastructure! Instead of clicking through the UI, we use **Terraform** to declare exactly what we want in code, and Terraform figures out how to build it in AWS.

This module uses the official `terraform-aws-modules` to provision a complete Virtual Private Cloud (VPC) and an Elastic Kubernetes Service (EKS) cluster.

## Prerequisites
1. Installed **Terraform**.
2. Installed the **AWS CLI** and configured your credentials (see `aws-configration/aws-cli-setup.md`).

---

## Deployment Steps

### 1. Initialize Terraform
This downloads the required providers (AWS) and modules (VPC, EKS).
```bash
terraform init
```

### 2. Plan the Infrastructure
This command tells you exactly what Terraform is going to create without actually creating anything yet. It's a dry run.
```bash
terraform plan
```
*You should see a massive list of resources that will be added (around 50+ resources including routing tables, subnets, node groups, etc).*

### 3. Apply the Changes (Deploy!)
This will provision the infrastructure.
```bash
terraform apply
```
*Type `yes` when prompted.*
⏳ **Note:** EKS clusters take about **15-20 minutes** to fully provision. Go grab a coffee!

---

## Connect to your Cluster
Once `terraform apply` finishes successfully, it will print out the command to configure your local `kubectl`. It looks like this:
```bash
aws eks update-kubeconfig --region us-east-1 --name my-prod-cluster
```

Run that command, then verify your nodes are ready:
```bash
kubectl get nodes
```

---

## 🛑 Destruction (Crucial Step)
The beauty of Terraform is that it tracks exactly what it created. You don't have to hunt through the AWS Console to delete hidden resources.

To delete **everything** and stop AWS from billing you:
```bash
terraform destroy
```
*Type `yes` when prompted.*

Terraform will systematically tear down the node group, the cluster, the subnets, and the VPC in the correct order.
