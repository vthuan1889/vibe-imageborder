package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	imgservice "vibe-imageborder/internal/image"
	"vibe-imageborder/internal/models"
	"vibe-imageborder/internal/template"
	"vibe-imageborder/internal/updater"
)

// MaxBatchSize limits the number of images that can be processed in one batch.
const MaxBatchSize = 1000

// Event name constants.
const (
	EventProgress  = "progress"
	EventComplete  = "complete"
	EventError     = "error"
	EventCancelled = "cancelled"
)

// App struct holds the application state and services.
type App struct {
	ctx            context.Context
	fonts          embed.FS
	templateSvc    *template.Service
	imageSvc       *imgservice.Service
	compositor     *imgservice.Compositor
	fontManager    *imgservice.FontManager
	textRenderer   *imgservice.TextRenderer
	cancelFunc     context.CancelFunc
	isProcessing   bool
	processingLock sync.Mutex
}

// NewApp creates a new App with all services initialized.
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

// startup is called when the app starts. The context is saved for runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// validatePath validates and cleans a file path.
func validatePath(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	cleanPath := filepath.Clean(path)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}
	return absPath, nil
}

// validatePaths validates multiple paths and returns cleaned paths.
func validatePaths(paths []string) ([]string, error) {
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		clean, err := validatePath(path)
		if err != nil {
			return nil, err
		}
		if clean != "" {
			result = append(result, clean)
		}
	}
	return result, nil
}

// sanitizeError removes sensitive information from error messages.
func sanitizeError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	// Remove absolute paths - just show filename
	if strings.Contains(msg, string(filepath.Separator)) {
		parts := strings.Split(msg, string(filepath.Separator))
		msg = parts[len(parts)-1]
	}
	// Limit length
	if len(msg) > 200 {
		msg = msg[:200] + "..."
	}
	return msg
}

// SelectProductFiles opens multi-file dialog for product images.
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
	return validatePaths(files)
}

// SelectFrameFile opens single file dialog for frame image.
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
	return validatePath(file)
}

// SelectTemplateFile opens dialog for template .txt file.
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
	return validatePath(file)
}

// SelectOutputFolder opens folder selection dialog.
func (a *App) SelectOutputFolder() (string, error) {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Output Folder",
	})
	if err != nil {
		return "", err
	}
	return validatePath(folder)
}

// GetDefaultOutputFolder returns the user's Downloads folder as default output.
func (a *App) GetDefaultOutputFolder() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	downloads := filepath.Join(home, "Downloads")
	if _, err := os.Stat(downloads); err == nil {
		return downloads
	}
	return home
}

// LoadTemplate loads template and returns field names.
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

// GetTemplateBackground returns background color from template.
func (a *App) GetTemplateBackground(path string) (string, error) {
	return a.templateSvc.GetBackground(path)
}

