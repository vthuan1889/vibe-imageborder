# Implementation Plan: Image Border Application

**Project:** vibe-imageborder
**Tech Stack:** Go + Wails v3 + React + TailwindCSS
**Created:** 2026-01-07
**Status:** Ready for Implementation

---

## Overview

Desktop app ghép ảnh sản phẩm với khung viền + text overlay động. Input: nhiều ảnh sản phẩm + 1 khung + template JSON. Output: N ảnh ghép với text fields.

**Reference Brainstorming:** [plans/reports/brainstormer-260107-0945-image-border-app-solution.md](plans/reports/brainstormer-260107-0945-image-border-app-solution.md)

---

## Architecture Summary

```
Frontend (React + TailwindCSS)
    ↓ Wails Bindings
Backend (Go)
    ├── TemplateService - Parse JSON, extract fields
    ├── ImageService - Composite, render text
    └── App - Wails bindings, events
```

**Core Libraries:**
- `disintegration/imaging` - Image loading, resizing
- `fogleman/gg` - 2D graphics, TTF text rendering
- `wails/v3` - Desktop framework

---

## Implementation Phases

### [Phase 1: Project Setup & Foundation](phase-01-project-setup.md)
**Goal:** Wails project initialized với Go + React frontend
**Status:** DONE (2026-01-07 11:15 AM)

**Tasks:**
1. Install Wails v3 CLI
2. Initialize project: `wails init -n vibe-imageborder -t react-ts`
3. Setup project structure
4. Configure TailwindCSS
5. Add Go dependencies
6. Test build & run

**Deliverables:**
- ✓ Running Wails app với React frontend
- ✓ TailwindCSS configured
- ✓ Dependencies installed
- ✓ Wails v3 installed: v3.0.0-alpha.57
- ✓ Go dependencies: imaging v1.6.2, gg v1.3.0
- ✓ Project structure created
- ✓ Frontend dev server working
- ✓ Bindings generated
- ✓ All tests passed

---

### [Phase 2: Template Service](phase-02-template-service.md)
**Goal:** Parse template JSON, extract dynamic fields
**Status:** DONE (2026-01-07 12:15 PM)

**Tasks:**
1. Define `TemplateField` và `Template` types
2. Implement `ParseTemplate(path)` - Load & parse JSON
3. Implement `ExtractDynamicFields(tmpl)` - Find `[fields]`
4. Implement `ReplaceVariables(text, values)` - Substitute values
5. Unit tests với reference templates

**Deliverables:**
- ✓ `internal/template/service.go`
- ✓ `internal/template/parser.go`
- ✓ `internal/models/types.go`
- ✓ Unit tests pass

---

### [Phase 3: Image Service - Core Processing](phase-03-image-service-core.md)
**Goal:** Load images, resize, composite without text

**Tasks:**
1. Define `CompositeRequest` type
2. Implement `LoadImage(path)` - Load JPEG/PNG
3. Implement `ResizeToFit(img, bounds)` - Contain mode
4. Implement `CompositeImages(product, frame)` - Overlay centered
5. Implement `SaveImage(img, path)` - PNG output
6. CLI test program

**Deliverables:**
- ✓ `internal/image/service.go`
- ✓ `internal/image/compositor.go`
- ✓ Basic image compositing working (no text)

---

### [Phase 4: Image Service - Text Rendering](phase-04-text-rendering.md)
**Goal:** Add text overlay với TTF fonts

**Tasks:**
1. Embed font: `assets/fonts/Roboto-Regular.ttf`
2. Implement `parseColor(colorStr)` - "white"/"#FFFFFF" → RGBA
3. Implement `parsePosition(posStr)` - "98,1720" → (x, y)
4. Implement `RenderText(ctx, field, values)` - Draw text
5. Integrate với `CompositeImages()`
6. Test với reference templates

