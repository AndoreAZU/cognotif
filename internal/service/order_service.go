package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"go.cognotif/internal/dto/request"
	"go.cognotif/internal/dto/response"
	"go.cognotif/internal/repository/constant"
	"go.cognotif/internal/repository/postgresql"
	pkg_error "go.cognotif/pkg/error"
	pkg_logger "go.cognotif/pkg/logger"
	pkg_util "go.cognotif/pkg/util"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req []request.ItemOrder, id_cust string) ([]*sql.Tx, error)
	GetOrder(ctx context.Context, query request.GetOrder) (response.GetOrder, error)
	CompletingOrder(ctx context.Context, id_cust, id_order string) (*sql.Tx, error)
	GenerateReport(ctx context.Context) ([]postgresql.OrderReport, error)
}

type orderService struct {
	*pkg_logger.Logger
	postgresql.Postgresql
	ProductOrderService
	ProductService
}

func NewOrderService(log *pkg_logger.Logger, psql postgresql.Postgresql, pos ProductOrderService, ps ProductService) OrderService {
	return &orderService{Logger: log, Postgresql: psql, ProductOrderService: pos, ProductService: ps}
}

func (os *orderService) CreateOrder(ctx context.Context, req []request.ItemOrder, id_cust string) ([]*sql.Tx, error) {
	var list_tx []*sql.Tx
	id_order := uuid.NewString()

	tx, err := os.Postgresql.InsertNewOrder(ctx, id_order, id_cust, pkg_util.Now(), constant.ORDER_STATUS_PENDING)
	if err != nil {
		os.Hashcode(ctx).Error("error when insert order: ", err)
		return nil, fmt.Errorf(pkg_error.GENERAL_ERROR)
	}
	list_tx = append(list_tx, tx)

	for _, p := range req {
		_, err := os.ProductService.GetProductById(ctx, p.IdProduct)
		if err != nil {
			os.Hashcode(ctx).Error("error when get product: ", err)
			return nil, postgresql.RollbackTxSql(list_tx, err)
		}

		tx, err := os.ProductOrderService.InsertProductOrder(ctx, id_order, p.IdProduct, p.Quantity)
		if err != nil {
			os.Hashcode(ctx).Error("error when insert product order: ", err)
			return nil, postgresql.RollbackTxSql(list_tx, fmt.Errorf(pkg_error.GENERAL_ERROR))
		}
		list_tx = append(list_tx, tx)
	}
	return list_tx, nil
}

func (os *orderService) GetOrder(ctx context.Context, query request.GetOrder) (response.GetOrder, error) {
	var list_order response.GetOrder

	orders, err := os.Postgresql.GetOrderByQuery(ctx, query.IdOrder, query.IdCust, query.Status)
	if err != nil {
		os.Hashcode(ctx).Error("error when get order: ", err)
		if err == sql.ErrNoRows {
			return response.GetOrder{}, fmt.Errorf(pkg_error.ORDER_NOT_EXIST)
		}
		return response.GetOrder{}, err
	}

	for _, o := range orders {
		list_po, err := os.ProductOrderService.GetProductOrderByIdOrder(ctx, o.Id)
		if err != nil {
			os.Hashcode(ctx).Error("error when get product order: ", err)
			return response.GetOrder{}, fmt.Errorf(pkg_error.GENERAL_ERROR)
		}

		var products []response.ProductDetail
		for _, po := range list_po {
			product, err := os.ProductService.GetProductById(ctx, po.IdProduct)
			if err != nil {
				os.Hashcode(ctx).Error("error when get product: ", err)
				return response.GetOrder{}, fmt.Errorf(pkg_error.GENERAL_ERROR)
			}

			products = append(products, response.ProductDetail{
				Id:          product.Id,
				Name:        product.Name,
				Price:       product.Price,
				Description: product.Description,
				Image:       product.Image,
				Quantity:    po.Quantity,
			})
		}

		list_order.Orders = append(list_order.Orders, response.Order{
			Id:      o.Id,
			IdCust:  o.IdCust,
			Date:    o.Date,
			Status:  o.Status,
			Product: products,
		})
	}

	return list_order, nil
}

func (os *orderService) CompletingOrder(ctx context.Context, id_cust, id_order string) (*sql.Tx, error) {
	_, err := os.GetOrder(ctx, request.GetOrder{IdCust: id_cust, Status: constant.ORDER_STATUS_PENDING, IdOrder: id_order})
	if err != nil {
		os.Hashcode(ctx).Error("error when get order: ", err)
		return nil, err
	}

	tx, err := os.Postgresql.UpdateOrderStatus(ctx, id_order, constant.ORDER_STATUS_COMPLETE)
	if err != nil {
		os.Hashcode(ctx).Error("error when update status order: ", err)
		return nil, fmt.Errorf(pkg_error.GENERAL_ERROR)
	}

	return tx, nil
}

func (os *orderService) GenerateReport(ctx context.Context) ([]postgresql.OrderReport, error) {
	report, err := os.Postgresql.GenerateDataReport(ctx)
	if err != nil {
		os.Hashcode(ctx).Error("error when generate report: ", err)
		return nil, fmt.Errorf(pkg_error.GENERAL_ERROR)
	}

	return report, nil
}
