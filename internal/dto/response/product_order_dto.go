package response

type GetProductOrder struct {
	Id        string `json:"id"`
	IdOrder   string `json:"id_order"`
	IdProduct int    `json:"id_product"`
	Quantity  int    `json:"quantity"`
}
