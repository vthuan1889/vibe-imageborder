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
