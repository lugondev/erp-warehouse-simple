import { PrismaClient } from '@prisma/client';
import { faker } from '@faker-js/faker';
import * as bcrypt from 'bcryptjs';

/**
 * Helper function to generate random decimal numbers with specified precision
 */
function randomDecimal(min: number, max: number, precision: number = 2): number {
	const value = Math.random() * (max - min) + min;
	const multiplier = Math.pow(10, precision);
	return Math.round(value * multiplier) / multiplier;
}

const prisma = new PrismaClient();

/**
 * Main seed function that orchestrates the data generation process
 */
async function main() {
	console.log('Starting database seeding...');

	// Clear existing data (in reverse order of dependencies)
	await clearExistingData();

	// Generate data in order of dependencies
	const roles = await createRoles();
	const users = await createUsers(roles);
	const warehouses = await createWarehouses(users);
	await createItemCategories();
	const customers = await createCustomers();
	await createCustomerAddresses(customers);
	await createSequences();

	// Create audit logs for various actions
	await createAuditLogs(users);

	// Create inventory-related data
	const inventories = await createInventories(warehouses);
	await createInventoryHistories(inventories, users);
	await createStockEntries(warehouses, users);

	console.log('Database seeding completed successfully!');
}

/**
 * Clear all existing data from the database in reverse order of dependencies
 */
async function clearExistingData() {
	console.log('Clearing existing data...');

	await prisma.audit_logs.deleteMany({});
	await prisma.inventory_histories.deleteMany({});
	await prisma.stock_entries.deleteMany({});
	await prisma.inventories.deleteMany({});
	await prisma.customer_addresses.deleteMany({});
	await prisma.customers.deleteMany({});
	await prisma.sequences.deleteMany({});
	await prisma.warehouses.deleteMany({});
	await prisma.users.deleteMany({});
	await prisma.roles.deleteMany({});
	await prisma.item_categories.deleteMany({});

	console.log('Existing data cleared.');
}

/**
 * Create roles with appropriate permissions
 */
async function createRoles() {
	console.log('Creating roles...');

	const roles = [
		{
			name: 'Admin',
			permissions: [
				'MANAGE_USERS', 'MANAGE_ROLES', 'MANAGE_WAREHOUSES',
				'MANAGE_INVENTORY', 'MANAGE_ITEMS', 'MANAGE_SUPPLIERS',
				'MANAGE_CUSTOMERS', 'MANAGE_ORDERS', 'MANAGE_FINANCES',
				'MANAGE_REPORTS', 'VIEW_ALL'
			]
		},
		{
			name: 'Warehouse Manager',
			permissions: [
				'MANAGE_WAREHOUSES', 'MANAGE_INVENTORY', 'VIEW_ITEMS',
				'VIEW_SUPPLIERS', 'VIEW_ORDERS'
			]
		},
		{
			name: 'Sales Representative',
			permissions: [
				'MANAGE_CUSTOMERS', 'MANAGE_ORDERS', 'VIEW_INVENTORY',
				'VIEW_ITEMS'
			]
		},
		{
			name: 'Inventory Clerk',
			permissions: [
				'MANAGE_INVENTORY', 'VIEW_WAREHOUSES', 'VIEW_ITEMS',
				'VIEW_SUPPLIERS'
			]
		},
		{
			name: 'Finance Officer',
			permissions: [
				'MANAGE_FINANCES', 'VIEW_ORDERS', 'VIEW_CUSTOMERS',
				'VIEW_SUPPLIERS'
			]
		}
	];

	const createdRoles: any[] = [];

	for (const role of roles) {
		const createdRole = await prisma.roles.create({
			data: {
				name: role.name,
				permissions: role.permissions,
				created_at: new Date(),
				updated_at: new Date()
			}
		});

		createdRoles.push(createdRole);
		console.log(`Created role: ${role.name}`);
	}

	return createdRoles;
}

/**
 * Create users with different roles
 */
