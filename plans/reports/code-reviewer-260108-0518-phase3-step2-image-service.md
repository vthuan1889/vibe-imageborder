# Code Review: Phase 3 Step 2 - Image Service Core

**Date:** 2026-01-08 05:18
**Reviewer:** code-reviewer (afddd1a)
**Scope:** internal/image package (service.go, compositor.go + tests)

---

## Summary

**Critical Issues:** 2
**Warnings:** 3
**Test Coverage:** 95.8%
**Build Status:** ✓ Pass

---

## Critical Issues

### 1. Path Traversal Vulnerability (service.go:29-30)
**Location:** `LoadImage()`

```go
func (s *Service) LoadImage(path string) (image.Image, error) {
    cleanPath := filepath.Clean(path)
    return imaging.Open(cleanPath)
}
```

**Issue:** `filepath.Clean` không ngăn được path traversal attacks (e.g., `../../etc/passwd`). Malicious input có thể đọc bất kỳ file nào user có quyền.

**Impact:** Security breach, file disclosure

**Fix:** Validate path nằm trong allowed directories:
```go
func (s *Service) LoadImage(path string) (image.Image, error) {
    cleanPath := filepath.Clean(path)
    absPath, err := filepath.Abs(cleanPath)
    if err != nil {
        return nil, fmt.Errorf("invalid path: %w", err)
    }
    // Validate path is within allowed directories
    // if !isAllowedPath(absPath) { return error }
    return imaging.Open(absPath)
}
```

### 2. No Image Size Limits (service.go:28-30, compositor.go:28-61)
**Location:** `LoadImage()`, `Composite()`

**Issue:** Không giới hạn dimensions/file size → OOM attacks. Attacker upload 100000x100000 image = crash server.

**Impact:** DoS, memory exhaustion

**Fix:** Add validation:
```go
const (
    MaxImageWidth  = 8192
    MaxImageHeight = 8192
    MaxFileSize    = 50 * 1024 * 1024 // 50MB
)

func (s *Service) LoadImage(path string) (image.Image, error) {
    // Check file size first
    fi, err := os.Stat(path)
    if err != nil || fi.Size() > MaxFileSize {
        return nil, fmt.Errorf("file too large")
    }

    img, err := imaging.Open(path)
    if err != nil {
        return nil, err
    }

    // Validate dimensions
    w, h := s.GetDimensions(img)
    if w > MaxImageWidth || h > MaxImageHeight {
        return nil, fmt.Errorf("image dimensions exceed limits")
    }

    return img, nil
}
```

---

## High Priority Warnings

### 3. Unchecked File Extension (service.go:34-60)
**Location:** `SaveImage()`

**Issue:** `SaveImage` chỉ dựa vào `format` param, không validate file extension thực tế → có thể ghi đè arbitrary files.

**Severity:** Medium-High

**Recommendation:** Validate output path:
```go
allowedExts := map[string]bool{".png": true, ".jpg": true, ".jpeg": true}
ext := filepath.Ext(outputPath)
if !allowedExts[strings.ToLower(ext)] {
    return fmt.Errorf("unsafe output path")
}
```

### 4. Silent Error in parseColor (service.go:88-98)
**Location:** `parseColor()`, `hexToByte()`

**Issue:** `fmt.Sscanf` error bị ignore → invalid hex returns 0, không log warning.

```go
func hexToByte(s string) uint8 {
    var val int
    fmt.Sscanf(s, "%x", &val) // error ignored
    return uint8(val)
}
```

**Impact:** Silent failures, debugging khó khăn

**Fix:**
```go
func parseColor(hex string) (color.Color, error) {
    hex = strings.TrimPrefix(hex, "#")
    if len(hex) != 6 {
        return nil, fmt.Errorf("invalid hex length")
    }
    // ... validate hex chars
}
```

### 5. Memory Not Released Explicitly
**Location:** `Composite()`, `CompositeWithPosition()`

**Issue:** Không gọi explicit cleanup. Go GC handle được, nhưng với large images hoặc batch processing → memory spikes.

**Recommendation:** Add context/cancellation support, document memory usage patterns.

---

## Medium Priority

### 6. DRY Violation in Compositor (compositor.go:28-103)
**Issue:** `Composite()` và `CompositeWithPosition()` duplicate 70% logic.

**Recommendation:** Refactor shared logic:
```go
func (c *Compositor) prepareCanvas(width, height int, bgColor string) *image.RGBA {
    canvas := image.NewRGBA(image.Rect(0, 0, width, height))
    if bgColor != "" {
        bg := c.service.CreateBlankCanvas(width, height, bgColor)
        draw.Draw(canvas, canvas.Bounds(), bg, image.Point{}, draw.Src)
    }
    return canvas
}
```

### 7. Webp Encoder Fallback Not Documented
**Location:** service.go:55-57

**Issue:** Comment says "falls back to PNG" nhưng user expect webp output → unexpected behavior.

**Recommendation:** Return error or log warning when webp requested.

---

## Low Priority

### 8. ToRGBA Utility Function Unused
**Location:** compositor.go:105-115

**Issue:** `ToRGBA()` defined nhưng không dùng trong codebase.

**Recommendation:** Remove hoặc document intended usage (YAGNI).

### 9. Hard-coded Quality Default
**Location:** service.go:54

**Issue:** JPEG quality param passed nhưng không có default/validation.

**Recommendation:** Add validation:
```go
if quality < 1 || quality > 100 {
    quality = 90 // default
}
```

---

## Positive Observations

✓ **Excellent test coverage** (95.8%)
✓ **Clean separation of concerns** (service vs compositor)
✓ **Good error wrapping** with `fmt.Errorf`
✓ **Proper use of `t.TempDir()`** in tests
✓ **Thread-safe** (Service stateless)

---

## Recommended Actions

**Priority 1 (Must Fix):**
1. Add image dimension/size limits in `LoadImage()`
2. Implement path validation to prevent traversal

**Priority 2 (Should Fix):**
3. Validate output paths in `SaveImage()`
4. Handle color parsing errors properly
5. Refactor DRY violations in compositor

**Priority 3 (Nice to Have):**
6. Remove unused `ToRGBA()` or document usage
7. Add JPEG quality validation
8. Document webp fallback behavior

---

## Metrics

- **Files Reviewed:** 4 (service.go, compositor.go, 2 test files)
- **LOC Analyzed:** ~350
- **Type Coverage:** N/A (Go)
- **Test Coverage:** 95.8%
- **Build Status:** Pass
- **Go Vet:** Clean

---

## Plan Status Update

**Phase 3 Step 2:** ⚠️ Partially Complete

**Todo List Status:**
- [x] Create `internal/image/service.go`
- [x] Create `internal/image/compositor.go`
- [x] Create unit tests
- [ ] **BLOCKED:** Security fixes required before production
- [ ] Test with real product + frame images (after fixes)
- [ ] Benchmark composite performance

**Next Steps:**
1. Fix Critical Issues #1-2
2. Address High Priority Warnings #3-5
3. Complete remaining todos
4. Update plan to "Complete"

---

## Unresolved Questions

1. Allowed directories for image loading?
2. Max file size limit (50MB reasonable cho product images)?
3. WebP encoding required hoặc fallback to PNG acceptable?
4. Batch processing planned? (affects memory strategy)
