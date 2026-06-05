# Image Border Tool

A desktop application for batch processing images with custom frames, text overlays, and templates. Built with Wails (Go + React).

## Features

- **Sidebar Navigation**: Multi-section UI with expandable menu structure
- **Create Frame**: Batch compose product images with frames and text overlays
- **Copy Image**: Recursively copy images with imperceptible fingerprint changes for social platforms
- **Batch Processing**: Process hundreds of images at once
- **Custom Frames**: Add PNG frames/borders to product images
- **Text Overlays**: Add dynamic text using template files
- **Multiple Formats**: Export to PNG, JPG, or WebP
- **Auto-Update**: Built-in update checker via GitHub Releases
- **Cross-Platform**: Windows support (macOS/Linux possible)

### Navigation (Sidebar)

| Menu | Status | Description |
|------|--------|-------------|
| **Create Frame** | Available | Main workflow: select images, frame, template, preview, and batch export |
| **Copy Image** | Available | Copy images recursively between folders with invisible fingerprint uniquification |
| **Frame Library** | Coming soon | Browse and manage saved frame templates |
| **Settings** | Coming soon | App preferences and default output options |

## Installation

### From Release (Recommended)

1. Go to [Releases](https://github.com/vthuan1889/vibe-imageborder/releases)
2. Download `ImageBorderTool-amd64-installer.exe`
3. Run the installer

### From Source

See [Build Instructions](#build-instructions) below.

## Usage

### Basic Workflow

Open the app and select **Create Frame** from the sidebar (default view).

1. **Select Product Images**: Click "Product Images" to choose images to process
2. **Select Frame**: Choose a PNG frame image (transparent areas show the product)
3. **Optional Template**: Load a `.txt` template file for text overlays
4. **Set Output**: Choose format (PNG/JPG/WebP), quality, and output folder
5. **Preview**: Click "Preview" to see the first image result
6. **Generate**: Click "Generate All" to process all images

The app remembers your last selected sidebar tab between sessions.

### Copy Image Workflow

Select **Copy Image** from the sidebar:

1. **Source Folder**: Choose the folder containing original images (scanned recursively)
2. **Destination Folder**: Choose where copies should be written (folder structure preserved)
3. **Copy**: Process all images — each copy is visually identical but has a unique file fingerprint

### Template Format

Create a `.txt` file with the following structure:

```
background=#FFFFFF

[field_name]
text=Your Text Here
x=100
y=50
size=24
color=#000000
font=BeVietnamPro
```

### Check for Updates

Click the "Check for Update" button in the top-right corner to check for new versions.

## Build Instructions

### Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- [NSIS](https://nsis.sourceforge.io/) (for Windows installer)

### Setup

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Clone repository
git clone https://github.com/vthuan1889/vibe-imageborder.git
cd vibe-imageborder

# Install frontend dependencies
cd frontend && npm install && cd ..
```

### Development

```bash
# Run in development mode with hot reload
wails dev
```

### Build

```bash
# Build executable only
wails build

# Build with NSIS installer (Windows)
wails build -nsis

# Build with version info
wails build -nsis -ldflags "-X 'main.version=v1.0.0'"
```

Output files will be in `build/bin/`:
- `ImageBorderTool.exe` - Standalone executable
- `ImageBorderTool-amd64-installer.exe` - NSIS installer

### Release

Push a tag to trigger automated release:

```bash
git tag v1.0.1
git push origin v1.0.1
```

GitHub Actions will build and upload the installer to Releases.

## Project Structure

```
vibe-imageborder/
├── app.go                 # Main app logic and Wails bindings
├── main.go                # Entry point with version info
├── wails.json             # Wails configuration
├── frontend/              # React frontend
│   ├── src/
│   │   ├── App.tsx        # Layout shell (sidebar + view switcher)
│   │   ├── views/         # Feature screens
│   │   │   └── CreateFrameView.tsx
│   │   ├── components/    # Reusable UI components
│   │   │   ├── Sidebar.tsx
│   │   │   ├── ComingSoonView.tsx
│   │   │   ├── FilePicker.tsx
│   │   │   └── ...
│   │   ├── types/
│   │   │   └── navigation.ts  # Sidebar menu definitions
│   │   └── utils/
│   │       └── storage.ts     # localStorage persistence
│   └── wailsjs/           # Auto-generated Wails bindings
├── internal/              # Go packages
│   ├── image/             # Image processing
│   ├── models/            # Data models
│   ├── template/          # Template parsing
│   └── updater/           # Auto-update logic
├── docs/                  # Project documentation
├── build/                 # Build assets
│   └── windows/           # Windows-specific (icon, NSIS)
└── .github/workflows/     # CI/CD
```

### Frontend Architecture

The UI uses view-based navigation (no `react-router`). `App.tsx` renders a fixed sidebar and switches content by `activeView` state.

To add a new feature:

1. Add an entry to `frontend/src/types/navigation.ts`
2. Create a view in `frontend/src/views/`
3. Register the view in `frontend/src/App.tsx`
4. Add Wails bindings in `app.go` if backend support is needed

## License

MIT License

Copyright (c) 2026 KKAuto.net

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
