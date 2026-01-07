# Code Review: Phase 5 Wails Backend Implementation

## Scope

- Files reviewed: app.go, main.go
- Lines of code: ~275 lines
- Review focus: Recent changes (HEAD vs HEAD~1)
- Plan: plans/260108-0121-image-border-app/phase-05-wails-backend.md
- Build status: Passing
- Test status: Passing with race detector

## Overall Assessment

Implementation is **functionally complete** và đáp ứng tất cả yêu cầu của Phase 5. Code structure sạch, logic rõ ràng, nhưng có **5 critical security issues** và **2 high-priority performance risks** cần fix ngay.

## Critical Issues

### 1. Path Traversal Vulnerability (CRITICAL - Security)

**Location:** app.go:58-108, 111-120, 129-169, 173-236

**Problem:** File selection dialogs không validate paths. User có thể select files ngoài boundaries dự kiến hoặc paths malicious.

```go
// app.go:58-68 - NO PATH VALIDATION
func (a *App) SelectProductFiles() ([]string, error) {
    files, err := runtime.OpenMultipleFilesDialog(a.ctx, ...)
    if err != nil {
        return nil, err
    }
    return files, nil  // ❌ Returns user paths unchecked
}
```

**Impact:**
- Path traversal attacks (e.g., `../../system/file`)
- Access to sensitive files outside intended directories
- Potential code execution if combined with other vulnerabilities

**Fix:** Validate paths before returning:

```go
func (a *App) SelectProductFiles() ([]string, error) {
    files, err := runtime.OpenMultipleFilesDialog(a.ctx, ...)
    if err != nil {
        return nil, err
    }

    // Validate each path
    for _, path := range files {
        cleanPath := filepath.Clean(path)
        absPath, err := filepath.Abs(cleanPath)
        if err != nil {
            return nil, fmt.Errorf("invalid path: %w", err)
        }

        // Optional: Check if within allowed directory
        // if !strings.HasPrefix(absPath, allowedDir) {
        //     return nil, fmt.Errorf("path outside allowed directory")
        // }
    }

    return files, nil
}
```

Apply to all 4 selection methods.

### 2. Missing Batch Size Validation Implementation (CRITICAL - DoS)

**Location:** app.go:21-22, 173-177

**Problem:** `MaxBatchSize` constant được define nhưng validation logic thiếu kiểm tra comprehensive.

```go
const MaxBatchSize = 1000

func (a *App) ProcessBatch(req models.ProcessRequest) error {
    // Validate batch size
    if len(req.ProductImages) > MaxBatchSize {
        return fmt.Errorf("batch size exceeds maximum of %d images", MaxBatchSize)
    }
    // ❌ What if len == 0? What if paths are duplicates?
```

**Impact:**
- DoS through empty batch (wastes resources)
- DoS through duplicate paths (processing same file 1000x)
- Memory exhaustion if individual images are huge

**Fix:**

```go
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
            return fmt.Errorf("duplicate path detected: %s", path)
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

    // Continue with existing code...
```

### 3. Race Condition in Cancel Logic (CRITICAL - Concurrency)

**Location:** app.go:179-190, 267-274

**Problem:** Lock được release trước khi goroutine hoàn tất, tạo window cho race condition.

```go
func (a *App) ProcessBatch(req models.ProcessRequest) error {
    a.processingLock.Lock()
    ctx, cancel := context.WithCancel(a.ctx)
    a.cancelFunc = cancel
    a.processingLock.Unlock()  // ❌ Lock released too early

    defer func() {
        a.processingLock.Lock()
        a.cancelFunc = nil
        a.processingLock.Unlock()
    }()

    // Processing happens here without lock protection
    // Meanwhile, CancelProcessing() can be called...
```

**Scenario:**
1. ProcessBatch locks, sets cancelFunc, unlocks
2. CancelProcessing locks, calls cancelFunc(), unlocks
3. ProcessBatch's defer runs, sets cancelFunc = nil
4. Second ProcessBatch call starts → old context still active

**Impact:**
- Cancel calls may affect wrong batch
- Concurrent batches may interfere
- Undefined behavior if multiple ProcessBatch calls

