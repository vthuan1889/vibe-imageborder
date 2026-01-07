# Phase 5: Wails Backend Integration

**Goal:** Expose services qua Wails bindings, implement file dialogs, progress events

**Duration:** ~2-3 hours

**Dependencies:** Phase 1-4 complete

---

## Overview

Integrate backend services với Wails framework:
1. Create `App` struct với service dependencies
2. Bind methods for file selection (dialogs)
3. Bind processing methods
4. Implement progress events
5. Error handling và validation

---

## Task 5.1: Update App Structure

Modify `app.go`:

```go
package main

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
	"vibe-imageborder/internal/image"
	"vibe-imageborder/internal/template"
)

//go:embed assets/fonts/*
var fontsFS embed.FS

// App struct
type App struct {
	ctx         context.Context
	app         *application.App
	imageSvc    *image.Service
	templateSvc *template.Service
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		imageSvc:    image.NewService(),
		templateSvc: template.NewService(),
	}
}

// Startup is called when the app starts
func (a *App) Startup(ctx context.Context, app *application.App) {
	a.ctx = ctx
	a.app = app
}

// GetFontPath returns embedded font path
func GetFontPath() string {
	return "assets/fonts/Roboto-Regular.ttf"
}
```

---

## Task 5.2: Implement File Selection Methods

Add to `app.go`:

```go
import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// FileSelectionResult contains selected file paths
type FileSelectionResult struct {
	Paths    []string `json:"paths"`
	Success  bool     `json:"success"`
	Error    string   `json:"error"`
}

// SelectProductImages opens multi-select file dialog for product images
func (a *App) SelectProductImages() FileSelectionResult {
	result := FileSelectionResult{Success: false}

	options := application.OpenDialogOptions{
		Title: "Select Product Images",
		Filters: []application.FileFilter{
			{
				DisplayName: "Images (*.jpg, *.png)",
				Pattern:     "*.jpg;*.jpeg;*.png",
			},
		},
		CanChooseFiles:      true,
		CanChooseDirectories: false,
		AllowsMultipleSelection: true,
	}

	paths, err := application.OpenFileDialog(a.app, options)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	if len(paths) == 0 {
		result.Error = "No files selected"
		return result
	}

	result.Paths = paths
	result.Success = true
	return result
}

// SelectFrameImage opens single-select file dialog for frame image
func (a *App) SelectFrameImage() FileSelectionResult {
	result := FileSelectionResult{Success: false}

	options := application.OpenDialogOptions{
		Title: "Select Frame Image",
		Filters: []application.FileFilter{
			{
				DisplayName: "Images (*.png)",
				Pattern:     "*.png",
			},
		},
		CanChooseFiles:      true,
		CanChooseDirectories: false,
		AllowsMultipleSelection: false,
	}

	paths, err := application.OpenFileDialog(a.app, options)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	if len(paths) == 0 {
		result.Error = "No file selected"
		return result
	}

	result.Paths = []string{paths[0]}
	result.Success = true
	return result
}

// SelectTemplateFile opens file dialog for template file
func (a *App) SelectTemplateFile() FileSelectionResult {
	result := FileSelectionResult{Success: false}

	options := application.OpenDialogOptions{
		Title: "Select Template File",
		Filters: []application.FileFilter{
			{
				DisplayName: "Template Files (*.txt)",
				Pattern:     "*.txt",
			},
		},
		CanChooseFiles:      true,
		CanChooseDirectories: false,
		AllowsMultipleSelection: false,
	}

	paths, err := application.OpenFileDialog(a.app, options)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	if len(paths) == 0 {
		result.Error = "No file selected"
		return result
	}

	result.Paths = []string{paths[0]}
	result.Success = true
	return result
}

// SelectOutputDirectory opens directory picker
func (a *App) SelectOutputDirectory() FileSelectionResult {
	result := FileSelectionResult{Success: false}

	options := application.OpenDialogOptions{
		Title:                "Select Output Directory",
		CanChooseFiles:       false,
		CanChooseDirectories: true,
		AllowsMultipleSelection: false,
	}

	paths, err := application.OpenFileDialog(a.app, options)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	if len(paths) == 0 {
		result.Error = "No directory selected"
		return result
	}

	result.Paths = []string{paths[0]}
	result.Success = true
	return result
}
```

---

## Task 5.3: Implement Template Parsing Method

Add to `app.go`:

```go
import (
	"vibe-imageborder/internal/models"
)

// TemplateInfo contains template và extracted fields
type TemplateInfo struct {
	Template models.Template `json:"template"`
	Fields   []string        `json:"fields"`
	Success  bool            `json:"success"`
	Error    string          `json:"error"`
}

// ParseTemplateFile loads template và extracts dynamic fields
func (a *App) ParseTemplateFile(path string) TemplateInfo {
	result := TemplateInfo{Success: false}

	// Load template
	tmpl, err := a.templateSvc.Load(path)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to load template: %v", err)
		return result
	}

	// Extract fields
	fields := a.templateSvc.GetDynamicFields(tmpl)

	result.Template = tmpl
	result.Fields = fields
	result.Success = true
	return result
}
```

---

## Task 5.4: Implement Batch Processing Method

Add to `app.go`:

