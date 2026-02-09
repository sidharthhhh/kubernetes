# Architecture Diagram

This diagram visualizes the flow of traffic from the client to the application, including TLS termination at the Ingress level.

```mermaid
graph TD
    Client["Client (curl / Browser)"] -->|1. HTTPS Request (api.local)| Ingress[Traefik Ingress Controller]
    
    subgraph "Kubernetes Cluster"
        direction TB
        Ingress -->|"2. Decrypt & Route"| Service[Service: whoami-service]
        Service -->|"3. Load Balance"| Pod[Pod: whoami]
        
        Secret[Secret: whoami-tls] -.->|"4. TLS Certificate"| Ingress
    end

    classDef k8s fill:#326ce5,stroke:#fff,stroke-width:2px,color:#fff;
    class Ingress,Service,Pod,Secret k8s;
```

## Flow Description
1.  **Client**: Sends an encrypted HTTPS request to `api.local` on port 443.
2.  **Ingress Controller**: 
    - Terminates TLS using the certificate stored in the `whoami-tls` Secret.
    - Inspects the Host header (`api.local`).
    - Routes the unencrypted traffic to the `whoami-service`.
3.  **Service**: Forwards the traffic to one of the available `whoami` pods.
4.  **Pod**: Processes the request and returns the response (acting as the backend).
