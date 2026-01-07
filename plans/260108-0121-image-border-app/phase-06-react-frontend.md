# Phase 6: React Frontend

## Context

- Plan: [plan.md](./plan.md)
- Previous: [Phase 5 - Wails Backend](./phase-05-wails-backend.md)

## Overview

| Field | Value |
|-------|-------|
| Priority | P1 - Critical Path |
| Status | Pending |
| Effort | 3h |

Build React frontend with 2-column layout: input controls (left) and preview/output (right). Use TailwindCSS for styling.

## Requirements

### Functional
- File pickers for products, frame, template
- Dynamic form for template fields
- Preview display with loading state
- Output settings (format, quality, folder)
- Progress bar during batch processing
- Generate button

### Non-functional
- Responsive 2-column layout
- Clean, modern UI with TailwindCSS
- Loading states for async operations
- Error display for user feedback

## UI Layout (2-Column)

```
┌────────────────────────────┬────────────────────────────────┐
│       COLUMN 1 (40%)       │         COLUMN 2 (60%)         │
├────────────────────────────┼────────────────────────────────┤
│  📁 Product Images         │  Preview                       │
│  [Drop/Browse]             │  [Image Preview Area]          │
│                            │  [Preview First Image]         │
│  🖼️ Frame Image            │                                │
│  [Select frame]            │  Output Settings               │
│                            │  Format: [PNG/JPG/WebP]        │
│  📄 Template (Optional)    │  Quality: [slider]             │
│  [Browse .txt]             │  Output: [folder path]         │
│                            │                                │
│  Text Fields               │  Progress                      │
│  (dynamic, show when       │  [████████░░░░] 67%            │
│   template loaded)         │                                │
│  - Barcode: [input]        │  [Generate All Button]         │
│  - Price: [input]          │                                │
└────────────────────────────┴────────────────────────────────┘
```

## Related Code Files

### Files to Create
| File | Purpose |
|------|---------|
| `frontend/src/App.tsx` | Main app component |
| `frontend/src/components/FilePicker.tsx` | File selection component |
| `frontend/src/components/TemplateFields.tsx` | Dynamic form for template fields |
| `frontend/src/components/Preview.tsx` | Image preview component |
| `frontend/src/components/ProgressBar.tsx` | Progress bar component |
| `frontend/src/components/OutputSettings.tsx` | Format/quality/folder settings |

### Files to Modify
| File | Change |
|------|--------|
| `frontend/src/style.css` | Add Tailwind + custom styles |
| `frontend/index.html` | Update title |

## Implementation Steps

### Step 1: Update style.css

```css
/* frontend/src/style.css */
@tailwind base;
@tailwind components;
@tailwind utilities;

body {
  @apply bg-gray-50 font-sans;
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  @apply bg-gray-100;
}

::-webkit-scrollbar-thumb {
  @apply bg-gray-300 rounded-full;
}

::-webkit-scrollbar-thumb:hover {
  @apply bg-gray-400;
}
```

### Step 2: Create FilePicker.tsx

```tsx
// frontend/src/components/FilePicker.tsx
import { FC } from 'react';

interface FilePickerProps {
  label: string;
  icon: string;
  files: string[];
  multiple?: boolean;
  onSelect: () => void;
  onClear?: () => void;
}

export const FilePicker: FC<FilePickerProps> = ({
  label,
  icon,
  files,
  multiple = false,
  onSelect,
  onClear,
}) => {
  const hasFiles = files.length > 0;

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-4">
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2">
          <span className="text-lg">{icon}</span>
          <span className="font-medium text-gray-700">{label}</span>
        </div>
        {hasFiles && onClear && (
          <button
            onClick={onClear}
            className="text-sm text-gray-400 hover:text-red-500"
          >
            Clear
          </button>
        )}
      </div>

      {hasFiles ? (
        <div className="text-sm text-gray-600">
          {multiple ? (
            <span className="bg-blue-100 text-blue-700 px-2 py-1 rounded">
              {files.length} file(s) selected
            </span>
          ) : (
            <span className="truncate block">{files[0].split('\\').pop()}</span>
          )}
        </div>
      ) : (
        <button
          onClick={onSelect}
          className="w-full py-3 border-2 border-dashed border-gray-300 rounded-lg
                     text-gray-500 hover:border-blue-400 hover:text-blue-500
                     transition-colors"
        >
          Click to browse
        </button>
      )}
    </div>
  );
};
```

