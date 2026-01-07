# Image Border Application - Solution Brainstorming

**Date:** 2026-01-07
**Project:** vibe-imageborder
**Tech Stack:** Go + Wails v3 + TailwindCSS

---

## Problem Statement

Xây dựng desktop app ghép ảnh sản phẩm với khung viền, kèm text overlay động:

- **Input:** Nhiều ảnh sản phẩm + 1 ảnh khung + template text file (.txt)
- **Processing:** Overlay từng sản phẩm lên khung, apply text fields từ template
- **Output:** N ảnh ghép (1 khung → N outputs)
- **Text handling:** Parse template JSON với dynamic fields `[barcode]`, `[size_dai]`, etc.

---

## Requirements Analysis

### Functional Requirements

1. **Image Selection**
   - Multi-select product images (JPEG/PNG)
   - Single frame/border image selection
   - Template file picker (.txt với JSON config)

2. **Image Compositing**
   - Overlay mode: product centered trong frame
   - Aspect ratio preservation (contain mode)
   - Maintain input image quality
   - PNG transparency support

3. **Text Rendering**
   - Parse template JSON: `{field: {text, position, fontsize, color}}`
   - Dynamic field replacement: `[barcode]` → user input
   - TTF font support (bundled)
   - Position-based placement (x,y coordinates)

4. **Batch Processing**
   - Sequential processing (memory-optimized)
   - Progress tracking UI
   - Error handling per image

5. **Output**
   - Save to user-selected directory
   - Naming convention: `{original_name}_framed.png`
   - Same quality as input

### Non-Functional Requirements

- **Performance:** Sequential processing, ~2-5s per image
- **UX:** Simple, clean interface với TailwindCSS
- **Reliability:** Graceful error handling, validation
- **Maintainability:** Clean architecture, modular code

---

## Reference Analysis

### Existing C# Project (D:\Code-Tool\Software\web-tool\UploadImage)

**Key Observations:**
- Uses `System.Drawing` for image manipulation
- Template structure: JSON với fields `{text, position, fontsize, color}`
- Multiple frame templates: khung-002-05.txt, khung-004-01.txt, etc.
- File upload functionality to web API

**Template Format Example (khung-002-05.txt):**
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

**Insights:**
- Position format: `"x,y"` strings
- Dynamic fields in `[brackets]`
- Common fields: barcode, size dimensions, product codes
- Multiple templates per frame type

---

## Technical Solution

### Architecture Overview

```
┌─────────────────────────────────────────┐
│          Frontend (Wails + React)       │
│  ┌────────────────────────────────────┐ │
│  │  - File pickers                    │ │
│  │  - Input fields (dynamic)          │ │
│  │  - Progress bar                    │ │
│  │  - Preview                         │ │
│  └────────────────────────────────────┘ │
└──────────────┬──────────────────────────┘
               │ Wails Bindings
┌──────────────▼──────────────────────────┐
│          Backend (Go)                   │
│  ┌────────────────────────────────────┐ │
│  │  ImageService                      │ │
│  │  - LoadImages()                    │ │
│  │  - CompositeImage()                │ │
│  │  - RenderText()                    │ │
│  │  - SaveOutput()                    │ │
│  └────────────────────────────────────┘ │
│  ┌────────────────────────────────────┐ │
│  │  TemplateService                   │ │
│  │  - ParseTemplate()                 │ │
│  │  - ExtractFields()                 │ │
│  │  - ReplaceVariables()              │ │
│  └────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

### Tech Stack Details

#### Go Libraries

**Image Processing:**
- **`disintegration/imaging`** - Image loading, resizing, transformations
  - Pros: Fast, simple API, production-stable
  - Use: Load images, resize product to fit frame

- **`fogleman/gg`** - 2D graphics, text rendering
  - Pros: Excellent TTF support, easy API
  - Use: Composite images, render text overlays

- **`golang/freetype`** (dependency of gg) - TTF font rendering

**Utilities:**
- `encoding/json` - Parse template files
- `path/filepath` - File operations
- `strings` - Text manipulation

**Wails v3:**
- Frontend-backend bindings
- Native file dialogs
- Event system for progress updates

#### Frontend

- **React** - UI components
- **TailwindCSS** - Styling
- **Wails Runtime** - Native APIs

---

## Implementation Approach

### 1. Project Structure

```
vibe-imageborder/
├── main.go                 # Wails app entry
├── app.go                  # App struct, bindings
├── frontend/
│   ├── src/
│   │   ├── App.jsx        # Main component
│   │   ├── components/
│   │   │   ├── FilePicker.jsx
│   │   │   ├── TemplateFields.jsx
│   │   │   └── ProgressBar.jsx
│   │   └── styles/
│   └── package.json
├── internal/
│   ├── image/
│   │   ├── service.go      # ImageService
│   │   └── compositor.go   # Compositing logic
│   ├── template/
│   │   ├── service.go      # TemplateService
│   │   └── parser.go       # JSON parsing
│   └── models/
│       └── types.go        # Shared types
├── assets/
│   └── fonts/
│       └── Roboto-Regular.ttf  # Bundled font
└── wails.json
```

### 2. Core Components

#### A. Template Service

```go
type TemplateField struct {
    Text     string `json:"text"`
    Position string `json:"position"` // "x,y"
    FontSize string `json:"fontsize"`
    Color    string `json:"color"`
}

