# Phase 4: Updater Frontend (React)

> Parent: [plan.md](./plan.md)

## Overview

| Field | Value |
|-------|-------|
| Date | 2026-01-08 |
| Priority | P1 |
| Effort | 2h |
| Status | pending |

## Requirements

1. "Check for Update" button in UI
2. Show current version
3. Display update available notification
4. Download progress indicator
5. Confirm before installing

## Implementation Steps

### 4.1 Create UpdateButton Component

Create `frontend/src/components/UpdateButton.tsx`:

```tsx
import { useState } from 'react';
import {
  GetVersion,
  CheckForUpdate,
  DownloadAndInstallUpdate,
} from '../../wailsjs/go/main/App';

interface UpdateInfo {
  available: boolean;
  current: string;
  latest: string;
  downloadUrl: string;
}

export function UpdateButton() {
  const [version, setVersion] = useState('');
  const [checking, setChecking] = useState(false);
  const [downloading, setDownloading] = useState(false);
  const [updateInfo, setUpdateInfo] = useState<UpdateInfo | null>(null);

  // Load version on mount
  useState(() => {
    GetVersion().then(setVersion);
  });

  const handleCheck = async () => {
    setChecking(true);
    try {
      const info = await CheckForUpdate();
      setUpdateInfo(info);
      if (!info.available) {
        alert('You are using the latest version!');
      }
    } catch (e) {
      alert('Failed to check for updates: ' + e);
    } finally {
      setChecking(false);
    }
  };

  const handleUpdate = async () => {
    if (!updateInfo?.downloadUrl) return;

    const confirmed = window.confirm(
      `Update to ${updateInfo.latest}?\n\nThe app will close and the installer will run.`
    );
    if (!confirmed) return;

    setDownloading(true);
    try {
      await DownloadAndInstallUpdate(updateInfo.downloadUrl);
      // App will quit after this
    } catch (e) {
      alert('Update failed: ' + e);
      setDownloading(false);
    }
  };

  return (
    <div className="flex items-center gap-2 text-sm">
      <span className="text-gray-500">v{version}</span>

      {updateInfo?.available ? (
        <button
          onClick={handleUpdate}
          disabled={downloading}
          className="px-3 py-1 bg-green-500 hover:bg-green-600 text-white
                     rounded text-xs disabled:opacity-50"
        >
          {downloading ? 'Downloading...' : `Update to ${updateInfo.latest}`}
        </button>
      ) : (
        <button
          onClick={handleCheck}
          disabled={checking}
          className="px-3 py-1 bg-gray-200 hover:bg-gray-300 text-gray-700
                     rounded text-xs disabled:opacity-50"
        >
          {checking ? 'Checking...' : 'Check for Update'}
        </button>
      )}
    </div>
  );
}
```

### 4.2 Add to App.tsx

Add import and render in header/footer area:

```tsx
import { UpdateButton } from './components/UpdateButton';

// In JSX, add to header or bottom of page
<div className="absolute top-2 right-4">
  <UpdateButton />
</div>
```

### 4.3 Wails Bindings

After adding Go methods, regenerate bindings:
```bash
wails generate module
```

New bindings will be created:
- `frontend/wailsjs/go/main/App.js`
- `frontend/wailsjs/go/main/App.d.ts`

## Files to Create/Modify

| File | Change |
|------|--------|
| `frontend/src/components/UpdateButton.tsx` | New component |
| `frontend/src/App.tsx` | Import and render UpdateButton |

## UI Placement Options

1. **Header right corner** (recommended) - Always visible
2. **Settings modal** - Less intrusive
3. **Footer** - Subtle placement

## Success Criteria

- [ ] Version displays correctly
- [ ] Check button triggers API call
- [ ] "Up to date" message when no update
- [ ] Update button appears when available
- [ ] Confirmation dialog before install
- [ ] App closes after starting download