**Deliverables:**
- ✓ `assets/fonts/Roboto-Regular.ttf` embedded
- ✓ Text rendering working
- ✓ Test outputs match expected positions

---

### [Phase 5: Wails Backend Integration](phase-05-wails-backend.md)
**Goal:** Expose services qua Wails bindings

**Tasks:**
1. Create `App` struct trong `app.go`
2. Bind methods:
   - `SelectProductImages()` → file picker (multi)
   - `SelectFrameImage()` → file picker (single)
   - `SelectTemplateFile()` → file picker (.txt)
   - `ParseTemplate(path)` → extract fields
   - `ProcessBatch(request)` → composite images
3. Implement progress events: `runtime.EventsEmit()`
4. Error handling và validation

**Deliverables:**
- ✓ `app.go` với bindings
- ✓ File dialogs working
- ✓ Progress events emitting

---

### [Phase 6: React Frontend - UI Components](phase-06-react-frontend.md)
**Goal:** Build UI với file pickers, dynamic form, progress

**Tasks:**
1. `FilePicker` component - Products, Frame, Template
2. `TemplateFields` component - Dynamic inputs từ template
3. `ProgressBar` component - Real-time progress
4. `App.jsx` - Main layout và state management
5. Style với TailwindCSS
6. Wire up Wails bindings

**Deliverables:**
- ✓ `frontend/src/components/FilePicker.jsx`
- ✓ `frontend/src/components/TemplateFields.jsx`
- ✓ `frontend/src/components/ProgressBar.jsx`
- ✓ `frontend/src/App.jsx`
- ✓ Functional UI

---

### [Phase 7: Integration & Testing](phase-07-integration-testing.md)
**Goal:** End-to-end workflow working

**Tasks:**
1. Test complete workflow: select files → fill fields → process
2. Test với reference templates (khung-002-05.txt, etc.)
3. Test batch processing (10+ images)
4. Validate output quality
5. Test error scenarios
6. Performance testing

**Deliverables:**
- ✓ E2E workflow functional
- ✓ Output quality validated
- ✓ Error handling robust

---

### [Phase 8: Polish & Production Readiness](phase-08-polish.md)
**Goal:** Production-quality UX và error handling

**Tasks:**
1. Output directory selection
2. Filename customization options
3. Clear error messages và validation
4. Loading states và UI polish
5. App icon và branding
6. Build executable
7. User testing

**Deliverables:**
- ✓ Polished UX
- ✓ Clear error messages
- ✓ Production build
- ✓ User documentation

---

## File Structure (Final)

```
vibe-imageborder/
├── main.go                     # Wails entry point
├── app.go                      # App struct, bindings
├── go.mod
├── go.sum
├── wails.json
├── build/                      # Build configs
├── frontend/
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   ├── src/
│   │   ├── App.jsx            # Main component
│   │   ├── main.jsx
│   │   ├── style.css
│   │   └── components/
│   │       ├── FilePicker.jsx
│   │       ├── TemplateFields.jsx
│   │       └── ProgressBar.jsx
│   └── dist/                   # Build output
├── internal/
│   ├── template/
│   │   ├── service.go         # TemplateService
│   │   └── parser.go          # JSON parsing logic
│   ├── image/
│   │   ├── service.go         # ImageService
│   │   └── compositor.go      # Compositing logic
│   └── models/
│       └── types.go           # Shared types
├── assets/
│   └── fonts/
│       └── Roboto-Regular.ttf # Bundled font
└── tests/
    ├── fixtures/              # Test images, templates
    └── *_test.go
```

---

## Technical Specifications

### Template JSON Format

```json
{
  "barcode": {
    "text": "[barcode]",
    "position": "98,1720",
    "fontsize": "45",
    "color": "white"
  },
  "size": {
    "text": "D[size_dai] x R[size_rong] x C[size_cao] cm",
    "position": "26,1852",
    "fontsize": "40",
    "color": "white"
  }
}
```

