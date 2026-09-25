// // package main

// // import (
// // 	"encoding/json"
// // 	"flag"
// // 	"fmt"
// // 	"os"
// // 	"os/exec"
// // 	"path/filepath"
// // 	"runtime"
// // 	"strings"
// // )

// // type Config struct {
// // 	Java       string `json:"java"`
// // 	Gradle     string `json:"gradle"`
// // 	AndroidSDK string `json:"androidSdk"`
// // 	BuildTools string `json:"buildTools"`
// // 	CompileSDK string `json:"compileSdk"`
// // }

// // func main() {
// // 	build := flag.Bool("build", false, "generate and build a real sample APK")
// // 	flag.Parse()

// // 	cfg, err := loadConfig()
// // 	if err != nil {
// // 		fail(err)
// // 	}

// // 	fmt.Println("=== ApkGenerator VS Code Toolchain Test ===")
// // 	fmt.Println("OS:", runtime.GOOS, runtime.GOARCH)

// // 	if err := checkConfig(&cfg); err != nil {
// // 		fail(err)
// // 	}

// // 	if err := checkJava(cfg.Java); err != nil {
// // 		fail(err)
// // 	}
// // 	if err := checkGradle(cfg.Gradle); err != nil {
// // 		fail(err)
// // 	}
// // 	if err := checkSDK(&cfg); err != nil {
// // 		fail(err)
// // 	}

// // 	fmt.Println()
// // 	fmt.Println("TOOLCHAIN CHECK: PASS")

// // 	if !*build {
// // 		fmt.Println()
// // 		fmt.Println("Run with -build to generate and build a real APK.")
// // 		return
// // 	}

// // 	project := filepath.Join("generated", "HelloApk")
// // 	if err := generateProject(project, &cfg); err != nil {
// // 		fail(err)
// // 	}

// // 	fmt.Println()
// // 	fmt.Println("=== REAL BUILD ===")
// // 	if err := runGradle(project, cfg.Gradle, cfg.AndroidSDK); err != nil {
// // 		fail(err)
// // 	}

// // 	apk := filepath.Join(project, "app", "build", "outputs", "apk", "debug", "app-debug.apk")
// // 	info, err := os.Stat(apk)
// // 	if err != nil || info.Size() == 0 {
// // 		fail(fmt.Errorf("build finished without a real APK: %s", apk))
// // 	}

// // 	fmt.Println()
// // 	fmt.Println("BUILD RESULT: SUCCESS")
// // 	fmt.Println("APK:", apk)
// // 	fmt.Printf("APK size: %d bytes\n", info.Size())
// // }

// // func loadConfig() (Config, error) {
// // 	data, err := os.ReadFile("config.json")
// // 	if err != nil {
// // 		return Config{}, err
// // 	}
// // 	var cfg Config
// // 	if err := json.Unmarshal(data, &cfg); err != nil {
// // 		return Config{}, err
// // 	}
// // 	return cfg, nil
// // }

// // func checkConfig(cfg *Config) error {
// // 	if strings.TrimSpace(cfg.Java) == "" ||
// // 		strings.TrimSpace(cfg.Gradle) == "" ||
// // 		strings.TrimSpace(cfg.AndroidSDK) == "" {
// // 		return fmt.Errorf("edit config.json and set java, gradle and androidSdk paths")
// // 	}
// // 	return nil
// // }

// // func checkJava(java string) error {
// // 	fmt.Println()
// // 	fmt.Println("[1/3] Java")
// // 	out, err := commandOutput(java, "-version")
// // 	if err != nil {
// // 		return fmt.Errorf("Java failed: %w\n%s", err, out)
// // 	}
// // 	fmt.Print(out)
// // 	return nil
// // }

// // func checkGradle(gradle string) error {
// // 	fmt.Println()
// // 	fmt.Println("[2/3] Gradle")
// // 	out, err := commandOutput(gradle, "--version")
// // 	if err != nil {
// // 		return fmt.Errorf("Gradle failed: %w\n%s", err, out)
// // 	}
// // 	fmt.Print(out)
// // 	return nil
// // }

