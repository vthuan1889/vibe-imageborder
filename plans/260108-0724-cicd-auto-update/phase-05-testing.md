# Phase 5: Testing & Validation

> Parent: [plan.md](./plan.md)

## Overview

| Field | Value |
|-------|-------|
| Date | 2026-01-08 |
| Priority | P2 |
| Effort | 1h |
| Status | pending |

## Requirements

1. Unit tests for version comparison
2. Manual testing workflow
3. CI workflow validation
4. End-to-end update flow test

## Implementation Steps

### 5.1 Unit Tests for Updater

Create `internal/updater/updater_test.go`:

```go
package updater

import "testing"

func TestCompareVersions(t *testing.T) {
    tests := []struct {
        v1, v2 string
        want   int
    }{
        {"v1.0.0", "v1.0.0", 0},
        {"v1.0.0", "v1.0.1", -1},
        {"v1.0.1", "v1.0.0", 1},
        {"v1.0.0", "v2.0.0", -1},
        {"v1.9.0", "v1.10.0", -1},
        {"dev", "v1.0.0", -1},
        {"v1.0.0", "dev", 1},
    }

    for _, tt := range tests {
        got := compareVersions(tt.v1, tt.v2)
        if got != tt.want {
            t.Errorf("compareVersions(%q, %q) = %d, want %d",
                tt.v1, tt.v2, got, tt.want)
        }
    }
}
```

### 5.2 Manual Testing Checklist

**Local Build Test:**
```bash
# Build with version
wails build --target windows/amd64 -ldflags "-X 'main.version=v0.9.0'"

# Run and check version displays
# Click "Check for Update"
# Verify update is detected (if v1.0.0+ exists)
```

**CI Workflow Test:**
```bash
# Create test tag
git tag v1.0.1-test
git push origin v1.0.1-test

# Monitor Actions tab
# Verify build completes
# Check release assets

# Cleanup
git push origin --delete v1.0.1-test
git tag -d v1.0.1-test
```

### 5.3 End-to-End Test Flow

1. Build app with version `v0.9.0`
2. Create GitHub release `v1.0.0` with installer
3. Run app, click "Check for Update"
4. Verify update detected
5. Click update, confirm dialog
6. Verify installer downloads and runs
7. Verify new version installed

### 5.4 Edge Cases to Test

| Scenario | Expected Behavior |
|----------|-------------------|
| No internet | Error message shown |
| No releases exist | "Up to date" message |
| Same version | "Up to date" message |
| Download interrupted | Error message shown |
| Invalid release format | Graceful error handling |

## Files to Create

| File | Description |
|------|-------------|
| `internal/updater/updater_test.go` | Unit tests |

## Run Tests

```bash
# Run updater tests
go test ./internal/updater/...

# Run all tests
go test ./...
```

## Success Criteria

- [ ] All unit tests pass
- [ ] Manual local build works
- [ ] CI workflow completes successfully
- [ ] Release assets uploaded correctly
- [ ] Full update flow works end-to-end
- [ ] Error cases handled gracefully
