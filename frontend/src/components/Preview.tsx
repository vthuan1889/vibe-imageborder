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
