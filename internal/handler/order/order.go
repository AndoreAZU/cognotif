package order

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"go.cognotif/internal/dto/request"
	"go.cognotif/internal/repository/constant"
	"go.cognotif/internal/repository/postgresql"
	"go.cognotif/internal/service"
	pkg_error "go.cognotif/pkg/error"
	pkg_util "go.cognotif/pkg/util"
)

type Order interface {
	CreateOrder(c echo.Context) error
	GetOrderByCust(c echo.Context) error
	GetOrderByAdmin(c echo.Context) error
	CompletingOrder(c echo.Context) error
	GenerateReport(c echo.Context) error
}

type Service struct {
	service.OrderService
}

type instance struct {
	Service
}

func NewOrder(srvc Service) Order {
	return &instance{Service: srvc}
}

func (x *instance) CreateOrder(c echo.Context) error {
	ctx := c.Request().Context()
	id_cust := c.Get(constant.CTX_CUST_ID).(string)

	var request request.CreateOrder
	if err := pkg_util.BindRequestAndValidate(c, &request); err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	list_tc, err := x.OrderService.CreateOrder(ctx, request.Items, id_cust)
	if err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	postgresql.CommitTxSql(list_tc)
	return c.JSON(200, constant.RESPONSE_SUCCESS)
}

func (x *instance) GetOrderByCust(c echo.Context) error {
	ctx := c.Request().Context()
	id_cust := c.Get(constant.CTX_CUST_ID).(string)

	var req request.GetOrderByCust
	if err := pkg_util.BindRequestAndValidate(c, &req); err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	orders, err := x.OrderService.GetOrder(ctx, request.GetOrder{IdCust: id_cust, IdOrder: req.IdOrder, Status: req.Status})
	if err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	return c.JSON(200, orders)
}

// admin handler
func (x *instance) GetOrderByAdmin(c echo.Context) error {
	ctx := c.Request().Context()

	var req request.GetOrderByAdmin
	if err := pkg_util.BindRequestAndValidate(c, &req); err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	orders, err := x.OrderService.GetOrder(ctx, request.GetOrder{IdCust: req.IdCust, IdOrder: req.IdOrder, Status: req.Status})
	if err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	return c.JSON(200, orders)
}

// admin handler
func (x *instance) GenerateReport(c echo.Context) error {
	ctx := c.Request().Context()

	orders, err := x.OrderService.GenerateReport(ctx)
	if err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	return c.JSON(200, orders)
}

func (x *instance) CompletingOrder(c echo.Context) error {
	ctx := c.Request().Context()
	id_cust := c.Get(constant.CTX_CUST_ID).(string)

	var requset request.CompletingOrder
	if err := pkg_util.BindRequestAndValidate(c, &requset); err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	tx, err := x.OrderService.CompletingOrder(ctx, id_cust, requset.Id)
	if err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	postgresql.CommitTxSql([]*sql.Tx{tx})
	return c.JSON(200, constant.RESPONSE_ORDER_COMPLETE)
}
