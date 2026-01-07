# Phase 1: Project Setup & Foundation

**Goal:** Initialize Wails v3 project với Go + React + TailwindCSS

**Duration:** ~2-3 hours

**Prerequisites:**
- Go 1.21+ installed
- Node.js 18+ installed
- Wails v3 CLI installed

---

## Tasks

### Task 1.1: Install Wails v3 CLI

```bash
# Install latest Wails v3
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# Verify installation
wails3 version
```

**Expected Output:** `Wails v3.0.0-alpha.x`

---

### Task 1.2: Initialize Project

```bash
# Create project với React template
wails3 init -n vibe-imageborder -t react-ts

# Navigate to project
cd vibe-imageborder

# Initial structure check
ls -la
```

**Expected Structure:**
```
vibe-imageborder/
├── main.go
├── app.go
├── go.mod
├── wails.json
├── build/
└── frontend/
    ├── package.json
    ├── vite.config.ts
    └── src/
```

---

### Task 1.3: Configure Project Structure

Create internal directories:

```bash
mkdir -p internal/{template,image,models}
mkdir -p assets/fonts
mkdir -p tests/fixtures/{templates,frames,products}
```

**Final Structure:**
```
vibe-imageborder/
├── main.go
├── app.go
├── go.mod
├── wails.json
├── internal/
│   ├── template/
│   ├── image/
│   └── models/
├── assets/
│   └── fonts/
├── tests/
│   └── fixtures/
│       ├── templates/
│       ├── frames/
│       └── products/
└── frontend/
```

---

### Task 1.4: Install Go Dependencies

Update `go.mod`:

```bash
go get github.com/disintegration/imaging@v1.6.2
go get github.com/fogleman/gg@v1.3.0
go mod tidy
```

Verify `go.mod`:
```go
module vibe-imageborder

go 1.21

require (
    github.com/wailsapp/wails/v3 v3.0.0-alpha.x
    github.com/disintegration/imaging v1.6.2
    github.com/fogleman/gg v1.3.0
)
```

---

### Task 1.5: Configure TailwindCSS

Install Tailwind trong frontend:

```bash
cd frontend

# Install Tailwind và dependencies
npm install -D tailwindcss@^3.4.0 postcss autoprefixer

# Initialize Tailwind config
npx tailwindcss init -p
```

Update `frontend/tailwind.config.js`:

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

Update `frontend/src/style.css`:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

/* Custom styles */
body {
  @apply bg-gray-50;
}
```

---

### Task 1.6: Test Build & Run

```bash
# Return to root
cd ..

# Development mode
wails3 dev

# Should open app window với React frontend
```

**Verify:**
- App window opens
- React app displays
- TailwindCSS styles apply
- No console errors

---

## Acceptance Criteria

- ✓ Wails v3 project initialized
- ✓ Go dependencies installed (`imaging`, `gg`)
- ✓ TailwindCSS configured và working
- ✓ Project structure matches spec
- ✓ `wails3 dev` runs successfully
- ✓ App window displays React frontend

---

## Deliverables

### Files Created/Modified

1. `go.mod` - Dependencies configured
2. `frontend/package.json` - Tailwind installed
3. `frontend/tailwind.config.js` - Tailwind config
4. `frontend/src/style.css` - Tailwind directives
5. Directory structure - `internal/`, `assets/`, `tests/`

### Validation

Run checklist:

```bash
# 1. Go dependencies installed
go list -m all | grep imaging
go list -m all | grep gg

# 2. Node dependencies installed
cd frontend && npm list tailwindcss

# 3. App runs
cd .. && wails3 dev
```

All checks pass → Phase 1 complete ✅

---

## Troubleshooting

### Issue: Wails v3 not found

**Solution:**
```bash
# Ensure Go bin in PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Reinstall Wails
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

### Issue: Node modules errors

**Solution:**
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Issue: Tailwind not applying styles

**Solution:**
- Check `tailwind.config.js` content paths
- Verify `@tailwind` directives trong `style.css`
- Restart dev server

---

## Next Phase

[Phase 2: Template Service](phase-02-template-service.md) - Parse template JSON, extract dynamic fields
