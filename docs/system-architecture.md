# Kiến trúc hệ thống

Dự án "vibe-imageborder" được xây dựng trên kiến trúc lai (hybrid architecture), tận dụng sức mạnh của Go cho logic backend và tính linh hoạt của các công nghệ web hiện đại cho giao diện người dùng. Nền tảng của kiến trúc này là khung công tác Wails, cung cấp một cầu nối liền mạch giữa hai môi trường này.

## 1. Tổng quan kiến trúc

Kiến trúc có thể được hình dung thành hai lớp chính:

-   **Backend Layer (Go)**: Chịu trách nhiệm về logic nghiệp vụ cốt lõi, xử lý hình ảnh, và quản lý tài nguyên hệ thống.
-   **Frontend Layer (Web Technologies)**: Cung cấp giao diện người dùng (UI) tương tác và phản hồi nhanh, nơi người dùng sẽ tương tác với ứng dụng.

<pre>
+---------------------+      +------------------------+
|   Backend (Go)      |      |   Frontend (Web)       |
|                     |      |                        |
| - Wails Runtime     |<---->|- Wails JavaScript API   |
| - Business Logic    |      | - React Application    |
| - Image Processing  |      | - Vite Build Tool      |
|   (imaging, gg)     |      | - Tailwind CSS         |
| - System Interaction|      | - User Interface       |
+---------------------+      +------------------------+
        ^
        | (Wails Bundler)
        v
+---------------------+
|   Desktop App       |
| (Wails Runtime)     |
+---------------------+
</pre>

## 2. Backend Layer (Go)

-   **Ngôn ngữ**: Go
-   **Khung công tác chính**: Wails v3
-   **Mô tả**:
    -   Backend là trung tâm của các hoạt động xử lý hình ảnh. Nó sử dụng các thư viện Go chuyên dụng như `github.com/disintegration/imaging` và `github.com/fogleman/gg` để thực hiện các thao tác như thêm đường viền, thay đổi kích thước, hoặc áp dụng bộ lọc cho hình ảnh.
    -   Wails runtime cho phép backend Go phơi bày các hàm Go trực tiếp cho frontend thông qua một API JavaScript được tạo tự động. Điều này giúp giao diện người dùng gọi các hàm backend một cách dễ dàng và an toàn.
    -   Các tệp trong `internal/models` sẽ chứa các định nghĩa kiểu dữ liệu được chia sẻ giữa các thành phần backend và có khả năng là frontend (thông qua việc chuyển đổi dữ liệu Wails).

## 3. Frontend Layer (Web Technologies)

-   **Ngôn ngữ**: TypeScript, JavaScript, CSS
-   **Khung công tác chính**: React
-   **Công cụ xây dựng**: Vite
-   **Thư viện UI/Styling**: Tailwind CSS
-   **Mô tả**:
    -   Frontend là một ứng dụng web được xây dựng bằng React và TypeScript, cung cấp giao diện đồ họa cho người dùng để tương tác với các tính năng của ứng dụng.
    -   Vite được sử dụng để xây dựng và phục vụ ứng dụng React, cung cấp trải nghiệm phát triển nhanh chóng với tính năng Hot Module Replacement (HMR).
    -   Tailwind CSS cung cấp một bộ khung tiện ích để tạo kiểu nhanh chóng và tùy chỉnh giao diện người dùng, đảm bảo một thiết kế sạch sẽ và nhất quán.
    -   Giao diện người dùng giao tiếp với backend Go thông qua Wails JavaScript API. Các sự kiện từ UI sẽ kích hoạt các hàm Go backend để xử lý dữ liệu hoặc thực hiện các tác vụ nặng.

## 4. Tương tác và Luồng dữ liệu

1.  **Khởi tạo ứng dụng**: Khi ứng dụng desktop Wails khởi chạy, nó tải backend Go và hiển thị giao diện người dùng web (Frontend).
2.  **Tương tác của người dùng**: Người dùng tương tác với giao diện người dùng React (ví dụ: chọn một hình ảnh, nhập tham số đường viền).
3.  **Gọi Backend**: Các hành động của người dùng sẽ kích hoạt các lệnh gọi đến các hàm Go được phơi bày bởi Wails JavaScript API.
4.  **Xử lý Backend**: Backend Go nhận các lệnh gọi này, thực hiện logic nghiệp vụ hoặc thao tác xử lý hình ảnh bằng các thư viện Go.
5.  **Phản hồi về Frontend**: Kết quả của hoạt động backend (ví dụ: hình ảnh đã xử lý, dữ liệu) được trả về frontend.
6.  **Cập nhật UI**: Frontend React cập nhật giao diện người dùng để hiển thị kết quả cho người dùng.

## 5. Các thành phần chính và mục đích của chúng

-   **`main.go`**: Điểm nhập cho ứng dụng Wails Go, khởi tạo backend và cấu hình ứng dụng desktop.
-   **`go.mod`**: Quản lý các phần phụ thuộc Go, bao gồm Wails và các thư viện xử lý hình ảnh.
-   **`internal/models/types.go`**: (Hiện tại là trình giữ chỗ) Sẽ chứa các định nghĩa kiểu dữ liệu dùng chung.
-   **`frontend/index.html`**: Tệp HTML cơ sở cho ứng dụng web frontend.
-   **`frontend/src/main.tsx`**: Điểm nhập cho ứng dụng React frontend.
-   **`frontend/src/App.tsx`**: Thành phần React cấp cao nhất.
-   **`frontend/src/style.css`**: Các kiểu CSS cơ sở, bao gồm các chỉ thị Tailwind.
-   **`frontend/package.json`**: Quản lý các phần phụ thuộc Node.js/frontend.
-   **`frontend/tailwind.config.js`**: Cấu hình tùy chỉnh cho Tailwind CSS.

Kiến trúc này cho phép phát triển hiệu quả, tận dụng sức mạnh của Go cho các tác vụ hiệu suất cao và sự nhanh nhẹn của React cho giao diện người dùng tương tác.
