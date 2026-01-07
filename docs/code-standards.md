# Tiêu chuẩn mã hóa

Tài liệu này phác thảo các tiêu chuẩn mã hóa và hướng dẫn thực hành tốt nhất cho dự án "vibe-imageborder". Việc tuân thủ các tiêu chuẩn này đảm bảo khả năng đọc, khả năng bảo trì và chất lượng mã tổng thể.

## 1. Tiêu chuẩn mã Go (Backend)

### Định dạng
-   Tuân thủ `gofmt` để định dạng mã tự động.
-   Sử dụng `goimports` để tự động thêm/xóa các gói import.

### Cấu trúc
-   Các gói nên được đặt tên ngắn gọn, rõ ràng và tất cả bằng chữ thường.
-   Tổ chức mã thành các gói và mô-đun hợp lý. Ví dụ: `internal/models` cho các định nghĩa cấu trúc dữ liệu, `internal/service` cho logic nghiệp vụ.

### Đặt tên
-   **Biến và Hàm**: `camelCase` cho các tên nội bộ, `PascalCase` cho các tên được export (có thể truy cập từ các gói khác).
-   **Hằng số**: `PascalCase` cho các hằng số được export, `camelCase` cho các hằng số nội bộ. Có thể sử dụng `UPPER_SNAKE_CASE` cho các hằng số không được export nếu chúng là các từ viết tắt hoặc hằng số toàn cục.
-   **Structs**: `PascalCase`.
-   **Interfaces**: `PascalCase`, thường kết thúc bằng `er` (ví dụ: `Reader`, `Writer`).

### Xử lý lỗi
-   Trả về lỗi một cách rõ ràng dưới dạng giá trị cuối cùng.
-   Sử dụng `fmt.Errorf` để bọc lỗi và thêm ngữ cảnh.
-   Tránh sử dụng `panic` trừ khi lỗi là không thể phục hồi (ví dụ: lỗi khởi tạo chương trình).

### Nhận xét
-   Mọi hàm và struct được export đều phải có nhận xét tài liệu.
-   Các nhận xét nên giải thích "cái gì" và "tại sao", không chỉ "cách thức".

## 2. Tiêu chuẩn mã React/TypeScript (Frontend)

### Định dạng
-   Sử dụng Prettier và ESLint để tự động định dạng và kiểm tra linting.
-   Tuân thủ các cài đặt được xác định trong `.eslintrc.*` và `prettier.config.*`.

### Cấu trúc thư mục
-   Tổ chức các thành phần theo tính năng hoặc miền.
-   Ví dụ: `src/components`, `src/pages`, `src/utils`, `src/assets`.

### Đặt tên
-   **Thành phần**: `PascalCase` (ví dụ: `App`, `Button`).
-   **Props**: `camelCase`.
-   **Hàm**: `camelCase`.
-   **Biến**: `camelCase`.
-   **Tệp**: `PascalCase.tsx` cho các thành phần, `camelCase.ts` cho các tệp tiện ích/logic.

### TypeScript
-   Sử dụng TypeScript một cách nhất quán để tận dụng các lợi ích về kiểu dữ liệu.
-   Định nghĩa các kiểu giao diện (`interface`) hoặc kiểu (`type`) rõ ràng cho các props, state và cấu trúc dữ liệu.
-   Tránh sử dụng `any` trừ khi thực sự cần thiết và có lý do chính đáng.

### Tailwind CSS
-   Sử dụng các lớp tiện ích của Tailwind CSS để tạo kiểu.
-   Tránh viết CSS tùy chỉnh trong `src/style.css` nếu có thể sử dụng các lớp Tailwind.
-   Sử dụng `@apply` một cách tiết kiệm để tạo các thành phần tùy chỉnh phức tạp hoặc các mẫu thường xuyên sử dụng lại.

### Thành phần
-   Giữ các thành phần nhỏ, tập trung và có thể sử dụng lại.
-   Sử dụng các hook của React (ví dụ: `useState`, `useEffect`, `useContext`) một cách thích hợp.
-   Tránh sao chép mã; ưu tiên các thành phần và hàm tiện ích có thể sử dụng lại.

### Hiệu suất
-   Tối ưu hóa việc hiển thị thành phần bằng cách sử dụng `React.memo` hoặc `useCallback`/`useMemo` khi cần thiết để tránh hiển thị lại không cần thiết.

## 3. Quản lý phụ thuộc

-   **Go**: Sử dụng `go mod tidy` để đảm bảo `go.mod` và `go.sum` được cập nhật.
-   **Frontend**: Sử dụng `npm install` hoặc `yarn install` để quản lý các phần phụ thuộc frontend. Giữ `package.json` sạch sẽ và được cập nhật.

## 4. Nhận xét và Tài liệu

-   Viết nhận xét rõ ràng, ngắn gọn khi mã phức tạp hoặc khi có quyết định thiết kế quan trọng.
-   Giữ tài liệu dự án (thư mục `./docs`) được cập nhật và đồng bộ với mã.
-   Đảm bảo `README.md` cung cấp tổng quan đầy đủ và hướng dẫn thiết lập.
