package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"go.cognotif/internal/dto/response"
	"go.cognotif/internal/repository/postgresql"
	pkg_error "go.cognotif/pkg/error"
	pkg_logger "go.cognotif/pkg/logger"
)

type ProductOrderService interface {
	InsertProductOrder(ctx context.Context, id_order string, id_product, qty int) (*sql.Tx, error)
	GetProductOrderByIdOrder(ctx context.Context, id_order string) ([]response.GetProductOrder, error)
}

type productOrderService struct {
	*pkg_logger.Logger
	postgresql.Postgresql
}

func NewProductOrderService(log *pkg_logger.Logger, psql postgresql.Postgresql) ProductOrderService {
	return &productOrderService{Logger: log, Postgresql: psql}
}

func (pos *productOrderService) InsertProductOrder(ctx context.Context, id_order string, id_product, qty int) (*sql.Tx, error) {
	return pos.Postgresql.InsertNewProductOrder(ctx, uuid.NewString(), id_order, id_product, qty)
}

func (pos *productOrderService) GetProductOrderByIdOrder(ctx context.Context, id_order string) ([]response.GetProductOrder, error) {
	var list_po []response.GetProductOrder

	po, err := pos.Postgresql.GetProductOrderByIdOrder(ctx, id_order)
	if err != nil {
		pos.Hashcode(ctx).Error("error when get product order")
		return nil, fmt.Errorf(pkg_error.GENERAL_ERROR)
	}

	for _, p := range po {
		list_po = append(list_po, response.GetProductOrder{
			Id:        p.Id,
			IdOrder:   p.IdOrder,
			IdProduct: p.IdProduct,
			Quantity:  p.Quantity,
		})
	}

	return list_po, nil
}
