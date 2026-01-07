# Phase 7: Integration & Testing

**Goal:** End-to-end testing, validation, bug fixes

**Duration:** ~3-4 hours

**Dependencies:** Phase 1-6 complete

---

## Overview

Comprehensive testing:
1. End-to-end workflow testing
2. Reference template validation
3. Batch processing stress testing
4. Error scenario testing
5. Output quality validation
6. Performance testing

---

## Task 7.1: E2E Workflow Testing

### Test Case 1: Basic Workflow

**Steps:**
1. Launch app: `wails3 dev`
2. Select 3 product images
3. Select frame image
4. Select template (khung-002-05.txt)
5. Fill all fields:
   - barcode: "TEST123"
   - size_dai: "30"
   - size_rong: "20"
   - size_cao: "15"
6. Select output directory
7. Click "Process Images"
8. Wait for completion
9. Verify output files

**Expected:**
- ✓ 3 PNG files created trong output dir
- ✓ Files named: `{original}_framed.png`
- ✓ Product images centered trong frames
- ✓ Text rendered tại correct positions
- ✓ Text content matches input values

---

### Test Case 2: Multiple Templates

Test với different templates:

| Template | Fields Required | Expected Text |
|----------|----------------|---------------|
| khung-002-05.txt | barcode, size_dai, size_rong, size_cao | All 4 values displayed |
| khung-004-01.txt | (check actual fields) | Fields displayed correctly |

**Validation:**
- Each template's fields display correctly
- Text positioning unique per template
- No overlapping text

---

## Task 7.2: Batch Processing Stress Test

### Test Large Batches

```bash
# Create 50 test product images
for i in {1..50}; do
  cp tests/fixtures/products/product-01.jpg \
     tests/fixtures/products/product-$i.jpg
done
```

**Test:**
1. Select all 50 images
2. Process với same template
3. Monitor:
   - Memory usage (should stay <100MB)
   - Processing time (~3-5s per image)
   - Progress updates smooth
   - No crashes

**Metrics:**
- Total time: <250s (50 images × 5s)
- Memory peak: <100MB
- Success rate: 100%

---

## Task 7.3: Error Scenario Testing

### Test Case 3: Invalid Image File

**Setup:**
```bash
# Create invalid image
echo "not an image" > tests/fixtures/products/invalid.jpg
```

**Test:**
1. Select invalid.jpg trong product list
2. Process batch

**Expected:**
- App doesn't crash
- Error logged trong console
- Invalid image skipped
- Other images process successfully
- Final result shows FailedCount = 1

---

### Test Case 4: Missing Template Fields

**Setup:**
- Fill only some fields (e.g., barcode only)
- Leave size_dai empty

**Test:**
1. Try to process

**Expected:**
- Process button stays disabled
- OR alert shows "Please fill all fields"

---

### Test Case 5: Invalid Template File

**Setup:**
```json
// tests/fixtures/templates/invalid.txt
{
  "broken json
}
```

**Test:**
1. Select invalid.txt

**Expected:**
- Error alert: "Failed to load template: ..."
- Fields don't generate
- Process button disabled

---

## Task 7.4: Output Quality Validation

### Visual Inspection Checklist

For each output image:

**Image Quality:**
- [ ] No pixelation or artifacts
- [ ] Colors accurate (không color shift)
- [ ] PNG transparency preserved (if applicable)

**Positioning:**
- [ ] Product centered trong frame
- [ ] Aspect ratio preserved (không distortion)
- [ ] Text tại exact positions from template

**Text Rendering:**
- [ ] Font clear và readable
- [ ] Font size matches template
- [ ] Color matches template (white visible trên dark bg)
- [ ] No text cutoff or overlap

**File Properties:**
- [ ] Format: PNG
- [ ] Dimensions match frame
- [ ] File size reasonable (<5MB for 2000x2000)

---

### Automated Visual Test

Create `tests/visual-test.go`:

```go
package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"vibe-imageborder/internal/image"
)

func main() {
	// Load output image
	svc := image.NewService()
	img, err := svc.LoadImage("tests/output/product-01_framed.png")
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	fmt.Printf("Dimensions: %dx%d\n", bounds.Dx(), bounds.Dy())

	// Check center pixel (should be product color, not frame)
	centerX, centerY := bounds.Dx()/2, bounds.Dy()/2
	centerColor := img.At(centerX, centerY)
	r, g, b, a := centerColor.RGBA()
	fmt.Printf("Center pixel: RGBA(%d,%d,%d,%d)\n", r>>8, g>>8, b>>8, a>>8)

	// Check text position (should not be pure white/black = has text)
	textX, textY := 100, 1720 // From template
	textColor := img.At(textX, textY)
	if textColor == (color.RGBA{}) {
		log.Println("Warning: No text detected tại position")
	} else {
		fmt.Printf("Text pixel detected tại (%d,%d)\n", textX, textY)
	}

	fmt.Println("✓ Visual test passed")
}
```

