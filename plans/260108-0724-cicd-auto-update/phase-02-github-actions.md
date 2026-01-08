# Phase 2: GitHub Actions Workflow

> Parent: [plan.md](./plan.md)

## Overview

| Field | Value |
|-------|-------|
| Date | 2026-01-08 |
| Priority | P1 |
| Effort | 2h |
| Status | pending |

## Requirements

1. Trigger on tag push matching `v*.*.*`
2. Build Windows amd64 executable
3. Create NSIS installer
4. Upload artifacts to GitHub Releases
5. Use GitHub-hosted Windows runner

## Implementation Steps

### 2.1 Create Workflow File

Create `.github/workflows/release.yml`:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*.*.*'

permissions:
  contents: write

jobs:
  build-windows:
    runs-on: windows-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest

      - name: Install NSIS
        run: choco install nsis -y

      - name: Install frontend dependencies
        run: npm install
        working-directory: frontend

      - name: Get version from tag
        id: version
        shell: bash
        run: echo "VERSION=${GITHUB_REF#refs/tags/}" >> $GITHUB_OUTPUT

      - name: Build with Wails
        run: |
          wails build --target windows/amd64 --nsis `
            -ldflags "-X 'main.version=${{ steps.version.outputs.VERSION }}'"

      - name: Upload to Release
        uses: softprops/action-gh-release@v2
        with:
          files: |
            build/bin/*.exe
          draft: false
          prerelease: false
```

### 2.2 NSIS Output Verification

Existing NSIS template outputs to:
- `build/bin/ImageBorderTool-amd64-installer.exe`

### 2.3 Release Assets

Expected assets per release:
- `ImageBorderTool.exe` - Standalone executable
- `ImageBorderTool-amd64-installer.exe` - NSIS installer

## Files to Create

| File | Description |
|------|-------------|
| `.github/workflows/release.yml` | CI/CD workflow |

## Workflow Trigger

```bash
# Create and push tag
git tag v1.0.1
git push origin v1.0.1
```

## Success Criteria

- [ ] Workflow triggers on tag push
- [ ] Go 1.24 and Node.js 20 installed
- [ ] Wails CLI available
- [ ] NSIS creates installer
- [ ] Both exe files uploaded to release
- [ ] Release is public (not draft)
