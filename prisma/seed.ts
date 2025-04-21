import { PrismaClient } from '@prisma/client';
import { faker } from '@faker-js/faker';
import * as bcrypt from 'bcryptjs';
import { v4 as uuidv4 } from 'uuid';

/**
 * Helper function to generate random decimal numbers with specified precision
 */
function randomDecimal(min: number, max: number, precision: number = 2): number {
	const multipleOf = Math.pow(10, -precision);
	return Number(faker.number.float({ min, max, multipleOf }));
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
	const vendors = await createVendors();
	const skus = await createSKUs(vendors);
	const stores = await createStores(users);
	const stocks = await createStocks(skus, stores);
	await createStockHistories(stocks);
	await createStockEntries(skus, stores, users);
	const clients = await createClients();
	await createClientAddresses(clients);
	await createAuditLogs(users);

	console.log('Database seeding completed successfully!');
}

/**
 * Clear all existing data from the database in reverse order of dependencies
 */
async function clearExistingData() {
	console.log('Clearing existing data...');

	await prisma.audit_logs.deleteMany({});
	await prisma.stock_histories.deleteMany({});
	await prisma.stock_entries.deleteMany({});
	await prisma.stocks.deleteMany({});
	await prisma.client_addresses.deleteMany({});
	await prisma.clients.deleteMany({});
	await prisma.stores.deleteMany({});
	await prisma.skus.deleteMany({});
	await prisma.vendors.deleteMany({});
	await prisma.users.deleteMany({});
	await prisma.roles.deleteMany({});

	console.log('Existing data cleared.');
}

/**
 * Create roles with appropriate permissions
 */
