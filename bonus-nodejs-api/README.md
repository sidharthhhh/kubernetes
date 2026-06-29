# Deploying a Node.js API to AWS EKS

Unlike running Kubernetes locally on your laptop, deploying to a cloud cluster like EKS requires that your Docker image is hosted in a public (or authenticated) registry like **Docker Hub** or **AWS ECR**, because the EKS worker nodes need to pull the image over the internet.

Here is exactly how to deploy the Node.js API we just created.

---

## 1. Build and Push the Docker Image

You will need a free [Docker Hub](https://hub.docker.com/) account for this.

1. Open your terminal and navigate to the `app` folder:
   ```bash
   cd bonus-nodejs-api/app
   ```

2. Login to Docker Hub from your terminal:
   ```bash
   docker login
   ```

3. Build the Docker image (replace `<YOUR_DOCKERHUB_USERNAME>` with your actual username):
   ```bash
   docker build -t sidharthhhh/nodejs-api:v1 .
   ```

4. Push the image to Docker Hub so EKS can download it:
   ```bash
   docker push sidharthhhh/nodejs-api:v1
   ```

---

## 2. Update the Deployment Manifest

Now we need to tell Kubernetes to use *your* newly uploaded image.

1. Open `bonus-nodejs-api/k8s/deployment.yaml`.
2. Find the line that says `image: <YOUR_DOCKERHUB_USERNAME>/nodejs-api:v1`.
3. Replace `<YOUR_DOCKERHUB_USERNAME>` with your actual username.

*(If you just want to test if EKS works without building your own image right now, you can change the image line to `image: paulbouwer/hello-kubernetes:1.10.1` and `containerPort: 8080`).*

---

## 3. Deploy to EKS!

Apply the manifests to your cluster:
```bash
cd ../k8s
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

Check if your pods are running:
```bash
kubectl get pods
```

---

## 4. Access your API via AWS Load Balancer

Because we set the Service type to `LoadBalancer` in `service.yaml`, AWS has automatically provisioned a real, public Load Balancer for you!

Get the public URL:
```bash
kubectl get svc nodejs-api-service
```

Look under the **EXTERNAL-IP** column. You will see a long AWS URL (e.g., `a1b2c3d4e5f6...us-east-1.elb.amazonaws.com`). 
*Note: It takes AWS about 2-3 minutes to fully provision the Load Balancer, so if the site doesn't load immediately, wait a minute and refresh.*

Copy that URL and paste it into your browser! You should see your Node.js API responding with a success message and the hostname of the exact Pod that served the request.
