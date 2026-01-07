package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"vibe-imageborder/internal/image"
	"vibe-imageborder/internal/models"
	"vibe-imageborder/internal/template"
)

// App struct
type App struct {
	imageSvc    *image.Service
	templateSvc *template.Service
}

// NewApp creates a new App
func NewApp() *App {
	return &App{
		imageSvc:    image.NewService(),
		templateSvc: template.NewService(),
	}
}

// TemplateInfo contains template and extracted fields
type TemplateInfo struct {
	Template models.Template `json:"template"`
	Fields   []string        `json:"fields"`
	Success  bool            `json:"success"`
	Error    string          `json:"error"`
}

// ParseTemplateFile loads template and extracts fields
func (a *App) ParseTemplateFile(path string) TemplateInfo {
	result := TemplateInfo{Success: false}

	tmpl, err := a.templateSvc.Load(path)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to load template: %v", err)
		return result
	}

	fields := a.templateSvc.GetDynamicFields(tmpl)

	result.Template = tmpl
	result.Fields = fields
	result.Success = true
	return result
}

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
	Success        bool     `json:"success"`
	ProcessedCount int      `json:"processedCount"`
	FailedCount    int      `json:"failedCount"`
	OutputPaths    []string `json:"outputPaths"`
	Error          string   `json:"error"`
}

// ProgressUpdate sent during processing
type ProgressUpdate struct {
	Current  int    `json:"current"`
	Total    int    `json:"total"`
	Filename string `json:"filename"`
	Status   string `json:"status"` // "processing", "success", "error"
}

// ProcessBatch processes all product images
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

		time.Sleep(10 * time.Millisecond)
	}

	result.ProcessedCount = len(outputPaths)
	result.FailedCount = failedCount
	result.OutputPaths = outputPaths
	result.Success = true

	return result
}

// emitProgress sends progress to frontend
// TODO: Implement with Wails v3 application.EmitEvent when API is stable
func (a *App) emitProgress(update ProgressUpdate) {
	// Placeholder - events will be implemented in Phase 7
	log.Printf("Progress: %d/%d - %s (%s)", update.Current, update.Total, update.Filename, update.Status)
}