**Fix:** Prevent concurrent processing:

```go
type App struct {
    // ... existing fields
    isProcessing bool  // Add this field
}

func (a *App) ProcessBatch(req models.ProcessRequest) error {
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

    // ... rest of code
}
```

### 4. Silent Error Swallowing (HIGH - Reliability)

**Location:** app.go:153-154, 203-204

**Problem:** Template errors được ignore với `_` instead of reporting to user.

```go
if req.TemplatePath != "" {
    bgColor, _ = a.templateSvc.GetBackground(req.TemplatePath)
    overlays, _ = a.templateSvc.GetOverlays(req.TemplatePath, req.FieldValues)
    // ❌ Errors silently ignored
}
```

**Impact:**
- Template parsing errors không được report
- User không biết tại sao text không hiển thị
- Debug khó vì không có error messages

**Fix:**

```go
var bgColor string
var overlays map[string]models.TextOverlay

if req.TemplatePath != "" {
    var err error
    bgColor, err = a.templateSvc.GetBackground(req.TemplatePath)
    if err != nil {
        return "", fmt.Errorf("failed to get template background: %w", err)
    }

    overlays, err = a.templateSvc.GetOverlays(req.TemplatePath, req.FieldValues)
    if err != nil {
        return "", fmt.Errorf("failed to get template overlays: %w", err)
    }
}
```

### 5. Event Data Not Sanitized (MEDIUM - Security)

**Location:** app.go:195, 213, 227, 235

**Problem:** Event payloads chứa raw user data chưa sanitize.

```go
runtime.EventsEmit(a.ctx, "error", map[string]string{"message": err.Error()})
// ❌ err.Error() có thể chứa sensitive path info
```

**Impact:**
- File paths leak trong error messages
- Stack traces expose internal structure
- XSS potential nếu frontend render HTML

**Fix:**

```go
// Add sanitization helper
func sanitizeError(err error) string {
    msg := err.Error()
    // Remove absolute paths
    msg = filepath.Base(msg)
    // Limit length
    if len(msg) > 200 {
        msg = msg[:200] + "..."
    }
    return msg
}

// Usage
runtime.EventsEmit(a.ctx, "error", map[string]string{
    "message": sanitizeError(err),
})
```

## High Priority Findings

### 6. Memory Leak Risk in Batch Processing (HIGH - Performance)

**Location:** app.go:209-233

**Problem:** Images được load vào memory nhưng không được explicitly released. Go GC sẽ handle, nhưng với 1000 images có thể cause memory spike.

```go
for i, productPath := range req.ProductImages {
    err := a.processSingleImage(productPath, frame, bgColor, overlays, req)
    // ❌ No explicit cleanup of loaded product image
}
```

**Current memory profile:**
- Frame: Loaded once (good!)
- Products: Loaded sequentially, held until GC
- With 1000 x 5MB images = potential 5GB memory usage

**Fix:** Not critical vì Go GC works, nhưng add monitoring:

```go
// After processing each image
if (i+1) % 100 == 0 {
    runtime.GC()  // Hint to GC to clean up
}
```

Or better: Process in batches of 100.

### 7. Output Path Collision (HIGH - Data Loss)

**Location:** app.go:256-261

**Problem:** Output filename không check for existing files, có thể overwrite.

```go
outputName := nameWithoutExt + "_framed"
outputPath := filepath.Join(req.OutputDir, outputName)
// ❌ No check if file exists
return a.imageSvc.SaveImage(result.Image, outputPath, req.Format, req.Quality)
```

**Impact:**
- Overwrite existing processed images
- Data loss nếu user re-process same batch

**Fix:**

```go
// Generate unique output path
outputName := nameWithoutExt + "_framed"
outputPath := filepath.Join(req.OutputDir, outputName)

// Add suffix if file exists
counter := 1
for {
    testPath := outputPath + ext
    if _, err := os.Stat(testPath); os.IsNotExist(err) {
        outputPath = testPath
        break
    }
    outputPath = fmt.Sprintf("%s_%d", outputName, counter)
    counter++
    if counter > 1000 {
        return fmt.Errorf("too many duplicate files")
    }
}

return a.imageSvc.SaveImage(result.Image, outputPath, req.Format, req.Quality)
```