async function createUsers(roles: any[]) {
	console.log('Creating users...');

	const users: any[] = [];
	const saltRounds = 10;

	// Create one admin user
	const adminPassword = await bcrypt.hash('admin123', saltRounds);
	const adminUser = await prisma.users.create({
		data: {
			username: 'admin',
			email: 'admin@example.com',
			password: adminPassword,
			role_id: roles.find(r => r.name === 'Admin')!.id,
			status: 'active',
			last_login: faker.date.recent(),
			created_at: new Date(),
			updated_at: new Date()
		}
	});

	users.push(adminUser);
	console.log('Created admin user');

	// Create 15 regular users with various roles
	for (let i = 0; i < 15; i++) {
		const firstName = faker.person.firstName();
		const lastName = faker.person.lastName();
		const username = faker.internet.userName({ firstName, lastName }).toLowerCase();
		const email = faker.internet.email({ firstName, lastName }).toLowerCase();
		const password = await bcrypt.hash('password123', saltRounds);

		// Assign a random role (excluding admin for most users)
		const roleIndex = i < 2 ? 0 : Math.floor(Math.random() * (roles.length - 1)) + 1;
		const role = roles[roleIndex];

		const user = await prisma.users.create({
			data: {
				username,
				email,
				password,
				role_id: role.id,
				status: faker.helpers.arrayElement(['active', 'inactive', 'suspended']),
				last_login: faker.date.recent(),
				created_at: faker.date.past(),
				updated_at: faker.date.recent(),
			}
		});

		users.push(user);
		console.log(`Created user: ${username} with role ${role.name}`);
	}

	return users;
}

/**
 * Create warehouses with different types and statuses
 */
async function createWarehouses(users: any[]) {
	console.log('Creating warehouses...');

	const warehouseTypes = ['RAW', 'FINISHED', 'GENERAL'];
	const warehouseStatuses = ['ACTIVE', 'INACTIVE'];
	const warehouses: any[] = [];

	// Get users with warehouse management permissions
	const warehouseManagers = users.filter(user =>
		user.roles?.permissions?.includes('MANAGE_WAREHOUSES') ||
		user.roles?.name === 'Admin' ||
		user.roles?.name === 'Warehouse Manager'
	);

	// If no specific warehouse managers, use any user
	const managerPool = warehouseManagers.length > 0 ? warehouseManagers : users;

	for (let i = 0; i < 5; i++) {
		const warehouseId = faker.string.uuid();
		const warehouseType = faker.helpers.arrayElement(warehouseTypes);
		const manager = faker.helpers.arrayElement(managerPool);

		const warehouse = await prisma.warehouses.create({
			data: {
				id: warehouseId,
				name: `${faker.location.city()} ${warehouseType} Warehouse`,
				address: faker.location.streetAddress({ useFullAddress: true }),
				type: warehouseType,
				manager_id: manager.id.toString(),
				contact: faker.phone.number(),
				status: faker.helpers.arrayElement(warehouseStatuses),
				created_at: faker.date.past(),
				updated_at: faker.date.recent()
			}
		});

		warehouses.push(warehouse);
		console.log(`Created warehouse: ${warehouse.name}`);
	}

	return warehouses;
}

/**
 * Create item categories with hierarchical structure
 */
async function createItemCategories() {
	console.log('Creating item categories...');

	const categories = [
		{ name: 'Electronics', description: 'Electronic devices and components' },
		{ name: 'Furniture', description: 'Office and home furniture' },
		{ name: 'Clothing', description: 'Apparel and accessories' },
		{ name: 'Food', description: 'Food products and ingredients' },
		{ name: 'Raw Materials', description: 'Materials for manufacturing' }
	];

	const subCategories = [
		{ name: 'Computers', description: 'Desktop and laptop computers', parent: 'Electronics' },
		{ name: 'Smartphones', description: 'Mobile phones and accessories', parent: 'Electronics' },
		{ name: 'Chairs', description: 'Office and dining chairs', parent: 'Furniture' },
		{ name: 'Tables', description: 'Office and dining tables', parent: 'Furniture' },
		{ name: 'T-shirts', description: 'Casual and formal t-shirts', parent: 'Clothing' },
		{ name: 'Pants', description: 'Casual and formal pants', parent: 'Clothing' },
		{ name: 'Dairy', description: 'Milk and dairy products', parent: 'Food' },
		{ name: 'Grains', description: 'Rice, wheat, and other grains', parent: 'Food' },
		{ name: 'Metals', description: 'Metal sheets and components', parent: 'Raw Materials' },
		{ name: 'Plastics', description: 'Plastic materials and components', parent: 'Raw Materials' }
	];

	const createdCategories: Record<string, any> = {};

	// Create main categories
	for (const category of categories) {
		const created = await prisma.item_categories.create({
			data: {
				name: category.name,
				description: category.description,
				created_at: new Date(),
				updated_at: new Date()
			}
		});

		createdCategories[category.name] = created;
		console.log(`Created category: ${category.name}`);
	}

	// Create subcategories
	for (const subCategory of subCategories) {
		const parentCategory = createdCategories[subCategory.parent];

		if (parentCategory) {
			const created = await prisma.item_categories.create({
				data: {
					name: subCategory.name,
					description: subCategory.description,
					parent_id: parentCategory.id,
					created_at: new Date(),
					updated_at: new Date()
				}
			});

			createdCategories[subCategory.name] = created;
			console.log(`Created subcategory: ${subCategory.name} under ${subCategory.parent}`);
		}
	}

	return Object.values(createdCategories);
}

