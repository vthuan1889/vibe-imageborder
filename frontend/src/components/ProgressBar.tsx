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
