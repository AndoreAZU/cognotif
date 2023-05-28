package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.cognotif/internal/handler/customer"
	"go.cognotif/internal/handler/order"
	"go.cognotif/internal/handler/product"
	"go.cognotif/internal/middleware"
	"go.cognotif/internal/repository/postgresql"
	"go.cognotif/internal/repository/validator"
	"go.cognotif/internal/service"
	pkg_logger "go.cognotif/pkg/logger"
	pkg_util "go.cognotif/pkg/util"
	"go.cognotif/server"
)

//go:embed files/*
var env embed.FS

// LoadEnv doing process split string from file .env
// and extract each key and value to os environment
func LoadEnv(env string) {
	s := strings.Split(env, "\n")

	for _, v := range s {
		if len(v) == 0 || !strings.Contains(v, "=") {
			continue
		}

		vS := strings.SplitN(v, "=", 2)
		os.Setenv(vS[0], vS[1])
	}
}

func Initialize() {
	profile := os.Getenv("APP_ENV")
	logrus.Info("start load config profile ", profile)
	key_value, err := env.ReadFile("files/" + profile + ".env")
	if err != nil {
		logrus.Error("failed init env from file : ", err.Error())
	}
	LoadEnv(string(key_value))
}

type Depedency struct {
	Server struct {
		Server server.Depedency
	}
}

type Configuration struct {
	Server struct {
		Server server.Configuration
	}
}

type Service struct {
	Handler struct {
		Customer customer.Customer
		Product  product.Product
		Order    order.Order

		Server struct {
			ServerHandler server.Handler
		}
	}
}

type _service struct {
	CustomerService     service.CustomerService
	ProductService      service.ProductService
	ProductOrderSerivce service.ProductOrderService
	OrderService        service.OrderService
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Initialize()

	// initialize pkg
	log := pkg_logger.NewLogger()
	psql := postgresql.New(postgresql.Configuration{
		DSN_POSTGRES: os.Getenv("DSN_COGNOTIF_PSQL"),
		DSN_MYSQL:    os.Getenv("DSN_COGNOTIF_MYSQL"),
	}, postgresql.Dependency{
		Logger: log,
	})
	validator := validator.New()

	c := new(Configuration)
	d := new(Depedency)
	s := new(Service)
	_s := new(_service)

	{
		// initialize service
		_s.CustomerService = service.NewCustomerService(log, psql)
		_s.ProductService = service.NewProductService(log, psql)
		_s.ProductOrderSerivce = service.NewProductOrderService(log, psql)
		_s.OrderService = service.NewOrderService(log, psql, _s.ProductOrderSerivce, _s.ProductService)
	}

	middleware := middleware.New(log)

	{
		c.Server.Server = server.Configuration{
			Port:         pkg_util.Atoi(os.Getenv("PORT")),
			ReadTimeout:  time.Duration(pkg_util.Atoi(os.Getenv("READ_TIMEOUT"))) * time.Second,
			WriteTimeout: time.Duration(pkg_util.Atoi(os.Getenv("WRITE_TIMEOUT"))) * time.Second,
		}
	}

	{
		// initialize handler
		s.Handler.Customer = customer.NewCustomer(customer.Service{
			CustomerService: _s.CustomerService,
		})

		s.Handler.Product = product.NewProduct(product.Service{
			ProductService: _s.ProductService,
		})

		s.Handler.Order = order.NewOrder(order.Service{
			OrderService: _s.OrderService,
		})
	}

	{
		d.Server.Server = server.Depedency{
			Logger:         log,
			Validator:      validator,
			MiddlewareHTTP: middleware,
		}
	}

	{
		// initialize server handler
		s.Handler.Server.ServerHandler = server.Handler{
			Customer: s.Handler.Customer,
			Product:  s.Handler.Product,
			Order:    s.Handler.Order,
		}
	}

	svc := server.New(c.Server.Server, d.Server.Server, s.Handler.Server.ServerHandler)

	if err := svc.Serve(); err != nil {
		log.Fatal(fmt.Sprintf("failed start service: %s", err))
	}

	if err := svc.Shutdown(ctx); err != nil {
		log.Fatal(fmt.Sprintf("failed shutdown service: %s", err))
	}
}
