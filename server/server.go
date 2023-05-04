package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.cognotif/internal/handler/customer"
	"go.cognotif/internal/handler/order"
	"go.cognotif/internal/handler/product"
	cognotif_middleware "go.cognotif/internal/middleware"
	"go.cognotif/internal/repository/validator"
	pkg_logger "go.cognotif/pkg/logger"
)

type Server interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

type Configuration struct {
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

type Handler struct {
	customer.Customer
	product.Product
	order.Order
}

type Depedency struct {
	*validator.Validator
	*pkg_logger.Logger
	cognotif_middleware.MiddlewareHTTP
}

type instance struct {
	Configuration
	Depedency
	handler *http.Server
}

func New(cfg Configuration, dep Depedency, han Handler) Server {
	handler := echo.New()
	handler.Validator = dep.Validator
	handler.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 100, Burst: 10, ExpiresIn: 1 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}))

	handler.Use(echo.WrapMiddleware(dep.Middleware))
	handler.Use(middleware.Recover())
	handler.Use(middleware.CORS())

	root := handler.Group("/v1")
	customer := root.Group("", dep.MiddlewareCustomer)
	admin := root.Group("", dep.MiddlewareAdmin)

	// non auth
	root.POST("/user", han.Customer.RegisterCustomer)
	root.POST("/user/login", han.Customer.LoginCustomer)
	root.POST("/admin/login", han.Customer.LoginAdmin)
	root.GET("/product", han.Product.GetListProduct)

	// customer
	customer.GET("/user/profile", han.Customer.GetProfile)
	customer.GET("/user/order", han.Order.GetOrderByCust)
	customer.POST("/order", han.Order.CreateOrder)
	customer.GET("/order/complete", han.Order.CompletingOrder)

	// admin
	admin.GET("/admin/order", han.Order.GetOrderByAdmin)
	admin.GET("/admin/order/report", han.Order.GenerateReport)

	svc := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		Handler:      handler,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &instance{cfg, dep, &svc}
}

func (x *instance) Serve() error {
	x.Logger.Info(fmt.Sprintf("service started at port %d ...", x.Port))
	return x.handler.ListenAndServe()
}

func (x *instance) Shutdown(ctx context.Context) error {
	x.Logger.Info("service shutdown ...", x.Port)
	return x.handler.Shutdown(ctx)
}
