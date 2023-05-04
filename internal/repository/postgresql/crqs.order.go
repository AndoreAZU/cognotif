package postgresql

import (
	"context"
	"database/sql"
)

type TableOrder interface {
	InsertNewOrder(ctx context.Context, id, id_cust, date, status string) (*sql.Tx, error)
	UpdateOrderStatus(ctx context.Context, id_order, status string) (*sql.Tx, error)
	GetOrderByQuery(ctx context.Context, id, id_cust, status string) ([]Order, error)
	GenerateDataReport(ctx context.Context) ([]OrderReport, error)
}

type OrderReport struct {
	IdOrder string  `json:"id_order"`
	Name    string  `json:"customer_name"`
	Date    string  `json:"order_date"`
	Sum     float64 `json:"total_price"`
	Status  string  `json:"status_order"`
}

type Order struct {
	Id     string `json:"id"`
	IdCust string `json:"id_customer"`
	Date   string `json:"date"`
	Status string `json:"status"`
}

func (x *instance) InsertNewOrder(ctx context.Context, id, id_cust, date, status string) (*sql.Tx, error) {
	return x.ExecContext(ctx, x.GetSQL("order.insert"), id, id_cust, date, status)
}

func (x *instance) UpdateOrderStatus(ctx context.Context, id_order, status string) (*sql.Tx, error) {
	return x.ExecContext(ctx, x.GetSQL("order.update_status"), status, id_order)
}

func (x *instance) GetOrderByQuery(ctx context.Context, id, id_cust, status string) ([]Order, error) {
	stmt, err := x.QueryContext(ctx, x.GetSQL("order.select_by_field"), id, id_cust, status)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var orders []Order
	for stmt.Next() {
		var order Order
		if err = stmt.Scan(&order.Id, &order.IdCust, &order.Date, &order.Status); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return nil, sql.ErrNoRows
	}

	return orders, nil
}

func (x *instance) GenerateDataReport(ctx context.Context) ([]OrderReport, error) {
	stmt, err := x.QueryContext(ctx, x.GetSQL("order.generate_report"))
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var orders []OrderReport
	for stmt.Next() {
		var order OrderReport
		if err = stmt.Scan(&order.IdOrder, &order.Name, &order.Date, &order.Sum, &order.Status); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return nil, sql.ErrNoRows
	}

	return orders, nil
}