### Step 3: Create TemplateFields.tsx

```tsx
// frontend/src/components/TemplateFields.tsx
import { FC } from 'react';

interface TemplateFieldsProps {
  fields: string[];
  values: Record<string, string>;
  onChange: (field: string, value: string) => void;
}

export const TemplateFields: FC<TemplateFieldsProps> = ({
  fields,
  values,
  onChange,
}) => {
  if (fields.length === 0) return null;

  // Format field name for display
  const formatLabel = (field: string) => {
    return field
      .replace(/_/g, ' ')
      .replace(/\b\w/g, (c) => c.toUpperCase());
  };

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-4">
      <div className="flex items-center gap-2 mb-3">
        <span className="text-lg">📝</span>
        <span className="font-medium text-gray-700">Text Fields</span>
      </div>

      <div className="space-y-3">
        {fields.map((field) => (
          <div key={field}>
            <label className="block text-sm text-gray-600 mb-1">
              {formatLabel(field)}
            </label>
            <input
              type="text"
              value={values[field] || ''}
              onChange={(e) => onChange(field, e.target.value)}
              placeholder={`Enter ${field}`}
              className="w-full px-3 py-2 border border-gray-300 rounded-md
                         focus:ring-2 focus:ring-blue-500 focus:border-transparent
                         text-sm"
            />
          </div>
        ))}
      </div>
    </div>
  );
};
```

### Step 4: Create Preview.tsx

```tsx
// frontend/src/components/Preview.tsx
import { FC } from 'react';

interface PreviewProps {
  imageData: string | null;
  isLoading: boolean;
  onPreview: () => void;
  canPreview: boolean;
}

export const Preview: FC<PreviewProps> = ({
  imageData,
  isLoading,
  onPreview,
  canPreview,
}) => {
  return (
    <div className="bg-white rounded-lg border border-gray-200 p-4 flex-1">
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <span className="text-lg">🖼️</span>
          <span className="font-medium text-gray-700">Preview</span>
        </div>
        <button
          onClick={onPreview}
          disabled={!canPreview || isLoading}
          className="px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 rounded
                     disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isLoading ? 'Loading...' : 'Preview First'}
        </button>
      </div>

      <div className="aspect-video bg-gray-100 rounded-lg flex items-center justify-center overflow-hidden">
        {isLoading ? (
          <div className="animate-pulse text-gray-400">Generating preview...</div>
        ) : imageData ? (
          <img
            src={imageData}
            alt="Preview"
            className="max-w-full max-h-full object-contain"
          />
        ) : (
          <span className="text-gray-400">No preview yet</span>
        )}
      </div>
    </div>
  );
};
```

### Step 5: Create ProgressBar.tsx

```tsx
// frontend/src/components/ProgressBar.tsx
import { FC } from 'react';

interface ProgressBarProps {
  current: number;
  total: number;
  currentFile: string;
  isProcessing: boolean;
}

export const ProgressBar: FC<ProgressBarProps> = ({
  current,
  total,
  currentFile,
  isProcessing,
}) => {
  const percentage = total > 0 ? Math.round((current / total) * 100) : 0;

  if (!isProcessing && current === 0) {
    return null;
  }

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-4">
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm font-medium text-gray-700">
          {isProcessing ? 'Processing...' : 'Complete!'}
        </span>
        <span className="text-sm text-gray-500">
          {current}/{total} ({percentage}%)
        </span>
      </div>

      <div className="w-full bg-gray-200 rounded-full h-2 mb-2">
        <div
          className="bg-blue-500 h-2 rounded-full transition-all duration-300"
          style={{ width: `${percentage}%` }}
        />
      </div>

      {currentFile && (
        <div className="text-xs text-gray-500 truncate">
          {currentFile}
        </div>
      )}
    </div>
  );
};
```

### Step 6: Create OutputSettings.tsx

