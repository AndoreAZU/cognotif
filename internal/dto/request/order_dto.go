package request

type GetOrderByAdmin struct {
	IdCust  string `query:"id-customer"`
	IdOrder string `query:"id_order"`
	Status  string `query:"status"`
}

type GetOrderByCust struct {
	IdOrder string `query:"id-order"`
	Status  string `query:"status"`
}

type CompletingOrder struct {
	Id string `query:"id" validate:"required"`
}

type GetOrder struct {
	IdOrder string `json:"id_order"`
	IdCust  string `json:"id_customer"`
	Status  string `json:"status"`
}

type CreateOrder struct {
	Items []ItemOrder `json:"items" validate:"required,min=1,dive,required"`
}

type ItemOrder struct {
	IdProduct int `json:"id_product" validate:"required"`
	Quantity  int `json:"quantity" validate:"required"`
}
