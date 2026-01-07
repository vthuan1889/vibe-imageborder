---
title: "Image Border App - Go + Wails 2"
description: "Desktop app ghép hình sản phẩm + khung + text overlay"
status: completed
priority: P1
effort: 16h
branch: main
tags: [feature, desktop, go, wails, image-processing]
created: 2026-01-08
---

# Image Border App - Implementation Plan

## Overview

Build desktop app to composite product images with frame overlays and text annotations using Go + Wails 2. Port core functionality from existing C# WinForms app.

## Architecture

```
┌─────────────────────────────────────────────┐
│              Wails 2 App                    │
├─────────────────────────────────────────────┤
│  Frontend (React + TailwindCSS)             │
│  ├── FilePicker                             │
│  ├── TemplateFields (dynamic form)          │
│  ├── Preview                                │
│  └── ProgressBar                            │
├─────────────────────────────────────────────┤
│  Backend (Go)                               │
│  ├── ImageService (composite, resize)       │
│  ├── TemplateService (parse JSON)           │
│  └── FileService (browse, save)             │
└─────────────────────────────────────────────┘
```

## Phases

| # | Phase | Status | Effort | Link |
|---|-------|--------|--------|------|
| 1 | Project Setup & Foundation | Completed | 2h | [phase-01](./phase-01-project-setup.md) |
| 2 | Template Service | Completed | 2h | [phase-02](./phase-02-template-service.md) |
| 3 | Image Service - Core | Completed | 3h | [phase-03](./phase-03-image-service-core.md) |
| 4 | Image Service - Text Rendering | Completed | 2h | [phase-04](./phase-04-text-rendering.md) |
| 5 | Wails Backend Integration | Completed | 2h | [phase-05](./phase-05-wails-backend.md) |
| 6 | React Frontend | Completed | 3h | [phase-06](./phase-06-react-frontend.md) |
| 7 | Integration & Testing | Completed | 1.5h | [phase-07](./phase-07-integration-testing.md) |
| 8 | Polish & Optimization | Completed | 0.5h | [phase-08](./phase-08-polish.md) |

## Tech Stack

| Component | Technology |
|-----------|------------|
| Framework | Wails 2 |
| Backend | Go 1.21+ |
| Frontend | React + TypeScript |
| Styling | TailwindCSS |
| Image resize | `disintegration/imaging` |
| Text render | `fogleman/gg` |
| Font loading | `golang/freetype` |

## Dependencies

- Wails CLI installed
- Go 1.21+
- Node.js 18+
- Bundle fonts: Be Vietnam Pro, Roboto

## Key Files

```
vibe-imageborder/
├── main.go
├── app.go
├── internal/
│   ├── image/
│   │   ├── compositor.go
│   │   └── service.go
│   ├── template/
│   │   ├── parser.go
│   │   └── service.go
│   └── models/
│       └── types.go
├── frontend/
│   └── src/
│       ├── App.tsx
│       └── components/
├── assets/
│   └── fonts/
└── build/
```

## Success Criteria

1. Process 100 images < 30 seconds
2. Vietnamese diacritics render correctly
3. Output quality matches input resolution
4. UI remains responsive during batch processing

## Related

- Brainstorm: [brainstormer-260108-0121-image-border-app-solution.md](../reports/brainstormer-260108-0121-image-border-app-solution.md)
- Reference app: `D:\Code-Tool\Software\web-tool\UploadImage`
