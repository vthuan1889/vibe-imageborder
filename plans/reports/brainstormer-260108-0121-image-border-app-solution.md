# Brainstorm Report: Image Border Application

**Date:** 2026-01-08
**Status:** Solution Agreed

---

## Problem Statement

Build desktop app to composite product images with frame overlays and text annotations. Port from existing C# WinForms app (UploadImage) to Go + Wails 2.

### Requirements Summary
- Multi-select product images (batch processing)
- Single frame selection (overlay)
- Template parsing (.txt JSON format) → extract `[placeholder]` fields
- Dynamic form generation for text input
- Single image preview before batch export
- Export formats: PNG, JPG, WebP
- Data source: Manual input (not API)
- Bundle 2-3 fonts for Vietnamese text

---

## Solution Architecture

### Chosen Approach: Monolithic + Async Goroutines

```
┌─────────────────────────────────────────────┐
│              Wails 2 App                    │
├─────────────────────────────────────────────┤
│  Frontend (React + TailwindCSS)             │
│  ├── FilePicker                             │
│  ├── TemplateFields (dynamic form)          │
│  ├── Preview                                │
│  └── ProgressBar                            │
├─────────────────────────────────────────────┤
│  Backend (Go)                               │
│  ├── ImageService (composite, resize)       │
│  ├── TemplateService (parse JSON)           │
│  └── FileService (browse, save)             │
└─────────────────────────────────────────────┘
```

### Rationale
- Simple architecture sufficient for ~100 image batches
- Go handles image processing efficiently
- Wails 2 provides native bindings Go ↔ JS
- Async goroutines prevent UI blocking during batch process

---

## Tech Stack

| Component | Technology | Reason |
|-----------|------------|--------|
| Framework | Wails 2 | Stable, production-ready |
| Backend | Go 1.21+ | Fast image processing |
| Frontend | React + TypeScript | Modern, type-safe |
| Styling | TailwindCSS | Utility-first, rapid dev |
| Image resize | `disintegration/imaging` | Mature, fast |
| Text render | `fogleman/gg` | 2D graphics + freetype |
| Font loading | `golang/freetype` | TTF/OTF support |

### Bundle Fonts
1. **Be Vietnam Pro** - Native Vietnamese glyphs
2. **Roboto** - Clean, modern fallback
3. **SF Pro** (optional) - Premium look

---

## Project Structure

```
vibe-imageborder/
├── main.go                    # Entry point
├── app.go                     # Wails bindings
├── internal/
│   ├── image/
│   │   ├── compositor.go      # Image compositing
│   │   └── service.go         # Load, save, resize
│   ├── template/
│   │   ├── parser.go          # JSON parsing
│   │   └── service.go         # Field extraction
│   └── models/
│       └── types.go           # Shared types
├── frontend/
│   ├── src/
│   │   ├── App.tsx
│   │   └── components/
│   │       ├── FilePicker.tsx
│   │       ├── TemplateFields.tsx
│   │       ├── Preview.tsx
│   │       └── ProgressBar.tsx
├── assets/
│   └── fonts/
│       ├── BeVietnamPro-Regular.ttf
│       └── Roboto-Regular.ttf
└── build/
```

---

## UI Layout (2-Column)

```
┌────────────────────────────┬────────────────────────────────┐
│       COLUMN 1 (40%)       │         COLUMN 2 (60%)         │
├────────────────────────────┼────────────────────────────────┤
│  📁 Product Images         │  Preview                       │
│  [Drop/Browse]             │  [Image Preview Area]          │
│                            │  [Preview First Image]         │
│  🖼️ Frame Image            │                                │
│  [Select frame]            │  Output Settings               │
│                            │  Format: [PNG/JPG/WebP]        │
│  📄 Template (Optional)    │  Quality: [slider]             │
│  [Browse .txt]             │  Output: [folder path]         │
│                            │                                │
│  Text Fields               │  Progress                      │
│  (dynamic, show when       │  [████████░░░░] 67%            │
│   template loaded)         │                                │
│  - Barcode: [input]        │  [Generate All Button]         │
│  - Price: [input]          │                                │
│  - Size: [input]           │                                │
└────────────────────────────┴────────────────────────────────┘
```

---

## Core Features (MVP)

1. **File Selection**
   - Multi-select products (drag & drop + browse)
   - Single frame selection
   - Template .txt selection

2. **Template Parsing**
   - Auto-detect `[field_name]` from JSON text property
   - Generate dynamic input form

3. **Image Compositing**
   - Resize product to fit frame dimensions
   - Overlay frame with alpha transparency
   - Render text at JSON-defined positions

4. **Preview**
   - Single image preview before batch
   - Basic zoom controls

5. **Batch Export**
   - Progress bar with percentage
   - Format/quality selection
   - Output folder selection

---

## Template Format (JSON)

```json
{
  "background": "#f1eeea",
  "barcode": {
    "text": "[barcode]",
    "position": "90,1852",
    "fontsize": "50",
    "color": "white"
  },
  "price": {
    "text": "Giá [price]K",
    "position": "10,1712",
    "fontsize": "50",
    "color": "white"
  },
  "size": {
    "text": "D[size_dai] x R[size_rong] x C[size_cao] CM",
    "position": "1100,10",
    "fontsize": "60",
    "color": "white"
  }
}
```

Fields extracted: `barcode`, `price`, `size_dai`, `size_rong`, `size_cao`

---

## Processing Workflow

```
1. User selects product images (multi)
2. User selects frame image
3. User selects template .txt (optional)
   → Backend parses → returns field list
   → Frontend renders dynamic form
4. User fills in field values
5. User clicks "Preview First"
   → Backend composites first image
   → Returns base64 → Frontend displays
6. User clicks "Generate All"
   → Backend processes each image sequentially
   → Emits progress events → Frontend updates bar
7. Complete → Files saved to output folder
```

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Vietnamese text rendering issues | High | Bundle fonts with full glyph support, test early |
| Large batch memory usage | Medium | Process sequentially, dispose after each |
| Template format variations | Medium | Strict JSON validation, clear error messages |
| UI freeze during processing | Medium | Async goroutines + progress events |

---

## Success Criteria

1. Process 100 images in < 30 seconds
2. Vietnamese diacritics render correctly
3. Output quality matches input resolution
4. UI remains responsive during batch processing

---

## Next Steps

- [ ] Create detailed implementation plan with phases
- [ ] Setup Wails 2 project structure
- [ ] Implement template parser first (critical path)
- [ ] Add image compositing service
- [ ] Build React frontend components
- [ ] Integration testing with real templates

---

## Unresolved Questions

None at this time. All major decisions agreed upon.
