# Phase 3: Updater Backend (Go)

> Parent: [plan.md](./plan.md)

## Overview

| Field | Value |
|-------|-------|
| Date | 2026-01-08 |
| Priority | P1 |
| Effort | 2h |
| Status | pending |

## Requirements

1. Query GitHub Releases API for latest version
2. Compare versions using semver
3. Download NSIS installer to temp directory
4. Execute installer and quit app
5. Handle errors gracefully

## Implementation Steps

### 3.1 Create Updater Package

Create `internal/updater/updater.go`:

```go
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
)

const (
    GitHubOwner     = "vthuan1889"
    GitHubRepo      = "vibe-imageborder"
    InstallerSuffix = "-amd64-installer.exe"
)

type ReleaseInfo struct {
    TagName string  `json:"tag_name"`
    Assets  []Asset `json:"assets"`
}

type Asset struct {
    Name               string `json:"name"`
    BrowserDownloadURL string `json:"browser_download_url"`
}

type UpdateInfo struct {
    Available   bool   `json:"available"`
    Current     string `json:"current"`
    Latest      string `json:"latest"`
    DownloadURL string `json:"downloadUrl"`
}

// CheckUpdate checks GitHub for newer version
func CheckUpdate(currentVersion string) (*UpdateInfo, error) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest",
        GitHubOwner, GitHubRepo)

    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to check update: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
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

    info.Available = compareVersions(currentVersion, release.TagName) < 0
    return info, nil
}

// DownloadAndInstall downloads installer and runs it
func DownloadAndInstall(downloadURL string) error {
    tempDir := os.TempDir()
    installerPath := filepath.Join(tempDir, "ImageBorderTool-update.exe")

    // Download
    resp, err := http.Get(downloadURL)
    if err != nil {
        return fmt.Errorf("download failed: %w", err)
    }
    defer resp.Body.Close()

    out, err := os.Create(installerPath)
    if err != nil {
        return fmt.Errorf("cannot create temp file: %w", err)
    }
    defer out.Close()

    if _, err := io.Copy(out, resp.Body); err != nil {
        return fmt.Errorf("download incomplete: %w", err)
    }
    out.Close()

    // Run installer (detached)
    cmd := exec.Command(installerPath)
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("cannot start installer: %w", err)
    }

    return nil
}

// compareVersions compares v1 and v2
// Returns: -1 if v1 < v2, 0 if equal, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
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
```

### 3.2 Add Methods to App

Add to `app.go`:

```go
import "vibe-imageborder/internal/updater"

// CheckForUpdate queries GitHub for updates
func (a *App) CheckForUpdate() (*updater.UpdateInfo, error) {
    return updater.CheckUpdate(version)
}

// DownloadAndInstallUpdate downloads and runs installer
func (a *App) DownloadAndInstallUpdate(downloadURL string) error {
    if err := updater.DownloadAndInstall(downloadURL); err != nil {
        return err
    }
    // Quit app after starting installer
    runtime.Quit(a.ctx)
    return nil
}
```

## Files to Create/Modify

| File | Change |
|------|--------|
| `internal/updater/updater.go` | New - update logic |
| `app.go` | Add CheckForUpdate, DownloadAndInstallUpdate |

## API Endpoints Used

- `GET https://api.github.com/repos/{owner}/{repo}/releases/latest`
- No auth required for public repos

## Success Criteria

- [ ] CheckUpdate returns correct version info
- [ ] Semver comparison handles all cases
- [ ] Installer downloads to temp folder
- [ ] Installer starts successfully
- [ ] App quits after starting installer
