# Path to OpenSSL (usually found in Git installation on Windows)
$opensslPath = "C:\Program Files\Git\usr\bin\openssl.exe"

if (-not (Test-Path $opensslPath)) {
    Write-Error "OpenSSL not found at $opensslPath. Please install Git for Windows or OpenSSL."
    exit 1
}

Write-Host "Generating Self-Signed Certificate using OpenSSL..."

# Generate Private Key and Certificate
& $opensslPath req -x509 -newkey rsa:4096 -sha256 -nodes -keyout tls.key -out tls.crt -subj "/CN=api.local" -days 365

if ($LASTEXITCODE -eq 0) {
    Write-Host "Certificate (tls.crt) and Private Key (tls.key) generated successfully."
    
    # Create Kubernetes Secret
    Write-Host "Creating Kubernetes TLS Secret 'whoami-tls'..."
    kubectl create secret tls whoami-tls --cert=tls.crt --key=tls.key --dry-run=client -o yaml | kubectl apply -f -
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Secret 'whoami-tls' created/updated."
    }
    else {
        Write-Error "Failed to create Kubernetes secret."
    }
}
else {
    Write-Error "Failed to generate certificate."
}