## Medium Priority Improvements

### 8. Missing Input Validation (MEDIUM - Robustness)

**Location:** app.go:129-136

**Problem:** GeneratePreview không validate req fields thoroughly.

```go
if len(req.ProductImages) == 0 {
    return "", fmt.Errorf("no product images selected")
}
if req.FrameImage == "" {
    return "", fmt.Errorf("no frame image selected")
}
// ❌ No validation of Format, Quality, TemplatePath
```

**Fix:**

```go
// Validate format
validFormats := map[string]bool{"png": true, "jpg": true, "jpeg": true, "webp": true}
if req.Format != "" && !validFormats[strings.ToLower(req.Format)] {
    return "", fmt.Errorf("invalid format: %s", req.Format)
}

// Validate quality
if req.Quality < 0 || req.Quality > 100 {
    req.Quality = 90  // Default
}

// Validate template path exists if provided
if req.TemplatePath != "" {
    if _, err := os.Stat(req.TemplatePath); err != nil {
        return "", fmt.Errorf("template file not found: %w", err)
    }
}
```

### 9. Error Handling Continues Processing (MEDIUM - UX)

**Location:** app.go:229-232

**Problem:** Errors during batch được log nhưng batch continues. User không biết có images failed.

```go
if err != nil {
    // Log error but continue
    fmt.Printf("Error processing %s: %v\n", productPath, err)
}
// ❌ No tracking of failures, no summary report
```

**Fix:**

```go
type ProcessResult struct {
    TotalProcessed int      `json:"totalProcessed"`
    TotalFailed    int      `json:"totalFailed"`
    Failures       []string `json:"failures,omitempty"`
}

// In ProcessBatch
var failures []string

for i, productPath := range req.ProductImages {
    err := a.processSingleImage(...)

    if err != nil {
        failures = append(failures, filepath.Base(productPath))
        fmt.Printf("Error processing %s: %v\n", productPath, err)
    }

    // Emit progress with status
    progress := models.ProcessProgress{
        Current: i + 1,
        Total:   total,
        File:    filepath.Base(productPath),
        Success: err == nil,
    }
    runtime.EventsEmit(a.ctx, "progress", progress)
}

// Final summary
result := ProcessResult{
    TotalProcessed: total - len(failures),
    TotalFailed:    len(failures),
    Failures:       failures,
}
runtime.EventsEmit(a.ctx, "complete", result)
```

### 10. Missing Context Timeout (MEDIUM - Resource Management)

**Location:** app.go:182

**Problem:** Batch processing không có timeout, có thể run forever.

```go
ctx, cancel := context.WithCancel(a.ctx)
// ❌ No timeout protection
```

**Fix:**

```go
// Add configurable timeout
const BatchProcessingTimeout = 30 * time.Minute

ctx, cancel := context.WithTimeout(a.ctx, BatchProcessingTimeout)
```

## Low Priority Suggestions

### 11. Magic Strings (LOW - Maintainability)

Event names "progress", "complete", "error", "cancelled" là magic strings. Consider constants:

```go
const (
    EventProgress  = "progress"
    EventComplete  = "complete"
    EventError     = "error"
    EventCancelled = "cancelled"
)
```

### 12. Missing Metrics/Logging (LOW - Observability)

No structured logging, chỉ có `fmt.Printf`. Consider structured logger:

```go
log.Info("batch processing started",
    "total", total,
    "outputDir", req.OutputDir)
```

## Positive Observations

1. **Clean separation of concerns**: App layer calls services correctly, không có business logic leak
2. **Frame reuse optimization**: Frame loaded once per batch instead of per image - excellent!
3. **Proper mutex usage**: Lock pattern correct (mặc dù có race condition cần fix)
4. **Context-based cancellation**: Using Go idioms correctly
5. **Service initialization**: Dependency injection pattern good
6. **Event-based progress**: Non-blocking async pattern appropriate for Wails
7. **Error wrapping**: Using `fmt.Errorf` with `%w` for error chains

## Architecture Assessment

