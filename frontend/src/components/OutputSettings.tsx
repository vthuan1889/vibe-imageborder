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
