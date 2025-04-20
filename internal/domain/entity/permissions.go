package entity

// Supplier permissions
const (
	SupplierCreate Permission = "supplier:create"
	SupplierRead   Permission = "supplier:read"
	SupplierUpdate Permission = "supplier:update"
	SupplierDelete Permission = "supplier:delete"

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

// Customer permissions
const (
	CustomerCreate Permission = "customer:create"
	CustomerRead   Permission = "customer:read"
	CustomerUpdate Permission = "customer:update"
	CustomerDelete Permission = "customer:delete"

	CustomerAddressCreate Permission = "customer:address:create"
	CustomerAddressRead   Permission = "customer:address:read"
	CustomerAddressUpdate Permission = "customer:address:update"
	CustomerAddressDelete Permission = "customer:address:delete"

	CustomerDebtRead   Permission = "customer:debt:read"
	CustomerDebtUpdate Permission = "customer:debt:update"

	CustomerLoyaltyRead   Permission = "customer:loyalty:read"
	CustomerLoyaltyUpdate Permission = "customer:loyalty:update"
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
