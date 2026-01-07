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

### Giai đoạn 2: Dịch vụ Mẫu (Template Service)

Giai đoạn 2 tập trung vào việc triển khai một dịch vụ mẫu mạnh mẽ để xử lý các mẫu (templates), cho phép trích xuất trường và thay thế biến động.

-   **Backend (Go)**:
    -   Tệp `internal/models/types.go` được cập nhật với các định nghĩa kiểu liên quan đến mẫu.
    -   `internal/template/parser.go`: Triển khai logic phân tích mẫu để trích xuất các trường và biến.
    -   `internal/template/service.go`: Cung cấp các chức năng để quản lý, xử lý và áp dụng các mẫu.
    -   `internal/template/parser_test.go`: Các bài kiểm tra đơn vị cho trình phân tích mẫu, đảm bảo độ chính xác và độ tin cậy.

-   **Kiểm tra**:
    -   Độ bao phủ kiểm thử đạt 95.7%, đảm bảo chất lượng cao của mã.
    -   Tất cả các bài kiểm tra đều vượt qua, xác nhận tính đúng đắn của việc triển khai dịch vụ mẫu.

-   **Tính năng chính**:
    -   Phân tích mẫu: Khả năng phân tích các chuỗi mẫu để xác định các phần tử động.
    -   Trích xuất trường: Trích xuất chính xác các trường dữ liệu từ các mẫu đã cho.
    -   Thay thế biến: Thay thế các biến trong mẫu bằng các giá trị được cung cấp một cách linh hoạt.

### Lịch sử thay đổi phiên bản

-   **Giai đoạn 1 hoàn thành**: Thiết lập dự án ban đầu, cấu hình backend Go với Wails, tích hợp frontend React/TypeScript/Tailwind. (2026-01-07)
-   **Giai đoạn 2 hoàn thành**: Triển khai Template Service với tính năng phân tích mẫu, trích xuất trường và thay thế biến. Độ bao phủ kiểm thử 95.7%. (2026-01-07)

