# Phase 7: Integration & Testing

## Context

- Plan: [plan.md](./plan.md)
- Previous: [Phase 6 - React Frontend](./phase-06-react-frontend.md)

## Overview

| Field | Value |
|-------|-------|
| Priority | P2 |
| Status | Pending |
| Effort | 1.5h |

End-to-end testing with real templates and images from reference app. Verify all features work together.

## Requirements

### Functional
- Test with real template files from `D:\Code-Tool\Software\web-tool\UploadImage\file\`
- Test batch processing with multiple images
- Verify Vietnamese text renders correctly
- Test all output formats

### Non-functional
- No crashes or freezes
- Memory usage stays reasonable
- Processing time within targets

## Test Cases

### TC1: Basic Compositing
| Step | Action | Expected |
|------|--------|----------|
| 1 | Select 1 product image | File shown in picker |
| 2 | Select frame image | File shown in picker |
| 3 | Click Preview | Preview displays composited image |
| 4 | Select output folder | Folder path shown |
| 5 | Click Generate | Single output file created |

### TC2: Template with Vietnamese Text
| Step | Action | Expected |
|------|--------|----------|
| 1 | Select products | Files shown |
| 2 | Select frame (khung-002-05.png) | Frame shown |
| 3 | Select template (khung-002-05.txt) | Fields appear: barcode, size_dai, size_rong, size_cao |
| 4 | Enter values (e.g., "SP001", "50", "30", "20") | Values saved |
| 5 | Click Preview | Text visible at correct positions |
| 6 | Verify Vietnamese diacritics | "Giá", "CM" render correctly |

### TC3: Batch Processing
| Step | Action | Expected |
|------|--------|----------|
| 1 | Select 10+ product images | Count shown |
| 2 | Select frame and template | Files shown |
| 3 | Fill in field values | Values saved |
| 4 | Click Generate All | Progress bar updates |
| 5 | Wait for completion | All files created in output |
| 6 | Verify file names | `*_framed.png` naming |

### TC4: Cancel Processing
| Step | Action | Expected |
|------|--------|----------|
| 1 | Start batch with 20+ images | Processing starts |
| 2 | Click Cancel mid-way | Processing stops |
| 3 | Check output folder | Partial files created |

### TC5: Format Options
| Format | Quality | Expected |
|--------|---------|----------|
| PNG | N/A | Lossless, larger file |
| JPG | 90 | Compressed, good quality |
| JPG | 50 | Compressed, visible artifacts |
| WebP | 90 | Compressed, good quality |

### TC6: Edge Cases
| Case | Expected |
|------|----------|
| No product selected | Generate disabled |
| No frame selected | Preview disabled |
| No template | Fields hidden, no text overlay |
| Empty field values | Placeholder not replaced |
| Invalid template JSON | Error message shown |

## Test Files

Copy from reference app for testing:
```
D:\Code-Tool\Software\web-tool\UploadImage\file\
├── khung-002-05.png  (frame)
├── khung-002-05.txt  (template)
├── khung-004-01.png  (frame with different layout)
├── khung-004-01.txt  (template)
└── [product images from your collection]
```

## Implementation Steps

### Step 1: Create Test Directory

```bash
mkdir -p tests/fixtures/frames
mkdir -p tests/fixtures/products
mkdir -p tests/fixtures/templates
mkdir -p tests/output
```

### Step 2: Copy Test Files

Copy from reference app to test fixtures.

### Step 3: Create Integration Test

```go
// tests/integration_test.go
package tests

import (
    "os"
    "path/filepath"
    "testing"

    imgservice "vibe-imageborder/internal/image"
    "vibe-imageborder/internal/models"
    "vibe-imageborder/internal/template"
)

func TestEndToEnd(t *testing.T) {
    // Setup
    imageSvc := imgservice.NewService()
    templateSvc := template.NewService()

    // Paths
    productPath := "fixtures/products/product-01.jpg"
    framePath := "fixtures/frames/khung-002-05.png"
    templatePath := "fixtures/templates/khung-002-05.txt"
    outputPath := "output/test_e2e_output.png"

    // Load template
    fields, err := templateSvc.GetFields(templatePath)
    if err != nil {
        t.Fatalf("Failed to load template: %v", err)
    }

    if len(fields) == 0 {
        t.Fatal("Expected fields from template")
    }

    // Prepare values
    values := map[string]string{
        "barcode":   "TEST001",
        "size_dai":  "100",
        "size_rong": "50",
        "size_cao":  "30",
    }

    // Get overlays
    overlays, err := templateSvc.GetOverlays(templatePath, values)
    if err != nil {
        t.Fatalf("Failed to get overlays: %v", err)
    }

    // Load images
    product, err := imageSvc.LoadImage(productPath)
    if err != nil {
        t.Fatalf("Failed to load product: %v", err)
    }

    frame, err := imageSvc.LoadImage(framePath)
    if err != nil {
        t.Fatalf("Failed to load frame: %v", err)
    }

    // Note: Font manager requires embedded FS, skip text in unit test
    // Full text rendering tested manually

    // Composite
    compositor := imgservice.NewCompositor(imageSvc)
    result := compositor.Composite(product, frame, "#f1eeea")

    if result.Width == 0 || result.Height == 0 {
        t.Fatal("Result has zero dimensions")
    }

    // Save
    err = imageSvc.SaveImage(result.Image, outputPath, "png", 90)
    if err != nil {
        t.Fatalf("Failed to save: %v", err)
    }

    // Verify output exists
    if _, err := os.Stat(outputPath); os.IsNotExist(err) {
        t.Fatal("Output file not created")
    }
}
```

### Step 4: Create Benchmark Test

```go
// tests/benchmark_test.go
package tests

import (
    "testing"

    imgservice "vibe-imageborder/internal/image"
)

func BenchmarkComposite(b *testing.B) {
    imageSvc := imgservice.NewService()
    compositor := imgservice.NewCompositor(imageSvc)

    product, _ := imageSvc.LoadImage("fixtures/products/product-01.jpg")
    frame, _ := imageSvc.LoadImage("fixtures/frames/khung-002-05.png")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        compositor.Composite(product, frame, "#f1eeea")
    }
}
```

### Step 5: Manual Testing Checklist

Run `wails dev` and test:

- [ ] App starts without errors
- [ ] File dialogs open on Windows
- [ ] Multiple product selection works
- [ ] Template fields appear dynamically
- [ ] Preview generates correctly
- [ ] Vietnamese text renders (diacritics)
- [ ] Progress bar updates
- [ ] Cancel button stops processing
- [ ] Output files are created
- [ ] PNG/JPG/WebP formats work
- [ ] Quality slider affects JPG output

## Todo List

- [ ] Create test fixtures directory
- [ ] Copy test files from reference app
- [ ] Create integration test
- [ ] Create benchmark test
- [ ] Run `go test ./tests/...`
- [ ] Manual testing with real data
- [ ] Fix any issues found

## Success Criteria

1. All unit tests pass
2. Integration test passes
3. Benchmark shows < 300ms per image
4. Manual testing checklist complete
5. Vietnamese text renders correctly

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Test files not available | High | Document required files |
| Font loading in tests | Medium | Skip font tests or embed test font |
| Slow tests | Low | Use smaller test images |

## Next Steps

After completion, proceed to [Phase 8: Polish](./phase-08-polish.md)
