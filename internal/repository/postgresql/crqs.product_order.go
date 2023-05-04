package postgresql

import (
	"context"
	"database/sql"
)

type TableProductOrder interface {
	InsertNewProductOrder(ctx context.Context, id, id_order string, id_product, qty int) (*sql.Tx, error)
	GetProductOrderByIdOrder(ctx context.Context, id_order string) ([]ProductOrder, error)
}

type ProductOrder struct {
	Id        string `json:"id"`
	IdOrder   string `json:"id_order"`
	IdProduct int    `json:"id_product"`
	Quantity  int    `json:"quantity"`
}

func (x *instance) InsertNewProductOrder(ctx context.Context, id, id_order string, id_product, qty int) (*sql.Tx, error) {
	return x.ExecContext(ctx, x.GetSQL("product_order.insert"), id, id_order, id_product, qty)
}

func (x *instance) GetProductOrderByIdOrder(ctx context.Context, id_order string) ([]ProductOrder, error) {
	stmt, err := x.QueryContext(ctx, x.GetSQL("product_order.select_by_id_order"), id_order)
	if err != nil {
		return []ProductOrder{}, err
	}

	defer stmt.Close()

	var products []ProductOrder
	for stmt.Next() {
		var product ProductOrder
		if err = stmt.Scan(&product.Id, &product.IdOrder, &product.IdProduct, &product.Quantity); err != nil {
			return []ProductOrder{}, err
		}
		products = append(products, product)
	}

	return products, nil
}
