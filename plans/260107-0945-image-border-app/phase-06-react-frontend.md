# Phase 6: React Frontend - UI Components

**Goal:** Build React UI với file pickers, dynamic form, progress bar

**Duration:** ~4-5 hours

**Dependencies:** Phase 1-5 complete

---

## Overview

Create React components for:
1. File selection (products, frame, template, output dir)
2. Dynamic form generation từ template fields
3. Progress tracking during processing
4. Result display

---

## Task 6.1: Setup Wails Runtime Binding

Create `frontend/src/wails-runtime.js`:

```js
// Wails v3 runtime helpers
export const SelectProductImages = window.wails.SelectProductImages;
export const SelectFrameImage = window.wails.SelectFrameImage;
export const SelectTemplateFile = window.wails.SelectTemplateFile;
export const SelectOutputDirectory = window.wails.SelectOutputDirectory;
export const ParseTemplateFile = window.wails.ParseTemplateFile;
export const ProcessBatch = window.wails.ProcessBatch;

// Event listener
export const EventsOn = (eventName, callback) => {
  window.wails.EventsOn(eventName, callback);
};

export const EventsOff = (eventName) => {
  window.wails.EventsOff(eventName);
};
```

---

## Task 6.2: Create FilePicker Component

Create `frontend/src/components/FilePicker.jsx`:

```jsx
import React from 'react';

export default function FilePicker({ label, type, value, onChange, disabled }) {
  const [displayValue, setDisplayValue] = React.useState('');

  React.useEffect(() => {
    if (type === 'products' && value?.length > 0) {
      setDisplayValue(`${value.length} file(s) selected`);
    } else if (value) {
      setDisplayValue(value);
    } else {
      setDisplayValue('');
    }
  }, [value, type]);

  const handleClick = async () => {
    let result;

    switch (type) {
      case 'products':
        result = await window.wails.SelectProductImages();
        if (result.Success) {
          onChange(result.Paths);
        }
        break;
      case 'frame':
        result = await window.wails.SelectFrameImage();
        if (result.Success) {
          onChange(result.Paths[0]);
        }
        break;
      case 'template':
        result = await window.wails.SelectTemplateFile();
        if (result.Success) {
          onChange(result.Paths[0]);
        }
        break;
      case 'output':
        result = await window.wails.SelectOutputDirectory();
        if (result.Success) {
          onChange(result.Paths[0]);
        }
        break;
    }
  };

  return (
    <div className="mb-4">
      <label className="block text-sm font-medium text-gray-700 mb-2">
        {label}
      </label>
      <div className="flex gap-2">
        <input
          type="text"
          value={displayValue}
          readOnly
          placeholder="No file selected"
          className="flex-1 px-3 py-2 border border-gray-300 rounded-md bg-gray-50 text-gray-700"
        />
        <button
          onClick={handleClick}
          disabled={disabled}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          Browse...
        </button>
      </div>
    </div>
  );
}
```

---

## Task 6.3: Create TemplateFields Component

Create `frontend/src/components/TemplateFields.jsx`:

```jsx
import React from 'react';

export default function TemplateFields({ fields, values, onChange, disabled }) {
  if (!fields || fields.length === 0) {
    return null;
  }

  const handleChange = (field, value) => {
    onChange({ ...values, [field]: value });
  };

  return (
    <div className="bg-white p-4 rounded-lg border border-gray-200">
      <h3 className="text-lg font-semibold text-gray-800 mb-3">
        Template Fields
      </h3>
      <div className="grid grid-cols-2 gap-4">
        {fields.map((field) => (
          <div key={field}>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {field.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())}
            </label>
            <input
              type="text"
              value={values[field] || ''}
              onChange={(e) => handleChange(field, e.target.value)}
              disabled={disabled}
              placeholder={`Enter ${field}`}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
            />
          </div>
        ))}
      </div>
    </div>
  );
}
```

---

## Task 6.4: Create ProgressBar Component

Create `frontend/src/components/ProgressBar.jsx`:

```jsx
import React from 'react';

export default function ProgressBar({ current, total, filename, status }) {
  if (!total || total === 0) {
    return null;
  }

  const percentage = Math.round((current / total) * 100);

  const statusColors = {
    processing: 'bg-blue-600',
    success: 'bg-green-600',
    error: 'bg-red-600',
  };

  const statusText = {
    processing: 'Processing...',
    success: 'Complete!',
    error: 'Error',
  };

  return (
    <div className="bg-white p-4 rounded-lg border border-gray-200">
      <div className="flex justify-between items-center mb-2">
        <span className="text-sm font-medium text-gray-700">
          {current} / {total} images
        </span>
        <span className="text-sm font-medium text-gray-600">
          {percentage}%
        </span>
      </div>

      <div className="w-full bg-gray-200 rounded-full h-3 mb-2">
        <div
          className={`h-3 rounded-full transition-all duration-300 ${
            statusColors[status] || statusColors.processing
          }`}
          style={{ width: `${percentage}%` }}
        />
      </div>

      {filename && (
        <div className="text-sm text-gray-600">
          <span className="font-medium">{statusText[status] || ''}</span>{' '}
          {filename}
        </div>
      )}
    </div>
  );
}
```

---

## Task 6.5: Create Main App Component

Update `frontend/src/App.jsx`:

