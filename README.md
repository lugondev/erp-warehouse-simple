# ERP Warehouse Simple

A Go-based ERP Warehouse System built with Clean Architecture principles and Domain-Driven Design. It includes an API Gateway, comprehensive authentication, authorization, audit logging, and uses PostgreSQL with GORM and Prisma for database interactions.

## Technology Stack

- **Backend:** Go 1.21+ with Gin Web Framework
- **API Gateway:** Custom Go Gateway
- **Database:** PostgreSQL
- **ORM:** GORM
- **Database Tooling:** Prisma (for seeding)
- **Authentication:** JWT (Access & Refresh Tokens)
- **Authorization:** Role-based Access Control (RBAC)
- **Containerization:** Docker & Docker Compose
- **Package Manager:** Bun (for TypeScript scripts like seeding)

## Features

- User Authentication with JWT and refresh tokens
- Role-based Authorization with granular permissions
- Password reset functionality
- Comprehensive audit logging
- API authentication between modules
- Warehouse and Inventory Management
- Supplier Management
- Manufacturing Process Management
- Product/SKU Management with categorization
- Purchase Management with workflow (request → approval → order → receipt → payment)
- Customer Management with loyalty program and debt tracking
- Sales Order Management with delivery and invoicing
- Finance Management with invoices, payments, accounts receivable/payable, and financial reporting
- Reports and Analytics with inventory reports, sales reports, purchase reports, profit and loss reports, and dashboard metrics

## Project Structure

```
.
├── cmd
│   ├── gateway             # API Gateway entry point
│   └── server              # Main application entry point
├── internal
│   ├── domain
│   │   └── entity          # Core domain models
│   ├── application
│   │   └── usecase         # Application-specific business logic
│   └── infrastructure
│       ├── auth            # JWT handling, context
│       ├── config          # Configuration loading
│       ├── database        # Database connection (GORM), migrations
│       ├── gateway         # API Gateway logic (proxy, middleware)
│       ├── repository      # Data access layer implementations
│       ├── server          # Main application HTTP server (Gin), handlers
│       └── service         # Infrastructure-level services (e.g., Audit)
├── prisma                  # Prisma schema and seeding script
├── client                  # (Placeholder or future client-side code)
├── docs                    # API documentation (Swagger)
├── Dockerfile              # Dockerfile for the main application
├── Dockerfile.gateway      # Dockerfile for the API Gateway
├── docker-compose.yml      # Docker Compose configuration
├── go.mod / go.sum         # Go module dependencies
├── package.json / bun.lock # Node.js dependencies (for Prisma/TS scripts)
└── README.md               # This file
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21 or higher (for local development)
- Bun (for local development, includes Node.js)

### Running with Docker (Recommended)

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/yourusername/erp-warehouse-simple.git # Replace with actual repo URL if known
    cd erp-warehouse-simple
    ```
2.  **Copy environment variables:**
    ```bash
    cp .env.example .env
    # Review and adjust .env variables if necessary (secrets, ports, etc.)
    ```
3.  **Build and start services:**
    ```bash
    docker-compose up --build -d
    ```
    This command builds the images for the main application (`app`) and the API Gateway (`api-gateway`), starts them along with the PostgreSQL database (`postgres`), and runs them in detached mode.

4.  **Access the application:**
    The API Gateway is exposed on `http://localhost:8000`. All API requests should go through the gateway. The main application runs internally on port 8080.

5.  **(Optional) Seed the database:**
    If you need initial data (like the admin user), run the Prisma seed script using Docker:
    ```bash
    # First, install dependencies inside a temporary container if you haven't built the main 'app' service yet
    # docker-compose run --rm app bun install

    # Then run the seed script using the 'app' service definition
    docker-compose run --rm app bun run seed
    ```
    Alternatively, if you have Bun installed locally:
    ```bash
    bun install # Install dependencies for the script locally
    bun run seed  # Run seed script locally (ensure .env points to Docker DB)
    ```

### Running Locally (Advanced)

Running locally requires manual setup of Go, PostgreSQL, Bun/Node.js, and environment variables.

1.  **Install Prerequisites:**
    *   Go 1.21+
    *   PostgreSQL Server
    *   Bun (includes Node.js)
2.  **Setup Database:**
    *   Create a PostgreSQL database (e.g., `erp_db`).
    *   Configure connection details.
3.  **Environment Variables:**
    *   Copy `.env.example` to `.env`.
    *   Update `.env` with your local database credentials, JWT secrets, and desired ports. Ensure `ERP_DATABASE_HOST` points to your local Postgres instance (e.g., `localhost`).
