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
- `App.tsx`: Layout shell (sidebar + view switcher)
- `views/CreateFrameView.tsx`: Main batch compose workflow
- `views/CopyImageView.tsx`: Recursive image copy with uniquify fingerprint
- `components/Sidebar.tsx`, `types/navigation.ts`: Sidebar navigation
- View-based routing via `activeView` state (no react-router)

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

## Cursor Cloud specific instructions

**Product:** Wails v2 desktop app (Go backend + React frontend). Primary target is Windows; Linux dev is partially supported.

**One-time Linux system packages** (not in the VM update script): `libgtk-3-dev`, `libwebkit2gtk-4.1-dev`, `build-essential`, `pkg-config`.

**PATH:** Ensure Wails CLI is on `PATH`: `export PATH="$PATH:$(go env GOPATH)/bin"` (install once with `go install github.com/wailsapp/wails/v2/cmd/wails@latest`).

**Frontend embed:** `main.go` embeds `frontend/dist`. Run `cd frontend && npm run build` before `go build` / `go test` on the root module, or let `wails build` / `wails dev` build it automatically.

**Lint:** No ESLint or golangci-lint config in repo. Use `cd frontend && npm run build` (runs `tsc`) for frontend type-checking.

**Tests on Linux (recommended):**
- `go test ./internal/image ./tests -v` — image compositing, copy pipeline, templates (all pass headless)
- `go test ./...` currently fails on Linux due to pre-existing issues: `internal/updater` uses Windows-only `syscall.SysProcAttr.HideWindow`, and `internal/template/service_test.go` references removed cache APIs

**Wails GUI on Linux:** `wails dev` / `wails build` fail until `internal/updater/updater.go` is split with `//go:build windows` (or equivalent) for `DownloadAndInstall`. Use Windows or fix that package for full GUI E2E. `wails doctor` may warn about `libwebkit` even when `libwebkit2gtk-4.1-dev` is installed (Wails looks for 4.0 pkg-config name on some distros).

**Standalone Vite preview** (`npm run preview`) serves static assets but the UI stays blank without the Wails runtime (`window.go.main`); do not use it as an app smoke test.

**Core workflow smoke test (no GUI):** `go test ./tests -run TestIntegration_BasicComposite -v` writes `tests/output/test-composite.png` (product + frame composite).
