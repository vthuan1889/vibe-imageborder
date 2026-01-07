# Phase 5: Wails Backend Integration

## Context

- Plan: [plan.md](./plan.md)
- Previous: [Phase 4 - Text Rendering](./phase-04-text-rendering.md)

## Overview

| Field | Value |
|-------|-------|
| Priority | P1 - Critical Path |
| Status | ⚠️ Complete with Issues |
| Effort | 2h (actual) |

Expose Go services to frontend via Wails bindings. Implement file dialogs, template loading, preview generation, and batch processing with progress events.

## Requirements

### Functional
- File/folder selection dialogs
- Load template and return field list
- Generate single preview image (base64)
- Batch process with progress events
- Cancel batch processing

### Non-functional
- Async batch processing (non-blocking UI)
- Progress updates every image
- Error handling with user-friendly messages

## Architecture

```
Wails Bindings (app.go)
├── SelectProductFiles()     → []string
├── SelectFrameFile()        → string
├── SelectTemplateFile()     → string
├── SelectOutputFolder()     → string
├── LoadTemplate()           → []string (field names)
├── GeneratePreview()        → string (base64 PNG)
├── ProcessBatch()           → void (emits progress)
└── CancelProcessing()       → void

Events:
├── progress → {current, total, file}
├── complete → {outputDir}
└── error → {message}
```

## Related Code Files

### Files to Modify
| File | Change |
|------|--------|
| `app.go` | Add all Wails-exposed methods |
| `main.go` | Initialize services, pass to App |

## Implementation Steps

### Step 1: Update App struct

```go
// app.go
package main

import (
    "context"
    "embed"
    "encoding/base64"
    "fmt"
    "path/filepath"
    "sync"

    "github.com/wailsapp/wails/v2/pkg/runtime"

    imgservice "vibe-imageborder/internal/image"
    "vibe-imageborder/internal/models"
    "vibe-imageborder/internal/template"
)

// App struct
type App struct {
    ctx            context.Context
    fonts          embed.FS
    templateSvc    *template.Service
    imageSvc       *imgservice.Service
    compositor     *imgservice.Compositor
    fontManager    *imgservice.FontManager
    textRenderer   *imgservice.TextRenderer
    cancelFunc     context.CancelFunc
    processingLock sync.Mutex
}

// NewApp creates new App
func NewApp(fonts embed.FS) *App {
    imageSvc := imgservice.NewService()
    fontManager := imgservice.NewFontManager(fonts)

    return &App{
        fonts:        fonts,
        templateSvc:  template.NewService(),
        imageSvc:     imageSvc,
        compositor:   imgservice.NewCompositor(imageSvc),
        fontManager:  fontManager,
        textRenderer: imgservice.NewTextRenderer(fontManager),
    }
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
}
```

### Step 2: File Selection Methods

```go
// app.go - File selection methods

// SelectProductFiles opens multi-file dialog for product images
func (a *App) SelectProductFiles() ([]string, error) {
    files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
        Title: "Select Product Images",
        Filters: []runtime.FileFilter{
            {DisplayName: "Images", Pattern: "*.jpg;*.jpeg;*.png;*.webp"},
        },
    })
    if err != nil {
        return nil, err
    }
    return files, nil
}

// SelectFrameFile opens single file dialog for frame image
func (a *App) SelectFrameFile() (string, error) {
    file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
        Title: "Select Frame Image",
        Filters: []runtime.FileFilter{
            {DisplayName: "Images", Pattern: "*.png;*.jpg;*.jpeg"},
        },
    })
    if err != nil {
        return "", err
    }
    return file, nil
}

// SelectTemplateFile opens dialog for template .txt file
func (a *App) SelectTemplateFile() (string, error) {
    file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
        Title: "Select Template File",
        Filters: []runtime.FileFilter{
            {DisplayName: "Template", Pattern: "*.txt"},
        },
    })
    if err != nil {
        return "", err
    }
    return file, nil
}

// SelectOutputFolder opens folder selection dialog
func (a *App) SelectOutputFolder() (string, error) {
    folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
        Title: "Select Output Folder",
    })
    if err != nil {
        return "", err
    }
    return folder, nil
}
```

### Step 3: Template Methods

```go
// app.go - Template methods

// LoadTemplate loads template and returns field names
func (a *App) LoadTemplate(path string) ([]string, error) {
    if path == "" {
        return []string{}, nil
    }

    fields, err := a.templateSvc.GetFields(path)
    if err != nil {
        return nil, fmt.Errorf("failed to load template: %w", err)
    }
    return fields, nil
}

// GetTemplateBackground returns background color from template
func (a *App) GetTemplateBackground(path string) (string, error) {
    return a.templateSvc.GetBackground(path)
}
```

### Step 4: Preview Method

```go
// app.go - Preview generation

// GeneratePreview creates preview of first image
func (a *App) GeneratePreview(req models.ProcessRequest) (string, error) {
    if len(req.ProductImages) == 0 {
        return "", fmt.Errorf("no product images selected")
    }
    if req.FrameImage == "" {
        return "", fmt.Errorf("no frame image selected")
    }

    // Load images
    product, err := a.imageSvc.LoadImage(req.ProductImages[0])
    if err != nil {
        return "", fmt.Errorf("failed to load product: %w", err)
    }

    frame, err := a.imageSvc.LoadImage(req.FrameImage)
    if err != nil {
        return "", fmt.Errorf("failed to load frame: %w", err)
    }

    // Get background and overlays
    bgColor := ""
    var overlays map[string]models.TextOverlay

    if req.TemplatePath != "" {
        bgColor, _ = a.templateSvc.GetBackground(req.TemplatePath)
        overlays, _ = a.templateSvc.GetOverlays(req.TemplatePath, req.FieldValues)
    }

    // Composite
    result, err := a.compositor.CompositeWithText(product, frame, bgColor, overlays, a.textRenderer)
    if err != nil {
        return "", fmt.Errorf("failed to composite: %w", err)
    }

    // Encode to base64 PNG
    var buf bytes.Buffer
    if err := png.Encode(&buf, result.Image); err != nil {
        return "", fmt.Errorf("failed to encode: %w", err)
    }

    return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
```