4.  **Run Migrations (if applicable):**
    *   This project uses GORM and has SQL migration files in `internal/infrastructure/database/migrations/`. You'll need a tool like `golang-migrate` or execute them manually against your local database.
    ```bash
    # Example using golang-migrate (install it first: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
    # migrate -database "postgres://user:password@host:port/dbname?sslmode=disable" -path internal/infrastructure/database/migrations up
    ```
5.  **Install Node Dependencies:**
    ```bash
    bun install
    ```
6.  **Run Database Seeding:**
    ```bash
    bun run seed
    ```
7.  **Run the Main Application:**
    ```bash
    # Ensure environment variables from .env are loaded (e.g., using direnv or source .env)
    go run cmd/server/main.go
    ```
8.  **Run the API Gateway (Optional):**
    ```bash
    # Ensure environment variables from .env are loaded
    go run cmd/gateway/main.go
    ```
    *Note: Adjust `ERP_APIGATEWAY_SERVICES_WAREHOUSE_URL` in `.env` to point to your locally running main application (e.g., `http://localhost:8080` if using default port).*

The API Gateway will be available at `http://localhost:8000` (or the port configured in `.env`).

## Default Admin Account

The system creates a default admin account when seeded:
- Username: admin
- Email: admin@example.com
- Password: admin123

Please change these credentials after first login.

## API Endpoints

All endpoints are accessed via the API Gateway (e.g., `http://localhost:8000/api/v1/...`).

### Authentication Endpoints

- `POST /api/v1/auth/register` - Register a new user
  ```json
  {
    "username": "string",
    "email": "string",
    "password": "string",
    "role_id": "number"
  }
  ```

- `POST /api/v1/auth/login` - Login
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```

- `POST /api/v1/auth/refresh-token` - Refresh access token
  ```json
  {
    "refresh_token": "string"
  }
  ```

- `POST /api/v1/auth/forgot-password` - Request password reset
  ```json
  {
    "email": "string"
  }
  ```

- `POST /api/v1/auth/reset-password` - Reset password
  ```json
  {
    "token": "string",
    "new_password": "string"
  }
  ```

### Protected Endpoints

All protected endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <your-token>
```

#### User Management

- `GET /api/v1/users` - List all users (Admin only)
- `GET /api/v1/users/:id` - Get user details
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user (Admin only)
- `POST /api/v1/users/logout` - Logout user

#### Role Management

- `POST /api/v1/roles` - Create role (Admin only)
- `GET /api/v1/roles` - List all roles
- `GET /api/v1/roles/:id` - Get role details
- `PUT /api/v1/roles/:id` - Update role (Admin only)
- `DELETE /api/v1/roles/:id` - Delete role (Admin only)

#### Audit Logs

- `GET /api/v1/audit/logs` - List all audit logs (Admin only)
- `GET /api/v1/audit/logs/user/:id` - Get user audit logs (Admin only)

#### Product/SKU Management

- `POST /api/v1/items` - Create a new item
- `GET /api/v1/items` - List items with filters
- `GET /api/v1/items/search` - Search items by term
- `GET /api/v1/items/:id` - Get item details
- `GET /api/v1/items/sku/:sku` - Get item by SKU
- `PUT /api/v1/items/:id` - Update item
- `DELETE /api/v1/items/:id` - Delete item
- `POST /api/v1/items/bulk` - Bulk create items
- `PUT /api/v1/items/bulk` - Bulk update items

#### Item Category Management

- `POST /api/v1/item-categories` - Create a new category
- `GET /api/v1/item-categories` - List all categories
- `GET /api/v1/item-categories/tree` - Get categories in tree structure
- `GET /api/v1/item-categories/:id` - Get category details
- `PUT /api/v1/item-categories/:id` - Update category
- `DELETE /api/v1/item-categories/:id` - Delete category
- `GET /api/v1/item-categories/:id/items` - Get items in a category

#### Customer Management

- `POST /api/v1/customers` - Create a new customer
- `GET /api/v1/customers` - List customers with filters
- `GET /api/v1/customers/:id` - Get customer details
- `PUT /api/v1/customers/:id` - Update customer
- `DELETE /api/v1/customers/:id` - Delete customer

- `POST /api/v1/customers/:id/addresses` - Add address to customer
- `GET /api/v1/customers/:id/addresses` - Get customer addresses
- `PUT /api/v1/customers/addresses/:addressId` - Update address
- `DELETE /api/v1/customers/addresses/:addressId` - Delete address