**Field Extraction:**
- Regex: `\[([^\]]+)\]` → extracts: `barcode`, `size_dai`, `size_rong`, `size_cao`
- Generate input fields dynamically trong UI

### Image Processing Pipeline

```
1. Load product image → imaging.Open(productPath)
2. Load frame image → imaging.Open(framePath)
3. Resize product to fit → imaging.Fit(product, frameW, frameH, Lanczos)
4. Create gg.Context → gg.NewContextForImage(frame)
5. Calculate center position → (frameW - productW)/2, (frameH - productH)/2
6. Draw product → ctx.DrawImage(product, centerX, centerY)
7. For each template field:
   a. Parse position → parsePosition("98,1720") = (98, 1720)
   b. Parse fontsize → strconv.Atoi("45") = 45
   c. Parse color → parseColor("white") = RGBA{255,255,255,255}
   d. Replace variables → ReplaceVariables("D[size_dai] cm", {"size_dai": "30"}) = "D30 cm"
   e. Load font → ctx.LoadFontFace(fontPath, fontSize)
   f. Set color → ctx.SetColor(color)
   g. Draw text → ctx.DrawString(text, x, y)
8. Save output → imaging.Save(ctx.Image(), outputPath)
```

### Color Parsing Logic

```go
func parseColor(colorStr string) color.Color {
    switch strings.ToLower(colorStr) {
    case "white":
        return color.RGBA{255, 255, 255, 255}
    case "black":
        return color.RGBA{0, 0, 0, 255}
    case "red":
        return color.RGBA{255, 0, 0, 255}
    default:
        // Parse hex: #RRGGBB
        if strings.HasPrefix(colorStr, "#") {
            // Implementation: hex to RGBA
        }
        return color.RGBA{255, 255, 255, 255} // Default white
    }
}
```

### Progress Events

```go
// Backend
runtime.EventsEmit(ctx, "progress", map[string]interface{}{
    "current": i + 1,
    "total": total,
    "filename": filepath.Base(productPath),
})

// Frontend
runtime.EventsOn("progress", (data) => {
    setProgress({
        current: data.current,
        total: data.total,
        filename: data.filename
    });
});
```

---

## Dependencies

### Go Modules

```go
module github.com/yourusername/vibe-imageborder

go 1.21

require (
    github.com/wailsapp/wails/v3 v3.0.0-alpha.0
    github.com/disintegration/imaging v1.6.2
    github.com/fogleman/gg v1.3.0
)
```

### Frontend Packages

```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "@vitejs/plugin-react": "^4.0.0",
    "autoprefixer": "^10.4.14",
    "postcss": "^8.4.24",
    "tailwindcss": "^3.4.0",
    "vite": "^5.0.0"
  }
}
```

---

## Testing Strategy

### Unit Tests (Go)

```
internal/template/
  ├── parser_test.go
  │   ├── TestParseTemplate
  │   ├── TestExtractDynamicFields
  │   └── TestReplaceVariables
internal/image/
  └── compositor_test.go
      ├── TestLoadImage
      ├── TestResizeToFit
      ├── TestCompositeImages
      └── TestRenderText
```

### Integration Tests

```
tests/
  └── integration_test.go
      ├── TestCompleteWorkflow
      ├── TestBatchProcessing
      └── TestErrorHandling
```

### Test Fixtures

```
tests/fixtures/
  ├── templates/
  │   ├── khung-002-05.txt
  │   └── khung-004-01.txt
  ├── frames/
  │   ├── frame-01.png
  │   └── frame-02.png
  └── products/
      ├── product-01.jpg
      └── product-02.jpg
```

---

## Error Handling

### Validation Rules

1. **Image Files:**
   - Format: JPEG/PNG only
   - Size: Max 10000x10000px (prevent OOM)
   - Exists: File must exist

2. **Template File:**
   - Format: Valid JSON
   - Required fields: text, position, fontsize, color
   - Position format: "x,y"

