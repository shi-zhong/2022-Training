package model

var (
	UserTypeCustomer   uint8 = 1
	UserTypeMerchant   uint8 = 2
//	UserTypeAdmin      uint8 = 3
//	UserTypeSuperAdmin uint8 = 4
)

var (
	AddressNotExsit   uint8 = 0
	AddressDefault    uint8 = 1
	AddressNotDefault uint8 = 2
	AddressDelete     uint8 = 3
)

var (
	CommodityStatusNotCreate    uint = 0
	CommodityStatusOnShelf      uint = 1
	CommodityStatusOffShelf     uint = 2
	CommodityStatusForceOnShelf uint = 3
	CommodityStatusNotEnought   uint = 4
	CommodityStatusDeleted      uint = 5
)

var (
	OrderCreated   uint = 1
	OrderDue       uint = 2
	OrderCommodity uint = 3
	OrderFinish    uint = 4
	OrderCancel    uint = 5
)

var (
	// ShareBillWaitingForTwo 主人入团， 还差俩
	ShareBillWaitingForTwo uint8 = 1
	// ShareBillWaitingForOne 还差一个
	ShareBillWaitingForOne uint8 = 2
	ShareBillSuccess       uint8 = 3
	ShareBillFailed        uint8 = 4
)