/**
 * Create customers with different types
 */
async function createCustomers() {
	console.log('Creating customers...');

	const customerTypes = ['INDIVIDUAL', 'BUSINESS', 'GOVERNMENT'];
	const loyaltyTiers = ['STANDARD', 'SILVER', 'GOLD', 'PLATINUM'];
	const customers: any[] = [];

	for (let i = 0; i < 20; i++) {
		const isBusiness = i % 3 === 0;
		const customerType = faker.helpers.arrayElement(customerTypes);

		let name, email, contacts;

		if (isBusiness) {
			name = faker.company.name();
			email = faker.internet.email({ firstName: name.split(' ')[0], lastName: '', provider: 'example.com' }).toLowerCase();

			// Create 1-3 contacts for business customers
			const contactsCount = faker.number.int({ min: 1, max: 3 });
			contacts = [];

			for (let j = 0; j < contactsCount; j++) {
				contacts.push({
					name: faker.person.fullName(),
					position: faker.person.jobTitle(),
					email: faker.internet.email().toLowerCase(),
					phone: faker.phone.number()
				});
			}
		} else {
			const firstName = faker.person.firstName();
			const lastName = faker.person.lastName();
			name = `${firstName} ${lastName}`;
			email = faker.internet.email({ firstName, lastName }).toLowerCase();
			contacts = null;
		}

		const customer = await prisma.customers.create({
			data: {
				code: `CUST${(i + 1).toString().padStart(5, '0')}`,
				name,
				type: customerType,
				email,
				phone_number: faker.phone.number(),
				tax_id: isBusiness ? faker.finance.accountNumber(9) : null,
				contacts: contacts as any,
				credit_limit: randomDecimal(1000, 50000, 2),
				current_debt: randomDecimal(0, 10000, 2),
				loyalty_tier: faker.helpers.arrayElement(loyaltyTiers),
				loyalty_points: faker.number.int({ min: 0, max: 10000 }),
				notes: faker.helpers.maybe(() => faker.lorem.paragraph(), { probability: 0.7 }),
				created_at: faker.date.past(),
				updated_at: faker.date.recent()
			}
		});

		customers.push(customer);
		console.log(`Created customer: ${customer.name} (${customer.code})`);
	}

	return customers;
}

/**
 * Create customer addresses for each customer
 */
async function createCustomerAddresses(customers: any[]) {
	console.log('Creating customer addresses...');

	const addressTypes = ['BILLING', 'SHIPPING', 'BOTH'];
	const addresses: any[] = [];

	for (const customer of customers) {
		// Each customer gets 1-3 addresses
		const addressCount = faker.number.int({ min: 1, max: 3 });

		for (let i = 0; i < addressCount; i++) {
			const isDefault = i === 0; // First address is default
			const addressType = addressCount === 1 ? 'BOTH' : faker.helpers.arrayElement(addressTypes);

			const address = await prisma.customer_addresses.create({
				data: {
					customer_id: customer.id,
					type: addressType,
					street: faker.location.streetAddress(),
					city: faker.location.city(),
					state: faker.location.state(),
					postal_code: faker.location.zipCode(),
					country: faker.location.country(),
					is_default: isDefault,
					created_at: faker.date.past(),
					updated_at: faker.date.recent()
				}
			});

			addresses.push(address);
		}

		console.log(`Created ${addressCount} addresses for customer: ${customer.name}`);
	}

	return addresses;
}

/**
 * Create sequences for various document types
 */
async function createSequences() {
	console.log('Creating sequences...');

	const sequenceTypes = [
		'PURCHASE_ORDER', 'PURCHASE_REQUEST', 'PURCHASE_RECEIPT',
		'SALES_ORDER', 'DELIVERY_ORDER', 'INVOICE',
		'PAYMENT', 'PRODUCTION_ORDER', 'TRANSFER_ORDER'
	];

	const sequences: any[] = [];

	for (const type of sequenceTypes) {
		const sequence = await prisma.sequences.create({
			data: {
				id: type,
				value: 1
			}
		});

		sequences.push(sequence);
		console.log(`Created sequence for: ${type}`);
	}

	return sequences;
}

/**
 * Create audit logs for various user actions
 */
