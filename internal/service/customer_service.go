package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.cognotif/internal/dto/request"
	"go.cognotif/internal/dto/response"
	"go.cognotif/internal/repository/postgresql"
	"go.cognotif/internal/repository/token"
	pkg_error "go.cognotif/pkg/error"
	pkg_logger "go.cognotif/pkg/logger"
	pkg_util "go.cognotif/pkg/util"
)

type CustomerService interface {
	GetCustomerByEmail(ctx context.Context, email string) (response.GetCustomer, error)
	RegisterNewCustomer(ctx context.Context, request request.RegisterCustomer) (*sql.Tx, error)
	CustomerLogin(ctx context.Context, email, password string) (response.Login, error)
	AdminLogin(ctx context.Context, email, password string) (response.Login, error)
	GetCustomerById(ctx context.Context, id string) (response.GetCustomer, error)
}

type customerService struct {
	*pkg_logger.Logger
	postgresql.Postgresql
}

func NewCustomerService(log *pkg_logger.Logger, psql postgresql.Postgresql) CustomerService {
	return &customerService{Logger: log, Postgresql: psql}
}

func (cs *customerService) GetCustomerByEmail(ctx context.Context, email string) (response.GetCustomer, error) {
	customer, err := cs.Postgresql.GetCustomerByEmail(ctx, email)
	if err != nil {
		cs.Hashcode(ctx).Error("error when get customer: ", err)
		if err == sql.ErrNoRows {
			return response.GetCustomer{}, fmt.Errorf(pkg_error.CUSTOMER_NOT_EXISTS)
		}
		return response.GetCustomer{}, fmt.Errorf(pkg_error.GENERAL_ERROR)
	}

	return response.GetCustomer{
		Id:       customer.Id,
		Name:     customer.Name,
		Email:    customer.Email,
		Password: customer.Password,
	}, nil
}

func (cs *customerService) GetCustomerById(ctx context.Context, id string) (response.GetCustomer, error) {
	customer, err := cs.Postgresql.GetCustomerById(ctx, id)
	if err != nil {
		cs.Hashcode(ctx).Error("error when get customer: ", err)
		if err == sql.ErrNoRows {
			return response.GetCustomer{}, fmt.Errorf(pkg_error.CUSTOMER_NOT_EXISTS)
		}
		return response.GetCustomer{}, fmt.Errorf(pkg_error.GENERAL_ERROR)
	}

	return response.GetCustomer{
		Id:       customer.Id,
		Name:     customer.Name,
		Email:    customer.Email,
		Password: customer.Password,
	}, nil
}

func (cs *customerService) RegisterNewCustomer(ctx context.Context, request request.RegisterCustomer) (*sql.Tx, error) {
	encrypt_pass, err := pkg_util.HashPassword(request.Password)
	if err != nil {
		cs.Hashcode(ctx).Error("error when hashing password: ", err)
		return nil, fmt.Errorf(pkg_error.GENERAL_ERROR)
	}

	customer, err := cs.GetCustomerByEmail(ctx, request.Email)
	if err != nil && strings.Compare(err.Error(), pkg_error.CUSTOMER_NOT_EXISTS) != 0 {
		cs.Hashcode(ctx).Error("error when check customer: ", err)
		return nil, err
	}

	if len(customer.Id) != 0 {
		cs.Hashcode(ctx).Info("customer with email ", request.Email, " already exists")
		return nil, fmt.Errorf(pkg_error.CUSTOMER_ALREADY_EXISTS)
	}

	tx, err := cs.Postgresql.InsertNewCustomer(ctx, uuid.NewString(), request.Name, request.Email, encrypt_pass)
	if err != nil {
		cs.Hashcode(ctx).Error("error when insert new customer: ", err)
		return nil, postgresql.RollbackTxSql([]*sql.Tx{tx}, fmt.Errorf(pkg_error.GENERAL_ERROR))
	}

	return tx, nil
}

func (cs *customerService) CustomerLogin(ctx context.Context, email, password string) (response.Login, error) {
	t := token.Token

	customer, err := cs.GetCustomerByEmail(ctx, email)
	if err != nil {
		cs.Hashcode(ctx).Error("error when login: ", err)
		if strings.Compare(err.Error(), pkg_error.CUSTOMER_NOT_EXISTS) == 0 {
			return response.Login{}, fmt.Errorf(pkg_error.INVALID_CREDENTIAL)
		}
		return response.Login{}, err
	}

	if !pkg_util.CheckPasswordHash(password, customer.Password) {
		cs.Hashcode(ctx).Error("password not match")
		return response.Login{}, fmt.Errorf(pkg_error.INVALID_CREDENTIAL)
	}

	expired := time.Now().Add(time.Duration(pkg_util.Atoi(os.Getenv("TOKEN_DURATION"))) * time.Second)
	token, _ := t.Create(
		t.WithKeypair(pkg_util.EkstractKeypairCognotif()),
		t.WithDuration(time.Duration(pkg_util.Atoi(os.Getenv("TOKEN_DURATION")))*time.Second),
		t.WithIDCustomer(customer.Id),
		t.WithIsAdmin(false),
	)

	return response.Login{
		Token:     token,
		ExpiredAt: expired,
	}, nil
}

func (cs *customerService) AdminLogin(ctx context.Context, email, password string) (response.Login, error) {
	t := token.Token

	admin, err := cs.GetAdminByEmail(ctx, email)
	if err != nil {
		cs.Hashcode(ctx).Error("error when login: ", err)
		if strings.Compare(err.Error(), pkg_error.CUSTOMER_NOT_EXISTS) == 0 {
			return response.Login{}, fmt.Errorf(pkg_error.INVALID_CREDENTIAL)
		}
		return response.Login{}, err
	}

	if !pkg_util.CheckPasswordHash(password, admin.Password) {
		cs.Hashcode(ctx).Error("password not match")
		return response.Login{}, fmt.Errorf(pkg_error.INVALID_CREDENTIAL)
	}

	expired := time.Now().Add(time.Duration(pkg_util.Atoi(os.Getenv("TOKEN_DURATION"))) * time.Second)
	token, _ := t.Create(
		t.WithKeypair(pkg_util.EkstractKeypairCognotif()),
		t.WithDuration(time.Duration(pkg_util.Atoi(os.Getenv("TOKEN_DURATION")))*time.Second),
		t.WithIDCustomer(admin.Id),
		t.WithIsAdmin(true),
	)

	return response.Login{
		Token:     token,
		ExpiredAt: expired,
	}, nil
}
