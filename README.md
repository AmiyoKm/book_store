# BookBond: A Full-Stack E-commerce Platform for Books

![BookBond Logo Placeholder](https://via.placeholder.com/150/007bff/ffffff?text=BookBond)

BookBond is a comprehensive full-stack e-commerce application designed for buying and selling books. It features a modern React frontend and a robust Go backend, providing a seamless experience for users to browse, purchase, and manage their book collections.

## 🌐 Live Demo

*   **Frontend:** [https://bookbond.netlify.app](https://bookbond.netlify.app)
*   **Backend API:** [https://book-bond.onrender.com/api/v1](https://book-bond.onrender.com/api/v1)

## ✨ Features

### Frontend (Client)
*   **User Authentication:** Secure sign-up, sign-in, password reset, and account activation flows.
*   **Book Catalog:** Browse a wide range of books with detailed views.
*   **Shopping Cart:** Add, update, and remove items from your cart.
*   **Wishlist:** Save books for later.
*   **Order Management:** View past orders and their details.
*   **User Profile:** Manage personal information.
*   **Responsive Design:** Optimized for various screen sizes.
*   **Theme Toggling:** Light and Dark mode support.

### Backend (API)
*   **RESTful API:** Clean and well-structured API endpoints for all functionalities.
*   **Authentication & Authorization:** JWT-based authentication for secure access.
*   **Database Management:** PostgreSQL integration for data persistence.
*   **Email Notifications:** For account activation and password resets.
*   **Scalable Architecture:** Built with Go for high performance and concurrency.

## 🚀 Technologies Used

### Frontend
*   **React:** A JavaScript library for building user interfaces.
*   **Vite:** A fast build tool for modern web projects.
*   **TypeScript:** A typed superset of JavaScript that compiles to plain JavaScript.
*   **Tailwind CSS:** A utility-first CSS framework for rapid UI development.
*   **Shadcn/ui:** Reusable components built with Radix UI and Tailwind CSS.
*   **React Query:** For efficient data fetching, caching, and synchronization.
*   **React Router:** For declarative routing in React applications.
*   **Axios:** Promise-based HTTP client for the browser and Node.js.
*   **Sonner:** A toast library for notifications.

### Backend
*   **Go:** A statically typed, compiled programming language designed for building simple, reliable, and efficient software.
*   **Chi:** A lightweight, idiomatic, and composable router for building HTTP services in Go.
*   **PostgreSQL:** A powerful, open-source object-relational database system.
*   **JWT (JSON Web Tokens):** For secure authentication.
*   **go-playground/validator:** For request payload validation.
*   **go.uber.org/zap:** A fast, structured, leveled logging in Go.
*   **joho/godotenv:** For loading environment variables from `.env` files.
*   **gomail/mail:** For sending emails.

### Other
*   **Docker & Docker Compose:** For containerization and easy environment setup.
*   **Git:** Version control system.
*   **GitHub Actions:** For Continuous Integration/Continuous Deployment (CI/CD).

## 🛠️ Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Make sure you have the following installed:
*   [Go](https://golang.org/doc/install) (version 1.24.2 or higher)
*   [Node.js](https://nodejs.org/en/download/) (LTS version recommended)
*   [npm](https://www.npmjs.com/get-npm) (comes with Node.js)
*   [Docker](https://www.docker.com/get-started) & [Docker Compose](https://docs.docker.com/compose/install/)
*   [PostgreSQL](https://www.postgresql.org/download/) (or use Docker for the database)

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/AmiyoKm/book_store.git
    cd book_store
    ```

2.  **Set up Environment Variables:**
    Create a `.env` file in the root directory of the project and populate it with necessary environment variables.
    Example `.env` (adjust values as needed):
    ```
    DB_DSN="host=localhost port=5432 user=bookstore password=password dbname=bookstore sslmode=disable"
    JWT_SECRET="your_jwt_secret_key"
    SMTP_HOST="smtp.mailtrap.io"
    SMTP_PORT=2525
    SMTP_USERNAME="your_mailtrap_username"
    SMTP_PASSWORD="your_mailtrap_password"
    SMTP_SENDER="BookBond <no-reply@bookbond.com>"
    ```
    For the frontend, create a `.env` file in the `client/` directory:
    ```
    VITE_BACKEND_PROD_ENDPOINT="http://localhost:8080/api/v1"
    ```

3.  **Run with Docker Compose (Recommended for full setup):**
    This will set up the PostgreSQL database, run migrations, and start both the Go backend and React frontend.
    ```bash
    docker-compose up --build
    ```
    The frontend will be accessible at `http://localhost:5173` (or similar, check Docker logs), and the backend API at `http://localhost:8080/api/v1`.

4.  **Manual Setup (Alternative)**

    #### Backend Setup
    ```bash
    # Navigate to the backend directory
    cd cmd/api

    # Install Go dependencies
    go mod tidy

    # Run database migrations (ensure PostgreSQL is running and DB_DSN is correct)
    # You might need to install 'migrate' tool: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    migrate -path ../migrate/migrations -database "${DB_DSN}" up

    # Run the backend API
    go run .
    ```
    The backend API will run on `http://localhost:8080`.

    #### Frontend Setup
    ```bash
    # Navigate to the client directory
    cd client

    # Install Node.js dependencies
    npm install

    # Start the development server
    npm run dev
    ```
    The frontend development server will typically run on `http://localhost:5173`.

## 📂 Project Structure

```
.
├── .air.toml             # Air (Go live-reloading) configuration
├── .gitignore            # Git ignore file
├── docker-compose.yml    # Docker Compose configuration
├── go.mod                # Go module dependencies
├── go.sum                # Go module checksums
├── makefile              # Makefile for common commands
├── .git/                 # Git repository metadata
├── .github/              # GitHub Actions workflows
│   └── workflows/
│       └── go.yml        # CI/CD for Go backend
├── bin/                  # Compiled binaries (if any)
├── client/               # React Frontend Application
│   ├── public/           # Static assets
│   ├── src/              # Frontend source code
│   │   ├── components/   # Reusable UI components
│   │   ├── config/       # Axios, API service configurations
│   │   ├── lib/          # Utility functions
│   │   ├── pages/        # Page-specific components (e.g., books, cart, auth)
│   │   └── types/        # TypeScript type definitions
│   ├── package.json      # Frontend dependencies and scripts
│   └── vite.config.ts    # Vite build configuration
├── cmd/                  # Main application entry points
│   ├── api/              # Go Backend API application
│   └── migrate/          # Database migration files
├── docs/                 # API documentation (Swagger/OpenAPI)
├── internal/             # Internal Go packages/modules
│   ├── auth/             # Authentication logic
│   ├── db/               # Database connection and utilities
│   ├── env/              # Environment variable loading
│   ├── mail/             # Email sending utilities and templates
│   └── store/            # Business logic for different entities (books, users, carts, etc.)
└── script/               # Utility scripts (e.g., for development, deployment)
```

## 🤝 Contributing

Contributions are welcome! If you have suggestions for improvements or new features, please open an issue or submit a pull request.

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the `LICENSE` file for details. (Note: You might need to create a `LICENSE` file if one doesn't exist).

## 📞 Contact

For any questions or support, please open an issue in the GitHub repository.