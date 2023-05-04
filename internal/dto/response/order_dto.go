package response

type GetOrder struct {
	Orders []Order `json:"order"`
}

type Order struct {
	Id      string          `json:"id"`
	IdCust  string          `json:"id_customer"`
	Date    string          `json:"date"`
	Status  string          `json:"status"`
	Product []ProductDetail `json:"products"`
}

type ProductDetail struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Quantity    int     `json:"quantity"`
}