type Template map[string]TemplateField

func ParseTemplate(path string) (Template, error)
func ExtractDynamicFields(tmpl Template) []string
func ReplaceVariables(text string, values map[string]string) string
```

**Logic:**
1. Load template JSON file
2. Parse into Template struct
3. Extract unique field names from `[brackets]`
4. Return field list for UI to generate inputs
5. Replace `[field]` với user input values

#### B. Image Service

```go
type CompositeRequest struct {
    ProductPaths []string
    FramePath    string
    Template     Template
    FieldValues  map[string]string
    OutputDir    string
}

func (s *ImageService) ProcessBatch(req CompositeRequest, onProgress func(int, int)) error
```

**Processing Pipeline:**

```
For each product image:
  1. Load product image → imaging.Open()
  2. Load frame image → imaging.Open()
  3. Calculate product size → imaging.Fit() to contain mode
  4. Create gg.Context from frame
  5. Draw resized product centered → dc.DrawImage()
  6. For each template field:
     - Parse position (x,y)
     - Parse fontsize
     - Parse color (white/black/hex)
     - Replace [variables] in text
     - Load font face → dc.LoadFontFace()
     - Draw text → dc.DrawString()
  7. Save output → imaging.Save()
  8. Emit progress event
```

#### C. Frontend UI Flow

**Initial State:**
- 3 file pickers: Products (multi), Frame (single), Template (single)
- Process button (disabled)

**After Template Selected:**
- Parse template
- Show dynamic input fields for `[barcode]`, `[size_dai]`, etc.
- Enable Process button

**Processing:**
- Show progress bar
- Disable inputs
- Stream progress updates from backend

**Complete:**
- Show success message với output path
- Re-enable inputs for next batch

---

## Key Technical Decisions

### 1. Image Compositing Strategy

**Decision:** Overlay product centered on frame

**Implementation:**
```go
// Resize product to fit frame with contain mode
productFit := imaging.Fit(productImg, frameWidth, frameHeight, imaging.Lanczos)

// Calculate center position
x := (frameWidth - productFit.Bounds().Dx()) / 2
y := (frameHeight - productFit.Bounds().Dy()) / 2

// Create gg context from frame
dc := gg.NewContextForImage(frameImg)

// Draw product at center
dc.DrawImage(productFit, x, y)
```

**Rationale:**
- Preserves product aspect ratio
- Simple, predictable layout
- No data loss from cropping

**Alternative considered:** Template-defined product area
- Rejected: Over-engineering for current needs (YAGNI)
- Can add later if needed

### 2. Text Rendering Approach

**Decision:** Use `fogleman/gg` với bundled TTF fonts

**Implementation:**
```go
// Load bundled font
fontPath := filepath.Join("assets", "fonts", "Roboto-Regular.ttf")
dc.LoadFontFace(fontPath, fontSize)

// Parse color
color := parseColor(field.Color) // "white" → color.RGBA{255,255,255,255}
dc.SetColor(color)