// // func normalizeCompileSDK(version string) string {
// // 	version = strings.TrimSpace(version)
// // 	if version == "" {
// // 		return version
// // 	}
// // 	return strings.TrimSuffix(version, ".0")
// // }

// // func resolvePlatformDir(androidSDK, compileSDK string) (string, error) {
// // 	platformsDir := filepath.Join(androidSDK, "platforms")
// // 	variants := []string{
// // 		compileSDK,
// // 		normalizeCompileSDK(compileSDK),
// // 		compileSDK + ".0",
// // 		compileSDK + ".1",
// // 		normalizeCompileSDK(compileSDK) + ".0",
// // 		normalizeCompileSDK(compileSDK) + ".1",
// // 	}

// // 	seen := map[string]bool{}
// // 	for _, v := range variants {
// // 		if v == "" || seen[v] {
// // 			continue
// // 		}
// // 		seen[v] = true

// // 		candidate := filepath.Join(platformsDir, "android-"+v)
// // 		if _, err := os.Stat(candidate); err == nil {
// // 			return candidate, nil
// // 		}
// // 	}

// // 	return "", fmt.Errorf("unable to find an Android platform directory for compileSdk %q under %s", compileSDK, platformsDir)
// // }

// // func checkSDK(cfg *Config) error {
// // 	fmt.Println()
// // 	fmt.Println("[3/3] Android SDK")

// // 	platformDir, err := resolvePlatformDir(cfg.AndroidSDK, cfg.CompileSDK)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	platform := filepath.Join(platformDir, "android.jar")
// // 	buildTools := filepath.Join(cfg.AndroidSDK, "build-tools", cfg.BuildTools)

// // 	required := []string{
// // 		platform,
// // 		filepath.Join(buildTools, "aapt2"),
// // 		filepath.Join(buildTools, "d8"),
// // 		filepath.Join(buildTools, "zipalign"),
// // 		filepath.Join(buildTools, "apksigner"),
// // 	}

// // 	if runtime.GOOS == "windows" {
// // 		required = []string{
// // 			platform,
// // 			filepath.Join(buildTools, "aapt2.exe"),
// // 			filepath.Join(buildTools, "d8.bat"),
// // 			filepath.Join(buildTools, "zipalign.exe"),
// // 			filepath.Join(buildTools, "apksigner.bat"),
// // 		}
// // 	}

// // 	for _, p := range required {
// // 		if _, err := os.Stat(p); err != nil {
// // 			return fmt.Errorf("missing SDK tool: %s", p)
// // 		}
// // 		fmt.Println("OK:", p)
// // 	}
// // 	return nil
// // }

// // func generateProject(project string, cfg *Config) error {
// // 	fmt.Println()
// // 	fmt.Println("Generating:", project)

// // 	if err := os.RemoveAll(project); err != nil {
// // 		return err
// // 	}

// // 	dirs := []string{
// // 		filepath.Join(project, "app", "src", "main", "java", "com", "example", "generated"),
// // 		filepath.Join(project, "app", "src", "main", "res", "values"),
// // 	}
// // 	for _, d := range dirs {
// // 		if err := os.MkdirAll(d, 0755); err != nil {
// // 			return err
// // 		}
// // 	}

// // 	write := func(name, content string) error {
// // 		return os.WriteFile(filepath.Join(project, name), []byte(content), 0644)
// // 	}

// // 	settings := `pluginManagement { repositories { google(); mavenCentral(); gradlePluginPortal() } }
// // dependencyResolutionManagement { repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS); repositories { google(); mavenCentral() } }
// // rootProject.name = "HelloApk"
// // include(":app")
// // `
// // 	if err := write("settings.gradle", settings); err != nil {
// // 		return err
// // 	}

// // 	if err := write("local.properties", fmt.Sprintf("sdk.dir=%s\n", filepath.ToSlash(cfg.AndroidSDK))); err != nil {
// // 		return err
// // 	}

