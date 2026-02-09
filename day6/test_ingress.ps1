Write-Host "Testing Ingress connectivity using curl.exe and bypassing DNS..."
Write-Host "Target: https://api.local/ (resolving to 127.0.0.1)"

# Check if curl.exe is available
if (Get-Command curl.exe -ErrorAction SilentlyContinue) {
    # Use curl.exe with -k (insecure) and --resolve to force api.local to 127.0.0.1
    # This bypasses the need for hosts file modification for this specific test
    curl.exe -k --resolve api.local:443:127.0.0.1 https://api.local/ -v
}
else {
    Write-Error "curl.exe not found. Please ensure curl is installed (included in Windows 10/11)."
}

Write-Host "`n---------------------------------------------------"
Write-Host "NOTE: If you see a successful connection above,"
Write-Host "but 'curl https://api.local' fails, it means your"
Write-Host "HOSTS file does not have the entry: 127.0.0.1 api.local"
Write-Host "---------------------------------------------------"
