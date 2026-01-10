# Hướng dẫn CI/CD và Auto-Update cho phần mềm Open Source

## Prompt Template

```
Thiết lập CI/CD và Auto-Update cho [TÊN_APP]:

1. CI/CD với GitHub Actions:
   - Trigger khi push tag v*.*.*
   - Build Windows executable + NSIS installer
   - Upload lên GitHub Releases

2. Auto-Update trong app:
   - Nút "Kiểm tra cập nhật"
   - So sánh version với GitHub Releases API
   - Download installer và chạy với admin rights
   - App tự đóng sau khi chạy installer

Thông tin:
- Repo: [OWNER]/[REPO]
- Tech: [Wails/Electron/Tauri]
- Version hiện tại: v1.0.0
```

---

## Cấu trúc files cần tạo

```
.github/workflows/release.yml   # GitHub Actions workflow
internal/updater/updater.go     # Backend: check/download/install
internal/updater/updater_test.go # Unit tests
frontend/.../UpdateButton.tsx   # Frontend: UI component
main.go                         # Thêm version variable
app.go                          # Thêm bindings cho frontend
```

---

## Các lưu ý quan trọng

### GitHub Actions

| Vấn đề | Giải pháp |
|--------|-----------|
| Wails build flag | Dùng `-platform windows/amd64` (không phải `--target`) |
| NSIS not found | Thêm `C:\Program Files (x86)\NSIS` vào PATH sau khi install |
| Version injection | Dùng ldflags: `-ldflags "-X 'main.version=v1.0.0'"` |
| Upload assets | Dùng `softprops/action-gh-release@v2` |

### Auto-Update Backend

| Vấn đề | Giải pháp |
|--------|-----------|
| GitHub API | `GET https://api.github.com/repos/{owner}/{repo}/releases/latest` |
| Version compare | So sánh semver, strip prefix "v" trước khi compare |
| Download path | Dùng `os.TempDir()` để tránh permission issues |
| Admin rights | PowerShell `Start-Process -FilePath "..." -Verb RunAs` |
| Ẩn cửa sổ PS | `syscall.SysProcAttr{HideWindow: true}` |
| App quit | Gọi `runtime.Quit()` sau khi start installer |

### Auto-Update Frontend

| Vấn đề | Giải pháp |
|--------|-----------|
| Load version | Gọi `GetVersion()` khi component mount |
| Check update | Gọi backend, hiển thị loading state |
| Confirm update | Dùng `confirm()` trước khi download |
| Download state | Hiển thị "Đang tải..." khi downloading |

---

## Flow hoạt động

```
[Push tag v*.*.*]
    ↓
[GitHub Actions trigger]
    ↓
[Build executable + NSIS installer]
    ↓
[Upload to GitHub Releases]
    ↓
[User mở app → Click "Kiểm tra cập nhật"]
    ↓
[App gọi GitHub API → So sánh version]
    ↓
[Có bản mới → User confirm → Download installer]
    ↓
[Chạy installer với admin rights → App quit]
```

---

## Checklist triển khai

- [ ] Thêm version variable vào main entry point
- [ ] Tạo updater package với 3 functions: CheckUpdate, DownloadAndInstall, CompareVersions
- [ ] Thêm bindings: GetVersion, CheckForUpdate, DownloadAndInstallUpdate
- [ ] Tạo UpdateButton component với states: version, checking, downloading, updateInfo
- [ ] Tạo GitHub Actions workflow với: checkout, setup go/node, install wails/nsis, build, upload
- [ ] Test local: build với `-nsis` flag
- [ ] Push code và tạo tag để test CI/CD

---

## Các lỗi thường gặp

| Lỗi | Nguyên nhân | Fix |
|-----|-------------|-----|
| `flag: -target undefined` | Wails v2 syntax khác | Dùng `-platform` |
| `makensis not found` | NSIS chưa trong PATH | Thêm vào PATH trong workflow |
| `requires elevation` | NSIS installer cần admin | Dùng PowerShell RunAs |
| `syscall undefined` | Cross-platform issue | Build tag cho Windows only |
| `404 từ GitHub API` | Chưa có release nào | Tạo release đầu tiên |

---

## Mở rộng (tùy chọn)

- **Multi-platform**: Thêm jobs cho macOS/Linux
- **Checksum verification**: Download SHA256 và verify trước khi install
- **Silent install**: Thêm `/S` flag cho NSIS
- **Progress bar**: Tracking download progress
- **Auto-check on startup**: Kiểm tra update khi app khởi động
