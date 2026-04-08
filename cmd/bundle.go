package cmd

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"canopy/pkg/registry"
)

// RunBundle handles "canopy bundle <app-path> [flags]".
func RunBundle(args []string) {
	fs := flag.NewFlagSet("bundle", flag.ExitOnError)
	output := fs.String("output", "", "Output path for the .app bundle (default: ./<AppName>.app)")
	name := fs.String("name", "", "Override app name (default: from canopy.json or directory name)")
	icon := fs.String("icon", "", "Path to .icns file for the app icon")
	bundleID := fs.String("bundle-id", "", "Override bundle identifier (default: com.canopy.app.<name>)")
	sign := fs.Bool("sign", false, "Codesign the bundle with hardened runtime")
	identity := fs.String("identity", "", "Signing identity (default: auto-detect Developer ID Application)")
	notarize := fs.Bool("notarize", false, "Notarize the bundle with Apple (implies --sign)")
	appleID := fs.String("apple-id", "", "Apple ID email for notarization")
	teamID := fs.String("team-id", "", "Team ID for notarization")
	password := fs.String("password", "", "App-specific password for notarization")
	keychainProfile := fs.String("keychain-profile", "", "Keychain profile for notarization (alternative to apple-id/team-id/password)")
	fs.StringVar(output, "o", "", "Output path for the .app bundle (short)")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: canopy bundle <app-path> [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Create a standalone macOS .app bundle from a Canopy app.\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  <app-path>    App directory, or owner/repo for an installed package\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
		os.Exit(1)
	}

	appPath := resolveAppPath(fs.Arg(0))

	info, err := os.Stat(appPath)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "error: %s is not a directory\n", appPath)
		os.Exit(1)
	}

	manifest := readManifestOrDefaults(appPath)

	// Apply flag overrides
	if *name != "" {
		manifest.Name = *name
	}
	if *bundleID != "" {
		manifest.BundleID = *bundleID
	}
	if *icon != "" {
		manifest.Icon = *icon
	}

	// Derive defaults
	if manifest.Name == "" {
		manifest.Name = titleCase(filepath.Base(appPath))
	}
	if manifest.Version == "" {
		manifest.Version = "1.0.0"
	}
	if manifest.BundleID == "" {
		manifest.BundleID = "com.canopy.app." + sanitizeID(manifest.Name)
	}

	// Determine output path
	outPath := *output
	if outPath == "" {
		outPath = manifest.Name + ".app"
	}
	outPath, _ = filepath.Abs(outPath)

	// Remove existing bundle
	os.RemoveAll(outPath)

	// Create bundle structure
	macosDir := filepath.Join(outPath, "Contents", "MacOS")
	resDir := filepath.Join(outPath, "Contents", "Resources")
	appResDir := filepath.Join(resDir, "app")

	if err := os.MkdirAll(macosDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(appResDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Copy canopy binary
	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot find own executable: %v\n", err)
		os.Exit(1)
	}
	exePath, _ = filepath.EvalSymlinks(exePath)
	if err := copyFilePerm(exePath, filepath.Join(macosDir, "canopy"), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error copying binary: %v\n", err)
		os.Exit(1)
	}

	// Copy app files
	if err := copyDir(appPath, appResDir); err != nil {
		fmt.Fprintf(os.Stderr, "error copying app files: %v\n", err)
		os.Exit(1)
	}

	// Copy icon if provided
	hasIcon := false
	iconPath := resolveIconPath(manifest.Icon, appPath)
	if iconPath != "" {
		if err := copyFilePerm(iconPath, filepath.Join(resDir, "AppIcon.icns"), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not copy icon: %v\n", err)
		} else {
			hasIcon = true
		}
	}

	// Generate Info.plist
	if err := writeInfoPlist(filepath.Join(outPath, "Contents", "Info.plist"), manifest, hasIcon); err != nil {
		fmt.Fprintf(os.Stderr, "error writing Info.plist: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Built %s\n", outPath)

	// Notarize implies sign
	if *notarize {
		*sign = true
	}

	// Codesign
	if *sign {
		fmt.Printf("Signing %s...\n", outPath)
		if err := signBundle(outPath, *identity); err != nil {
			fmt.Fprintf(os.Stderr, "error signing: %v\n", err)
			os.Exit(1)
		}
	}

	// Notarize
	if *notarize {
		creds, err := resolveNotarizeCredentials(*appleID, *teamID, *password, *keychainProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Notarizing %s (this may take several minutes)...\n", outPath)
		if err := notarizeBundle(outPath, creds); err != nil {
			fmt.Fprintf(os.Stderr, "error notarizing: %v\n", err)
			os.Exit(1)
		}
	}
}

// resolveAppPath resolves an app path argument. If it looks like owner/repo,
// check the installed packages directory.
func resolveAppPath(arg string) string {
	// Check if it looks like owner/repo (exactly one slash, no path separators)
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`, arg); matched {
		// Try installed package
		reg, err := registry.New()
		if err == nil {
			if entry := reg.Get(arg); entry != nil && entry.Path != "" {
				return entry.Path
			}
		}
		// Try direct path under ~/.canopy/apps/
		home, err := os.UserHomeDir()
		if err == nil {
			pkgPath := filepath.Join(home, ".canopy", "apps", arg)
			if info, err := os.Stat(pkgPath); err == nil && info.IsDir() {
				return pkgPath
			}
		}
	}
	// Use as-is (relative or absolute path)
	abs, err := filepath.Abs(arg)
	if err != nil {
		return arg
	}
	return abs
}

func readManifestOrDefaults(appPath string) registry.Manifest {
	data, err := os.ReadFile(filepath.Join(appPath, "canopy.json"))
	if err != nil {
		return registry.Manifest{}
	}
	var m registry.Manifest
	if json.Unmarshal(data, &m) != nil {
		return registry.Manifest{}
	}
	return m
}

func resolveIconPath(iconField, appDir string) string {
	if iconField == "" {
		return ""
	}
	// If it's an absolute path to an .icns file
	if filepath.IsAbs(iconField) && strings.HasSuffix(iconField, ".icns") {
		if _, err := os.Stat(iconField); err == nil {
			return iconField
		}
	}
	// If it's a relative path to an .icns file
	if strings.HasSuffix(iconField, ".icns") {
		abs := filepath.Join(appDir, iconField)
		if _, err := os.Stat(abs); err == nil {
			return abs
		}
	}
	// SF Symbol names or other non-file values — can't embed
	return ""
}

func sanitizeID(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return strings.ToLower(re.ReplaceAllString(name, "-"))
}

func titleCase(name string) string {
	label := strings.ReplaceAll(name, "_", " ")
	label = strings.ReplaceAll(label, "-", " ")
	words := strings.Fields(label)
	for j, w := range words {
		if len(w) > 0 {
			words[j] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func copyFilePerm(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFilePerm(path, target, info.Mode())
	})
}

type plistData struct {
	Name     string
	BundleID string
	Version  string
	HasIcon  bool
}

func writeInfoPlist(path string, m registry.Manifest, hasIcon bool) error {
	tmpl, err := template.New("plist").Parse(infoPlistTmpl)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, plistData{
		Name:     m.Name,
		BundleID: m.BundleID,
		Version:  m.Version,
		HasIcon:  hasIcon,
	})
}

// signBundle codesigns the .app with hardened runtime and embedded entitlements.
func signBundle(appPath, identity string) error {
	if _, err := exec.LookPath("codesign"); err != nil {
		return fmt.Errorf("codesign not found; install Xcode Command Line Tools: xcode-select --install")
	}

	if identity == "" {
		found, err := findSigningIdentity()
		if err != nil {
			return err
		}
		identity = found
	}

	// Write embedded entitlements to temp file
	tmpFile, err := os.CreateTemp("", "canopy-entitlements-*.plist")
	if err != nil {
		return fmt.Errorf("create temp entitlements: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write(entitlementsData); err != nil {
		tmpFile.Close()
		return fmt.Errorf("write temp entitlements: %w", err)
	}
	tmpFile.Close()

	// Sign with hardened runtime
	cmd := exec.Command("codesign",
		"--sign", identity,
		"--options", "runtime",
		"--entitlements", tmpFile.Name(),
		"--timestamp",
		"--force",
		"--deep",
		appPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("codesign failed: %s\n%s", err, out)
	}

	// Verify
	cmd = exec.Command("codesign", "--verify", "--verbose=2", appPath)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("codesign verification failed: %s\n%s", err, out)
	}

	fmt.Printf("Signed with: %s\n", identity)
	return nil
}

// findSigningIdentity searches the keychain for a Developer ID Application identity.
func findSigningIdentity() (string, error) {
	cmd := exec.Command("security", "find-identity", "-v", "-p", "codesigning")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to list signing identities: %w", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		// Lines look like: 1) ABC123 "Developer ID Application: Name (TEAMID)"
		idx := strings.Index(line, `"Developer ID Application:`)
		if idx < 0 {
			continue
		}
		end := strings.LastIndex(line, `"`)
		if end > idx {
			return line[idx+1 : end], nil
		}
	}

	return "", fmt.Errorf("no Developer ID Application identity found in keychain; specify --identity or install a signing certificate")
}

type notarizeCredentials struct {
	keychainProfile string
	appleID         string
	teamID          string
	password        string
}

func resolveNotarizeCredentials(flagAppleID, flagTeamID, flagPassword, flagProfile string) (notarizeCredentials, error) {
	// Keychain profile takes priority
	profile := flagProfile
	if profile == "" {
		profile = os.Getenv("CANOPY_KEYCHAIN_PROFILE")
	}
	if profile != "" {
		return notarizeCredentials{keychainProfile: profile}, nil
	}

	// Collect individual credentials from flags then env
	appleID := flagAppleID
	if appleID == "" {
		appleID = os.Getenv("CANOPY_APPLE_ID")
	}
	teamID := flagTeamID
	if teamID == "" {
		teamID = os.Getenv("CANOPY_TEAM_ID")
	}
	pwd := flagPassword
	if pwd == "" {
		pwd = os.Getenv("CANOPY_APP_PASSWORD")
	}

	var missing []string
	if appleID == "" {
		missing = append(missing, "apple-id (or CANOPY_APPLE_ID)")
	}
	if teamID == "" {
		missing = append(missing, "team-id (or CANOPY_TEAM_ID)")
	}
	if pwd == "" {
		missing = append(missing, "password (or CANOPY_APP_PASSWORD)")
	}
	if len(missing) > 0 {
		return notarizeCredentials{}, fmt.Errorf("notarization requires credentials; missing: %s\n  use --keychain-profile as an alternative", strings.Join(missing, ", "))
	}

	return notarizeCredentials{appleID: appleID, teamID: teamID, password: pwd}, nil
}

// notarizeBundle submits the .app for Apple notarization and staples the ticket.
func notarizeBundle(appPath string, creds notarizeCredentials) error {
	if _, err := exec.LookPath("xcrun"); err != nil {
		return fmt.Errorf("xcrun not found; install Xcode Command Line Tools: xcode-select --install")
	}
	if _, err := exec.LookPath("ditto"); err != nil {
		return fmt.Errorf("ditto not found; install Xcode Command Line Tools: xcode-select --install")
	}

	// Create temp zip
	tmpDir, err := os.MkdirTemp("", "canopy-notarize-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, "bundle.zip")
	cmd := exec.Command("ditto", "-c", "-k", "--keepParent", appPath, zipPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ditto failed: %s\n%s", err, out)
	}

	// Submit for notarization
	submitArgs := []string{"notarytool", "submit", zipPath, "--wait"}
	if creds.keychainProfile != "" {
		submitArgs = append(submitArgs, "--keychain-profile", creds.keychainProfile)
	} else {
		submitArgs = append(submitArgs, "--apple-id", creds.appleID, "--team-id", creds.teamID, "--password", creds.password)
	}

	cmd = exec.Command("xcrun", submitArgs...)
	cmd.Stdout = os.Stderr // stream progress to stderr
	cmd.Stderr = os.Stderr
	var outBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stderr, &outBuf)

	if err := cmd.Run(); err != nil {
		output := outBuf.String()
		// Try to extract submission ID for user to check logs
		if id := extractSubmissionID(output); id != "" {
			return fmt.Errorf("notarization failed: %w\n  check details: xcrun notarytool log %s", err, id)
		}
		return fmt.Errorf("notarization failed: %w\n%s", err, output)
	}

	// Staple the ticket
	cmd = exec.Command("xcrun", "stapler", "staple", appPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("stapler failed: %s\n%s", err, out)
	}

	fmt.Printf("Notarized and stapled: %s\n", appPath)
	return nil
}

// extractSubmissionID tries to find a submission UUID in notarytool output.
func extractSubmissionID(output string) string {
	re := regexp.MustCompile(`id:\s*([0-9a-f-]{36})`)
	m := re.FindStringSubmatch(output)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

const infoPlistTmpl = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleName</key>
	<string>{{.Name}}</string>
	<key>CFBundleDisplayName</key>
	<string>{{.Name}}</string>
	<key>CFBundleIdentifier</key>
	<string>{{.BundleID}}</string>
	<key>CFBundleVersion</key>
	<string>{{.Version}}</string>
	<key>CFBundleShortVersionString</key>
	<string>{{.Version}}</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleExecutable</key>
	<string>canopy</string>
	{{- if .HasIcon}}
	<key>CFBundleIconFile</key>
	<string>AppIcon</string>
	{{- end}}
	<key>LSMinimumSystemVersion</key>
	<string>13.0</string>
	<key>NSHighResolutionCapable</key>
	<true/>
	<key>NSCameraUsageDescription</key>
	<string>This app needs camera access to display live camera preview and capture photos.</string>
	<key>NSMicrophoneUsageDescription</key>
	<string>This app needs microphone access to record audio.</string>
	<key>NSSupportsAutomaticTermination</key>
	<false/>
</dict>
</plist>
`
