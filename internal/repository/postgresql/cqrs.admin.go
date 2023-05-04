package postgresql

import (
	"context"
)

type TableAdmin interface {
	GetAdminByEmail(ctx context.Context, email string) (Customer, error)
	GetAdminById(ctx context.Context, id string) (Customer, error)
}

type Admin struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (x *instance) GetAdminByEmail(ctx context.Context, email string) (Customer, error) {
	stmt, err := x.PrepareContext(ctx, x.GetSQL("admin.select_by_email"))
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

func (x *instance) GetAdminById(ctx context.Context, id string) (Customer, error) {
	stmt, err := x.PrepareContext(ctx, x.GetSQL("admin.select_by_id"))
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