### Step 5: Batch Processing

```go
// app.go - Batch processing

// ProcessBatch processes all images with progress events
func (a *App) ProcessBatch(req models.ProcessRequest) error {
    a.processingLock.Lock()

    // Create cancellable context
    ctx, cancel := context.WithCancel(a.ctx)
    a.cancelFunc = cancel
    a.processingLock.Unlock()

    defer func() {
        a.processingLock.Lock()
        a.cancelFunc = nil
        a.processingLock.Unlock()
    }()

    // Load frame once
    frame, err := a.imageSvc.LoadImage(req.FrameImage)
    if err != nil {
        runtime.EventsEmit(a.ctx, "error", map[string]string{"message": err.Error()})
        return err
    }

    // Get template data
    bgColor := ""
    var overlays map[string]models.TextOverlay
    if req.TemplatePath != "" {
        bgColor, _ = a.templateSvc.GetBackground(req.TemplatePath)
        overlays, _ = a.templateSvc.GetOverlays(req.TemplatePath, req.FieldValues)
    }

    total := len(req.ProductImages)

    for i, productPath := range req.ProductImages {
        // Check for cancellation
        select {
        case <-ctx.Done():
            runtime.EventsEmit(a.ctx, "cancelled", nil)
            return nil
        default:
        }

        // Process single image
        err := a.processSingleImage(productPath, frame, bgColor, overlays, req)

        // Emit progress
        progress := models.ProcessProgress{
            Current: i + 1,
            Total:   total,
            File:    filepath.Base(productPath),
        }
        runtime.EventsEmit(a.ctx, "progress", progress)

        if err != nil {
            // Log error but continue
            fmt.Printf("Error processing %s: %v\n", productPath, err)
        }
    }

    runtime.EventsEmit(a.ctx, "complete", map[string]string{"outputDir": req.OutputDir})
    return nil
}

func (a *App) processSingleImage(
    productPath string,
    frame image.Image,
    bgColor string,
    overlays map[string]models.TextOverlay,
    req models.ProcessRequest,
) error {
    product, err := a.imageSvc.LoadImage(productPath)
    if err != nil {
        return err
    }

    result, err := a.compositor.CompositeWithText(product, frame, bgColor, overlays, a.textRenderer)
    if err != nil {
        return err
    }

    // Generate output filename
    baseName := filepath.Base(productPath)
    ext := filepath.Ext(baseName)
    nameWithoutExt := baseName[:len(baseName)-len(ext)]
    outputName := nameWithoutExt + "_framed"
    outputPath := filepath.Join(req.OutputDir, outputName)

    return a.imageSvc.SaveImage(result.Image, outputPath, req.Format, req.Quality)
}

// CancelProcessing cancels ongoing batch processing
func (a *App) CancelProcessing() {
    a.processingLock.Lock()
    defer a.processingLock.Unlock()

    if a.cancelFunc != nil {
        a.cancelFunc()
    }
}
```

### Step 6: Update main.go

```go
// main.go
package main

import (
    "embed"

    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed assets/fonts/*
var fonts embed.FS

func main() {
    app := NewApp(fonts)

    err := wails.Run(&options.App{
        Title:  "Image Border Tool",
        Width:  1200,
        Height: 800,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
        OnStartup:        app.startup,
        Bind: []interface{}{
            app,
        },
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
```

## Todo List

- [x] Update App struct with services
- [x] Implement file selection dialogs
- [x] Implement template loading
- [x] Implement preview generation
- [x] Implement batch processing with events
- [x] Implement cancel functionality
- [x] Update main.go with initialization
- [ ] Test all Wails bindings (requires running app)
- [ ] **FIX SECURITY ISSUES** (see review report)

## Security Issues Found (Code Review 260108-0542)

**Must fix before Phase 6:**

1. **[CRITICAL]** Path traversal vulnerability - no path validation in file selection methods
2. **[CRITICAL]** Incomplete batch size validation - missing zero-length and duplicate checks
3. **[CRITICAL]** Race condition in cancel logic - needs isProcessing flag
4. **[CRITICAL]** Silent error swallowing in template operations
5. **[HIGH]** Output path collision - files may be overwritten

**Full report:** `plans/reports/code-reviewer-260108-0542-phase5-wails-backend.md`

## Success Criteria

1. File dialogs work on Windows
2. Template loading returns correct fields
3. Preview generates valid base64 image
4. Progress events emit correctly
5. Cancel stops processing immediately

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Dialog not working | High | Test on Windows early |
| Memory leak in batch | Medium | Dispose images after each |
| Race condition | Medium | Use mutex for cancel |

## Security Considerations

- Validate file paths before processing
- Limit max batch size (e.g., 1000 images)
- Sanitize event data

## Next Steps

After completion, proceed to [Phase 6: React Frontend](./phase-06-react-frontend.md)