```tsx
// frontend/src/components/OutputSettings.tsx
import { FC } from 'react';

interface OutputSettingsProps {
  format: string;
  quality: number;
  outputFolder: string;
  onFormatChange: (format: string) => void;
  onQualityChange: (quality: number) => void;
  onSelectFolder: () => void;
}

export const OutputSettings: FC<OutputSettingsProps> = ({
  format,
  quality,
  outputFolder,
  onFormatChange,
  onQualityChange,
  onSelectFolder,
}) => {
  return (
    <div className="bg-white rounded-lg border border-gray-200 p-4">
      <div className="flex items-center gap-2 mb-3">
        <span className="text-lg">⚙️</span>
        <span className="font-medium text-gray-700">Output Settings</span>
      </div>

      <div className="space-y-4">
        {/* Format */}
        <div>
          <label className="block text-sm text-gray-600 mb-1">Format</label>
          <select
            value={format}
            onChange={(e) => onFormatChange(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
          >
            <option value="png">PNG</option>
            <option value="jpg">JPG</option>
            <option value="webp">WebP</option>
          </select>
        </div>

        {/* Quality (only for JPG/WebP) */}
        {format !== 'png' && (
          <div>
            <label className="block text-sm text-gray-600 mb-1">
              Quality: {quality}%
            </label>
            <input
              type="range"
              min="1"
              max="100"
              value={quality}
              onChange={(e) => onQualityChange(parseInt(e.target.value))}
              className="w-full"
            />
          </div>
        )}

        {/* Output Folder */}
        <div>
          <label className="block text-sm text-gray-600 mb-1">Output Folder</label>
          <div className="flex gap-2">
            <input
              type="text"
              value={outputFolder}
              readOnly
              placeholder="Select output folder"
              className="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm
                         bg-gray-50 truncate"
            />
            <button
              onClick={onSelectFolder}
              className="px-3 py-2 bg-gray-100 hover:bg-gray-200 rounded-md text-sm"
            >
              📁
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
```

### Step 7: Create App.tsx