```jsx
import React, { useState, useEffect } from 'react';
import FilePicker from './components/FilePicker';
import TemplateFields from './components/TemplateFields';
import ProgressBar from './components/ProgressBar';

function App() {
  const [productPaths, setProductPaths] = useState([]);
  const [framePath, setFramePath] = useState('');
  const [templatePath, setTemplatePath] = useState('');
  const [outputDir, setOutputDir] = useState('');

  const [template, setTemplate] = useState(null);
  const [fields, setFields] = useState([]);
  const [fieldValues, setFieldValues] = useState({});

  const [processing, setProcessing] = useState(false);
  const [progress, setProgress] = useState({ current: 0, total: 0 });
  const [result, setResult] = useState(null);

  // Load template when selected
  useEffect(() => {
    if (!templatePath) return;

    const loadTemplate = async () => {
      const info = await window.wails.ParseTemplateFile(templatePath);
      if (info.Success) {
        setTemplate(info.Template);
        setFields(info.Fields);
        setFieldValues({});
      } else {
        alert(`Template error: ${info.Error}`);
      }
    };

    loadTemplate();
  }, [templatePath]);

  // Listen for progress events
  useEffect(() => {
    const handleProgress = (data) => {
      setProgress(data);
    };

    window.wails.EventsOn('processing:progress', handleProgress);

    return () => {
      window.wails.EventsOff('processing:progress');
    };
  }, []);

  const canProcess = () => {
    if (!productPaths.length || !framePath || !templatePath || !outputDir) {
      return false;
    }

    // Check all fields filled
    for (const field of fields) {
      if (!fieldValues[field]) {
        return false;
      }
    }

    return true;
  };

  const handleProcess = async () => {
    if (!canProcess()) {
      alert('Please fill all required fields');
      return;
    }

    setProcessing(true);
    setResult(null);
    setProgress({ current: 0, total: productPaths.length });

    const req = {
      ProductPaths: productPaths,
      FramePath: framePath,
      Template: template,
      FieldValues: fieldValues,
      OutputDir: outputDir,
    };

    const processResult = await window.wails.ProcessBatch(req);

    setProcessing(false);
    setResult(processResult);

    if (processResult.Success) {
      alert(
        `Processing complete!\nProcessed: ${processResult.ProcessedCount}\nFailed: ${processResult.FailedCount}`
      );
    } else {
      alert(`Processing failed: ${processResult.Error}`);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">
          Image Border Application
        </h1>

        {/* File Selection */}
        <div className="bg-white p-6 rounded-lg shadow-sm mb-6">
          <h2 className="text-xl font-semibold text-gray-800 mb-4">
            1. Select Files
          </h2>

          <FilePicker
            label="Product Images"
            type="products"
            value={productPaths}
            onChange={setProductPaths}
            disabled={processing}
          />

          <FilePicker
            label="Frame Image"
            type="frame"
            value={framePath}
            onChange={setFramePath}
            disabled={processing}
          />

          <FilePicker
            label="Template File"
            type="template"
            value={templatePath}
            onChange={setTemplatePath}
            disabled={processing}
          />

          <FilePicker
            label="Output Directory"
            type="output"
            value={outputDir}
            onChange={setOutputDir}
            disabled={processing}
          />
        </div>

        {/* Template Fields */}
        {fields.length > 0 && (
          <div className="mb-6">
            <h2 className="text-xl font-semibold text-gray-800 mb-4">
              2. Fill Template Fields
            </h2>
            <TemplateFields
              fields={fields}
              values={fieldValues}
              onChange={setFieldValues}
              disabled={processing}
            />
          </div>
        )}

        {/* Process Button */}
        <div className="bg-white p-6 rounded-lg shadow-sm mb-6">
          <button
            onClick={handleProcess}
            disabled={!canProcess() || processing}
            className="w-full py-3 bg-green-600 text-white text-lg font-semibold rounded-md hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            {processing ? 'Processing...' : 'Process Images'}
          </button>
        </div>

        {/* Progress */}
        {processing && (
          <ProgressBar
            current={progress.Current}
            total={progress.Total}
            filename={progress.Filename}
            status={progress.Status}
          />
        )}

        {/* Result */}
        {result && result.Success && (
          <div className="bg-green-50 border border-green-200 p-4 rounded-lg">
            <h3 className="text-lg font-semibold text-green-800 mb-2">
              ✓ Processing Complete
            </h3>
            <p className="text-green-700">
              Processed: {result.ProcessedCount} images
              {result.FailedCount > 0 && ` (${result.FailedCount} failed)`}
            </p>
            <p className="text-sm text-green-600 mt-2">
              Output location: {outputDir}
            </p>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
```

---

## Task 6.6: Update Styles

Update `frontend/src/style.css`:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

body {
  @apply bg-gray-50;
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  height: 100vh;
  overflow-y: auto;
}
```

---

## Acceptance Criteria

- ✓ File pickers open native dialogs
- ✓ Selected files display trong UI
- ✓ Template fields generate dynamically
- ✓ All fields required before processing enabled
- ✓ Progress bar updates in real-time
- ✓ Result message displays after processing
- ✓ UI responsive và clean (TailwindCSS)

---

## Deliverables

### Files Created/Modified

1. `frontend/src/components/FilePicker.jsx`
2. `frontend/src/components/TemplateFields.jsx`
3. `frontend/src/components/ProgressBar.jsx`
4. `frontend/src/App.jsx`
5. `frontend/src/style.css`

### Validation

```bash
# Run dev mode
wails3 dev

# Manual test checklist:
# ✓ Click "Browse" for each file picker
# ✓ Select template → fields appear
# ✓ Fill all fields → Process button enables
# ✓ Click Process → progress bar shows
# ✓ Complete → success message displays
```

---

## Next Phase

[Phase 7: Integration & Testing](phase-07-integration-testing.md)