// // 	rootBuild := `plugins {
// //     id 'com.android.application' version '9.2.0' apply false
// // }
// // `
// // 	if err := write("build.gradle", rootBuild); err != nil {
// // 		return err
// // 	}

// // 	appBuild := fmt.Sprintf(`plugins {
// //     id 'com.android.application'
// // }

// // android {
// //     namespace 'com.example.generated'
// //     compileSdk %s

// //     defaultConfig {
// //         applicationId 'com.example.generated'
// //         minSdk 24
// //         targetSdk 34
// //         versionCode 1
// //         versionName '1.0'
// //     }
// // }
// // `, normalizeCompileSDK(cfg.CompileSDK))
// // 	if err := write(filepath.Join("app", "build.gradle"), appBuild); err != nil {
// // 		return err
// // 	}

// // 	manifest := `<?xml version="1.0" encoding="utf-8"?>
// // <manifest xmlns:android="http://schemas.android.com/apk/res/android">
// //     <application android:theme="@style/AppTheme" android:label="Hello APK">
// //         <activity android:name=".MainActivity" android:exported="true">
// //             <intent-filter>
// //                 <action android:name="android.intent.action.MAIN"/>
// //                 <category android:name="android.intent.category.LAUNCHER"/>
// //             </intent-filter>
// //         </activity>
// //     </application>
// // </manifest>
// // `
// // 	if err := write(filepath.Join("app", "src", "main", "AndroidManifest.xml"), manifest); err != nil {
// // 		return err
// // 	}

// // 	activity := `package com.example.generated;

// // import android.app.Activity;
// // import android.os.Bundle;
// // import android.graphics.Color;
// // import android.widget.TextView;

// // public class MainActivity extends Activity {
// //     @Override
// //     protected void onCreate(Bundle state) {
// //         super.onCreate(state);
// //         TextView view = new TextView(this);
// //         view.setText("Hello from a real generated APK");
// //         view.setTextSize(24);
// //         view.setTextColor(Color.WHITE);
// //         view.setGravity(17);
// //         view.setBackgroundColor(Color.rgb(25, 25, 25));
// //         setContentView(view);
// //     }
// // }
// // `
// // 	return write(filepath.Join("app", "src", "main", "java", "com", "example", "generated", "MainActivity.java"), activity)
// // }

// // func runGradle(project, gradle, androidSDK string) error {
// // 	cmd := exec.Command(gradle, "assembleDebug")
// // 	cmd.Dir = project
// // 	cmd.Env = append(os.Environ(),
// // 		"ANDROID_HOME="+androidSDK,
// // 		"ANDROID_SDK_ROOT="+androidSDK,
// // 	)
// // 	cmd.Stdout = os.Stdout
// // 	cmd.Stderr = os.Stderr
// // 	return cmd.Run()
// // }

// // func commandOutput(program string, args ...string) (string, error) {
// // 	cmd := exec.Command(program, args...)
// // 	out, err := cmd.CombinedOutput()
// // 	return string(out), err
// // }

// // func fail(err error) {
// // 	fmt.Fprintln(os.Stderr, "\nBUILD/TOOLCHAIN TEST FAILED:")
// // 	fmt.Fprintln(os.Stderr, err)
// // 	os.Exit(1)
// // }
// package main

// import (
// 	"encoding/json"
// 	"flag"
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"runtime"
// 	"strings"
// )

// type Config struct {
// 	Java       string `json:"java"`
// 	Gradle     string `json:"gradle"`
// 	AndroidSDK string `json:"androidSdk"`
// 	BuildTools string `json:"buildTools"`
// 	CompileSDK string `json:"compileSdk"`
// }

// func main() {
// 	build := flag.Bool("build", false, "generate and build a real sample APK")
// 	flag.Parse()

// 	cfg, err := loadConfig()
// 	if err != nil {
// 		fail(err)
// 	}

// 	fmt.Println("=== ApkGenerator VS Code Toolchain Test ===")
// 	fmt.Println("OS:", runtime.GOOS, runtime.GOARCH)

