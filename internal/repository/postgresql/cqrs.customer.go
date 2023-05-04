package postgresql

import (
	"context"
	"database/sql"
)

type TableCustomer interface {
	InsertNewCustomer(ctx context.Context, id, name, email, password string) (*sql.Tx, error)
	GetCustomerByEmail(ctx context.Context, email string) (Customer, error)
	GetCustomerById(ctx context.Context, id string) (Customer, error)
	GetAllCustomer(ctx context.Context) ([]Customer, error)
}

type Customer struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (x *instance) InsertNewCustomer(ctx context.Context, id, name, email, password string) (*sql.Tx, error) {
	return x.ExecContext(ctx, x.GetSQL("customer.insert"), id, name, email, password)
}

func (x *instance) GetCustomerByEmail(ctx context.Context, email string) (Customer, error) {
	stmt, err := x.PrepareContext(ctx, x.GetSQL("customer.select_by_email"))
	if err != nil {
		return Customer{}, err
	}

	defer stmt.Close()

	var c Customer
	err = stmt.QueryRow(email).Scan(&c.Id, &c.Name, &c.Email, &c.Password)
	if err != nil {
		return Customer{}, err
	}

	return c, nil
}

func (x *instance) GetCustomerById(ctx context.Context, id string) (Customer, error) {
	stmt, err := x.PrepareContext(ctx, x.GetSQL("customer.select_by_id"))
	if err != nil {
		return Customer{}, err
	}

	defer stmt.Close()

	var c Customer
	err = stmt.QueryRow(id).Scan(&c.Id, &c.Name, &c.Email, &c.Password)
	if err != nil {
		return Customer{}, err
	}

	return c, nil
}

func (x *instance) GetAllCustomer(ctx context.Context) ([]Customer, error) {
	stmt, err := x.QueryContext(ctx, x.GetSQL("customer.select_all"))
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var customer []Customer
	for stmt.Next() {
		var c Customer
		if err = stmt.Scan(&c.Id, &c.Name, &c.Email); err != nil {
			return nil, err
		}
		customer = append(customer, c)
	}

	return customer, nil
}
