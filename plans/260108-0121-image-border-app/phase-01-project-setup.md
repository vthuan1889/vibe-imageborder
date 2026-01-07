# Phase 1: Project Setup & Foundation

## Context

- Plan: [plan.md](./plan.md)
- Brainstorm: [brainstormer-260108-0121-image-border-app-solution.md](../reports/brainstormer-260108-0121-image-border-app-solution.md)

## Overview

| Field | Value |
|-------|-------|
| Priority | P1 - Critical Path |
| Status | Completed (2026-01-08) |
| Effort | 2h |

Initialize Wails 2 project with Go backend + React frontend. Setup project structure, dependencies, and bundled fonts.

## Requirements

### Functional
- Wails 2 project initialized with React template
- Go module with required dependencies
- Frontend with TailwindCSS configured
- Bundle fonts embedded in assets

### Non-functional
- Clean project structure following Go conventions
- TypeScript strict mode enabled
- Hot reload working for development

## Architecture

```
vibe-imageborder/
├── main.go                    # Entry point
├── app.go                     # Wails app struct + methods
├── go.mod
├── go.sum
├── internal/
│   ├── image/
│   ├── template/
│   └── models/
│       └── types.go
├── assets/
│   └── fonts/
│       ├── BeVietnamPro-Regular.ttf
│       └── Roboto-Regular.ttf
├── frontend/
│   ├── package.json
│   ├── tsconfig.json
│   ├── tailwind.config.js
│   ├── postcss.config.js
│   ├── vite.config.ts
│   ├── index.html
│   └── src/
│       ├── main.tsx
│       ├── App.tsx
│       └── style.css
└── build/
    └── appicon.png
```

## Related Code Files

### Files to Create
| File | Purpose |
|------|---------|
| `main.go` | Wails app entry point |
| `app.go` | App struct with exposed methods |
| `internal/models/types.go` | Shared type definitions |
| `assets/fonts/*.ttf` | Bundled fonts |
| `frontend/tailwind.config.js` | TailwindCSS config |
| `frontend/postcss.config.js` | PostCSS config |

### Files to Modify
| File | Change |
|------|--------|
| `frontend/package.json` | Add tailwindcss deps |
| `frontend/src/style.css` | Add Tailwind directives |

## Implementation Steps

### Step 1: Initialize Wails Project (if not exists)
```bash
# If starting fresh
wails init -n vibe-imageborder -t react-ts

# Or if project exists, ensure wails.json is correct
```

### Step 2: Install Go Dependencies
```bash
go get github.com/disintegration/imaging
go get github.com/fogleman/gg
go get golang.org/x/image/font/opentype
```

### Step 3: Create Project Structure
```bash
mkdir -p internal/image internal/template internal/models
mkdir -p assets/fonts
```

### Step 4: Create types.go
```go
// internal/models/types.go
package models

// TextOverlay represents text to draw on image
type TextOverlay struct {
    Text     string `json:"text"`
    Position string `json:"position"` // "x,y"
    FontSize int    `json:"fontsize"`
    Color    string `json:"color"`
}

// TemplateConfig represents parsed template JSON
type TemplateConfig struct {
    Background string                 `json:"background,omitempty"`
    Fields     map[string]TextOverlay `json:"-"`
    Raw        map[string]interface{} `json:"-"`
}

// ProcessRequest represents batch processing request
type ProcessRequest struct {
    ProductImages []string          `json:"productImages"`
    FrameImage    string            `json:"frameImage"`
    TemplatePath  string            `json:"templatePath"`
    FieldValues   map[string]string `json:"fieldValues"`
    OutputDir     string            `json:"outputDir"`
    Format        string            `json:"format"` // png, jpg, webp
    Quality       int               `json:"quality"`
}

// ProcessProgress represents progress update
type ProcessProgress struct {
    Current int    `json:"current"`
    Total   int    `json:"total"`
    File    string `json:"file"`
}
```

### Step 5: Setup app.go
```go
// app.go
package main

import (
    "context"
)

// App struct
type App struct {
    ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
    return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
    return "Hello " + name + "!"
}
```

### Step 6: Setup main.go
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
    app := NewApp()

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

### Step 7: Download Bundle Fonts
Download from Google Fonts:
- Be Vietnam Pro: https://fonts.google.com/specimen/Be+Vietnam+Pro
- Roboto: https://fonts.google.com/specimen/Roboto

Place in `assets/fonts/`:
- `BeVietnamPro-Regular.ttf`
- `Roboto-Regular.ttf`

### Step 8: Setup TailwindCSS in Frontend
```bash
cd frontend
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

Update `tailwind.config.js`:
```js
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

Update `src/style.css`:
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

### Step 9: Verify Build
```bash
wails dev
```

## Todo List

- [x] Initialize Wails project (if needed)
- [x] Install Go dependencies
- [x] Create internal folder structure
- [x] Create models/types.go
- [x] Setup app.go with context
- [x] Setup main.go with embed
- [x] Download and add bundle fonts
- [x] Configure TailwindCSS
- [x] Verify wails dev runs successfully

## Success Criteria

1. `wails dev` starts without errors
2. Browser opens with React app
3. TailwindCSS classes work
4. Go modules resolve correctly
5. Fonts embedded in binary

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Wails CLI not installed | High | Document installation steps |
| Font file too large | Low | Fonts are ~200KB each, acceptable |
| Go version mismatch | Medium | Require Go 1.21+ in docs |

## Security Considerations

- Fonts bundled in binary, no external downloads
- No network calls in this phase

## Next Steps

After completion, proceed to [Phase 2: Template Service](./phase-02-template-service.md)