// 	if err := checkConfig(&cfg); err != nil {
// 		fail(err)
// 	}

// 	if err := checkJava(cfg.Java); err != nil {
// 		fail(err)
// 	}
// 	if err := checkGradle(cfg.Gradle); err != nil {
// 		fail(err)
// 	}
// 	if err := checkSDK(&cfg); err != nil {
// 		fail(err)
// 	}

// 	fmt.Println()
// 	fmt.Println("TOOLCHAIN CHECK: PASS")

// 	if !*build {
// 		fmt.Println()
// 		fmt.Println("Run with -build to generate and build a real APK.")
// 		return
// 	}

// 	project := filepath.Join("generated", "HelloApk")
// 	if err := generateProject(project, &cfg); err != nil {
// 		fail(err)
// 	}

// 	fmt.Println()
// 	fmt.Println("=== REAL BUILD ===")
// 	if err := runGradle(project, cfg.Gradle); err != nil {
// 		fail(err)
// 	}

// 	apk := filepath.Join(project, "app", "build", "outputs", "apk", "debug", "app-debug.apk")
// 	info, err := os.Stat(apk)
// 	if err != nil || info.Size() == 0 {
// 		fail(fmt.Errorf("build finished without a real APK: %s", apk))
// 	}

// 	fmt.Println()
// 	fmt.Println("BUILD RESULT: SUCCESS")
// 	fmt.Println("APK:", apk)
// 	fmt.Printf("APK size: %d bytes\n", info.Size())
// }

// func loadConfig() (Config, error) {
// 	data, err := os.ReadFile("config.json")
// 	if err != nil {
// 		return Config{}, err
// 	}
// 	var cfg Config
// 	if err := json.Unmarshal(data, &cfg); err != nil {
// 		return Config{}, err
// 	}
// 	return cfg, nil
// }

// func checkConfig(cfg *Config) error {
// 	if strings.TrimSpace(cfg.Java) == "" ||
// 		strings.TrimSpace(cfg.Gradle) == "" ||
// 		strings.TrimSpace(cfg.AndroidSDK) == "" {
// 		return fmt.Errorf("edit config.json and set java, gradle and androidSdk paths")
// 	}
// 	return nil
// }

// func checkJava(java string) error {
// 	fmt.Println()
// 	fmt.Println("[1/3] Java")
// 	out, err := commandOutput(java, "-version")
// 	if err != nil {
// 		return fmt.Errorf("Java failed: %w\n%s", err, out)
// 	}
// 	fmt.Print(out)
// 	return nil
// }

// func checkGradle(gradle string) error {
// 	fmt.Println()
// 	fmt.Println("[2/3] Gradle")
// 	out, err := commandOutput(gradle, "--version")
// 	if err != nil {
// 		return fmt.Errorf("Gradle failed: %w\n%s", err, out)
// 	}
// 	fmt.Print(out)
// 	return nil
// }

// func checkSDK(cfg *Config) error {
// 	fmt.Println()
// 	fmt.Println("[3/3] Android SDK")

// 	platform := filepath.Join(cfg.AndroidSDK, "platforms", "android-"+cfg.CompileSDK, "android.jar")
// 	buildTools := filepath.Join(cfg.AndroidSDK, "build-tools", cfg.BuildTools)

// 	required := []string{
// 		platform,
// 		filepath.Join(buildTools, "aapt2"),
// 		filepath.Join(buildTools, "d8"),
// 		filepath.Join(buildTools, "zipalign"),
// 		filepath.Join(buildTools, "apksigner"),
// 	}

// 	if runtime.GOOS == "windows" {
// 		required = []string{
// 			platform,
// 			filepath.Join(buildTools, "aapt2.exe"),
// 			filepath.Join(buildTools, "d8.bat"),
// 			filepath.Join(buildTools, "zipalign.exe"),
// 			filepath.Join(buildTools, "apksigner.bat"),
// 		}
// 	}

