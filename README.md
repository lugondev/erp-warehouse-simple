# ERP Warehouse System

A Go-based ERP Warehouse System built with Clean Architecture principles and Domain-Driven Design, featuring comprehensive authentication, authorization, audit logging, and an API Gateway for microservices integration.

## Technology Stack

- Go 1.21+
- Gin Web Framework
- GORM with PostgreSQL
- JWT Authentication with refresh tokens
- Role-based Access Control (RBAC)
- Audit Logging
- API Gateway with rate limiting and circuit breaking
- WebSocket support for real-time updates
- Docker & Docker Compose

## Features

- User Authentication with JWT and refresh tokens
- Role-based Authorization with granular permissions
- Password reset functionality
- Comprehensive audit logging
- API Gateway for microservices integration
  - Single entry point for all API requests
  - Routing and load balancing
  - Centralized authentication and authorization
  - Rate limiting and circuit breaking
  - Logging and monitoring
  - WebSocket support for real-time updates
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
│   ├── server              # Main application entry point
│   └── gateway             # API Gateway entry point
├── internal
│   ├── domain
│   │   └── entity         # Domain entities (User, Role, AuditLog)
│   ├── application
│   │   └── usecase        # Application business rules
│   └── infrastructure
│       ├── auth           # Authentication services
│       ├── config         # Configuration management
│       ├── database       # Database connection
│       ├── gateway        # API Gateway implementation
│       │   ├── middleware # Gateway middleware (rate limiting, circuit breaking)
│       │   ├── proxy      # Service proxy for routing requests
│       │   └── websocket  # WebSocket support for real-time updates
│       ├── repository     # Data persistence
│       ├── service        # Application services
│       └── server         # HTTP server and handlers
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21 or higher

### Running the Application

1. Clone the repository:
```bash
git clone https://github.com/yourusername/erp-warehouse-simple.git
cd erp-warehouse-simple
```

2. Start the application using Docker Compose:
```bash
docker-compose up --build
```

The main application will be available at `http://localhost:8080`
The API Gateway will be available at `http://localhost:8000`

### Running the API Gateway Separately

You can also run the API Gateway separately:

```bash
go run cmd/gateway/main.go
```

The API Gateway will be available at `http://localhost:8000`

## Default Admin Account

The system creates a default admin account on first run:
- Username: admin
- Email: admin@example.com
- Password: admin123

Please change these credentials after first login.

## API Gateway

The API Gateway serves as a single entry point for all API requests and provides the following features:

### Features

- **Centralized Routing**: All API requests go through a single entry point
- **Load Balancing**: Distributes requests across multiple service instances
- **Authentication & Authorization**: Centralized security for all services
- **Rate Limiting**: Prevents abuse by limiting request rates
- **Circuit Breaking**: Prevents cascading failures
- **Logging & Monitoring**: Centralized logging and monitoring
- **WebSocket Support**: Real-time updates and notifications

### Configuration

The API Gateway is configured in the application configuration:

```yaml
apigateway:
  enabled: true
  port: 8000
  tracing: true
  logging: true
  ratelimit:
    requests_per_second: 100
    burst: 50
  circuitbreak:
    max_requests: 100
    interval: 60
    timeout: 30
    consecutive_error: 5
  services:
    warehouse:
      url: http://localhost:8081
      timeout: 30
      retry_count: 3
      health_check: /health
    inventory:
      url: http://localhost:8082
      timeout: 30
      retry_count: 3
      health_check: /health
    # Other services...
```

### WebSocket Support

The API Gateway provides WebSocket support for real-time updates:

```javascript
// Connect to WebSocket
const socket = new WebSocket('ws://localhost:8000/ws');

// Handle messages
socket.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('Received:', data);
};

// Send a message
socket.send(JSON.stringify({
  type: 'subscribe',
  channel: 'inventory_updates'
}));
```

## API Endpoints

All API endpoints are accessible through the API Gateway at `http://localhost:8000`.

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