- `GET /api/v1/customers/:id/orders` - Get customer order history
- `GET /api/v1/customers/:id/debt` - Get customer debt information
- `PUT /api/v1/customers/:id/debt` - Update customer debt

- `PUT /api/v1/customers/:id/loyalty/points` - Update loyalty points
- `PUT /api/v1/customers/:id/loyalty/tier` - Update loyalty tier
- `GET /api/v1/customers/:id/loyalty/calculate-tier` - Calculate loyalty tier

#### Finance Management

- `POST /api/v1/finance/invoices` - Create a new invoice
- `GET /api/v1/finance/invoices` - List invoices with filters
- `GET /api/v1/finance/invoices/:id` - Get invoice details
- `PUT /api/v1/finance/invoices/:id` - Update invoice
- `PATCH /api/v1/finance/invoices/:id/status` - Update invoice status
- `POST /api/v1/finance/invoices/:id/cancel` - Cancel invoice

- `POST /api/v1/finance/payments` - Create a new payment
- `GET /api/v1/finance/payments` - List payments with filters
- `GET /api/v1/finance/payments/:id` - Get payment details
- `PUT /api/v1/finance/payments/:id` - Update payment
- `POST /api/v1/finance/payments/:id/confirm` - Confirm payment
- `POST /api/v1/finance/payments/:id/cancel` - Cancel payment
- `POST /api/v1/finance/payments/:id/refund` - Refund payment

- `GET /api/v1/finance/reports/accounts-receivable` - Get accounts receivable report
- `GET /api/v1/finance/reports/accounts-payable` - Get accounts payable report
- `GET /api/v1/finance/reports/finance` - Get financial report

#### Reports and Analytics

- `POST /api/v1/reports` - Create a new report
- `GET /api/v1/reports` - List reports with filters
- `GET /api/v1/reports/:id` - Get report details
- `DELETE /api/v1/reports/:id` - Delete report
- `POST /api/v1/reports/:id/export` - Export report to CSV, Excel, or PDF

- `POST /api/v1/reports/schedules` - Create a new report schedule
- `GET /api/v1/reports/schedules` - List report schedules
- `GET /api/v1/reports/schedules/:id` - Get report schedule details
- `PUT /api/v1/reports/schedules/:id` - Update report schedule
- `DELETE /api/v1/reports/schedules/:id` - Delete report schedule

- `GET /api/v1/reports/inventory/value` - Get inventory value report
- `GET /api/v1/reports/inventory/age` - Get inventory age report
- `GET /api/v1/reports/sales/products` - Get product sales report
- `GET /api/v1/reports/sales/customers` - Get customer sales report
- `GET /api/v1/reports/purchases/suppliers` - Get supplier purchase report
- `GET /api/v1/reports/financial/profit-loss` - Get profit and loss report
- `GET /api/v1/reports/dashboard/metrics` - Get dashboard metrics

## Available Permissions

- User Management: `user:create`, `user:read`, `user:update`, `user:delete`
- Role Management: `role:create`, `role:read`, `role:update`, `role:delete`
- Audit Logs: `audit:read`
- Module Integration: `module:integrate`
- Product Management: `product:create`, `product:read`, `product:update`, `product:delete`
- Customer Management: `customer:create`, `customer:read`, `customer:update`, `customer:delete`
- Customer Address: `customer:address:create`, `customer:address:read`, `customer:address:update`, `customer:address:delete`
- Customer Debt: `customer:debt:read`, `customer:debt:update`
- Customer Loyalty: `customer:loyalty:read`, `customer:loyalty:update`
- Finance Management: `finance:invoice:create`, `finance:invoice:read`, `finance:invoice:update`, `finance:invoice:delete`
- Payment Management: `finance:payment:create`, `finance:payment:read`, `finance:payment:update`, `finance:payment:process`
- Financial Reporting: `finance:report:read`
- Report Management: `report:create`, `report:read`, `report:update`, `report:delete`, `report:export`
- Report Schedule Management: `report:schedule:create`, `report:schedule:read`, `report:schedule:update`, `report:schedule:delete`

## Development

### Adding New Permissions

1. Add new permission constant in `internal/domain/entity/role.go`
2. Update role validation in `RoleUseCase`
3. Implement required handlers and middleware

### Audit Logging

The system automatically logs:
- User authentication events (login, logout)
- Resource modifications (create, update, delete)
- Access to protected resources

Audit logs include:
- User ID and action
- Resource type and ID
- Timestamp
- IP address
- User agent

## License

[MIT License](LICENSE)
