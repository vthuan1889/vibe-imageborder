import { FC } from 'react';
import { AppView, NAV_ITEMS } from '../types/navigation';

interface SidebarProps {
  activeView: AppView;
  onNavigate: (view: AppView) => void;
}

export const Sidebar: FC<SidebarProps> = ({ activeView, onNavigate }) => {
  return (
    <aside className="w-60 shrink-0 bg-white border-r border-gray-200 flex flex-col">
      <div className="p-4 border-b border-gray-200">
        <h1 className="text-lg font-semibold text-gray-800">Vibe Image Border</h1>
        <p className="text-xs text-gray-500 mt-1">Product frame composer</p>
      </div>

      <nav className="flex-1 p-3 flex flex-col gap-1">
        {NAV_ITEMS.map((item) => {
          const isActive = activeView === item.id;

          return (
            <button
              key={item.id}
              onClick={() => onNavigate(item.id)}
              className={`
                w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-left
                transition-colors
                ${isActive
                  ? 'bg-blue-50 text-blue-700 border-l-4 border-blue-500 pl-2'
                  : 'text-gray-700 hover:bg-gray-50 border-l-4 border-transparent pl-2'
                }
                ${item.disabled && !isActive ? 'opacity-70' : ''}
              `}
            >
              <span className="text-lg">{item.icon}</span>
              <span className="flex-1 font-medium text-sm">{item.label}</span>
              {item.disabled && (
                <span className="text-xs bg-gray-100 text-gray-500 px-1.5 py-0.5 rounded">
                  Soon
                </span>
              )}
            </button>
          );
        })}
      </nav>
    </aside>
  );
};
