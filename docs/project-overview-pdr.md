# Tổng quan dự án và PDR

## Tổng quan dự án

Dự án "vibe-imageborder" là một ứng dụng lai (hybrid application) được phát triển bằng khung công tác Wails, kết hợp các khả năng của Go backend với giao diện người dùng frontend dựa trên web. Mục tiêu chính là cung cấp một công cụ để xử lý hình ảnh, cụ thể là thêm đường viền vào hình ảnh.

### Giai đoạn 1: Thiết lập nền tảng

Giai đoạn 1 của dự án tập trung vào việc thiết lập kiến trúc cơ bản, tích hợp các công nghệ backend và frontend, và đảm bảo một môi trường phát triển chức năng. Các thành phần chính được thiết lập trong giai đoạn này bao gồm:

-   **Backend (Go)**:
    -   Sử dụng `go.mod` để quản lý các phần phụ thuộc Go.
    -   Tích hợp Wails v3 để kết nối backend Go với frontend.
    -   Sử dụng thư viện `github.com/disintegration/imaging` và `github.com/fogleman/gg` cho các hoạt động xử lý hình ảnh.
    -   Tệp `internal/models/types.go` đóng vai trò giữ chỗ cho các định nghĩa kiểu backend, chuẩn bị cho các tính năng trong tương lai.

-   **Frontend (React/TypeScript)**:
    -   Được xây dựng bằng React và TypeScript để tạo giao diện người dùng tương tác.
    -   Vite được sử dụng làm công cụ xây dựng để phát triển nhanh.
    -   Tailwind CSS được tích hợp để tạo kiểu nhanh chóng và tùy chỉnh.
    -   Tệp `frontend/package.json` quản lý các phần phụ thuộc frontend.
    -   `frontend/tailwind.config.js` cấu hình Tailwind CSS.
    -   `frontend/src/style.css` bao gồm các chỉ thị Tailwind cơ bản và kiểu tùy chỉnh.
    -   `frontend/src/main.tsx` là điểm nhập của ứng dụng React.

## Yêu cầu phát triển sản phẩm (PDRs)

### 1. Kiến trúc ứng dụng

-   **Yêu cầu chức năng**: Ứng dụng phải hoạt động như một ứng dụng máy tính để bàn sử dụng kiến trúc lai (Go backend, web frontend).
-   **Tiêu chí chấp nhận**: Ứng dụng khởi chạy thành công và hiển thị giao diện người dùng.
-   **Ràng buộc kỹ thuật**: Wails v3 phải được sử dụng làm khung công tác tích hợp.

### 2. Xử lý hình ảnh (Backend)

-   **Yêu cầu chức năng**: Backend Go phải có khả năng thực hiện các hoạt động xử lý hình ảnh cơ bản.
-   **Tiêu chí chấp nhận**: Các thư viện xử lý hình ảnh (ví dụ: `imaging`, `gg`) được tích hợp và sẵn sàng sử dụng.

### 3. Giao diện người dùng (Frontend)

-   **Yêu cầu chức năng**: Giao diện người dùng phải được xây dựng bằng công nghệ web hiện đại.
-   **Tiêu chí chấp nhận**: React, TypeScript, Vite và Tailwind CSS được thiết lập và cấu hình chính xác. Giao diện người dùng được hiển thị chính xác.
-   **Ràng buộc kỹ thuật**: Phải tuân thủ các phương pháp hay nhất của React và TypeScript.

### 4. Môi trường phát triển

-   **Yêu cầu chức năng**: Cung cấp một môi trường phát triển nhất quán và hiệu quả.
-   **Tiêu chí chấp nhận**: Các phần phụ thuộc của Go và Node.js được quản lý thông qua `go.mod` và `package.json` tương ứng.
-   **Ràng buộc kỹ thuật**: Các bản dựng phát triển và sản xuất phải hoạt động như mong đợi.

### Lịch sử thay đổi phiên bản

-   **Giai đoạn 1 hoàn thành**: Thiết lập dự án ban đầu, cấu hình backend Go với Wails, tích hợp frontend React/TypeScript/Tailwind. (2026-01-07)