// 	for _, p := range required {
// 		if _, err := os.Stat(p); err != nil {
// 			return fmt.Errorf("missing SDK tool: %s", p)
// 		}
// 		fmt.Println("OK:", p)
// 	}
// 	return nil
// }

// func generateProject(project string, cfg *Config) error {
// 	fmt.Println()
// 	fmt.Println("Generating:", project)

// 	if err := os.RemoveAll(project); err != nil {
// 		return err
// 	}

// 	dirs := []string{
// 		filepath.Join(project, "app", "src", "main", "java", "com", "example", "generated"),
// 		filepath.Join(project, "app", "src", "main", "res", "values"),
// 	}
// 	for _, d := range dirs {
// 		if err := os.MkdirAll(d, 0755); err != nil {
// 			return err
// 		}
// 	}

// 	write := func(name, content string) error {
// 		return os.WriteFile(filepath.Join(project, name), []byte(content), 0644)
// 	}

// 	settings := `pluginManagement { repositories { google(); mavenCentral(); gradlePluginPortal() } }
// dependencyResolutionManagement { repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS); repositories { google(); mavenCentral() } }
// rootProject.name = "HelloApk"
// include(":app")
// `
// 	if err := write("settings.gradle", settings); err != nil {
// 		return err
// 	}

// 	rootBuild := `plugins {
//     id 'com.android.application' version '9.2.0' apply false
// }
// `
// 	if err := write("build.gradle", rootBuild); err != nil {
// 		return err
// 	}

// 	appBuild := fmt.Sprintf(`plugins {
//     id 'com.android.application'
// }

// android {
//     namespace 'com.example.generated'
//     compileSdk = %s

//     defaultConfig {
//         applicationId 'com.example.generated'
//         minSdk = 24
//         targetSdk = 34
//         versionCode = 1
//         versionName = '1.0'
//     }

//     compileOptions {
//         sourceCompatibility = JavaVersion.VERSION_17
//         targetCompatibility = JavaVersion.VERSION_17
//     }
// }

// java {
//     toolchain {
//         languageVersion = JavaLanguageVersion.of(17)
//     }
// }
// `, cfg.CompileSDK)
// 	if err := write(filepath.Join("app", "build.gradle"), appBuild); err != nil {
// 		return err
// 	}

// 	manifest := `<?xml version="1.0" encoding="utf-8"?>
// <manifest xmlns:android="http://schemas.android.com/apk/res/android">
//     <application android:theme="@style/AppTheme" android:label="Hello APK">
//         <activity android:name=".MainActivity" android:exported="true">
//             <intent-filter>
//                 <action android:name="android.intent.action.MAIN"/>
//                 <category android:name="android.intent.category.LAUNCHER"/>
//             </intent-filter>
//         </activity>
//     </application>
// </manifest>
// `
// 	if err := write(filepath.Join("app", "src", "main", "AndroidManifest.xml"), manifest); err != nil {
// 		return err
// 	}

// 	activity := `package com.example.generated;

// import android.app.Activity;
// import android.os.Bundle;
// import android.graphics.Color;
// import android.widget.TextView;

// public class MainActivity extends Activity {
//     @Override
//     protected void onCreate(Bundle state) {
//         super.onCreate(state);
//         TextView view = new TextView(this);
//         view.setText("Hello from a real generated APK");
//         view.setTextSize(24);
//         view.setTextColor(Color.WHITE);
//         view.setGravity(17);
//         view.setBackgroundColor(Color.rgb(25, 25, 25));
//         setContentView(view);
//     }
// }
// `
// 	return write(filepath.Join("app", "src", "main", "java", "com", "example", "generated", "MainActivity.java"), activity)
// }

// func runGradle(project, gradle string) error {
// 	args := []string{"assembleDebug"}
// 	cmd := exec.Command(gradle, args...)
// 	cmd.Dir = project
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	return cmd.Run()
// }

// func commandOutput(program string, args ...string) (string, error) {
// 	cmd := exec.Command(program, args...)
// 	out, err := cmd.CombinedOutput()
// 	return string(out), err
// }

