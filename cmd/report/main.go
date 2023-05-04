package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"go.cognotif/internal/repository/constant"
	"go.cognotif/internal/repository/postgresql"
	pkg_logger "go.cognotif/pkg/logger"
	pkg_util "go.cognotif/pkg/util"
)

// generate report order
func main() {
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

	order_data, err := psql.GenerateDataReport(ctx)
	if err != nil {
		log.Hashcode(ctx).Fatal("Error when generate report: ", err)
	}

	report_csv, err := os.Create("report-order " + pkg_util.Now() + ".csv")
	if err != nil {
		log.Hashcode(ctx).Fatal("error when create file: ", err)
	}
	defer report_csv.Close()

	w := csv.NewWriter(report_csv)
	defer w.Flush()

	var data [][]string
	row := []string{"order_id", "customer_name", "order_date", "order_status", "total_price"}
	data = append(data, row)
	for _, order := range order_data {
		row := []string{order.IdOrder, order.Name, order.Date, order.Status, fmt.Sprintf("%f", order.Sum)}
		data = append(data, row)
	}
	w.WriteAll(data)

}