```tsx
// frontend/src/App.tsx
import { useState, useEffect } from 'react';
import { FilePicker } from './components/FilePicker';
import { TemplateFields } from './components/TemplateFields';
import { Preview } from './components/Preview';
import { ProgressBar } from './components/ProgressBar';
import { OutputSettings } from './components/OutputSettings';

import {
  SelectProductFiles,
  SelectFrameFile,
  SelectTemplateFile,
  SelectOutputFolder,
  LoadTemplate,
  GeneratePreview,
  ProcessBatch,
  CancelProcessing,
} from '../wailsjs/go/main/App';
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime';

function App() {
  // File state
  const [productFiles, setProductFiles] = useState<string[]>([]);
  const [frameFile, setFrameFile] = useState<string>('');
  const [templateFile, setTemplateFile] = useState<string>('');

  // Template state
  const [templateFields, setTemplateFields] = useState<string[]>([]);
  const [fieldValues, setFieldValues] = useState<Record<string, string>>({});

  // Output state
  const [format, setFormat] = useState('png');
  const [quality, setQuality] = useState(90);
  const [outputFolder, setOutputFolder] = useState('');

  // Preview state
  const [previewImage, setPreviewImage] = useState<string | null>(null);
  const [isPreviewLoading, setIsPreviewLoading] = useState(false);

  // Processing state
  const [isProcessing, setIsProcessing] = useState(false);
  const [progress, setProgress] = useState({ current: 0, total: 0, file: '' });

  // Event listeners
  useEffect(() => {
    EventsOn('progress', (data: any) => {
      setProgress({ current: data.current, total: data.total, file: data.file });
    });

    EventsOn('complete', () => {
      setIsProcessing(false);
    });

    EventsOn('error', (data: any) => {
      alert('Error: ' + data.message);
      setIsProcessing(false);
    });

    return () => {
      EventsOff('progress');
      EventsOff('complete');
      EventsOff('error');
    };
  }, []);

  // Handlers
  const handleSelectProducts = async () => {
    const files = await SelectProductFiles();
    if (files) setProductFiles(files);
  };

  const handleSelectFrame = async () => {
    const file = await SelectFrameFile();
    if (file) setFrameFile(file);
  };

  const handleSelectTemplate = async () => {
    const file = await SelectTemplateFile();
    if (file) {
      setTemplateFile(file);
      const fields = await LoadTemplate(file);
      setTemplateFields(fields || []);
      setFieldValues({});
    }
  };

  const handleSelectOutput = async () => {
    const folder = await SelectOutputFolder();
    if (folder) setOutputFolder(folder);
  };

  const handleFieldChange = (field: string, value: string) => {
    setFieldValues((prev) => ({ ...prev, [field]: value }));
  };

  const handlePreview = async () => {
    setIsPreviewLoading(true);
    try {
      const data = await GeneratePreview({
        productImages: productFiles,
        frameImage: frameFile,
        templatePath: templateFile,
        fieldValues,
        outputDir: outputFolder,
        format,
        quality,
      });
      setPreviewImage(data);
    } catch (e: any) {
      alert('Preview error: ' + e);
    } finally {
      setIsPreviewLoading(false);
    }
  };

  const handleGenerate = async () => {
    if (!outputFolder) {
      alert('Please select output folder');
      return;
    }

    setIsProcessing(true);
    setProgress({ current: 0, total: productFiles.length, file: '' });

    try {
      await ProcessBatch({
        productImages: productFiles,
        frameImage: frameFile,
        templatePath: templateFile,
        fieldValues,
        outputDir: outputFolder,
        format,
        quality,
      });
    } catch (e: any) {
      alert('Error: ' + e);
      setIsProcessing(false);
    }
  };

  const handleCancel = () => {
    CancelProcessing();
    setIsProcessing(false);
  };

  const canPreview = productFiles.length > 0 && frameFile !== '';
  const canGenerate = canPreview && outputFolder !== '';

  return (
    <div className="h-screen p-4 flex gap-4">
      {/* Column 1: Inputs */}
      <div className="w-2/5 flex flex-col gap-4 overflow-y-auto">
        <FilePicker
          label="Product Images"
          icon="📁"
          files={productFiles}
          multiple
          onSelect={handleSelectProducts}
          onClear={() => setProductFiles([])}
        />

        <FilePicker
          label="Frame Image"
          icon="🖼️"
          files={frameFile ? [frameFile] : []}
          onSelect={handleSelectFrame}
          onClear={() => setFrameFile('')}
        />

        <FilePicker
          label="Template (Optional)"
          icon="📄"
          files={templateFile ? [templateFile] : []}
          onSelect={handleSelectTemplate}
          onClear={() => {
            setTemplateFile('');
            setTemplateFields([]);
            setFieldValues({});
          }}
        />

        <TemplateFields
          fields={templateFields}
          values={fieldValues}
          onChange={handleFieldChange}
        />
      </div>

      {/* Column 2: Preview & Output */}
      <div className="w-3/5 flex flex-col gap-4 overflow-y-auto">
        <Preview
          imageData={previewImage}
          isLoading={isPreviewLoading}
          onPreview={handlePreview}
          canPreview={canPreview}
        />

        <OutputSettings
          format={format}
          quality={quality}
          outputFolder={outputFolder}
          onFormatChange={setFormat}
          onQualityChange={setQuality}
          onSelectFolder={handleSelectOutput}
        />

        <ProgressBar
          current={progress.current}
          total={progress.total}
          currentFile={progress.file}
          isProcessing={isProcessing}
        />

        {/* Generate Button */}
        <div className="flex gap-2">
          {isProcessing ? (
            <button
              onClick={handleCancel}
              className="flex-1 py-3 bg-red-500 hover:bg-red-600 text-white
                         font-medium rounded-lg"
            >
              Cancel
            </button>
          ) : (
            <button
              onClick={handleGenerate}
              disabled={!canGenerate}
              className="flex-1 py-3 bg-blue-500 hover:bg-blue-600 text-white
                         font-medium rounded-lg disabled:opacity-50
                         disabled:cursor-not-allowed"
            >
              Generate All ({productFiles.length} images)
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

export default App;
```

## Todo List

- [ ] Update `style.css` with Tailwind
- [ ] Create `FilePicker.tsx`
- [ ] Create `TemplateFields.tsx`
- [ ] Create `Preview.tsx`
- [ ] Create `ProgressBar.tsx`
- [ ] Create `OutputSettings.tsx`
- [ ] Create `App.tsx` with full state management
- [ ] Generate Wails bindings (`wails generate module`)
- [ ] Test UI with backend

## Success Criteria

1. 2-column layout renders correctly
2. File pickers open dialogs
3. Template fields appear when template loaded
4. Preview shows image
5. Progress bar updates during batch
6. Generate button triggers processing

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Wails bindings not generated | High | Run wails generate module |
| Events not received | Medium | Verify event names match backend |
| Layout broken on resize | Low | Test with different window sizes |

## Security Considerations

- Sanitize displayed file names
- No user input rendered as HTML

## Next Steps

After completion, proceed to [Phase 7: Integration Testing](./phase-07-integration-testing.md)
