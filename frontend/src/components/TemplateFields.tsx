import { FC } from 'react';

interface TemplateFieldsProps {
  fields: string[];
  values: Record<string, string>;
  onChange: (field: string, value: string) => void;
}

export const TemplateFields: FC<TemplateFieldsProps> = ({
  fields,
  values,
  onChange,
}) => {
  if (fields.length === 0) return null;

  // Format field name for display
  const formatLabel = (field: string) => {
    return field
      .replace(/_/g, ' ')
      .replace(/\b\w/g, (c) => c.toUpperCase());
  };

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-4">
      <div className="flex items-center gap-2 mb-3">
        <span className="text-lg">📝</span>
        <span className="font-medium text-gray-700">Text Fields</span>
      </div>

      <div className="space-y-3">
        {fields.map((field) => (
          <div key={field}>
            <label className="block text-sm text-gray-600 mb-1">
              {formatLabel(field)}
            </label>
            <input
              type="text"
              value={values[field] || ''}
              onChange={(e) => onChange(field, e.target.value)}
              placeholder={`Enter ${field}`}
              className="w-full px-3 py-2 border border-gray-300 rounded-md
                         focus:ring-2 focus:ring-blue-500 focus:border-transparent
                         text-sm"
            />
          </div>
        ))}
      </div>
    </div>
  );
};
