$ErrorActionPreference = "Stop"

Remove-Item -Recurse -Force .\generated -ErrorAction SilentlyContinue
Remove-Item -Recurse -Force .\build -ErrorAction SilentlyContinue
Write-Host "Clean complete."
