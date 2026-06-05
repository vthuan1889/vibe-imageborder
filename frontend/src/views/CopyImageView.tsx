import { useState, useEffect } from 'react';
import { FolderPicker } from '../components/FolderPicker';
import { ProgressBar } from '../components/ProgressBar';
import { StorageService } from '../utils/storage';

import {
  SelectOutputFolder,
  CopyImages,
  CancelProcessing,
} from '../../wailsjs/go/main/App';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

export function CopyImageView() {
  const [sourceFolder, setSourceFolder] = useState(
    StorageService.loadCopySourceFolder()
  );
  const [destFolder, setDestFolder] = useState(
    StorageService.loadCopyDestFolder()
  );

  const [isProcessing, setIsProcessing] = useState(false);
  const [progress, setProgress] = useState({ current: 0, total: 0, file: '' });
  const [lastResult, setLastResult] = useState<string | null>(null);

  useEffect(() => {
    EventsOn('progress', (data: { current: number; total: number; file: string }) => {
      setProgress({ current: data.current, total: data.total, file: data.file });
    });

    EventsOn('complete', (data: { totalProcessed?: number; totalFailed?: number; outputDir?: string }) => {
      setIsProcessing(false);
      const processed = data?.totalProcessed ?? 0;
      const failed = data?.totalFailed ?? 0;
      setLastResult(
        `Đã copy ${processed} ảnh${failed > 0 ? `, ${failed} lỗi` : ''} → ${data?.outputDir ?? destFolder}`
      );
    });

    EventsOn('error', (data: { message: string }) => {
      alert('Lỗi: ' + data.message);
      setIsProcessing(false);
    });

    EventsOn('cancelled', () => {
      setIsProcessing(false);
      setLastResult('Đã hủy copy.');
    });

    return () => {
      EventsOff('progress');
      EventsOff('complete');
      EventsOff('error');
      EventsOff('cancelled');
    };
  }, [destFolder]);

  const handleSelectSource = async () => {
    try {
      const folder = await SelectOutputFolder();
      if (folder) {
        setSourceFolder(folder);
        StorageService.saveCopySourceFolder(folder);
      }
    } catch (e) {
      console.error('Failed to select source folder:', e);
    }
  };

  const handleSelectDest = async () => {
    try {
      const folder = await SelectOutputFolder();
      if (folder) {
        setDestFolder(folder);
        StorageService.saveCopyDestFolder(folder);
      }
    } catch (e) {
      console.error('Failed to select destination folder:', e);
    }
  };

  const handleCopy = async () => {
    if (!sourceFolder) {
      alert('Vui lòng chọn thư mục hình gốc');
      return;
    }
    if (!destFolder) {
      alert('Vui lòng chọn thư mục đích');
      return;
    }
    if (sourceFolder === destFolder) {
      alert('Thư mục gốc và thư mục đích phải khác nhau');
      return;
    }

    setIsProcessing(true);
    setProgress({ current: 0, total: 0, file: '' });
    setLastResult(null);

    try {
      await CopyImages({ sourceDir: sourceFolder, destDir: destFolder });
    } catch (e) {
      alert('Lỗi: ' + e);
      setIsProcessing(false);
    }
  };

  const handleCancel = () => {
    CancelProcessing();
    setIsProcessing(false);
  };

  const canCopy = sourceFolder !== '' && destFolder !== '' && sourceFolder !== destFolder;

  return (
    <div className="h-full flex gap-4 overflow-hidden">
      <div className="w-2/5 flex flex-col gap-4 overflow-y-auto">
        <FolderPicker
          label="Thư mục hình gốc"
          icon="📂"
          folder={sourceFolder}
          placeholder="Chọn thư mục chứa ảnh gốc"
          onSelect={handleSelectSource}
        />

        <FolderPicker
          label="Thư mục đích"
          icon="📁"
          folder={destFolder}
          placeholder="Chọn thư mục copy ảnh tới"
          onSelect={handleSelectDest}
        />

        <div className="bg-blue-50 border border-blue-100 rounded-lg p-4 text-sm text-blue-800">
          <p className="font-medium mb-1">Cách hoạt động</p>
          <ul className="list-disc list-inside space-y-1 text-blue-700">
            <li>Quét đệ quy tất cả ảnh trong thư mục gốc</li>
            <li>Giữ nguyên cấu trúc subfolder ở thư mục đích</li>
            <li>Chỉnh sửa vô hình để social nhận diện là ảnh khác</li>
            <li>File .webp sẽ được lưu dưới dạng .png</li>
          </ul>
        </div>
      </div>

      <div className="w-3/5 flex flex-col gap-4 overflow-y-auto">
        <div className="bg-white rounded-lg border border-gray-200 p-6 flex-1 min-h-[200px] flex flex-col justify-center items-center text-center">
          {isProcessing ? (
            <div>
              <p className="text-lg font-medium text-gray-700 mb-2">Đang copy ảnh...</p>
              <p className="text-sm text-gray-500">
                Ảnh được chỉnh sửa nhẹ để đổi fingerprint, mắt người không thấy khác biệt.
              </p>
            </div>
          ) : lastResult ? (
            <div>
              <p className="text-lg font-medium text-green-700 mb-2">Hoàn tất</p>
              <p className="text-sm text-gray-600">{lastResult}</p>
            </div>
          ) : (
            <div>
              <p className="text-lg font-medium text-gray-700 mb-2">Copy hình</p>
              <p className="text-sm text-gray-500">
                Chọn thư mục gốc và đích, sau đó nhấn Copy để bắt đầu.
              </p>
            </div>
          )}
        </div>

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
              Hủy
            </button>
          ) : (
            <button
              onClick={handleCopy}
              disabled={!canCopy}
              className="flex-1 py-3 bg-blue-500 hover:bg-blue-600 text-white
                         font-medium rounded-lg disabled:opacity-50
                         disabled:cursor-not-allowed"
            >
              Copy
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
