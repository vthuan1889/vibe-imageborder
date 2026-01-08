// Package updater provides auto-update functionality via GitHub Releases.
package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	GitHubOwner     = "vthuan1889"
	GitHubRepo      = "vibe-imageborder"
	InstallerSuffix = "-amd64-installer.exe"
)

// ReleaseInfo represents GitHub release API response.
type ReleaseInfo struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

// Asset represents a release asset.
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// UpdateInfo contains update check result.
type UpdateInfo struct {
	Available   bool   `json:"available"`
	Current     string `json:"current"`
	Latest      string `json:"latest"`
	DownloadURL string `json:"downloadUrl"`
}

// CheckUpdate queries GitHub for newer version.
func CheckUpdate(currentVersion string) (*UpdateInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest",
		GitHubOwner, GitHubRepo)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to check update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	info := &UpdateInfo{
		Current: currentVersion,
		Latest:  release.TagName,
	}

	// Find installer asset
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, InstallerSuffix) {
			info.DownloadURL = asset.BrowserDownloadURL
			break
		}
	}

	info.Available = CompareVersions(currentVersion, release.TagName) < 0
	return info, nil
}

// DownloadAndInstall downloads installer and runs it.
func DownloadAndInstall(downloadURL string) error {
	tempDir := os.TempDir()
	installerPath := filepath.Join(tempDir, "ImageBorderTool-update.exe")

	// Download installer
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %d", resp.StatusCode)
	}

	out, err := os.Create(installerPath)
	if err != nil {
		return fmt.Errorf("cannot create temp file: %w", err)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		return fmt.Errorf("download incomplete: %w", err)
	}
	out.Close()

	// Run installer with elevation via PowerShell on Windows
	cmd := exec.Command("powershell.exe", "-Command", fmt.Sprintf(`Start-Process -FilePath "%s" -Verb RunAs`, installerPath))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cannot start installer: %w", err)
	}

	return nil
}

// CompareVersions compares two semantic versions.
// Returns: -1 if v1 < v2, 0 if equal, 1 if v1 > v2.
func CompareVersions(v1, v2 string) int {
	// Strip 'v' prefix
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	for i := 0; i < 3; i++ {
		var n1, n2 int
		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &n1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &n2)
		}
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}
	return 0
}
