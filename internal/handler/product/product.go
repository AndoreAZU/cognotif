package product

import (
	"github.com/labstack/echo/v4"
	"go.cognotif/internal/service"
	pkg_error "go.cognotif/pkg/error"
)

type Product interface {
	GetListProduct(c echo.Context) error
}

type Service struct {
	service.ProductService
}

type instance struct {
	Service
}

func NewProduct(srvc Service) Product {
	return &instance{Service: srvc}
}

func (x *instance) GetListProduct(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := x.ProductService.GetListProduct(ctx)
	if err != nil {
		return pkg_error.CreateError(c, err.Error())
	}

	return c.JSON(200, products)
}
