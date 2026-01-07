# Phase 8: Polish & Optimization

## Context

- Plan: [plan.md](./plan.md)
- Previous: [Phase 7 - Integration Testing](./phase-07-integration-testing.md)

## Overview

| Field | Value |
|-------|-------|
| Priority | P3 |
| Status | Pending |
| Effort | 0.5h |

Final polish: error handling, loading states, UI refinements, build configuration.

## Requirements

### Functional
- Clear error messages for users
- Loading indicators for all async operations
- Window title shows app name

### Non-functional
- Clean build output
- App icon configured
- No console errors in production

## Polish Items

### 1. Error Handling

Add user-friendly error messages:

```tsx
// In App.tsx, wrap async calls:
try {
  await SomeOperation();
} catch (e: any) {
  // Show toast or alert
  alert(`Error: ${e.message || e}`);
}
```

### 2. Loading States

Already implemented in Preview component. Verify all buttons disable during loading.

### 3. App Icon

Replace `build/appicon.png` with custom icon (512x512 PNG).

### 4. Window Configuration

```go
// In main.go
err := wails.Run(&options.App{
    Title:     "Image Border Tool",
    Width:     1200,
    Height:    800,
    MinWidth:  800,
    MinHeight: 600,
    // ...
})
```

### 5. Build Configuration

Update `wails.json`:
```json
{
  "name": "Image Border Tool",
  "outputfilename": "ImageBorderTool",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "author": {
    "name": "Your Name"
  }
}
```

### 6. Production Build

```bash
wails build
```

Output: `build/bin/ImageBorderTool.exe`

## Final Checklist

### Code Quality
- [ ] No `console.log` in production code
- [ ] No TODO comments left
- [ ] Error boundaries in React
- [ ] All imports used

### UI/UX
- [ ] Consistent spacing
- [ ] Focus states on inputs
- [ ] Disabled states clear
- [ ] Loading indicators present

### Build
- [ ] App icon set
- [ ] Window title correct
- [ ] Min window size set
- [ ] Production build works
- [ ] Executable runs standalone

### Documentation
- [ ] README.md with usage instructions
- [ ] Build instructions documented

## Todo List

- [ ] Review error handling
- [ ] Verify loading states
- [ ] Replace app icon
- [ ] Update wails.json
- [ ] Test production build
- [ ] Final QA pass

## Success Criteria

1. Production build runs without errors
2. No console warnings
3. All async operations have loading states
4. Errors show user-friendly messages

## Completion

After this phase, the MVP is complete and ready for use.