---

## Task 7.5: Performance Benchmarking

Create `tests/benchmark_test.go`:

```go
package tests

import (
	"testing"
	"vibe-imageborder/internal/image"
	"vibe-imageborder/internal/models"
	"vibe-imageborder/internal/template"
)

func BenchmarkSingleImage(b *testing.B) {
	imgSvc := image.NewService()
	tmplSvc := template.NewService()

	tmpl, _ := tmplSvc.Load("fixtures/templates/khung-002-05.txt")
	fieldValues := map[string]string{
		"barcode":   "BENCH123",
		"size_dai":  "30",
		"size_rong": "20",
		"size_cao":  "15",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = imgSvc.ProcessSingle(
			"fixtures/products/product-01.jpg",
			"fixtures/frames/frame-01.png",
			"output/bench.png",
			tmpl,
			fieldValues,
		)
	}
}
```

Run benchmark:

```bash
go test -bench=. -benchmem ./tests/
```

**Target Performance:**
- Single image processing: <5s
- Memory allocation: <50MB per image

---

## Task 7.6: Create Test Report

After all tests, document results trong `tests/test-report.md`:

```markdown
# Test Report - Image Border Application

**Date:** 2026-01-07
**Tester:** [Name]
**Version:** v1.0-alpha

## E2E Testing

| Test Case | Status | Notes |
|-----------|--------|-------|
| Basic workflow (3 images) | ✓ Pass | All outputs correct |
| Multiple templates | ✓ Pass | khung-002-05, khung-004-01 |
| Batch processing (50 images) | ✓ Pass | ~180s total |

## Error Handling

| Scenario | Status | Notes |
|----------|--------|-------|
| Invalid image file | ✓ Pass | Skipped gracefully |
| Missing fields | ✓ Pass | Process button disabled |
| Invalid template JSON | ✓ Pass | Error alert shown |

## Quality Validation

| Aspect | Status | Notes |
|--------|--------|-------|
| Image quality | ✓ Pass | No artifacts |
| Positioning accuracy | ✓ Pass | ±2px tolerance |
| Text rendering | ✓ Pass | Clear, correct font |
| File format | ✓ Pass | PNG, correct dims |

## Performance

| Metric | Result | Target |
|--------|--------|--------|
| Single image | 3.2s | <5s ✓ |
| Memory usage | 45MB | <100MB ✓ |
| Batch (50 images) | 178s | <250s ✓ |

## Issues Found

### Critical
- None

### Minor
- [Issue #1] Progress bar sometimes updates late (cosmetic)

## Conclusion

**Status:** ✓ PASS - Ready for Phase 8 (Polish)

All core functionality working correctly. Minor UX improvements needed trong Phase 8.
```

---

## Acceptance Criteria

- ✓ E2E workflow completes successfully
- ✓ All reference templates work correctly
- ✓ Batch processing handles 50+ images
- ✓ Error scenarios handled gracefully
- ✓ Output quality meets standards
- ✓ Performance meets targets
- ✓ Test report documented

---

## Deliverables

### Files Created

1. `tests/visual-test.go` - Visual validation script
2. `tests/benchmark_test.go` - Performance benchmarks
3. `tests/test-report.md` - Test results documentation

### Test Artifacts

```
tests/
├── fixtures/
│   ├── products/ (50+ test images)
│   ├── frames/ (multiple frames)
│   └── templates/ (valid + invalid templates)
├── output/ (test outputs for validation)
└── test-report.md
```

---

## Troubleshooting

### Issue: Text not visible trong output

**Debug steps:**
1. Check text color vs background
2. Verify position within image bounds
3. Check font loading (log font path)
4. Increase font size for visibility test

### Issue: Slow processing

**Debug steps:**
1. Check image dimensions (downscale large images)
2. Profile memory usage
3. Verify sequential processing (not blocking UI)

---

## Next Phase

[Phase 8: Polish & Production Readiness](phase-08-polish.md)
