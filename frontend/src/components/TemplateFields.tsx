import React from 'react';

interface TemplateFieldsProps {
  fields: string[];
  values: Record<string, string>;
  onChange: (values: Record<string, string>) => void;
  disabled?: boolean;
}

export const TemplateFields: React.FC<TemplateFieldsProps> = ({
  fields,
  values,
  onChange,
  disabled = false,
}) => {
  if (!fields || fields.length === 0) {
    return null;
  }

  const handleChange = (field: string, value: string) => {
    onChange({ ...values, [field]: value });
  };

  // Format field name: size_dai → Size Dai
  const formatFieldName = (field: string): string => {
    return field
      .replace(/_/g, ' ')
      .replace(/\b\w/g, (c) => c.toUpperCase());
  };

  return (
    <div className="bg-white p-4 rounded-lg border border-gray-200">
      <h3 className="text-lg font-semibold text-gray-800 mb-3">
        Template Fields
      </h3>
      <div className="grid grid-cols-2 gap-4">
        {fields.map((field) => (
          <div key={field}>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {formatFieldName(field)}
            </label>
            <input
              type="text"
              value={values[field] || ''}
              onChange={(e) => handleChange(field, e.target.value)}
              disabled={disabled}
              placeholder={`Enter ${field}`}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
            />
          </div>
        ))}
      </div>
    </div>
  );
};

export default TemplateFields;
