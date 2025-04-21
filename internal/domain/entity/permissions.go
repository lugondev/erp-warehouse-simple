package entity

// Permission represents a permission string
type Permission string

// User permissions
const (
	UserCreate Permission = "user:create"
	UserRead   Permission = "user:read"
	UserUpdate Permission = "user:update"
	UserDelete Permission = "user:delete"
)

// Role permissions
const (
	RoleCreate Permission = "role:create"
	RoleRead   Permission = "role:read"
	RoleUpdate Permission = "role:update"
	RoleDelete Permission = "role:delete"
)

// Store permissions
const (
	StoreCreate Permission = "store:create"
	StoreRead   Permission = "store:read"
	StoreUpdate Permission = "store:update"
	StoreDelete Permission = "store:delete"
)

// Stock permissions
const (
	StockRead   Permission = "stock:read"
	StockUpdate Permission = "stock:update"

	StockEntryCreate Permission = "stock:entry:create"
	StockEntryRead   Permission = "stock:entry:read"
)

// Vendor permissions
const (
	VendorCreate Permission = "vendor:create"
	VendorRead   Permission = "vendor:read"
	VendorUpdate Permission = "vendor:update"
	VendorDelete Permission = "vendor:delete"

	ProductCreate Permission = "product:create"
	ProductRead   Permission = "product:read"
	ProductUpdate Permission = "product:update"
	ProductDelete Permission = "product:delete"

	ContractCreate Permission = "contract:create"
	ContractRead   Permission = "contract:read"
	ContractUpdate Permission = "contract:update"
	ContractDelete Permission = "contract:delete"

	RatingCreate Permission = "rating:create"
	RatingRead   Permission = "rating:read"
)

// Manufacturing permissions
const (
	ManufacturingFacilityCreate Permission = "manufacturing:facility:create"
	ManufacturingFacilityRead   Permission = "manufacturing:facility:read"
	ManufacturingFacilityUpdate Permission = "manufacturing:facility:update"
	ManufacturingFacilityDelete Permission = "manufacturing:facility:delete"

	ProductionOrderCreate Permission = "manufacturing:order:create"
	ProductionOrderRead   Permission = "manufacturing:order:read"
	ProductionOrderUpdate Permission = "manufacturing:order:update"
	ProductionOrderDelete Permission = "manufacturing:order:delete"

	BOMCreate Permission = "manufacturing:bom:create"
	BOMRead   Permission = "manufacturing:bom:read"
	BOMUpdate Permission = "manufacturing:bom:update"
	BOMDelete Permission = "manufacturing:bom:delete"
)

// Purchase permissions
const (
	PurchaseRequestCreate  Permission = "purchase:request:create"
	PurchaseRequestRead    Permission = "purchase:request:read"
	PurchaseRequestUpdate  Permission = "purchase:request:update"
	PurchaseRequestDelete  Permission = "purchase:request:delete"
	PurchaseRequestApprove Permission = "purchase:request:approve"

	PurchaseOrderCreate  Permission = "purchase:order:create"
	PurchaseOrderRead    Permission = "purchase:order:read"
	PurchaseOrderUpdate  Permission = "purchase:order:update"
	PurchaseOrderDelete  Permission = "purchase:order:delete"
	PurchaseOrderApprove Permission = "purchase:order:approve"

	PurchaseReceiptCreate Permission = "purchase:receipt:create"
	PurchaseReceiptRead   Permission = "purchase:receipt:read"
	PurchaseReceiptUpdate Permission = "purchase:receipt:update"

	PurchasePaymentCreate Permission = "purchase:payment:create"
	PurchasePaymentRead   Permission = "purchase:payment:read"
	PurchasePaymentUpdate Permission = "purchase:payment:update"
)

// Client permissions
const (
	ClientCreate Permission = "client:create"
	ClientRead   Permission = "client:read"
	ClientUpdate Permission = "client:update"
	ClientDelete Permission = "client:delete"

	ClientAddressCreate Permission = "client:address:create"
	ClientAddressRead   Permission = "client:address:read"
	ClientAddressUpdate Permission = "client:address:update"
	ClientAddressDelete Permission = "client:address:delete"

	ClientDebtRead   Permission = "client:debt:read"
	ClientDebtUpdate Permission = "client:debt:update"

	ClientLoyaltyRead   Permission = "client:loyalty:read"
	ClientLoyaltyUpdate Permission = "client:loyalty:update"
)

// Sales Order permissions
const (
	SalesOrderCreate  Permission = "sales:order:create"
	SalesOrderRead    Permission = "sales:order:read"
	SalesOrderUpdate  Permission = "sales:order:update"
	SalesOrderDelete  Permission = "sales:order:delete"
	SalesOrderConfirm Permission = "sales:order:confirm"
	SalesOrderCancel  Permission = "sales:order:cancel"

	DeliveryOrderCreate  Permission = "delivery:order:create"
	DeliveryOrderRead    Permission = "delivery:order:read"
	DeliveryOrderUpdate  Permission = "delivery:order:update"
	DeliveryOrderProcess Permission = "delivery:order:process"

	InvoiceCreate Permission = "invoice:create"
	InvoiceRead   Permission = "invoice:read"
	InvoiceUpdate Permission = "invoice:update"
	InvoiceIssue  Permission = "invoice:issue"
	InvoicePay    Permission = "invoice:pay"
)

// Finance permissions
const (
	FinanceInvoiceCreate Permission = "finance:invoice:create"
	FinanceInvoiceRead   Permission = "finance:invoice:read"
	FinanceInvoiceUpdate Permission = "finance:invoice:update"
	FinanceInvoiceDelete Permission = "finance:invoice:delete"

	FinancePaymentCreate  Permission = "finance:payment:create"
	FinancePaymentRead    Permission = "finance:payment:read"
	FinancePaymentUpdate  Permission = "finance:payment:update"
	FinancePaymentProcess Permission = "finance:payment:process"

	FinanceReportRead Permission = "finance:report:read"
)

// Report permissions
const (
	ReportCreate Permission = "report:create"
	ReportRead   Permission = "report:read"
	ReportUpdate Permission = "report:update"
	ReportDelete Permission = "report:delete"
	ReportExport Permission = "report:export"

	ReportScheduleCreate Permission = "report:schedule:create"
	ReportScheduleRead   Permission = "report:schedule:read"
	ReportScheduleUpdate Permission = "report:schedule:update"
	ReportScheduleDelete Permission = "report:schedule:delete"
)

// Audit permissions
const (
	AuditLogRead Permission = "audit:log:read"
)
