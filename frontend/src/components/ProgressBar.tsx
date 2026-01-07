import React from 'react';

interface ProgressBarProps {
  current: number;
  total: number;
  filename?: string;
  status?: 'processing' | 'success' | 'error';
}

export const ProgressBar: React.FC<ProgressBarProps> = ({
  current,
  total,
  filename,
  status = 'processing',
}) => {
  if (!total || total === 0) {
    return null;
  }

  const percentage = Math.round((current / total) * 100);

  const statusColors: Record<string, string> = {
    processing: 'bg-blue-600',
    success: 'bg-green-600',
    error: 'bg-red-600',
  };

  const statusText: Record<string, string> = {
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
        <span className="text-sm font-medium text-gray-600">{percentage}%</span>
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
};

export default ProgressBar;
