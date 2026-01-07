import React, { useState, useEffect } from 'react';

interface FilePickerProps {
  label: string;
  type: 'products' | 'frame' | 'template' | 'output';
  value: string | string[];
  onChange: (value: string | string[]) => void;
  disabled?: boolean;
}

export const FilePicker: React.FC<FilePickerProps> = ({
  label,
  type,
  value,
  onChange,
  disabled = false,
}) => {
  const [displayValue, setDisplayValue] = useState<string>('');

  useEffect(() => {
    if (type === 'products' && Array.isArray(value) && value.length > 0) {
      setDisplayValue(`${value.length} file(s) selected`);
    } else if (typeof value === 'string' && value) {
      setDisplayValue(value);
    } else {
      setDisplayValue('');
    }
  }, [value, type]);

  const handleClick = () => {
    // For manual input instead of dialog (workaround)
    let result: string | string[] | null = null;

    if (type === 'products') {
      const input = prompt('Enter product image paths (comma-separated):');
      if (input) {
        result = input.split(',').map((p) => p.trim()).filter((p) => p);
        onChange(result);
      }
    } else if (type === 'output') {
      const input = prompt('Enter output directory path:');
      if (input) {
        result = input.trim();
        onChange(result);
      }
    } else {
      const input = prompt(`Enter ${label.toLowerCase()} path:`);
      if (input) {
        result = input.trim();
        onChange(result);
      }
    }
  };

  return (
    <div className="mb-4">
      <label className="block text-sm font-medium text-gray-700 mb-2">
        {label}
      </label>
      <div className="flex gap-2">
        <input
          type="text"
          value={displayValue}
          readOnly
          placeholder="No file selected"
          className="flex-1 px-3 py-2 border border-gray-300 rounded-md bg-gray-50 text-gray-700"
        />
        <button
          onClick={handleClick}
          disabled={disabled}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          Select...
        </button>
      </div>
    </div>
  );
};

export default FilePicker;
