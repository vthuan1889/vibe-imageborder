# Phase 1: Version Management

> Parent: [plan.md](./plan.md)

## Overview

| Field | Value |
|-------|-------|
| Date | 2026-01-08 |
| Priority | P1 |
| Effort | 1h |
| Status | pending |

## Requirements

1. Version injected at build time via ldflags
2. Version accessible from both Go and frontend
3. Format: `v1.0.0` (semantic with 'v' prefix)
4. Fallback for development builds

## Implementation Steps

### 1.1 Add Version Variable to main.go

```go
// Version info - injected at build time via ldflags
var (
    version   = "dev"
    buildTime = "unknown"
)
```

### 1.2 Add Version Getter to app.go

```go
// GetVersion returns the current app version
func (a *App) GetVersion() string {
    return version
}
```

Note: Need to move `version` var to app.go or export it.

### 1.3 Update wails.json (optional)

Add version field for NSIS to read:
```json
{
  "info": {
    "productVersion": "1.0.0"
  }
}
```

### 1.4 Build Command Update

Normal build:
```bash
wails build --target windows/amd64 --nsis
```

Release build (with version):
```bash
wails build --target windows/amd64 --nsis \
  -ldflags "-X 'main.version=v1.0.0' -X 'main.buildTime=$(date -u)'"
```

## Files to Modify

| File | Change |
|------|--------|
| `main.go` | Add version/buildTime variables |
| `app.go` | Add GetVersion() method |
| `wails.json` | Add info.productVersion |

## Success Criteria

- [ ] `GetVersion()` returns correct version
- [ ] Dev builds show "dev" version
- [ ] Release builds show tag version
- [ ] NSIS installer shows correct version