// func fail(err error) {
// 	fmt.Fprintln(os.Stderr, "\nBUILD/TOOLCHAIN TEST FAILED:")
// 	fmt.Fprintln(os.Stderr, err)
// 	os.Exit(1)
// }

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	Java       string `json:"java"`
	Gradle     string `json:"gradle"`
	AndroidSDK string `json:"androidSdk"`
	BuildTools string `json:"buildTools"`
	CompileSDK string `json:"compileSdk"`
}

func main() {
	build := flag.Bool("build", false, "generate and build a real sample APK")
	flag.Parse()

	cfg, err := loadConfig()
	if err != nil {
		fail(err)
	}

	fmt.Println("=== ApkGenerator VS Code Toolchain Test ===")
	fmt.Println("OS:", runtime.GOOS, runtime.GOARCH)

	if err := checkConfig(&cfg); err != nil {
		fail(err)
	}

	if err := checkJava(cfg.Java); err != nil {
		fail(err)
	}

	if err := checkGradle(cfg.Gradle, cfg.Java); err != nil {
		fail(err)
	}

	if err := checkSDK(&cfg); err != nil {
		fail(err)
	}

	fmt.Println()
	fmt.Println("TOOLCHAIN CHECK: PASS")

	if !*build {
		fmt.Println()
		fmt.Println("Run with -build to generate and build a real APK.")
		return
	}

	project := filepath.Join("generated", "HelloApk")

	if err := generateProject(project, &cfg); err != nil {
		fail(err)
	}

	fmt.Println()
	fmt.Println("=== REAL BUILD ===")

	if err := runGradle(project, cfg.Gradle, cfg.Java, cfg.AndroidSDK); err != nil {
		fail(err)
	}

	apk := filepath.Join(
		project,
		"app",
		"build",
		"outputs",
		"apk",
		"debug",
		"app-debug.apk",
	)

	info, err := os.Stat(apk)
	if err != nil || info.Size() == 0 {
		fail(fmt.Errorf("build finished without a real APK: %s", apk))
	}

	fmt.Println()
	fmt.Println("BUILD RESULT: SUCCESS")
	fmt.Println("APK:", apk)
	fmt.Printf("APK size: %d bytes\n", info.Size())
}

