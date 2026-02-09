Write-Host "Cleaning up Day 6 resources..."

# Delete Kubernetes Resources
Write-Host "Deleting Ingress..."
kubectl delete -f ingress.yaml --ignore-not-found

Write-Host "Deleting Service..."
kubectl delete -f service.yaml --ignore-not-found

Write-Host "Deleting Deployment..."
kubectl delete -f deployment.yaml --ignore-not-found

Write-Host "Deleting TLS Secret..."
kubectl delete secret whoami-tls --ignore-not-found

# Delete generated local files
if (Test-Path "tls.key") {
    Remove-Item "tls.key"
    Write-Host "Deleted tls.key"
}

if (Test-Path "tls.crt") {
    Remove-Item "tls.crt"
    Write-Host "Deleted tls.crt"
}

Write-Host "Cleanup complete."
Write-Host "Note: You may want to remove '127.0.0.1 api.local' from your hosts file manually."
