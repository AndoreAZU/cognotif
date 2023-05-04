package postgresql

import (
	"context"
)

type TableProduct interface {
	GetListProduct(ctx context.Context) ([]Product, error)
	GetProductById(ctx context.Context, id int) (Product, error)
}

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
}

func (x *instance) GetListProduct(ctx context.Context) ([]Product, error) {
	stmt, err := x.QueryContext(ctx, x.GetSQL("product.select_all"))
	if err != nil {
		return []Product{}, err
	}

	defer stmt.Close()

	var products []Product
	for stmt.Next() {
		var product Product
		if err = stmt.Scan(&product.Id, &product.Name, &product.Price, &product.Description, &product.Image); err != nil {
			return []Product{}, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (x *instance) GetProductById(ctx context.Context, id int) (Product, error) {
	stmt, err := x.PrepareContext(ctx, x.GetSQL("product.select_by_id"))
	if err != nil {
		return Product{}, err
	}

	defer stmt.Close()

	var p Product
	err = stmt.QueryRow(id).Scan(&p.Id, &p.Name, &p.Price, &p.Description, &p.Image)
	if err != nil {
		return Product{}, err
	}

	return p, nil
}
