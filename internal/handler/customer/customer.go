package customer

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

type Customer interface {
	RegisterCustomer(c echo.Context) error
	LoginCustomer(c echo.Context) error
	LoginAdmin(c echo.Context) error
	GetProfile(c echo.Context) error
}

type Service struct {
	service.CustomerService
}

type instance struct {
	Service
}

func NewCustomer(srvc Service) Customer {
	return &instance{Service: srvc}
}

func (x *instance) RegisterCustomer(c echo.Context) error {
	ctx := c.Request().Context()

	var request request.RegisterCustomer
	if err := pkg_util.BindRequestAndValidate(c, &request); err != nil {
		return pkg_error.CreateError(c, err.Error())
	}
	tx, err := x.CustomerService.RegisterNewCustomer(ctx, request)
	if err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	postgresql.CommitTxSql([]*sql.Tx{tx})
	return c.JSON(200, constant.RESPONSE_SUCCESS)
}

func (x *instance) LoginCustomer(c echo.Context) error {
	ctx := c.Request().Context()

	var request request.Login
	if err := pkg_util.BindRequestAndValidate(c, &request); err != nil {
		return pkg_error.CreateError(c, err.Error())
	}
	login, err := x.CustomerService.CustomerLogin(ctx, request.Email, request.Password)
	if err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	return c.JSON(200, login)
}

func (x *instance) LoginAdmin(c echo.Context) error {
	ctx := c.Request().Context()

	var request request.Login
	if err := pkg_util.BindRequestAndValidate(c, &request); err != nil {
		return pkg_error.CreateError(c, err.Error())
	}
	login, err := x.CustomerService.AdminLogin(ctx, request.Email, request.Password)
	if err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	return c.JSON(200, login)
}

func (x *instance) GetProfile(c echo.Context) error {
	ctx := c.Request().Context()
	id_cust := c.Get(constant.CTX_CUST_ID).(string)

	customer, err := x.CustomerService.GetCustomerById(ctx, id_cust)
	if err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	return c.JSON(200, customer)
}
