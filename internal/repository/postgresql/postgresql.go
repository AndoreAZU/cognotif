package postgresql

import (
	"context"
	"database/sql"
	"embed"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	pkg_logger "go.cognotif/pkg/logger"
)

type Postgresql interface {
	TableCustomer
	TableProduct
	TableOrder
	TableProductOrder
	TableAdmin
}

type Dependency struct {
	*pkg_logger.Logger
}

type Configuration struct {
	DSN_POSTGRES string
	DSN_MYSQL    string
}

type instance struct {
	Dependency
	Configuration

	mysql *sql.DB
	*sql.DB
}

//go:embed query/*.sql
var files_sql embed.FS

func New(cfg Configuration, dep Dependency) Postgresql {
	dep.Logger.Info("starting connect to db...")
	conn, err := sql.Open("postgres", cfg.DSN_POSTGRES)
	if err != nil {
		dep.Logger.Fatal("failed connect to db: ", err)
	}

	conn2, err := sql.Open("mysql", cfg.DSN_MYSQL)
	if err != nil {
		dep.Logger.Fatal("failed connect to db: ", err)
	}

	dep.Logger.Info("connected to db...")
	return &instance{Configuration: cfg, Dependency: dep, DB: conn, mysql: conn2}
}

func (x *instance) ExecContext(ctx context.Context, query string, args ...any) (tx *sql.Tx, err error) {
	tx, err = x.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return
	}

	_, err = tx.ExecContext(ctx, query, args...)
	return
}

func CommitTxSql(list_tx []*sql.Tx) {
	for _, tx := range list_tx {
		if tx != nil {
			tx.Commit()
		}
	}
}

func RollbackTxSql(list_tx []*sql.Tx, err error) error {
	for _, tx := range list_tx {
		if tx != nil {
			tx.Rollback()
		}
	}

	return err
}

func (x *instance) GetSQL(sql_name string) string {
	sql_properties, _ := files_sql.ReadFile("query/" + sql_name + ".sql")
	return string(sql_properties)
}