```go
import (
	"log"
	"time"
)

// ProcessRequest contains all processing parameters
type ProcessRequest struct {
	ProductPaths []string          `json:"productPaths"`
	FramePath    string            `json:"framePath"`
	Template     models.Template   `json:"template"`
	FieldValues  map[string]string `json:"fieldValues"`
	OutputDir    string            `json:"outputDir"`
}

// ProcessResult contains processing outcome
type ProcessResult struct {
	Success       bool     `json:"success"`
	ProcessedCount int     `json:"processedCount"`
	FailedCount   int      `json:"failedCount"`
	OutputPaths   []string `json:"outputPaths"`
	Error         string   `json:"error"`
}

// ProgressUpdate sent during processing
type ProgressUpdate struct {
	Current  int    `json:"current"`
	Total    int    `json:"total"`
	Filename string `json:"filename"`
	Status   string `json:"status"` // "processing", "success", "error"
}

// ProcessBatch processes all product images với template
func (a *App) ProcessBatch(req ProcessRequest) ProcessResult {
	result := ProcessResult{Success: false}

	total := len(req.ProductPaths)
	if total == 0 {
		result.Error = "No product images selected"
		return result
	}

	log.Printf("Processing %d images...", total)

	var outputPaths []string
	var failedCount int

	for i, productPath := range req.ProductPaths {
		filename := filepath.Base(productPath)

		// Emit progress
		a.emitProgress(ProgressUpdate{
			Current:  i + 1,
			Total:    total,
			Filename: filename,
			Status:   "processing",
		})

		// Generate output path
		outputFilename := fmt.Sprintf("%s_framed.png",
			filename[:len(filename)-len(filepath.Ext(filename))])
		outputPath := filepath.Join(req.OutputDir, outputFilename)

		// Process single image
		err := a.imageSvc.ProcessSingle(
			productPath,
			req.FramePath,
			outputPath,
			req.Template,
			req.FieldValues,
		)

		if err != nil {
			log.Printf("Error processing %s: %v", filename, err)
			failedCount++

			a.emitProgress(ProgressUpdate{
				Current:  i + 1,
				Total:    total,
				Filename: filename,
				Status:   "error",
			})
		} else {
			outputPaths = append(outputPaths, outputPath)

			a.emitProgress(ProgressUpdate{
				Current:  i + 1,
				Total:    total,
				Filename: filename,
				Status:   "success",
			})
		}

		// Small delay to prevent UI freeze
		time.Sleep(10 * time.Millisecond)
	}

	result.ProcessedCount = len(outputPaths)
	result.FailedCount = failedCount
	result.OutputPaths = outputPaths
	result.Success = true

	return result
}

// emitProgress sends progress update to frontend
func (a *App) emitProgress(update ProgressUpdate) {
	a.app.EmitEvent("processing:progress", update)
}
```

---

## Task 5.5: Update main.go Bindings

Modify `main.go`:

```go
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	// Create application
	app := application.New(application.Options{
		Name:        "vibe-imageborder",
		Description: "Image Border Application",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create app instance
	appInstance := NewApp()

	// Bind methods
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Image Border App",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
		Bind: []interface{}{
			appInstance.SelectProductImages,
			appInstance.SelectFrameImage,
			appInstance.SelectTemplateFile,
			appInstance.SelectOutputDirectory,
			appInstance.ParseTemplateFile,
			appInstance.ProcessBatch,
		},
	})

	// Run application
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
```

---

## Task 5.6: Test Backend Methods

Create test file `test-backend.go`:

```go
package main

import (
	"encoding/json"
	"fmt"
	"vibe-imageborder/internal/models"
)

func main() {
	app := NewApp()

	// Test template parsing
	tmplInfo := app.ParseTemplateFile("tests/fixtures/templates/khung-002-05.txt")
	if !tmplInfo.Success {
		fmt.Printf("Template parse failed: %s\n", tmplInfo.Error)
		return
	}

	fmt.Printf("Template loaded successfully\n")
	fmt.Printf("Fields: %v\n", tmplInfo.Fields)

	// Test processing
	req := ProcessRequest{
		ProductPaths: []string{
			"tests/fixtures/products/product-01.jpg",
		},
		FramePath: "tests/fixtures/frames/frame-01.png",
		Template:  tmplInfo.Template,
		FieldValues: map[string]string{
			"barcode":   "TEST123",
			"size_dai":  "30",
			"size_rong": "20",
			"size_cao":  "15",
		},
		OutputDir: "tests/output",
	}

	result := app.ProcessBatch(req)
	if !result.Success {
		fmt.Printf("Processing failed: %s\n", result.Error)
		return
	}

	fmt.Printf("✓ Processed: %d images\n", result.ProcessedCount)
	fmt.Printf("✓ Failed: %d images\n", result.FailedCount)
	fmt.Printf("✓ Outputs: %v\n", result.OutputPaths)
}
```

Run test:

```bash
go run test-backend.go
```

---

## Acceptance Criteria

- ✓ File dialogs open correctly (products, frame, template, output dir)
- ✓ `ParseTemplateFile()` returns template + fields
- ✓ `ProcessBatch()` processes multiple images
- ✓ Progress events emit during processing
- ✓ Error handling robust (invalid paths, missing files)
- ✓ Backend test program passes

---

## Deliverables

### Files Created/Modified

1. `app.go` - App struct với all bindings
2. `main.go` - Wails app setup với bindings
3. `test-backend.go` - Backend testing script

### Validation

```bash
# 1. Test backend methods
go run test-backend.go

# 2. Build app
wails3 build

# 3. Run app (manual testing)
wails3 dev
# - Click file pickers (should open dialogs)
# - Select template (should show fields in console)
```

---

## Next Phase

[Phase 6: React Frontend - UI Components](phase-06-react-frontend.md)
