package constant

type Hashcode struct{}

const (
	CTX_CUST_ID           = "x-cust-id"
	ORDER_STATUS_COMPLETE = "COMPLETE"
	ORDER_STATUS_PENDING  = "PENDING"
)

var RESPONSE_SUCCESS = map[string]string{
	"message": "success",
}

var RESPONSE_ORDER_COMPLETE = map[string]string{
	"message": "your order already complete",
}