async function createAuditLogs(users: any[]) {
	console.log('Creating audit logs...');

	const actions = ['CREATE', 'UPDATE', 'DELETE', 'LOGIN', 'LOGOUT', 'EXPORT'];
	const resources = [
		'USER', 'ROLE', 'WAREHOUSE', 'INVENTORY', 'ITEM',
		'CUSTOMER', 'SUPPLIER', 'ORDER', 'INVOICE', 'PAYMENT'
	];

	const logs: any[] = [];
	const logCount = 50; // Generate 50 audit logs

	for (let i = 0; i < logCount; i++) {
		const user = faker.helpers.arrayElement(users);
		const action = faker.helpers.arrayElement(actions);
		const resource = faker.helpers.arrayElement(resources);

		let detail;
		switch (action) {
			case 'CREATE':
				detail = `Created new ${resource.toLowerCase()} record`;
				break;
			case 'UPDATE':
				detail = `Updated ${resource.toLowerCase()} record`;
				break;
			case 'DELETE':
				detail = `Deleted ${resource.toLowerCase()} record`;
				break;
			case 'LOGIN':
				detail = `User logged in`;
				break;
			case 'LOGOUT':
				detail = `User logged out`;
				break;
			case 'EXPORT':
				detail = `Exported ${resource.toLowerCase()} data`;
				break;
			default:
				detail = `Performed ${action} on ${resource}`;
		}

		const log = await prisma.audit_logs.create({
			data: {
				user_id: user.id,
				action,
				resource,
				detail,
				ip: faker.internet.ip(),
				user_agent: faker.internet.userAgent(),
				created_at: faker.date.recent()
			}
		});

		logs.push(log);
	}

	console.log(`Created ${logCount} audit logs`);
	return logs;
}

/**
 * Create inventory records for products in warehouses
 */
async function createInventories(warehouses: any[]) {
	console.log('Creating inventories...');

	const inventories: any[] = [];
	const productCount = 30; // Generate 30 products with inventory

	for (let i = 0; i < productCount; i++) {
		const productId = `PROD${(i + 1).toString().padStart(5, '0')}`;

		// Distribute products across warehouses
		for (const warehouse of warehouses) {
			if (faker.datatype.boolean(0.7)) { // 70% chance to have this product in this warehouse
				const inventoryId = faker.string.uuid();
				const quantity = faker.number.int({ min: 0, max: 1000 });

				// Some products have batch/lot tracking
				const hasBatchTracking = faker.datatype.boolean(0.4);
				const batchNumber = hasBatchTracking ? `BATCH-${faker.string.alphanumeric(6)}` : null;
				const lotNumber = hasBatchTracking ? `LOT-${faker.string.alphanumeric(4)}` : null;

				// Some products have expiry dates
				const hasExpiry = faker.datatype.boolean(0.3);
				const manufactureDate = hasExpiry ? faker.date.past({ years: 1 }) : undefined;
				const expiryDate = hasExpiry ? faker.date.future({ years: 2, refDate: manufactureDate }) : null;

				const inventory = await prisma.inventories.create({
					data: {
						id: inventoryId,
						product_id: productId,
						warehouse_id: warehouse.id,
						quantity,
						bin_location: faker.helpers.maybe(() => `BIN-${faker.string.alphanumeric(3)}`, { probability: 0.8 }),
						shelf_number: faker.helpers.maybe(() => `SHELF-${faker.string.alphanumeric(2)}`, { probability: 0.8 }),
						zone_code: faker.helpers.maybe(() => `ZONE-${faker.string.alpha({ length: 1, casing: 'upper' })}`, { probability: 0.8 }),
						batch_number: batchNumber,
						lot_number: lotNumber,
						manufacture_date: manufactureDate,
						expiry_date: expiryDate,
						created_at: faker.date.past(),
						updated_at: faker.date.recent()
					}
				});

				inventories.push(inventory);
			}
		}
	}

	console.log(`Created ${inventories.length} inventory records for ${productCount} products`);
	return inventories;
}

/**
 * Create inventory history records for tracking inventory changes
 */
