# ERP Warehouse Management System

This is a simple ERP Warehouse Management System that provides functionality for managing warehouses, inventory, customers, suppliers, and more.

## Database Seeding

The project includes a comprehensive seed script to generate sample data for all tables in the database. This is useful for development and testing purposes.

### Prerequisites

- Node.js and npm/bun installed
- PostgreSQL database
- Prisma CLI

### Setup

1. Clone the repository
2. Install dependencies:
   ```bash
   npm install
   # or
   bun install
   ```
3. Set up your database connection in `.env`:
   ```
   DATABASE_URL="postgresql://username:password@localhost:5432/erp_warehouse"
   ```
4. Run Prisma migrations:
   ```bash
   npx prisma migrate dev
   # or
   bunx prisma migrate dev
   ```

### Running the Seed Script

To populate your database with sample data:

```bash
npm run seed
# or
bun run seed
```

The seed script will generate data for:

- Users and roles (including admins and regular users)
- Warehouses and inventory management
- Items, categories, and inventory tracking
- Customers and their addresses
- Sequences for document numbering
- Audit logs for user actions

### Known Issues

- The seed script uses TypeScript and may show some type warnings due to the dynamic nature of some of the data generation.
- When running the seed script, you may need to install additional dependencies:
  ```bash
  npm install @faker-js/faker bcrypt
  # or
  bun add @faker-js/faker bcrypt
  ```

## Development

To start development:

1. Make sure your database is running
2. Run the development server:
   ```bash
   npm run dev
   # or
   bun run dev
   ```

## Project Structure

- `prisma/`: Contains Prisma schema and migrations
  - `schema.prisma`: Database schema definition
  - `seed.ts`: Database seeding script
- `src/`: Source code for the application

## License

This project is licensed under the MIT License.
