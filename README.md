# ApkGenerator VS Code Test Project

This is a standalone VS Code/GitHub test project for the real APK build pipeline.

## What this test does

It uses the build/code-generation logic extracted from H2APK as a reference and provides a small Go CLI that:

1. Generates a real Android project from HTML/CSS/JS input.
2. Uses the local Android SDK/JDK/Gradle tools if configured.
3. Runs the real build commands.
4. Verifies that an APK was actually produced.
5. Does NOT fake a successful build.

## Important

This repository does not contain Android/JDK binaries. Put your Android-compatible toolchain on the machine and configure `config.json`.

The Windows `aapt2.exe` and `zipalign.exe` files from the original H2APK repository are intentionally not bundled here because this project is intended to test the pipeline before moving it into the Android app.

## Requirements

- VS Code
- Go 1.22+
- JDK 17
- Android SDK with:
  - platform android-37
  - build-tools 36.0.0
- Gradle 9.4.1 or compatible Gradle for the selected project
- A working `aapt2`, `d8`, `zipalign`, and `apksigner`

## Test

PowerShell:

```powershell
go run .\cmd\toolchain-test
```

Then:

```powershell
go run .\cmd\toolchain-test -build
```

The first command checks the configured toolchain.

The second command generates a sample Android project and invokes the real build pipeline.

A successful test is only reported when an actual APK exists.

## Configuration

Edit `config.json` if your SDK/JDK/Gradle locations differ.

Example:

```json
{
  "java": "C:\\Java\\jdk-17\\bin\\java.exe",
  "gradle": "C:\\Gradle\\gradle-9.4.1\\bin\\gradle.bat",
  "androidSdk": "C:\\Android\\Sdk",
  "buildTools": "36.0.0",
  "compileSdk": "37"
}
```