async function createInventoryHistories(inventories: any[], users: any[]) {
	console.log('Creating inventory histories...');

	const historyTypes = ['RECEIPT', 'ISSUE', 'ADJUSTMENT', 'TRANSFER', 'COUNT'];
	const histories: any[] = [];

	// Create 2-5 history records for each inventory
	for (const inventory of inventories) {
		const historyCount = faker.number.int({ min: 2, max: 5 });

		let currentQty = parseFloat(inventory.quantity);

		for (let i = 0; i < historyCount; i++) {
			const type = faker.helpers.arrayElement(historyTypes);
			const user = faker.helpers.arrayElement(users);

			// Calculate quantities based on history type
			let quantity;
			let previousQty;
			let newQty;

			switch (type) {
				case 'RECEIPT':
					quantity = randomDecimal(10, 100, 2);
					previousQty = currentQty - quantity;
					newQty = currentQty;
					break;
				case 'ISSUE':
					quantity = randomDecimal(1, 50, 2);
					previousQty = currentQty + quantity;
					newQty = currentQty;
					break;
				case 'ADJUSTMENT':
					quantity = randomDecimal(-20, 20, 2);
					previousQty = currentQty - quantity;
					newQty = currentQty;
					break;
				case 'TRANSFER':
					quantity = randomDecimal(5, 30, 2);
					previousQty = currentQty - quantity;
					newQty = currentQty;
					break;
				case 'COUNT':
					quantity = 0; // No quantity change for count
					previousQty = currentQty;
					newQty = currentQty;
					break;
				default:
					quantity = randomDecimal(1, 50, 2);
					previousQty = currentQty - quantity;
					newQty = currentQty;
			}

			// Ensure we don't have negative previous quantities
			if (previousQty < 0) {
				previousQty = 0;
				newQty = quantity;
			}

			// Create reference based on type
			let reference;
			switch (type) {
				case 'RECEIPT':
					reference = `PO-${faker.string.numeric(5)}`;
					break;
				case 'ISSUE':
					reference = `SO-${faker.string.numeric(5)}`;
					break;
				case 'TRANSFER':
					reference = `TR-${faker.string.numeric(5)}`;
					break;
				case 'COUNT':
					reference = `COUNT-${faker.string.numeric(5)}`;
					break;
				default:
					reference = `ADJ-${faker.string.numeric(5)}`;
			}

			const history = await prisma.inventory_histories.create({
				data: {
					id: faker.string.uuid(),
					inventory_id: inventory.id,
					type,
					quantity,
					previous_qty: previousQty,
					new_qty: newQty,
					reference,
					note: faker.helpers.maybe(() => faker.lorem.sentence(), { probability: 0.7 }),
					created_at: faker.date.recent(),
					created_by: user.id.toString()
				}
			});

			histories.push(history);

			// Update current quantity for next history record
			currentQty = previousQty;
		}
	}

	console.log(`Created ${histories.length} inventory history records`);
	return histories;
}

/**
 * Create stock entry records for inventory movements
 */
async function createStockEntries(warehouses: any[], users: any[]) {
	console.log('Creating stock entries...');

	const entryTypes = ['RECEIPT', 'ISSUE', 'TRANSFER_IN', 'TRANSFER_OUT', 'ADJUSTMENT'];
	const entries: any[] = [];
	const entryCount = 50; // Generate 50 stock entries

	for (let i = 0; i < entryCount; i++) {
		const type = faker.helpers.arrayElement(entryTypes);
		const warehouse = faker.helpers.arrayElement(warehouses);
		const user = faker.helpers.arrayElement(users);
		const productId = `PROD${faker.number.int({ min: 1, max: 30 }).toString().padStart(5, '0')}`;

		// Create reference based on type
		let reference;
		switch (type) {
			case 'RECEIPT':
				reference = `PO-${faker.string.numeric(5)}`;
				break;
			case 'ISSUE':
				reference = `SO-${faker.string.numeric(5)}`;
				break;
			case 'TRANSFER_IN':
			case 'TRANSFER_OUT':
				reference = `TR-${faker.string.numeric(5)}`;
				break;
			default:
				reference = `ADJ-${faker.string.numeric(5)}`;
		}

		const entry = await prisma.stock_entries.create({
			data: {
				id: faker.string.uuid(),
				warehouse_id: warehouse.id,
				product_id: productId,
				type,
				quantity: randomDecimal(1, 100, 2),
				batch_number: faker.helpers.maybe(() => `BATCH-${faker.string.alphanumeric(6)}`, { probability: 0.4 }),
				lot_number: faker.helpers.maybe(() => `LOT-${faker.string.alphanumeric(4)}`, { probability: 0.4 }),
				reference,
				note: faker.helpers.maybe(() => faker.lorem.sentence(), { probability: 0.7 }),
				created_at: faker.date.recent(),
				created_by: user.id.toString()
			}
		});

		entries.push(entry);
	}

	console.log(`Created ${entries.length} stock entry records`);
	return entries;
}

// Execute the main function
main()
	.catch((e) => {
		console.error('Error during seeding:', e);
		process.exit(1);
	})
	.finally(async () => {
		// Close Prisma client connection
		await prisma.$disconnect();
	});
