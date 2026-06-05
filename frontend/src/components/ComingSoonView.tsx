import { FC } from 'react';

interface ComingSoonViewProps {
  title: string;
  description: string;
}

export const ComingSoonView: FC<ComingSoonViewProps> = ({ title, description }) => {
  return (
    <div className="h-full flex items-center justify-center">
      <div className="bg-white rounded-lg border border-gray-200 p-8 max-w-md text-center">
        <div className="text-4xl mb-4">🚧</div>
        <h2 className="text-xl font-semibold text-gray-800 mb-2">{title}</h2>
        <p className="text-gray-600 mb-4">{description}</p>
        <p className="text-sm text-gray-400">This feature is under development.</p>
      </div>
    </div>
  );
};
