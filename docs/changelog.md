# Nhật ký thay đổi dự án Vibe Image Border

## Unreleased

_(Chưa có thay đổi mới.)_

## 1.0.6 - 2026-06-05

### Đã thêm
- **Copy Image** — copy đệ quy ảnh từ thư mục gốc sang thư mục đích, giữ cấu trúc subfolder.
- Chỉnh sửa pixel vô hình (LSB noise) để đổi fingerprint; mắt người không thấy khác biệt.
- `CopyImageView`, `FolderPicker`, API `CopyImages` và `Uniquify` backend.
- Sidebar điều hướng bên trái với **Create Frame**, **Copy Image**, **Frame Library**, **Settings**.
- `CreateFrameView` — tách toàn bộ luồng ghép ảnh từ `App.tsx` sang view riêng.
- `ComingSoonView` — màn placeholder cho **Frame Library** và **Settings**.
- Lưu tab sidebar và đường dẫn folder copy vào `localStorage`.

### Đã thay đổi
- `App.tsx` trở thành layout shell (sidebar + header + view switcher).
- Cấu trúc frontend mở rộng: thêm thư mục `views/`, `types/navigation.ts`.
- File `.webp` khi copy được lưu dưới dạng `.png` (giới hạn encoder).

## 0.1.0 - 2026-01-08

### Đã thêm
- Khởi tạo lộ trình dự án và tài liệu nhật ký thay đổi.
- Giai đoạn 1: Thiết lập dự án & Nền tảng đã hoàn thành.
- Giai đoạn 2: Dịch vụ mẫu đã hoàn thành.
    - Triển khai phân tích cú pháp tệp mẫu JSON.
    - Triển khai trích xuất trường và thay thế giá trị chỗ dành sẵn.
    - Đã thêm các trường bảo mật và đồng thời quan trọng vào danh sách việc cần làm.
    - Đã thêm các kiểm thử đơn vị cơ bản cho phân tích cú pháp và áp dụng giá trị.