// Draw text at position
dc.DrawString(text, x, y)
```

**Rationale:**
- Consistent rendering across platforms
- No system font dependencies
- High-quality anti-aliased text

**Font embedding in Wails:**
```go
//go:embed assets/fonts/*
var fontsFS embed.FS
```

### 3. Processing Model

**Decision:** Sequential processing với progress callbacks

**Implementation:**
```go
func (s *ImageService) ProcessBatch(req CompositeRequest) error {
    total := len(req.ProductPaths)

    for i, productPath := range req.ProductPaths {
        err := s.processOne(productPath, req)
        if err != nil {
            log.Printf("Error processing %s: %v", productPath, err)
            continue // Skip failed images
        }

        // Emit progress event
        runtime.EventsEmit(s.ctx, "progress", i+1, total)
    }

    return nil
}
```

**Rationale:**
- Memory-efficient: process one image at a time
- Simple error handling: skip failed, continue
- Real-time progress feedback

**Alternative considered:** Parallel processing với worker pool
- Rejected: Overkill for desktop app, RAM concerns
- Sequential is fast enough (~2-5s/image)

### 4. Template Format

**Decision:** Keep C# JSON format compatibility

**Rationale:**
- User familiarity (existing templates)
- Simple, human-readable
- Easy to create new templates

**Enhancement:** Position validation
```go
func parsePosition(pos string) (int, int, error) {
    parts := strings.Split(pos, ",")
    if len(parts) != 2 {
        return 0, 0, fmt.Errorf("invalid position format: %s", pos)
    }
    x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
    y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
    return x, y, nil
}
```

---

## Risk Analysis

### Technical Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Large images cause OOM | Medium | Validate image size, downscale if needed |
| Font rendering quality issues | Low | Use proven `gg` + `freetype` combo |
| Template parsing errors | Medium | Strict validation, clear error messages |
| PNG transparency corruption | Low | Use `imaging.Save()` với PNG encoder |

### UX Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Unclear field labels | Medium | Use descriptive placeholders from template |
| No preview before processing | High | **Future:** Add preview panel |
| Lost progress on errors | Medium | Continue on error, log failures |

### Performance Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Slow processing for large batches | Medium | Show progress, allow cancel |
| UI freeze during processing | High | Run processing in goroutine |

---

## Trade-offs Analysis

### 1. Sequential vs Parallel Processing

**Chosen: Sequential**

| Aspect | Sequential | Parallel |
|--------|-----------|----------|
| Speed | ~3s/image | ~1s/image (với 4 workers) |
| Memory | Low (~50MB) | High (~200MB+) |
| Complexity | Simple | Worker pool, sync |
| Error handling | Easy | Complex |

**Decision:** Sequential wins cho desktop app. Speed adequate, simplicity valued.

### 2. Bundled vs System Fonts

**Chosen: Bundled**

| Aspect | Bundled | System |
|--------|---------|--------|
| Consistency | ✓ Perfect | ✗ Platform-dependent |
| App size | +2MB | 0 |
| Font quality | ✓ Controlled | ? Variable |
| Setup | Embed once | User config |

**Decision:** Bundled fonts. 2MB cost negligible, consistency critical.

### 3. React vs Vue vs Svelte

**Chosen: React**

**Rationale:**
- Wails v3 has excellent React templates
- Larger ecosystem, more resources
- Team familiarity (assumption)

**Note:** TailwindCSS works equally well với all 3.

---

## Success Metrics

### Functional Success
- ✓ Process 100+ images without crash
- ✓ Accurate text positioning (±2px)
- ✓ No quality degradation in output
- ✓ Handle all template variations

### Performance Success
- ✓ <5s per image (2000x2000px)
- ✓ <100MB memory usage
- ✓ <3s app startup time

### UX Success
- ✓ 0 clicks to understand workflow
- ✓ Clear error messages
- ✓ Progress visibility
- ✓ Output preview (future)

---

## Implementation Plan Overview

### Phase 1: Foundation (Core Processing)
1. Setup Wails v3 project với Go + React
2. Implement TemplateService: parse JSON, extract fields
3. Implement ImageService: load, resize, composite
4. Basic CLI testing (no UI)

### Phase 2: Text Rendering
1. Integrate `fogleman/gg` + font embedding
2. Implement text parsing và color handling
3. Position-based text rendering
4. Test với reference templates

### Phase 3: Frontend (UI/UX)
1. File pickers: products (multi), frame, template
2. Dynamic form generation từ template fields
3. Process button + progress bar
4. Result display với output path

### Phase 4: Polish
1. Error handling + validation
2. Output directory selection
3. Filename customization
4. Testing với production data

### Phase 5: Enhancement (Optional)
1. Preview panel
2. Batch template processing
3. CSV import for field values
4. Custom fonts folder

---

## Dependencies

### Go Modules
```
require (
    github.com/wailsapp/wails/v3 v3.0.0
    github.com/disintegration/imaging v1.6.2
    github.com/fogleman/gg v1.3.0
    github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
)
```

### Frontend
```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "tailwindcss": "^3.4.0",
    "vite": "^5.0.0"
  }
}
```

### Assets
- Roboto-Regular.ttf (bundled in `assets/fonts/`)

---

## Next Steps

1. **User Decision:** Proceed với implementation plan?
2. **If Yes:** Generate detailed implementation plan với `/plan` command
3. **If No:** Iterate on brainstorming, address concerns

---

## Unresolved Questions

1. **Image size constraints:** Max input image dimensions? Downscale strategy?
2. **Output format:** Always PNG? Support JPEG option?
3. **Font selection:** Single bundled font OK, or need multiple weights/styles?
4. **Field validation:** Require all template fields filled, or allow empty?
5. **Error recovery:** Retry failed images? Save error log file?
6. **Future roadmap:** Web version planned? API needed?

---

## Conclusion

**Recommended Solution:**
- **Backend:** Go với `imaging` + `gg` libraries
- **Frontend:** React + TailwindCSS trong Wails v3
- **Architecture:** Clean separation: TemplateService + ImageService
- **Processing:** Sequential với progress tracking
- **Quality:** Production-ready, maintainable, extensible

**Why this approach wins:**
- **YAGNI:** Solves current requirements without over-engineering
- **KISS:** Simple architecture, proven libraries
- **DRY:** Modular services, reusable logic
- **Maintainable:** Clear structure, standard Go patterns
- **Extensible:** Easy to add preview, CSV import, etc. later

**Estimated Effort:**
- Phase 1-2 (Core): ~8-12 hours
- Phase 3 (UI): ~6-8 hours
- Phase 4 (Polish): ~4-6 hours
- **Total:** ~20-26 hours for production-ready v1.0

**Risk Level:** Low
- Proven tech stack
- Clear requirements
- Reference implementation available

---

**Ready to proceed với detailed implementation plan?**