func loadConfig() (Config, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func checkConfig(cfg *Config) error {
	if strings.TrimSpace(cfg.Java) == "" ||
		strings.TrimSpace(cfg.Gradle) == "" ||
		strings.TrimSpace(cfg.AndroidSDK) == "" {
		return fmt.Errorf(
			"edit config.json and set java, gradle and androidSdk paths",
		)
	}

	if strings.TrimSpace(cfg.BuildTools) == "" {
		return fmt.Errorf("config.json is missing buildTools")
	}

	if strings.TrimSpace(cfg.CompileSDK) == "" {
		return fmt.Errorf("config.json is missing compileSdk")
	}

	return nil
}

func checkJava(java string) error {
	fmt.Println()
	fmt.Println("[1/3] Java")

	out, err := commandOutput(java, "-version")
	if err != nil {
		return fmt.Errorf("Java failed: %w\n%s", err, out)
	}

	fmt.Print(out)

	return nil
}

func checkGradle(gradle, java string) error {
	fmt.Println()
	fmt.Println("[2/3] Gradle")

	javaHome := javaHomeFromJava(java)

	out, err := commandOutputWithEnv(
		javaHome,
		gradle,
		"--version",
	)

	if err != nil {
		return fmt.Errorf("Gradle failed: %w\n%s", err, out)
	}

	fmt.Print(out)

	return nil
}

func checkSDK(cfg *Config) error {
	fmt.Println()
	fmt.Println("[3/3] Android SDK")

	platform, err := findCompileSDK(cfg.AndroidSDK, cfg.CompileSDK)
	if err != nil {
		return err
	}

	buildTools := filepath.Join(
		cfg.AndroidSDK,
		"build-tools",
		cfg.BuildTools,
	)

	required := []string{
		platform,
	}

	if runtime.GOOS == "windows" {
		required = append(required,
			filepath.Join(buildTools, "aapt2.exe"),
			filepath.Join(buildTools, "d8.bat"),
			filepath.Join(buildTools, "zipalign.exe"),
			filepath.Join(buildTools, "apksigner.bat"),
		)
	} else {
		required = append(required,
			filepath.Join(buildTools, "aapt2"),
			filepath.Join(buildTools, "d8"),
			filepath.Join(buildTools, "zipalign"),
			filepath.Join(buildTools, "apksigner"),
		)
	}

	for _, p := range required {
		if _, err := os.Stat(p); err != nil {
			return fmt.Errorf("missing SDK tool: %s", p)
		}

		fmt.Println("OK:", p)
	}

	return nil
}

func findCompileSDK(sdkRoot, compileSDK string) (string, error) {
	platformsDir := filepath.Join(
		sdkRoot,
		"platforms",
	)

	entries, err := os.ReadDir(platformsDir)
	if err != nil {
		return "", fmt.Errorf(
			"cannot read Android SDK platforms directory: %w",
			err,
		)
	}

	// First try exact directory:
	// android-37
	exact := filepath.Join(
		platformsDir,
		"android-"+compileSDK,
		"android.jar",
	)

	if fileExists(exact) {
		return exact, nil
	}

	// Then allow:
	// android-37.0
	// android-37.1
	// android-37.2
	prefix := "android-" + compileSDK + "."

	var candidates []string

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		if strings.HasPrefix(name, prefix) {
			jar := filepath.Join(
				platformsDir,
				name,
				"android.jar",
			)

			if fileExists(jar) {
				candidates = append(candidates, jar)
			}
		}
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf(
			"Android platform android-%s not found under %s",
			compileSDK,
			platformsDir,
		)
	}

	// Sort lexicographically and use the highest matching version.
	sortStrings(candidates)

	return candidates[len(candidates)-1], nil
}

