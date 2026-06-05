import { FC } from 'react';

interface FolderPickerProps {
  label: string;
  icon: string;
  folder: string;
  placeholder?: string;
  onSelect: () => void;
}

export const FolderPicker: FC<FolderPickerProps> = ({
  label,
  icon,
  folder,
  placeholder = 'Select folder',
  onSelect,
}) => {
  return (
    <div className="bg-white rounded-lg border border-gray-200 p-4">
      <div className="flex items-center gap-2 mb-3">
        <span className="text-lg">{icon}</span>
        <span className="font-medium text-gray-700">{label}</span>
      </div>

      <div className="flex gap-2">
        <input
          type="text"
          value={folder}
          readOnly
          placeholder={placeholder}
          className="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm
                     bg-gray-50 truncate"
        />
        <button
          onClick={onSelect}
          className="px-3 py-2 bg-gray-100 hover:bg-gray-200 rounded-md text-sm"
        >
          📁
        </button>
      </div>
    </div>
  );
};