// GeneratePreview creates preview of first image.
func (a *App) GeneratePreview(req models.ProcessRequest) (string, error) {
	if len(req.ProductImages) == 0 {
		return "", fmt.Errorf("no product images selected")
	}
	if req.FrameImage == "" {
		return "", fmt.Errorf("no frame image selected")
	}

	// Validate format
	validFormats := map[string]bool{"png": true, "jpg": true, "jpeg": true, "webp": true, "": true}
	if !validFormats[strings.ToLower(req.Format)] {
		return "", fmt.Errorf("invalid format: %s", req.Format)
	}

	// Validate quality
	if req.Quality < 0 || req.Quality > 100 {
		req.Quality = 90
	}

	// Validate template path exists if provided
	if req.TemplatePath != "" {
		if _, err := os.Stat(req.TemplatePath); err != nil {
			return "", fmt.Errorf("template file not found: %w", err)
		}
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

	// Get background and overlays with proper error handling
	var bgColor string
	var overlays map[string]models.TextOverlay

	if req.TemplatePath != "" {
		bgColor, err = a.templateSvc.GetBackground(req.TemplatePath)
		if err != nil {
			return "", fmt.Errorf("failed to get template background: %w", err)
		}

		overlays, err = a.templateSvc.GetOverlays(req.TemplatePath, req.FieldValues)
		if err != nil {
			return "", fmt.Errorf("failed to get template overlays: %w", err)
		}
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

// ProcessBatch processes all images with progress events.
func (a *App) ProcessBatch(req models.ProcessRequest) error {
	// Validate batch size
	if len(req.ProductImages) == 0 {
		return fmt.Errorf("no images to process")
	}
	if len(req.ProductImages) > MaxBatchSize {
		return fmt.Errorf("batch size %d exceeds maximum %d", len(req.ProductImages), MaxBatchSize)
	}

	// Deduplicate paths
	seen := make(map[string]bool)
	for _, path := range req.ProductImages {
		if seen[path] {
			return fmt.Errorf("duplicate path detected: %s", filepath.Base(path))
		}
		seen[path] = true
	}

	// Validate required fields
	if req.FrameImage == "" {
		return fmt.Errorf("frame image required")
	}
	if req.OutputDir == "" {
		return fmt.Errorf("output directory required")
	}

	// Check if already processing
	a.processingLock.Lock()
	if a.isProcessing {
		a.processingLock.Unlock()
		return fmt.Errorf("another batch is already processing")
	}
	a.isProcessing = true

	// Create cancellable context
	ctx, cancel := context.WithCancel(a.ctx)
	a.cancelFunc = cancel
	a.processingLock.Unlock()

	defer func() {
		a.processingLock.Lock()
		a.cancelFunc = nil
		a.isProcessing = false
		a.processingLock.Unlock()
	}()

	// Load frame once
	frame, err := a.imageSvc.LoadImage(req.FrameImage)
	if err != nil {
		runtime.EventsEmit(a.ctx, EventError, map[string]string{"message": sanitizeError(err)})
		return err
	}

	// Get template data with proper error handling
	var bgColor string
	var overlays map[string]models.TextOverlay
	if req.TemplatePath != "" {
		bgColor, err = a.templateSvc.GetBackground(req.TemplatePath)
		if err != nil {
			runtime.EventsEmit(a.ctx, EventError, map[string]string{"message": sanitizeError(err)})
			return fmt.Errorf("failed to get template background: %w", err)
		}
		overlays, err = a.templateSvc.GetOverlays(req.TemplatePath, req.FieldValues)
		if err != nil {
			runtime.EventsEmit(a.ctx, EventError, map[string]string{"message": sanitizeError(err)})
			return fmt.Errorf("failed to get template overlays: %w", err)
		}
	}

	total := len(req.ProductImages)
	var failures []string

	for i, productPath := range req.ProductImages {
		// Check for cancellation
		select {
		case <-ctx.Done():
			runtime.EventsEmit(a.ctx, EventCancelled, nil)
			return nil
		default:
		}

		// Process single image
		err := a.processSingleImage(productPath, frame, bgColor, overlays, req)

		success := err == nil
		if !success {
			failures = append(failures, filepath.Base(productPath))
			fmt.Printf("Error processing %s: %v\n", productPath, err)
		}

		// Emit progress
		progress := models.ProcessProgress{
			Current: i + 1,
			Total:   total,
			File:    filepath.Base(productPath),
			Success: success,
		}
		runtime.EventsEmit(a.ctx, EventProgress, progress)
	}

	// Emit completion with summary
	result := map[string]interface{}{
		"outputDir":      req.OutputDir,
		"totalProcessed": total - len(failures),
		"totalFailed":    len(failures),
		"failures":       failures,
	}
	runtime.EventsEmit(a.ctx, EventComplete, result)
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

	// Determine output format
	outputFormat := req.Format
	if outputFormat == "" {
		outputFormat = "png"
	}

	// Generate unique output path to avoid collisions
	outputPath := filepath.Join(req.OutputDir, outputName)
	counter := 0
	for {
		testPath := outputPath
		if counter > 0 {
			testPath = fmt.Sprintf("%s_%d", outputPath, counter)
		}
		fullPath := testPath + "." + outputFormat
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			outputPath = testPath
			break
		}
		counter++
		if counter > 1000 {
			return fmt.Errorf("too many duplicate files for %s", baseName)
		}
	}

	return a.imageSvc.SaveImage(result.Image, outputPath, req.Format, req.Quality)
}

// CopyImages copies all images from source to destination recursively,
// applying imperceptible modifications so each file has a unique fingerprint.
func (a *App) CopyImages(req models.CopyImagesRequest) error {
	sourceDir, err := validatePath(req.SourceDir)
	if err != nil {
		return fmt.Errorf("invalid source directory: %w", err)
	}
	if sourceDir == "" {
		return fmt.Errorf("source directory required")
	}

	destDir, err := validatePath(req.DestDir)
	if err != nil {
		return fmt.Errorf("invalid destination directory: %w", err)
	}
	if destDir == "" {
		return fmt.Errorf("destination directory required")
	}

	if sourceDir == destDir {
		return fmt.Errorf("source and destination folders must be different")
	}

	files, err := imgservice.CollectImageFiles(sourceDir, MaxBatchSize)
	if err != nil {
		return err
	}

	a.processingLock.Lock()
	if a.isProcessing {
		a.processingLock.Unlock()
		return fmt.Errorf("another batch is already processing")
	}
	a.isProcessing = true

	ctx, cancel := context.WithCancel(a.ctx)
	a.cancelFunc = cancel
	a.processingLock.Unlock()

	defer func() {
		a.processingLock.Lock()
		a.cancelFunc = nil
		a.isProcessing = false
		a.processingLock.Unlock()
	}()

	total := len(files)
	var failures []string

	for i, srcPath := range files {
		select {
		case <-ctx.Done():
			runtime.EventsEmit(a.ctx, EventCancelled, nil)
			return nil
		default:
		}

		err := a.copySingleImage(srcPath, sourceDir, destDir)
		success := err == nil
		if !success {
			failures = append(failures, filepath.Base(srcPath))
			fmt.Printf("Error copying %s: %v\n", srcPath, err)
		}

		progress := models.ProcessProgress{
			Current: i + 1,
			Total:   total,
			File:    filepath.Base(srcPath),
			Success: success,
		}
		runtime.EventsEmit(a.ctx, EventProgress, progress)
	}

	result := map[string]interface{}{
		"outputDir":      destDir,
		"totalProcessed": total - len(failures),
		"totalFailed":    len(failures),
		"failures":       failures,
	}
	runtime.EventsEmit(a.ctx, EventComplete, result)
	return nil
}

func (a *App) copySingleImage(srcPath, sourceDir, destDir string) error {
	rel, err := filepath.Rel(sourceDir, srcPath)
	if err != nil {
		return fmt.Errorf("failed to resolve relative path: %w", err)
	}

	ext := filepath.Ext(srcPath)
	format, quality, outputExt := imgservice.FormatFromExtension(ext)

	destRel := strings.TrimSuffix(rel, ext) + outputExt
	destPath := filepath.Join(destDir, destRel)

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	img, err := a.imageSvc.LoadImage(srcPath)
	if err != nil {
		return err
	}

	seed := uniquifySeed(srcPath)
	uniqueImg := imgservice.Uniquify(img, seed)

	outputBase := strings.TrimSuffix(destPath, outputExt)
	return a.imageSvc.SaveImage(uniqueImg, outputBase, format, quality)
}

func uniquifySeed(path string) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(path))
	return int64(h.Sum64()) ^ time.Now().UnixNano()
}

// CancelProcessing cancels ongoing batch processing.
func (a *App) CancelProcessing() {
	a.processingLock.Lock()
	defer a.processingLock.Unlock()

	if a.cancelFunc != nil {
		a.cancelFunc()
	}
}

// GetVersion returns the current app version.
func (a *App) GetVersion() string {
	return version
}

// CheckForUpdate queries GitHub for available updates.
func (a *App) CheckForUpdate() (*updater.UpdateInfo, error) {
	return updater.CheckUpdate(version)
}

// DownloadAndInstallUpdate downloads and runs the installer, then quits the app.
func (a *App) DownloadAndInstallUpdate(downloadURL string) error {
	if err := updater.DownloadAndInstall(downloadURL); err != nil {
		return err
	}
	// Quit app after starting installer
	runtime.Quit(a.ctx)
	return nil
}
