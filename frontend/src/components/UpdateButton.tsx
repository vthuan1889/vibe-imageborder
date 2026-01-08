import { useState, useEffect } from 'react';
import {
  GetVersion,
  CheckForUpdate,
  DownloadAndInstallUpdate,
} from '../../wailsjs/go/main/App';

interface UpdateInfo {
  available: boolean;
  current: string;
  latest: string;
  downloadUrl: string;
}

export function UpdateButton() {
  const [version, setVersion] = useState('');
  const [checking, setChecking] = useState(false);
  const [downloading, setDownloading] = useState(false);
  const [updateInfo, setUpdateInfo] = useState<UpdateInfo | null>(null);

  useEffect(() => {
    GetVersion().then(setVersion);
  }, []);

  const handleCheck = async () => {
    setChecking(true);
    try {
      const info = await CheckForUpdate();
      setUpdateInfo(info);
      if (!info.available) {
        alert('Bạn đang sử dụng phiên bản mới nhất!');
      }
    } catch (e) {
      alert('Không thể kiểm tra cập nhật: ' + e);
    } finally {
      setChecking(false);
    }
  };

  const handleUpdate = async () => {
    if (!updateInfo?.downloadUrl) return;

    const confirmed = window.confirm(
      `Cập nhật lên ${updateInfo.latest}?\n\nỨng dụng sẽ đóng và trình cài đặt sẽ chạy.`
    );
    if (!confirmed) return;

    setDownloading(true);
    try {
      await DownloadAndInstallUpdate(updateInfo.downloadUrl);
    } catch (e) {
      alert('Cập nhật thất bại: ' + e);
      setDownloading(false);
    }
  };

  return (
    <div className="flex items-center gap-2 text-sm">
      <span className="text-gray-500">{version}</span>

      {updateInfo?.available ? (
        <button
          onClick={handleUpdate}
          disabled={downloading}
          className="px-3 py-1 bg-green-500 hover:bg-green-600 text-white
                     rounded text-xs disabled:opacity-50 cursor-pointer"
        >
          {downloading ? 'Đang tải...' : `Cập nhật ${updateInfo.latest}`}
        </button>
      ) : (
        <button
          onClick={handleCheck}
          disabled={checking}
          className="px-3 py-1 bg-gray-200 hover:bg-gray-300 text-gray-700
                     rounded text-xs disabled:opacity-50 cursor-pointer"
        >
          {checking ? 'Đang kiểm tra...' : 'Kiểm tra cập nhật'}
        </button>
      )}
    </div>
  );
}
