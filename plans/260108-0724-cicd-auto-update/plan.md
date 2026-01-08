---
title: "CI/CD and Auto-Update for Wails v2 App"
description: "Implement GitHub Actions CI/CD pipeline and in-app auto-update feature for vibe-imageborder"
status: completed
priority: P1
effort: 8h
branch: main
tags: [ci-cd, auto-update, github-actions, nsis, wails]
created: 2026-01-08
completed: 2026-01-08
---

# CI/CD and Auto-Update Implementation Plan

## Overview

Implement complete CI/CD pipeline via GitHub Actions and in-app auto-update feature for vibe-imageborder (Wails v2 + React + TypeScript).

## Goals

1. Automated build/release on tag push (v*.*.*)
2. Windows amd64 executable + NSIS installer
3. GitHub Releases integration
4. In-app update check with semver comparison
5. Seamless update via installer download + execution

## Architecture

```
[Tag Push] --> [GitHub Actions] --> [Build Wails] --> [Create NSIS] --> [Upload Release]
                                                                              |
[App] --> [Check Update] --> [Compare Version] --> [Download Installer] --> [Run + Quit]
```

## Phases

| Phase | Description | Effort | Status |
|-------|-------------|--------|--------|
| 1 | Version Management | 1h | completed |
| 2 | GitHub Actions Workflow | 2h | completed |
| 3 | Updater Backend (Go) | 2h | completed |
| 4 | Updater Frontend (React) | 2h | completed |
| 5 | Testing & Validation | 1h | completed |

## Key Decisions

- **Version storage**: Use ldflags to inject version at build time
- **Update channel**: GitHub Releases API (public, no auth required)
- **Installer naming**: `ImageBorderTool-amd64-installer.exe`
- **Temp download path**: `%TEMP%/ImageBorderTool-update.exe`

## Files to Create/Modify

**New Files:**
- `.github/workflows/release.yml` - CI/CD workflow
- `internal/updater/updater.go` - Update logic
- `frontend/src/components/UpdateButton.tsx` - UI component

**Modified Files:**
- `main.go` - Add version variable
- `app.go` - Bind updater methods
- `frontend/src/App.tsx` - Add update button

## Dependencies

- Existing NSIS template: `build/windows/installer/project.nsi`
- GitHub repo: vthuan1889/vibe-imageborder
- Current version: v1.0.0

## Success Criteria

- [ ] Tag push triggers automated build
- [ ] NSIS installer uploaded to GitHub Releases
- [ ] App shows "Check for Update" button
- [ ] Version comparison works correctly
- [ ] Installer downloads and executes properly
