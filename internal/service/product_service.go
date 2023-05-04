package service

import (
	"context"
	"database/sql"
	"fmt"

	"go.cognotif/internal/dto/response"
	"go.cognotif/internal/repository/postgresql"
	pkg_error "go.cognotif/pkg/error"
	pkg_logger "go.cognotif/pkg/logger"
)

type ProductService interface {
	GetListProduct(ctx context.Context) ([]response.GetProduct, error)
	GetProductById(ctx context.Context, id int) (response.GetProduct, error)
}

type productService struct {
	*pkg_logger.Logger
	postgresql.Postgresql
}

func NewProductService(log *pkg_logger.Logger, psql postgresql.Postgresql) ProductService {
	return &productService{Logger: log, Postgresql: psql}
}

func (ps *productService) GetListProduct(ctx context.Context) ([]response.GetProduct, error) {
	var list_product []response.GetProduct

	products, err := ps.Postgresql.GetListProduct(ctx)
	if err != nil {
		ps.Hashcode(ctx).Error("error when get list product: ", err)
		return nil, fmt.Errorf(pkg_error.GENERAL_ERROR)
	}

	for _, p := range products {
		list_product = append(list_product, response.GetProduct{
			Id:          p.Id,
			Name:        p.Name,
			Price:       p.Price,
			Description: p.Description,
			Image:       p.Image,
		})
	}

	return list_product, nil
}

func (ps *productService) GetProductById(ctx context.Context, id int) (response.GetProduct, error) {
	product, err := ps.Postgresql.GetProductById(ctx, id)
	if err != nil {
		ps.Hashcode(ctx).Error("error when get product: ", err)
		if err == sql.ErrNoRows {
			return response.GetProduct{}, fmt.Errorf(pkg_error.PRODUCT_NOT_EXIST)
		}
		return response.GetProduct{}, fmt.Errorf(pkg_error.GENERAL_ERROR)
	}

	return response.GetProduct{
		Id:          product.Id,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Image:       product.Image,
	}, nil
}