async function createRoles() {
	console.log('Creating roles...');

	// Define all permissions from permissions.go
	// User permissions
	const userPermissions = [
		"user:create", "user:read", "user:update", "user:delete"
	];

	// Role permissions
	const rolePermissions = [
		"role:create", "role:read", "role:update", "role:delete"
	];

	// Store permissions
	const storePermissions = [
		"store:create", "store:read", "store:update", "store:delete"
	];

	// Stock permissions
	const stockPermissions = [
		"stock:read", "stock:update",
		"stock:entry:create", "stock:entry:read"
	];

	// Vendor permissions
	const vendorPermissions = [
		"vendor:create", "vendor:read", "vendor:update", "vendor:delete",
		"product:create", "product:read", "product:update", "product:delete",
		"contract:create", "contract:read", "contract:update", "contract:delete",
		"rating:create", "rating:read"
	];

	// Manufacturing permissions
	const manufacturingPermissions = [
		"manufacturing:facility:create", "manufacturing:facility:read", "manufacturing:facility:update", "manufacturing:facility:delete",
		"manufacturing:order:create", "manufacturing:order:read", "manufacturing:order:update", "manufacturing:order:delete",
		"manufacturing:bom:create", "manufacturing:bom:read", "manufacturing:bom:update", "manufacturing:bom:delete"
	];

	// Purchase permissions
	const purchasePermissions = [
		"purchase:request:create", "purchase:request:read", "purchase:request:update", "purchase:request:delete", "purchase:request:approve",
		"purchase:order:create", "purchase:order:read", "purchase:order:update", "purchase:order:delete", "purchase:order:approve",
		"purchase:receipt:create", "purchase:receipt:read", "purchase:receipt:update",
		"purchase:payment:create", "purchase:payment:read", "purchase:payment:update"
	];

	// Client permissions
	const clientPermissions = [
		"client:create", "client:read", "client:update", "client:delete",
		"client:address:create", "client:address:read", "client:address:update", "client:address:delete",
		"client:debt:read", "client:debt:update",
		"client:loyalty:read", "client:loyalty:update"
	];

	// Sales Order permissions
	const salesPermissions = [
		"sales:order:create", "sales:order:read", "sales:order:update", "sales:order:delete", "sales:order:confirm", "sales:order:cancel",
		"delivery:order:create", "delivery:order:read", "delivery:order:update", "delivery:order:process",
		"invoice:create", "invoice:read", "invoice:update", "invoice:issue", "invoice:pay"
	];

	// Finance permissions
	const financePermissions = [
		"finance:invoice:create", "finance:invoice:read", "finance:invoice:update", "finance:invoice:delete",
		"finance:payment:create", "finance:payment:read", "finance:payment:update", "finance:payment:process",
		"finance:report:read"
	];

	// Report permissions
	const reportPermissions = [
		"report:create", "report:read", "report:update", "report:delete", "report:export",
		"report:schedule:create", "report:schedule:read", "report:schedule:update", "report:schedule:delete"
	];

	// Audit permissions
	const auditPermissions = [
		"audit:log:read"
	];

	// All permissions combined
	const allPermissions = [
		...userPermissions, ...rolePermissions, ...storePermissions,
		...stockPermissions, ...vendorPermissions, ...manufacturingPermissions,
		...purchasePermissions, ...clientPermissions, ...salesPermissions,
		...financePermissions, ...reportPermissions, ...auditPermissions
	];

	const roles = [
		{
			name: 'Admin',
			permissions: allPermissions // Admin has all permissions
		},
		{
			name: 'Store Manager',
			permissions: [
				...storePermissions, ...stockPermissions,
				"product:read", "vendor:read",
				"sales:order:read", "delivery:order:read", "delivery:order:process",
				"manufacturing:facility:read", "manufacturing:order:read"
			]
		},
		{
			name: 'Sales Representative',
			permissions: [
				...clientPermissions,
				"sales:order:create", "sales:order:read", "sales:order:update", "sales:order:confirm",
				"delivery:order:create", "delivery:order:read",
				"invoice:create", "invoice:read",
				"stock:read", "product:read"
			]
		},
		{
			name: 'Inventory Clerk',
			permissions: [
				...stockPermissions,
				"store:read", "product:read", "vendor:read",
				"purchase:receipt:create", "purchase:receipt:read"
			]
		},
		{
			name: 'Finance Officer',
			permissions: [
				...financePermissions,
				"client:read", "client:debt:read", "client:debt:update",
				"vendor:read", "sales:order:read", "invoice:read", "invoice:update", "invoice:pay",
				"purchase:order:read", "purchase:payment:create", "purchase:payment:read",
				"report:read", "report:export"
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
 * 
 * Relationship: User -> Role (Many-to-One)
 * Each user has one role, and a role can be assigned to many users
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
 * Create vendors for product sourcing
 */
async function createVendors() {
	console.log('Creating vendors...');

	const vendorTypes = ['MANUFACTURER', 'DISTRIBUTOR', 'WHOLESALER', 'IMPORTER'];
	const countries = ['United States', 'China', 'Germany', 'Japan', 'France', 'Brazil', 'India', 'Vietnam', 'Thailand'];
	const paymentMethods = ['BANK_TRANSFER', 'CREDIT_CARD', 'CHECK', 'CASH', 'PAYPAL'];
	const currencies = ['USD', 'EUR', 'GBP', 'JPY', 'CNY', 'VND'];

	const vendors: any[] = [];

	for (let i = 0; i < 10; i++) {
		const vendorCode = `VND${(i + 1).toString().padStart(5, '0')}`;
		const companyName = faker.company.name();

		const vendor = await prisma.vendors.create({
			data: {
				code: vendorCode,
				name: companyName,
				type: faker.helpers.arrayElement(vendorTypes),
				address: faker.location.streetAddress({ useFullAddress: true }),
				country: faker.helpers.arrayElement(countries),
				email: faker.internet.email({ firstName: companyName.split(' ')[0], lastName: '' }).toLowerCase(),
				phone: faker.phone.number(),
				website: `https://www.${companyName.toLowerCase().replace(/[^a-zA-Z0-9]/g, '')}.com`,
				tax_id: faker.finance.accountNumber(9),
				payment_method: faker.helpers.arrayElement(paymentMethods),
				payment_days: faker.helpers.arrayElement([15, 30, 45, 60, 90]),
				currency: faker.helpers.arrayElement(currencies),
				rating: faker.number.float({ min: 1, max: 5, multipleOf: 0.01 }),
				created_at: faker.date.past(),
				updated_at: faker.date.recent()
			}
		});

		vendors.push(vendor);
		console.log(`Created vendor: ${vendor.name} (${vendor.code})`);
	}

	return vendors;
}

/**
 * Create SKUs (Stock Keeping Units) with relationships to vendors
 * 
 * Relationships:
 * 1. SKU -> Vendor (Many-to-One, as vendor_id)
 * 2. SKU -> Manufacturer (Many-to-One, as manufacturer_id)
 */
async function createSKUs(vendors: any[]) {
	console.log('Creating SKUs...');

	const skus: any[] = [];
	const categories = ['Electronics', 'Furniture', 'Office Supplies', 'Raw Materials', 'Tools', 'Packaging'];
	const unitMeasures = ['EA', 'KG', 'L', 'M', 'BOX', 'PALLET', 'ROLL', 'SET'];
	const statuses = ['ACTIVE', 'DISCONTINUED', 'PENDING'];

	for (let i = 0; i < 20; i++) {
		// Generate random manufacturer and vendor (they can be the same vendor)
		const manufacturerId = faker.helpers.arrayElement(vendors).id;
		const vendorId = faker.helpers.maybe(() => faker.helpers.arrayElement(vendors).id, { probability: 0.7 });

		const skuId = uuidv4();
		const skuCode = `SKU-${faker.string.alphanumeric(6).toUpperCase()}`;
		const name = faker.commerce.productName();
		const category = faker.helpers.arrayElement(categories);

		// Technical specifications as JSON
		const techSpecs = {
			dimensions: {
				length: faker.number.float({ min: 1, max: 100, multipleOf: 0.1 }),
				width: faker.number.float({ min: 1, max: 100, multipleOf: 0.1 }),
				height: faker.number.float({ min: 1, max: 100, multipleOf: 0.1 }),
				weight: faker.number.float({ min: 0.1, max: 50, multipleOf: 0.01 })
			},
			material: faker.commerce.productMaterial(),
			color: faker.color.human(),
			certifications: faker.helpers.arrayElements(
				['ISO', 'CE', 'FDA', 'ROHS', 'UL'],
				faker.number.int({ min: 0, max: 3 })
			)
		};

		const sku = await prisma.skus.create({
			data: {
				id: skuId,
				sku_code: skuCode,
				name: name,
				description: faker.commerce.productDescription(),
				unit_of_measure: faker.helpers.arrayElement(unitMeasures),
				price: faker.number.float({ min: 5, max: 1000, multipleOf: 0.01 }),
				category: category,
				technical_specs: techSpecs,
				manufacturer_id: manufacturerId,
				vendor_id: vendorId,
				image_url: faker.helpers.maybe(() => faker.image.url(), { probability: 0.7 }),
				status: faker.helpers.arrayElement(statuses),
				created_at: faker.date.past(),
				updated_at: faker.date.recent()
			}
		});

		skus.push(sku);
		console.log(`Created SKU: ${sku.name} (${sku.sku_code})`);
	}

	return skus;
}

/**
 * Create stores (warehouses) managed by users
 * 
 * Relationship: Store -> User (Many-to-One)
 * Each store has one manager (user), and a user can manage multiple stores
 */
async function createStores(users: any[]) {
	console.log('Creating stores...');

	const storeTypes = ['WAREHOUSE', 'RETAIL', 'DISTRIBUTION_CENTER', 'MANUFACTURING'];
	const storeStatuses = ['ACTIVE', 'INACTIVE', 'MAINTENANCE'];
	const stores: any[] = [];

	// Get users with store management permissions
	const storeManagers = users.filter(user =>
		user.role_id === 1 || user.role_id === 2 // Admin or Store Manager roles
	);

	// If no specific store managers, use any user
	const managerPool = storeManagers.length > 0 ? storeManagers : users;

	for (let i = 0; i < 5; i++) {
		const storeId = uuidv4();
		const storeType = faker.helpers.arrayElement(storeTypes);
		const manager = faker.helpers.arrayElement(managerPool);

		const storeName = `${faker.location.city()} ${storeType.replace('_', ' ')}`;
		const storeCode = `ST${i + 1}${storeType.substring(0, 3).toUpperCase()}`;

		const store = await prisma.stores.create({
			data: {
				id: storeId,
				name: storeName,
				code: storeCode,
				address: faker.location.streetAddress({ useFullAddress: true }),
				type: storeType,
				manager_id: manager.id,
				contact: faker.phone.number(),
				status: faker.helpers.arrayElement(storeStatuses),
				created_at: faker.date.past(),
				updated_at: faker.date.recent()
			}
		});

		stores.push(store);
		console.log(`Created store: ${store.name} (${store.code})`);
	}

	return stores;
}

/**
 * Create stocks for SKUs in different stores
 * 
 * Relationships:
 * 1. Stock -> SKU (Many-to-One)
 * 2. Stock -> Store (Many-to-One)
 */
async function createStocks(skus: any[], stores: any[]) {
	console.log('Creating stocks...');

	const stocks: any[] = [];

	// Create stock entries for various SKUs across stores
	for (const sku of skus) {
		// Distribute this SKU across some of the stores
		for (const store of stores) {
			if (faker.datatype.boolean(0.7)) { // 70% chance to have this SKU in this store
				const stockId = uuidv4();
				const quantity = faker.number.int({ min: 0, max: 1000 });

				// Some products have batch/lot tracking
				const hasBatchTracking = faker.datatype.boolean(0.4);
				const batchNumber = hasBatchTracking ? `BATCH-${faker.string.alphanumeric(6)}` : null;
				const lotNumber = hasBatchTracking ? `LOT-${faker.string.alphanumeric(4)}` : null;

				// Some products have expiry dates
				const hasExpiry = faker.datatype.boolean(0.3);
				const manufactureDate = hasExpiry ? faker.date.past({ years: 1 }) : null;
				const expiryDate = hasExpiry ? faker.date.future({ years: 2, refDate: manufactureDate! }) : null;

				const stock = await prisma.stocks.create({
					data: {
						id: stockId,
						sk_uid: sku.id,
						store_id: store.id,
						quantity: quantity,
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

				stocks.push(stock);
				console.log(`Created stock for SKU ${sku.sku_code} in store ${store.name}: ${quantity} units`);
			}
		}
	}

	return stocks;
}

/**
 * Create stock history records for tracking stock changes
 * 
 * Relationship: StockHistory -> Stock (Many-to-One)
 * Each stock history entry is associated with one stock record
 */
async function createStockHistories(stocks: any[]) {
	console.log('Creating stock histories...');

	const historyTypes = ['RECEIPT', 'ISSUE', 'ADJUSTMENT', 'TRANSFER', 'COUNT'];
	const histories: any[] = [];

	// Create 2-5 history records for each stock
	for (const stock of stocks) {
		const historyCount = faker.number.int({ min: 2, max: 5 });

		let currentQty = parseFloat(stock.quantity);

		for (let i = 0; i < historyCount; i++) {
			const type = faker.helpers.arrayElement(historyTypes);

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

			const history = await prisma.stock_histories.create({
				data: {
					id: uuidv4(),
					stock_id: stock.id,
					type,
					quantity,
					previous_qty: previousQty,
					new_qty: newQty,
					reference,
					note: faker.helpers.maybe(() => faker.lorem.sentence(), { probability: 0.7 }),
					created_at: faker.date.recent(),
					created_by: faker.string.uuid() // Just a placeholder for user ID
				}
			});

			histories.push(history);

			// Update current quantity for next history record
			currentQty = previousQty;
		}
	}

	console.log(`Created ${histories.length} stock history records`);
	return histories;
}

/**
 * Create stock entry records for inventory movements
 * 
 * Relationships:
 * 1. StockEntry -> SKU (Many-to-One)
 * 2. StockEntry -> Store (Many-to-One)
 */
async function createStockEntries(skus: any[], stores: any[], users: any[]) {
	console.log('Creating stock entries...');

	const entryTypes = ['RECEIPT', 'ISSUE', 'TRANSFER_IN', 'TRANSFER_OUT', 'ADJUSTMENT'];
	const entries: any[] = [];
	const entryCount = 50; // Generate 50 stock entries

	for (let i = 0; i < entryCount; i++) {
		const type = faker.helpers.arrayElement(entryTypes);
		const store = faker.helpers.arrayElement(stores);
		const sku = faker.helpers.arrayElement(skus);
		const user = faker.helpers.arrayElement(users);

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
				id: uuidv4(),
				sk_uid: sku.id,
				store_id: store.id,
				type,
				quantity: randomDecimal(1, 100, 2),
				batch_number: faker.helpers.maybe(() => `BATCH-${faker.string.alphanumeric(6)}`, { probability: 0.4 }),
				lot_number: faker.helpers.maybe(() => `LOT-${faker.string.alphanumeric(4)}`, { probability: 0.4 }),
				manufacture_date: faker.helpers.maybe(() => faker.date.past({ years: 1 }), { probability: 0.3 }),
				expiry_date: faker.helpers.maybe(() => faker.date.future({ years: 2 }), { probability: 0.3 }),
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

/**
 * Create clients (customers)
 */
async function createClients() {
	console.log('Creating clients...');

	const clientTypes = ['INDIVIDUAL', 'BUSINESS', 'GOVERNMENT'];
	const loyaltyTiers = ['STANDARD', 'SILVER', 'GOLD', 'PLATINUM'];
	const clients: any[] = [];

	for (let i = 0; i < 20; i++) {
		const isBusiness = i % 3 === 0;
		const clientType = faker.helpers.arrayElement(clientTypes);

		let name, email, contacts;

		if (isBusiness) {
			name = faker.company.name();
			email = faker.internet.email({ firstName: name.split(' ')[0], lastName: '', provider: 'example.com' }).toLowerCase();

			// Create 1-3 contacts for business clients
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
			contacts = undefined; // Use undefined instead of null for JSON fields in Prisma
		}

		const client = await prisma.clients.create({
			data: {
				code: `CLT${(i + 1).toString().padStart(5, '0')}`,
				name,
				type: clientType,
				email,
				phone_number: faker.phone.number(),
				tax_id: isBusiness ? faker.finance.accountNumber(9) : null,
				contacts: contacts,
				credit_limit: randomDecimal(1000, 50000, 2),
				current_debt: randomDecimal(0, 10000, 2),
				loyalty_tier: faker.helpers.arrayElement(loyaltyTiers),
				loyalty_points: faker.number.int({ min: 0, max: 10000 }),
				notes: faker.helpers.maybe(() => faker.lorem.paragraph(), { probability: 0.7 }),
				created_at: faker.date.past(),
				updated_at: faker.date.recent()
			}
		});

		clients.push(client);
		console.log(`Created client: ${client.name} (${client.code})`);
	}

	return clients;
}

/**
 * Create client addresses for each client
 * 
 * Relationship: ClientAddress -> Client (Many-to-One)
 * Each client can have multiple addresses
 */
async function createClientAddresses(clients: any[]) {
	console.log('Creating client addresses...');

	const addressTypes = ['BILLING', 'SHIPPING', 'BOTH'];
	const addresses: any[] = [];

	for (const client of clients) {
		// Each client gets 1-3 addresses
		const addressCount = faker.number.int({ min: 1, max: 3 });

		for (let i = 0; i < addressCount; i++) {
			const isDefault = i === 0; // First address is default
			const addressType = addressCount === 1 ? 'BOTH' : faker.helpers.arrayElement(addressTypes);

			const address = await prisma.client_addresses.create({
				data: {
					client_id: client.id,
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

		console.log(`Created ${addressCount} addresses for client: ${client.name}`);
	}

	return addresses;
}

/**
 * Create audit logs for various user actions
 * 
 * Relationship: AuditLog -> User (Many-to-One)
 * Each audit log is associated with one user, and a user can have many audit logs
 */
async function createAuditLogs(users: any[]) {
	console.log('Creating audit logs...');

	const actions = ['CREATE', 'UPDATE', 'DELETE', 'LOGIN', 'LOGOUT', 'EXPORT'];
	const resources = [
		'USER', 'ROLE', 'STORE', 'INVENTORY', 'SKU',
		'CLIENT', 'VENDOR', 'ORDER', 'INVOICE', 'PAYMENT'
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
