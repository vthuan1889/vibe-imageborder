# Image Border Tool

A desktop application for batch processing images with custom frames, text overlays, and templates. Built with Wails (Go + React).

## Features

- **Batch Processing**: Process hundreds of images at once
- **Custom Frames**: Add PNG frames/borders to product images
- **Text Overlays**: Add dynamic text using template files
- **Multiple Formats**: Export to PNG, JPG, or WebP
- **Auto-Update**: Built-in update checker via GitHub Releases
- **Cross-Platform**: Windows support (macOS/Linux possible)

## Installation

### From Release (Recommended)

1. Go to [Releases](https://github.com/vthuan1889/vibe-imageborder/releases)
2. Download `ImageBorderTool-amd64-installer.exe`
3. Run the installer

### From Source

See [Build Instructions](#build-instructions) below.

## Usage

### Basic Workflow

1. **Select Product Images**: Click "Product Images" to choose images to process
2. **Select Frame**: Choose a PNG frame image (transparent areas show the product)
3. **Optional Template**: Load a `.txt` template file for text overlays
4. **Set Output**: Choose format (PNG/JPG/WebP), quality, and output folder
5. **Preview**: Click "Preview" to see the first image result
6. **Generate**: Click "Generate All" to process all images

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
│   │   ├── App.tsx        # Main React component
│   │   └── components/    # UI components
│   └── wailsjs/           # Auto-generated Wails bindings
├── internal/              # Go packages
│   ├── image/             # Image processing
│   ├── models/            # Data models
│   ├── template/          # Template parsing
│   └── updater/           # Auto-update logic
├── build/                 # Build assets
│   └── windows/           # Windows-specific (icon, NSIS)
└── .github/workflows/     # CI/CD
```

## License

MIT License

Copyright (c) 2026 TT PC 2

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
