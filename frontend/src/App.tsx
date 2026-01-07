import { useState, useEffect } from 'react';
import { App as AppBindings, ProcessRequest, TemplateInfo, ProcessResult, ProgressUpdate } from '../bindings/vibe-imageborder';
import { Events } from '@wailsio/runtime';
import FilePicker from './components/FilePicker';
import TemplateFields from './components/TemplateFields';
import ProgressBar from './components/ProgressBar';
import type { Template } from '../bindings/vibe-imageborder/internal/models/models';

function App() {
  const [productPaths, setProductPaths] = useState<string[]>([]);
  const [framePath, setFramePath] = useState<string>('');
  const [templatePath, setTemplatePath] = useState<string>('');
  const [outputDir, setOutputDir] = useState<string>('');

  const [template, setTemplate] = useState<Template | null>(null);
  const [fields, setFields] = useState<string[]>([]);
  const [fieldValues, setFieldValues] = useState<Record<string, string>>({});

  const [processing, setProcessing] = useState<boolean>(false);
  const [progress, setProgress] = useState<ProgressUpdate>({ current: 0, total: 0, filename: '', status: '' });
  const [result, setResult] = useState<ProcessResult | null>(null);

  // Load template when selected
  useEffect(() => {
    if (!templatePath) return;

    const loadTemplate = async () => {
      try {
        const info: TemplateInfo = await AppBindings.ParseTemplateFile(templatePath);
        if (info.success) {
          setTemplate(info.template);
          setFields(info.fields);
          setFieldValues({});
        } else {
          alert(`Template error: ${info.error}`);
        }
      } catch (err) {
        alert(`Failed to load template: ${err}`);
      }
    };

    loadTemplate();
  }, [templatePath]);

  // Listen for progress events
  useEffect(() => {
    const handleProgress = (data: ProgressUpdate) => {
      setProgress(data);
    };

    Events.On('processing:progress', handleProgress);

    return () => {
      Events.Off('processing:progress', handleProgress);
    };
  }, []);

  const canProcess = (): boolean => {
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

    if (!template) {
      alert('Template not loaded');
      return;
    }

    setProcessing(true);
    setResult(null);
    setProgress({ current: 0, total: productPaths.length, filename: '', status: 'processing' });

    const req: ProcessRequest = {
      productPaths: productPaths,
      framePath: framePath,
      template: template,
      fieldValues: fieldValues,
      outputDir: outputDir,
    };

    try {
      const processResult: ProcessResult = await AppBindings.ProcessBatch(req);

      setProcessing(false);
      setResult(processResult);

      if (processResult.success) {
        alert(
          `Processing complete!\nProcessed: ${processResult.processedCount}\nFailed: ${processResult.failedCount}`
        );
      } else {
        alert(`Processing failed: ${processResult.error}`);
      }
    } catch (err) {
      setProcessing(false);
      alert(`Error during processing: ${err}`);
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
            label="Product Images (comma-separated paths)"
            type="products"
            value={productPaths}
            onChange={(value) => setProductPaths(value as string[])}
            disabled={processing}
          />

          <FilePicker
            label="Frame Image Path"
            type="frame"
            value={framePath}
            onChange={(value) => setFramePath(value as string)}
            disabled={processing}
          />

          <FilePicker
            label="Template File Path"
            type="template"
            value={templatePath}
            onChange={(value) => setTemplatePath(value as string)}
            disabled={processing}
          />

          <FilePicker
            label="Output Directory"
            type="output"
            value={outputDir}
            onChange={(value) => setOutputDir(value as string)}
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
        {processing && progress.total > 0 && (
          <ProgressBar
            current={progress.current}
            total={progress.total}
            filename={progress.filename}
            status={progress.status as 'processing' | 'success' | 'error'}
          />
        )}

        {/* Result */}
        {result && result.success && (
          <div className="bg-green-50 border border-green-200 p-4 rounded-lg">
            <h3 className="text-lg font-semibold text-green-800 mb-2">
              ✓ Processing Complete
            </h3>
            <p className="text-green-700">
              Processed: {result.processedCount} images
              {result.failedCount > 0 && ` (${result.failedCount} failed)`}
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
