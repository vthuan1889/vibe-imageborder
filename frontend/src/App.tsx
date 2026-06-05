import { useState } from 'react';
import { Sidebar } from './components/Sidebar';
import { ComingSoonView } from './components/ComingSoonView';
import { UpdateButton } from './components/UpdateButton';
import { CreateFrameView } from './views/CreateFrameView';
import { CopyImageView } from './views/CopyImageView';
import { AppView, NAV_ITEMS } from './types/navigation';
import { StorageService } from './utils/storage';

function App() {
  const [activeView, setActiveView] = useState<AppView>(StorageService.loadActiveView());

  const handleNavigate = (view: AppView) => {
    setActiveView(view);
    StorageService.saveActiveView(view);
  };

  const placeholderItem = NAV_ITEMS.find(
    (item) => item.id === activeView && item.disabled
  );

  return (
    <div className="h-screen flex">
      <Sidebar activeView={activeView} onNavigate={handleNavigate} />

      <div className="flex-1 flex flex-col p-4 overflow-hidden">
        <div className="flex justify-end mb-4">
          <UpdateButton />
        </div>

        <div className="flex-1 overflow-hidden">
          {activeView === 'create-frame' && <CreateFrameView />}
          {activeView === 'copy-image' && <CopyImageView />}
          {placeholderItem && (
            <ComingSoonView
              title={placeholderItem.label}
              description={placeholderItem.description || ''}
            />
          )}
        </div>
      </div>
    </div>
  );
}

export default App;
