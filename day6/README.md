# Day 6: Networking II - Ingress & TLS

## Goal
Learn how to expose your services to the outside world using an **Ingress Controller** and secure them with **TLS**.

We will use the simple `traefik/whoami` application, which prints out the request headers and IP addresses, making it perfect for verifying that our traffic is passing through the Ingress correctly.

## Prerequisites
- You have a running Kubernetes cluster (e.g., k3s, minikube).
- If using k3s, **Traefik** is likely already installed and enabled as the default Ingress controller.

## Steps

### 1. Deploy the Application and Service
First, we deploy the backend application and the internal Service.

```powershell
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

Verify they are running:
```powershell
kubectl get pods
kubectl get svc
```
You should see `whoami-...` pods and a `whoami-service`.

### 2. Configure Ingress
Now we create the Ingress resource which tells the Ingress Controller (Traefik) how to route traffic.

```powershell
kubectl apply -f ingress.yaml
```

Check the status of the Ingress:
```powershell
kubectl get ingress
```
Wait until you see an **ADDRESS** (it might be your node IP or localhost).

### 3. DNS Spoofing (The `hosts` File)
Since we are using a fake domain (`api.local`), we need to tell our computer how to find it.

1.  Open **Notepad** as **Administrator**.
2.  Open the file: `C:\Windows\System32\drivers\etc\hosts`
3.  Add the following line at the bottom:
    ```text
    127.0.0.1  api.local
    ```
    *(Note: If you are using k3s/docker-desktop, `127.0.0.1` usually works. If you are using a VM or a specific node IP, replace `127.0.0.1` with that IP.)*

4.  Save the file.

### 4. Test the Ingress
Open a terminal (PowerShell) and try to curl the endpoint:

```powershell
curl http://api.local
```

Or open `http://api.local` in your browser.

**Success:** You should see output starting with `Hostname: whoami-...` and details about your request.

### 5. (Optional) TLS Setup
To secure this with HTTPS (even with a fake certificate):

**Prerequisite:** You need `openssl`. 
- If you have **Git for Windows** installed, simply right-click in the `day6` folder and select **Git Bash Here**. Run the commands below in that window.
- If you only have PowerShell and no `openssl`, it is hard to generate the `.key` file in the correct format. It is highly recommended to install Git for Windows or OpenSSL.

1.  **Generate a generic self-signed certificate:**
    ```bash
    openssl req -x509 -newkey rsa:4096 -sha256 -nodes -keyout tls.key -out tls.crt -subj "/CN=api.local" -days 365
    ```

    *If you verify files exist:* `ls tls.*` should show `tls.key` and `tls.crt`.

2.  **Create a Kubernetes Secret:**
    ```powershell
    kubectl create secret tls whoami-tls --cert=tls.crt --key=tls.key
    ```

3.  **Update Ingress to use TLS:**
    Edit `ingress.yaml` to include the `tls` section:
    ```yaml
    spec:
      tls:
      - hosts:
        - api.local
        secretName: whoami-tls
      rules:
      ...
    ```
    *(Or just uncomment the TLS section if provided in a later step)*.

4.  **Apply changes:**
    ```powershell
    kubectl apply -f ingress.yaml
    ```

5.  **Test HTTPS:**
    ```powershell
    curl -k https://api.local
    ```
    (The `-k` flag tells curl to ignore the insecure self-signed certificate).

## Troubleshooting
- **404 Not Found:** Check if `kubectl get ingress` shows the correct Address. Verify your `hosts` file.
- **502 Bad Gateway:** The Ingress controller can't reach your Service. Check if the Pods are running and the Service port (80) is correct.
- **`curl` Error (PowerShell):** PowerShell's `curl` alias doesn't fully support all flags (like `-k`). Use `curl.exe` instead.
- **"Remote name could not be resolved":** add `127.0.0.1 api.local` to `C:\Windows\System32\drivers\etc\hosts`.
- **Verify without hosts file:** Run `.\test_ingress.ps1` to test connectivity by bypassing DNS.

## Architecture
See [architecture.md](architecture.md) for a visual diagram of the request flow and TLS termination.

## Cleanup
To remove all resources created in this exercise, run:
```powershell
.\cleanup.ps1
```
This will delete the Ingress, Service, Deployment, TLS Secret, and the generated certificate files.