func generateProject(project string, cfg *Config) error {
	fmt.Println()
	fmt.Println("Generating:", project)

	if err := os.RemoveAll(project); err != nil {
		return err
	}

	dirs := []string{
		filepath.Join(
			project,
			"app",
			"src",
			"main",
			"java",
			"com",
			"example",
			"generated",
		),
		filepath.Join(
			project,
			"app",
			"src",
			"main",
			"res",
			"values",
		),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	write := func(name, content string) error {
		return os.WriteFile(
			filepath.Join(project, name),
			[]byte(content),
			0644,
		)
	}

	settings := `pluginManagement {
    repositories {
        google()
        mavenCentral()
        gradlePluginPortal()
    }
}

dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)

    repositories {
        google()
        mavenCentral()
    }
}

rootProject.name = "HelloApk"
include(":app")
`

	if err := write("settings.gradle", settings); err != nil {
		return err
	}

	rootBuild := `plugins {
    id 'com.android.application' version '9.2.0' apply false
}
`

	if err := write("build.gradle", rootBuild); err != nil {
		return err
	}

	appBuild := fmt.Sprintf(`plugins {
    id 'com.android.application'
}

android {
    namespace 'com.example.generated'
    compileSdk = %s

    defaultConfig {
        applicationId 'com.example.generated'
        minSdk = 24
        targetSdk = 34
        versionCode = 1
        versionName = '1.0'
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }
}

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(17)
    }
}
`, cfg.CompileSDK)

	if err := write(
		filepath.Join("app", "build.gradle"),
		appBuild,
	); err != nil {
		return err
	}

	/*
		IMPORTANT:
		Gradle needs the Android SDK location.

		Using forward slashes also works on Windows and avoids
		escaping problems inside local.properties.
	*/
	sdkDir := filepath.ToSlash(cfg.AndroidSDK)

	localProperties := fmt.Sprintf(
		"sdk.dir=%s\n",
		sdkDir,
	)

	if err := write(
		"local.properties",
		localProperties,
	); err != nil {
		return err
	}

	manifest := `<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android">

    <application
        android:theme="@style/AppTheme"
        android:label="Hello APK">

        <activity
            android:name=".MainActivity"
            android:exported="true">

            <intent-filter>
                <action android:name="android.intent.action.MAIN" />
                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>

        </activity>

    </application>

</manifest>
`

	if err := write(
		filepath.Join(
			"app",
			"src",
			"main",
			"AndroidManifest.xml",
		),
		manifest,
	); err != nil {
		return err
	}

	styles := `<?xml version="1.0" encoding="utf-8"?>
<resources>

    <style name="AppTheme"
        parent="@android:style/Theme.Material.NoActionBar">
        <item name="android:fontFamily">sans</item>
        <item name="android:colorAccent">#FFFFFF</item>
    </style>

</resources>
`

	if err := write(
		filepath.Join(
			"app",
			"src",
			"main",
			"res",
			"values",
			"styles.xml",
		),
		styles,
	); err != nil {
		return err
	}

	activity := `package com.example.generated;

import android.app.Activity;
import android.os.Bundle;
import android.graphics.Color;
import android.widget.TextView;

public class MainActivity extends Activity {

    @Override
    protected void onCreate(Bundle state) {
        super.onCreate(state);

        TextView view = new TextView(this);

        view.setText("Hello from a real generated APK");
        view.setTextSize(24);
        view.setTextColor(Color.WHITE);
        view.setGravity(17);
        view.setBackgroundColor(Color.rgb(25, 25, 25));

        setContentView(view);
    }
}
`

	// return write(
	// 	filepath.Join(
	// 		project,
	// 		"app",
	// 		"src",
	// 		"main",
	// 		"java",
	// 		"com",
	// 		"example",
	// 		"generated",
	// 		"MainActivity.java",
	// 	),
	// 	activity,
	// )
	return write(
	filepath.Join(
		"app",
		"src",
		"main",
		"java",
		"com",
		"example",
		"generated",
		"MainActivity.java",
	),
	activity,
)
}

func runGradle(
	project string,
	gradle string,
	java string,
	androidSDK string,
) error {

	args := []string{
		"assembleDebug",
	}

	cmd := exec.Command(
		gradle,
		args...,
	)

	cmd.Dir = project

	/*
		Force Gradle itself to run using JDK 17.

		This prevents Gradle from using the system JDK 23.
	*/
	javaHome := javaHomeFromJava(java)

	env := os.Environ()

	env = append(
		env,
		"JAVA_HOME="+javaHome,
		"ANDROID_HOME="+androidSDK,
		"ANDROID_SDK_ROOT="+androidSDK,
	)

	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func commandOutput(program string, args ...string) (string, error) {
	cmd := exec.Command(program, args...)

	out, err := cmd.CombinedOutput()

	return string(out), err
}

func commandOutputWithEnv(
	javaHome string,
	program string,
	args ...string,
) (string, error) {

	cmd := exec.Command(program, args...)

	env := os.Environ()

	env = append(
		env,
		"JAVA_HOME="+javaHome,
	)

	cmd.Env = env

	out, err := cmd.CombinedOutput()

	return string(out), err
}

func javaHomeFromJava(java string) string {
	absolute, err := filepath.Abs(java)
	if err != nil {
		return filepath.Dir(filepath.Dir(java))
	}

	// java.exe:
	// C:\...\jdk-17\bin\java.exe
	//
	// javaHome:
	// C:\...\jdk-17
	return filepath.Dir(filepath.Dir(absolute))
}

func fileExists(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		return false
	}

	return !info.IsDir()
}

func sortStrings(values []string) {
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if values[j] < values[i] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}
}

func fail(err error) {
	fmt.Fprintln(
		os.Stderr,
		"\nBUILD/TOOLCHAIN TEST FAILED:",
	)

	fmt.Fprintln(
		os.Stderr,
		err,
	)

	os.Exit(1)
}