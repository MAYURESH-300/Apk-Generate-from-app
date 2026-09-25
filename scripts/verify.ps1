$ErrorActionPreference = "Stop"

Write-Host "=== VS Code Toolchain Verification ==="

if (-not (Test-Path ".\config.json")) {
    throw "config.json not found."
}

go run .\cmd\toolchain-test

if ($LASTEXITCODE -ne 0) {
    throw "Toolchain verification failed."
}

Write-Host "Verification complete."