3. **User Input:**
   - All template fields must be filled
   - No empty strings

### Error Messages

```go
// User-friendly errors
var (
    ErrInvalidImageFormat = errors.New("Invalid image format. Please use JPEG or PNG")
    ErrImageTooLarge = errors.New("Image too large. Max size: 10000x10000px")
    ErrInvalidTemplate = errors.New("Invalid template file. Please check JSON format")
    ErrMissingFields = errors.New("Please fill all required fields")
)
```

---

## Performance Targets

| Metric | Target | Rationale |
|--------|--------|-----------|
| Processing time | <5s per image (2000x2000px) | Acceptable for batch |
| Memory usage | <100MB | Desktop app constraint |
| Startup time | <3s | Good UX |
| Batch capacity | 100+ images | Production use |

---

## Risk Mitigation

### Technical Risks

1. **Large images cause OOM**
   - Mitigation: Validate max dimensions, downscale if needed
   - Implementation: Check `img.Bounds()` before processing

2. **Font rendering quality**
   - Mitigation: Use proven `fogleman/gg` library
   - Testing: Visual comparison với reference outputs

3. **Template parsing errors**
   - Mitigation: Strict validation, clear error messages
   - Implementation: Validate JSON structure before processing

### UX Risks

1. **Unclear field labels**
   - Mitigation: Use field name as label with placeholder
   - Example: `[barcode]` → Input label: "Barcode"

2. **No preview**
   - Mitigation: Phase 8 enhancement (optional)
   - Workaround: Fast processing, users can retry

---

## Success Criteria

### Must Have (v1.0)

- ✓ Load multiple product images
- ✓ Select frame image và template file
- ✓ Dynamic form generation từ template
- ✓ Batch processing với progress tracking
- ✓ Text overlay với correct positioning
- ✓ Output quality matches input
- ✓ Error handling và validation

### Should Have (v1.1)

- Preview panel
- Output directory selection
- Filename customization

### Nice to Have (v2.0)

- CSV import for field values
- Custom fonts folder
- Batch template processing
- Image quality adjustment

---

## Development Guidelines

### Code Style

- **Go:** Follow [Effective Go](https://go.dev/doc/effective_go)
- **React:** Functional components với hooks
- **Naming:** Clear, descriptive names (KISS principle)

### Git Workflow

```bash
# Feature branches
git checkout -b phase-01-project-setup
git checkout -b phase-02-template-service
...

# Commit per task
git commit -m "feat(template): implement ParseTemplate function"
git commit -m "test(template): add unit tests for ExtractDynamicFields"
```

### Testing Requirements

- All public functions must have tests
- Test coverage: >80%
- Integration test for complete workflow

---

## Next Steps

1. **Start Phase 1:** [Project Setup & Foundation](phase-01-project-setup.md)
2. Track progress qua phase completion
3. Update this plan nếu requirements change
4. Review after each phase

---

## Unresolved Questions

*To be addressed during implementation:*

1. Image size constraints: Max input dimensions? Downscale strategy?
2. Output format: Always PNG? Support JPEG option?
3. Font selection: Single bundled font OK, or need multiple weights/styles?
4. Field validation: Require all fields filled, or allow empty?
5. Error recovery: Retry failed images? Save error log file?

*Will clarify with user during respective phases.*

---

## Resources

- **Brainstorming Report:** [plans/reports/brainstormer-260107-0945-image-border-app-solution.md](plans/reports/brainstormer-260107-0945-image-border-app-solution.md)
- **Wails v3 Docs:** https://wails.io/docs
- **imaging Library:** https://github.com/disintegration/imaging
- **gg Library:** https://github.com/fogleman/gg
- **Reference Project:** D:\Code-Tool\Software\web-tool\UploadImage

---

**Plan Status:** ✅ Ready for Implementation
**Estimated Effort:** 20-26 hours for v1.0
**Risk Level:** Low (proven tech stack, clear requirements)
