# Project Roadmap: Image Border Application

## Overview
This roadmap outlines the key phases and milestones for the development of the Image Border Application. It tracks progress against the overall project plan and provides a high-level view of the development lifecycle.

## Phases & Milestones

### Phase 1: Project Setup & Foundation
- **Status:** Completed
- **Start Date:** 2026-01-07
- **End Date:** 2026-01-07
- **Description:** Wails project initialized with Go backend and React frontend.
- **Key Deliverables:**
    - Running Wails app with React frontend
    - TailwindCSS configured
    - Go and Frontend dependencies installed

### Phase 2: Template Service
- **Status:** In Progress
- **Start Date:** 2026-01-07
- **Description:** Implement template parsing and dynamic field extraction from JSON.
- **Key Deliverables:**
    - `internal/template/service.go` and `parser.go`
    - `internal/models/types.go`
    - Unit tests for template parsing

### Phase 3: Image Service - Core Processing
- **Status:** Pending
- **Start Date:** TBD
- **End Date:** TBD
- **Description:** Core image loading, resizing, and compositing functionality.
- **Key Deliverables:**
    - `internal/image/service.go` and `compositor.go`
    - Basic image compositing (no text)

### Phase 4: Image Service - Text Rendering
- **Status:** Pending
- **Start Date:** TBD
- **End Date:** TBD
- **Description:** Add text overlay with TTF font rendering.
- **Key Deliverables:**
    - Embedded `Roboto-Regular.ttf` font
    - Working text rendering on images

### Phase 5: Wails Backend Integration
- **Status:** Pending
- **Start Date:** TBD
- **End Date:** TBD
- **Description:** Expose Go services to the frontend via Wails bindings.
- **Key Deliverables:**
    - `app.go` with Wails bindings
    - File dialogs and progress events

### Phase 6: React Frontend - UI Components
- **Status:** Pending
- **Start Date:** TBD
- **End Date:** TBD
- **Description:** Develop user interface components for file selection, dynamic forms, and progress display.
- **Key Deliverables:**
    - `FilePicker`, `TemplateFields`, `ProgressBar` components
    - Functional `App.jsx` with UI

### Phase 7: Integration & Testing
- **Status:** Pending
- **Start Date:** TBD
- **End Date:** TBD
- **Description:** End-to-end testing of the application workflow.
- **Key Deliverables:**
    - Fully functional E2E workflow
    - Validated output quality

### Phase 8: Polish & Production Readiness
- **Status:** Pending
- **Start Date:** TBD
- **End Date:** TBD
- **Description:** Refine user experience, implement error handling, and prepare for production build.
- **Key Deliverables:**
    - Polished UX and error messages
    - Production-ready build

## Progress Summary
- Overall Progress: 12.5%
- Current Phase: Phase 2 (Template Service)

## Changelog
### 2026-01-07
- Initial roadmap created.
- Phase 1: Project Setup & Foundation completed.
- Progress updated to 12.5%.
- Started Phase 2: Template Service.
