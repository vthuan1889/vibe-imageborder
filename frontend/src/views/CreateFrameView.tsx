import { useState, useEffect } from 'react';
import { FilePicker } from '../components/FilePicker';
import { TemplateFields } from '../components/TemplateFields';
import { Preview } from '../components/Preview';
import { ProgressBar } from '../components/ProgressBar';
import { OutputSettings } from '../components/OutputSettings';
import { StorageService } from '../utils/storage';

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
} from '../../wailsjs/go/main/App';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

export function CreateFrameView() {
  const savedState = StorageService.loadState();

  const [productFiles, setProductFiles] = useState<string[]>(savedState.productFiles || []);
  const [frameFile, setFrameFile] = useState<string>(savedState.frameFile || '');
  const [templateFile, setTemplateFile] = useState<string>(savedState.templateFile || '');

  const [templateFields, setTemplateFields] = useState<string[]>(savedState.templateFields || []);
  const [fieldValues, setFieldValues] = useState<Record<string, string>>(savedState.fieldValues || {});

  const [format, setFormat] = useState(savedState.format || 'png');
  const [quality, setQuality] = useState(savedState.quality || 90);
  const [outputFolder, setOutputFolder] = useState(savedState.outputFolder || '');

  const [previewImage, setPreviewImage] = useState<string | null>(null);
  const [isPreviewLoading, setIsPreviewLoading] = useState(false);

  const [isProcessing, setIsProcessing] = useState(false);
  const [progress, setProgress] = useState({ current: 0, total: 0, file: '' });

  useEffect(() => {
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

  const handleSelectProducts = async () => {
    try {
      const files = await SelectProductFiles();
      if (files && files.length > 0) {
        setProductFiles(files);
        StorageService.saveProductFiles(files);
      }
    } catch (e) {
      console.error('Failed to select products:', e);
    }
  };

  const handleSelectFrame = async () => {
    try {
      const file = await SelectFrameFile();
      if (file) {
        setFrameFile(file);
        StorageService.saveFrameFile(file);
      }
    } catch (e) {
      console.error('Failed to select frame:', e);
    }
  };

  const handleSelectTemplate = async () => {
    try {
      const file = await SelectTemplateFile();
      if (file) {
        setTemplateFile(file);
        StorageService.saveTemplateFile(file);
        const fields = await LoadTemplate(file);
        setTemplateFields(fields || []);
        StorageService.saveTemplateFields(fields || []);
        setFieldValues({});
        StorageService.saveFieldValues({});
      }
    } catch (e) {
      console.error('Failed to load template:', e);
    }
  };

  const handleSelectOutput = async () => {
    try {
      const folder = await SelectOutputFolder();
      if (folder) {
        setOutputFolder(folder);
        StorageService.saveOutputFolder(folder);
      }
    } catch (e) {
      console.error('Failed to select output folder:', e);
    }
  };

  const handleFieldChange = (field: string, value: string) => {
    setFieldValues((prev) => {
      const updated = { ...prev, [field]: value };
      StorageService.saveFieldValues(updated);
      return updated;
    });
  };

  const handlePreview = async () => {
    setIsPreviewLoading(true);
    try {
      const filledValues = Object.fromEntries(
        Object.entries(fieldValues).filter(([, value]) => value.trim() !== '')
      );
      const data = await GeneratePreview({
        productImages: productFiles,
        frameImage: frameFile,
        templatePath: templateFile,
        fieldValues: filledValues,
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
      const filledValues = Object.fromEntries(
        Object.entries(fieldValues).filter(([, value]) => value.trim() !== '')
      );
      await ProcessBatch({
        productImages: productFiles,
        frameImage: frameFile,
        templatePath: templateFile,
        fieldValues: filledValues,
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
    <div className="h-full flex gap-4 overflow-hidden">
      <div className="w-2/5 flex flex-col gap-4 overflow-y-auto">
        <FilePicker
          label="Product Images"
          icon="📁"
          files={productFiles}
          multiple
          onSelect={handleSelectProducts}
          onClear={() => {
            setProductFiles([]);
            StorageService.saveProductFiles([]);
          }}
        />

        <FilePicker
          label="Frame Image"
          icon="🖼️"
          files={frameFile ? [frameFile] : []}
          onSelect={handleSelectFrame}
          onClear={() => {
            setFrameFile('');
            StorageService.saveFrameFile('');
          }}
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
            StorageService.saveTemplateFile('');
            StorageService.saveTemplateFields([]);
            StorageService.saveFieldValues({});
          }}
        />

        <TemplateFields
          fields={templateFields}
          values={fieldValues}
          onChange={handleFieldChange}
        />
      </div>

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
          onFormatChange={(newFormat) => {
            setFormat(newFormat);
            StorageService.saveFormat(newFormat);
          }}
          onQualityChange={(newQuality) => {
            setQuality(newQuality);
            StorageService.saveQuality(newQuality);
          }}
          onSelectFolder={handleSelectOutput}
        />

        <ProgressBar
          current={progress.current}
          total={progress.total}
          currentFile={progress.file}
          isProcessing={isProcessing}
        />

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
