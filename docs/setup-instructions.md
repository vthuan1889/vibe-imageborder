# Hướng dẫn thiết lập dự án

Tài liệu này cung cấp hướng dẫn thiết lập dự án "vibe-imageborder" trên môi trường phát triển cục bộ của bạn. Dự án này bao gồm cả phần backend Go và frontend React/TypeScript.

## 1. Yêu cầu hệ thống

Trước khi bắt đầu, hãy đảm bảo rằng bạn đã cài đặt các công cụ sau:

-   **Go**: Phiên bản 1.25 trở lên.
    -   [Tải xuống Go](https://go.dev/dl/)
-   **Node.js**: Phiên bản 18 trở lên (bao gồm npm hoặc yarn).
    -   [Tải xuống Node.js](https://nodejs.org/en/download/)
-   **Wails CLI**: Công cụ dòng lệnh Wails.
    -   Cài đặt bằng cách chạy: `go install github.com/wailsapp/wails/v3/cmd/wails@latest`
-   **Trình duyệt hỗ trợ WebView2** (Windows): Đảm bảo bạn có Microsoft Edge WebView2 Runtime.
    -   [Tải xuống WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)
-   **Trình soạn thảo mã**: Ví dụ: VS Code với các tiện ích mở rộng Go và React/TypeScript.

## 2. Thiết lập dự án

### 2.1. Clone Repository

```bash
git clone <URL_TO_YOUR_REPOSITORY>
cd vibe-imageborder
```

### 2.2. Thiết lập Backend (Go)

1.  **Chuyển đến thư mục gốc của dự án**:
    ```bash
    cd <YOUR_PROJECT_ROOT_DIRECTORY>
    ```
2.  **Tải xuống các phần phụ thuộc Go**:
    ```bash
    go mod tidy
    ```
    Lệnh này sẽ tải xuống tất cả các phần phụ thuộc được liệt kê trong `go.mod`.

### 2.3. Thiết lập Frontend (React/TypeScript)

1.  **Chuyển đến thư mục frontend**:
    ```bash
    cd frontend
    ```
2.  **Cài đặt các phần phụ thuộc Node.js**:
    ```bash
    npm install
    # Hoặc nếu bạn dùng yarn
    # yarn install
    ```
    Lệnh này sẽ cài đặt tất cả các phần phụ thuộc được liệt kê trong `frontend/package.json`.

## 3. Chạy ứng dụng

Sau khi thiết lập cả backend và frontend, bạn có thể chạy ứng dụng Wails:

1.  **Chuyển đến thư mục gốc của dự án**:
    ```bash
    cd <YOUR_PROJECT_ROOT_DIRECTORY>
    ```
2.  **Chạy ứng dụng ở chế độ phát triển**:
    ```bash
    wails dev
    ```
    Lệnh này sẽ khởi động backend Go, xây dựng frontend bằng Vite, và mở cửa sổ ứng dụng desktop. Nó cũng cung cấp tính năng Hot Reloading cho cả backend và frontend trong quá trình phát triển.

## 4. Xây dựng ứng dụng

Để xây dựng một phiên bản sản phẩm của ứng dụng:

1.  **Chuyển đến thư mục gốc của dự án**:
    ```bash
    cd <YOUR_PROJECT_ROOT_DIRECTORY>
    ```
2.  **Xây dựng ứng dụng**:
    ```bash
    wails build
    ```
    Lệnh này sẽ tạo một tệp thực thi cho hệ điều hành hiện tại của bạn trong thư mục `build/bin`.

## 5. Cấu hình Tailwind CSS

Tailwind CSS được tích hợp sẵn trong dự án frontend. Bạn có thể tùy chỉnh cấu hình Tailwind bằng cách chỉnh sửa tệp `frontend/tailwind.config.js`.

-   Tất cả các lớp tiện ích Tailwind có sẵn để sử dụng trong các thành phần React của bạn.
-   Các kiểu cơ sở và tùy chỉnh được đưa vào thông qua `frontend/src/style.css`.

## 6. Khắc phục sự cố

-   **Lỗi Wails CLI**: Đảm bảo bạn đã cài đặt Wails CLI đúng cách và nó có trong PATH của bạn.
-   **Lỗi phần phụ thuộc Go**: Chạy `go mod tidy` để đảm bảo tất cả các phần phụ thuộc đã được giải quyết.
-   **Lỗi phần phụ thuộc Node.js**: Xóa thư mục `node_modules` và tệp `package-lock.json` (hoặc `yarn.lock`), sau đó chạy `npm install` (hoặc `yarn install`) lại.
-   **Không thấy giao diện người dùng**: Đảm bảo rằng ứng dụng Wails đang chạy và không có lỗi nào trong console của trình duyệt (bạn có thể mở công cụ dành cho nhà phát triển trong cửa sổ Wails).
