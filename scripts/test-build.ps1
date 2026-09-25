$ErrorActionPreference = "Stop"

go run .\cmd\toolchain-test -build

if ($LASTEXITCODE -ne 0) {
    throw "Real APK build test failed."
}

Write-Host "Real APK build test completed."
