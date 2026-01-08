import { useState, useEffect } from 'react';
import { FilePicker } from './components/FilePicker';
import { TemplateFields } from './components/TemplateFields';
import { Preview } from './components/Preview';
import { ProgressBar } from './components/ProgressBar';
import { OutputSettings } from './components/OutputSettings';
import { UpdateButton } from './components/UpdateButton';

import {
  SelectProductFiles,
  SelectFrameFile,
  SelectTemplateFile,
  SelectOutputFolder,
  GetDefaultOutputFolder,
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
    // Load default output folder on startup
    GetDefaultOutputFolder().then((folder) => {
      if (folder) setOutputFolder(folder);
    });

    EventsOn('progress', (data: { current: number; total: number; file: string }) => {
      setProgress({ current: data.current, total: data.total, file: data.file });
    });

    EventsOn('complete', () => {
      setIsProcessing(false);
    });

    EventsOn('error', (data: { message: string }) => {
      alert('Error: ' + data.message);
      setIsProcessing(false);
    });

    EventsOn('cancelled', () => {
      setIsProcessing(false);
    });

    return () => {
      EventsOff('progress');
      EventsOff('complete');
      EventsOff('error');
      EventsOff('cancelled');
    };
  }, []);

  // Handlers
  const handleSelectProducts = async () => {
    try {
      const files = await SelectProductFiles();
      if (files && files.length > 0) setProductFiles(files);
    } catch (e) {
      console.error('Failed to select products:', e);
    }
  };

  const handleSelectFrame = async () => {
    try {
      const file = await SelectFrameFile();
      if (file) setFrameFile(file);
    } catch (e) {
      console.error('Failed to select frame:', e);
    }
  };

  const handleSelectTemplate = async () => {
    try {
      const file = await SelectTemplateFile();
      if (file) {
        setTemplateFile(file);
        const fields = await LoadTemplate(file);
        setTemplateFields(fields || []);
        setFieldValues({});
      }
    } catch (e) {
      console.error('Failed to load template:', e);
    }
  };

  const handleSelectOutput = async () => {
    try {
      const folder = await SelectOutputFolder();
      if (folder) setOutputFolder(folder);
    } catch (e) {
      console.error('Failed to select output folder:', e);
    }
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
    } catch (e) {
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
    } catch (e) {
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
    <div className="h-screen p-4 flex flex-col">
      {/* Header with version and update button */}
      <div className="flex justify-end mb-4">
        <UpdateButton />
      </div>

      <div className="flex-1 flex gap-4 overflow-hidden">
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
    </div>
  );
}

export default App;
