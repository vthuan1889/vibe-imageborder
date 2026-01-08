# Agent Guidelines for vibe-imageborder

## Build & Test Commands

```bash
# Development
wails dev                          # Run with hot reload
wails build                        # Build executable
wails build -nsis                  # Build Windows installer

# Tests
go test ./...                      # Run all tests
go test ./tests -v                 # Run integration tests
go test ./internal/image -run TestCompositor  # Run single test

# Frontend
cd frontend && npm install && npm run build   # Build frontend
```

## Architecture

**Go Backend (Wails v2)**
- `app.go`: Main application logic with Wails bindings exposed to frontend
- `internal/image/`: Image processing (load, save, composite, text rendering)
- `internal/template/`: Template file parsing (`.txt` format with text overlays)
- `internal/models/`: Data structures (ProcessRequest, TextOverlay, etc.)
- `internal/updater/`: GitHub-based auto-update checking

**React Frontend**
- `frontend/src/`: Vite + TypeScript, Tailwind CSS
- `frontend/wailsjs/`: Auto-generated Wails bindings
- Main entry: `App.tsx` with UI components

**Key APIs**
- Batch image processing with progress events (EventProgress, EventComplete)
- Context-based cancellation for long-running tasks
- Font embedding from `assets/fonts/`
- File dialogs via Wails runtime

## Code Style

**Go**
- PascalCase for exported functions/types, camelCase for unexported
- Error wrapping: `fmt.Errorf("context: %w", err)`
- Path validation via `filepath.Clean`, `filepath.Abs`
- Error sanitization (remove sensitive paths) before frontend transmission
- Constants for magic numbers (e.g., MaxBatchSize=1000)

**Frontend**
- TypeScript with strict mode, React hooks
- Tailwind CSS for styling
- Import from `wailsjs/go/main` for Go bindings

**Naming**
- Services: `*Service` interface pattern (template.Service, image.Service)
- Handlers: `Handle*` or verb-first (SelectProductFiles, ProcessBatch)
- Events: PascalCase constants (EventProgress, EventComplete)

**Error Handling**
- Go: Wrap errors with context, return early
- Frontend: Emit error events via `runtime.EventsEmit(ctx, EventError, ...)`
- Validate inputs (file existence, path safety, format support)
