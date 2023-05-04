package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"go.cognotif/internal/dto/request"
	"go.cognotif/internal/dto/response"
	"go.cognotif/internal/repository/constant"
	"go.cognotif/internal/repository/postgresql"
	"go.cognotif/internal/service"
	pkg_error "go.cognotif/pkg/error"
	pkg_logger "go.cognotif/pkg/logger"
)

type EmailData struct {
	Name   string        `json:"customer_name"`
	Orders []OrderDetail `json:"order_detail"`
}

type OrderDetail struct {
	Orders              []response.Order `json:"orders"`
	LinkCompletingOrder string           `json:"link_completing_order"`
}

func __sending_email(ctx context.Context, name, email string, order response.GetOrder) {
	var detail []OrderDetail
	for _, o := range order.Orders {
		detail = append(detail, OrderDetail{
			Orders:              order.Orders,
			LinkCompletingOrder: "http://localhost:8082/v1/order/complete?id=" + o.Id,
		})
	}

	data := EmailData{
		Name:   name,
		Orders: detail,
	}

	json, _ := json.Marshal(data)

	fmt.Println("sending email to: ", email)
	fmt.Println("order detail: ", string(json))
}

func _check_pending_order() {
	ctx, cancel := context.WithCancel(context.Background())
	hashcode := strings.Split(uuid.NewString(), "-")[0]
	ctx = context.WithValue(ctx, constant.Hashcode{}, hashcode)
	defer cancel()

	log := pkg_logger.NewLogger()
	psql := postgresql.New(postgresql.Configuration{
		DSN_POSTGRES: "postgresql://cognotif:c0gn0t1f@localhost:5431/ecommerce?sslmode=disable",
		DSN_MYSQL:    "cognotif:c0gn0t1f@tcp(localhost:3306)/ecommerce",
	}, postgresql.Dependency{
		Logger: log,
	})

	product_order_service := service.NewProductOrderService(log, psql)
	product_service := service.NewProductService(log, psql)
	order_service := service.NewOrderService(log, psql, product_order_service, product_service)

	// get all customer
	customer, err := psql.GetAllCustomer(ctx)
	if err != nil {
		log.Hashcode(ctx).Fatal("error when get customer")
	}

	for _, c := range customer {
		order, err := order_service.GetOrder(ctx, request.GetOrder{IdCust: c.Id, Status: constant.ORDER_STATUS_PENDING})
		if err != nil && err.Error() != pkg_error.ORDER_NOT_EXIST {
			log.Hashcode(ctx).Fatal("error when get order: ", err)
		}

		__sending_email(ctx, c.Name, c.Email, order)
	}
}

func main() {
	// Worker Scheduler for inserting data to backend via api
	var cron = gocron.NewScheduler(time.Local)
	cron.Every(10).Seconds().Tag("insert_request").Do(_check_pending_order)
	cron.StartBlocking()
}