**Strengths:**
- Wails bindings exposed correctly
- Services injected properly via NewApp
- Async processing pattern appropriate
- Progress events non-blocking

**Concerns:**
- No abstraction for file I/O (harder to test)
- Hard dependency on Wails runtime (harder to unit test)
- No retry logic for transient failures
- No batch size tuning based on memory

**YAGNI Compliance:** ✅ Good - no over-engineering detected. Implementation focused on requirements.

## Recommended Actions

**Must Fix Before Phase 6:**

1. **[CRITICAL]** Add path validation in all 4 selection methods (Issue #1)
2. **[CRITICAL]** Implement comprehensive batch size validation (Issue #2)
3. **[CRITICAL]** Fix race condition in cancel logic (Issue #3)
4. **[CRITICAL]** Handle template errors properly (Issue #4)
5. **[HIGH]** Add output path collision detection (Issue #7)
6. **[HIGH]** Add input validation in GeneratePreview (Issue #8)

**Should Fix:**

7. **[MEDIUM]** Sanitize event data (Issue #5)
8. **[MEDIUM]** Track and report failures in batch (Issue #9)
9. **[MEDIUM]** Add context timeout (Issue #10)

**Nice to Have:**

10. **[LOW]** Extract event name constants (Issue #11)

## Plan Todo List Status

**File:** plans/260108-0121-image-border-app/phase-05-wails-backend.md

Current status từ plan:

```markdown
- [ ] Update App struct with services
- [ ] Implement file selection dialogs
- [ ] Implement template loading
- [ ] Implement preview generation
- [ ] Implement batch processing with events
- [ ] Implement cancel functionality
- [ ] Update main.go with initialization
- [ ] Test all Wails bindings
```

**Actual implementation status:**

- ✅ Update App struct with services - DONE (app.go:25-50)
- ✅ Implement file selection dialogs - DONE (app.go:58-108)
- ✅ Implement template loading - DONE (app.go:111-126)
- ✅ Implement preview generation - DONE (app.go:129-170)
- ✅ Implement batch processing with events - DONE (app.go:173-264)
- ✅ Implement cancel functionality - DONE (app.go:267-274)
- ✅ Update main.go with initialization - DONE (main.go:24)
- ⚠️ Test all Wails bindings - **CANNOT VERIFY** (requires running app)

**Success Criteria từ plan:**

1. ❓ File dialogs work on Windows - **CANNOT VERIFY** without running
2. ✅ Template loading returns correct fields - Logic correct
3. ✅ Preview generates valid base64 image - Logic correct
4. ✅ Progress events emit correctly - Logic correct
5. ⚠️ Cancel stops processing immediately - Has race condition, needs fix

## Test Coverage

**Build:** ✅ PASS
**Tests:** ✅ PASS (all internal packages)
**Race detector:** ✅ PASS (but app.go not covered by tests)

**Missing tests:**
- app.go has NO test file
- Wails bindings not unit testable without mocking runtime
- Integration tests needed for full validation

## Security Summary

**Vulnerabilities found:** 3 CRITICAL, 1 MEDIUM

1. Path traversal (no validation)
2. DoS via batch size
3. Race condition in concurrency
4. Data leak in events

**Security checklist from plan:**

- ❌ Validate file paths before processing - **NOT DONE**
- ⚠️ Limit max batch size (e.g., 1000 images) - **PARTIALLY DONE** (constant defined, validation incomplete)
- ❌ Sanitize event data - **NOT DONE**

## Metrics

- **Type Coverage:** N/A (Go không track)
- **Build Status:** ✅ PASS
- **Test Coverage:** Unknown (app.go untested)
- **Security Issues:** 4 (3 critical, 1 medium)
- **Performance Issues:** 2 (1 high, 1 medium)
- **Code Quality Issues:** 3 (1 high, 2 medium)

## Unresolved Questions

1. Should we implement max file size limit per image? (e.g., 50MB)
2. Should we add disk space check before batch processing?
3. Should we implement batch processing queue instead of blocking?
4. Should we add telemetry/analytics for usage patterns?
5. What's the fallback behavior if output directory becomes unavailable mid-batch?
6. Should we implement graceful shutdown handling for app close during batch?
