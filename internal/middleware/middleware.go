package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.cognotif/internal/repository/constant"
	"go.cognotif/internal/repository/token"
	pkg_error "go.cognotif/pkg/error"
	pkg_logger "go.cognotif/pkg/logger"
	pkg_util "go.cognotif/pkg/util"
)

type MiddlewareHTTP interface {
	Middleware(next http.Handler) http.Handler
	MiddlewareCustomer(next echo.HandlerFunc) echo.HandlerFunc
	MiddlewareAdmin(next echo.HandlerFunc) echo.HandlerFunc
}

type middlewareHTTP struct {
	*pkg_logger.Logger
}

func New(log *pkg_logger.Logger) MiddlewareHTTP {
	return &middlewareHTTP{Logger: log}
}

func (h *middlewareHTTP) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hc := generateHashcode()
		r = r.WithContext(context.WithValue(r.Context(), constant.Hashcode{}, hc))

		w.Header().Set("x-hashcode", hc)
		next.ServeHTTP(w, r)

	})
}

func (h *middlewareHTTP) MiddlewareCustomer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if len(c.Request().Header.Get("Authorization")) == 0 {
			return pkg_error.CreateError(c, pkg_error.MISSING_HEADER_AUTHORIZATION)
		}

		t := token.Token
		duration := time.Duration(pkg_util.Atoi(os.Getenv("TOKEN_DURATION"))) * time.Second

		// validate token
		token_data, err := t.VerifyRequest(c.Request(),
			t.WithKeypair(pkg_util.EkstractKeypairCognotif()),
			t.WithDuration(duration))

		if err != nil {
			return pkg_error.CreateError(c, err.Error())
		}

		if token_data.IsAdmin {
			return pkg_error.CreateError(c, pkg_error.ADMIN_FORBIDDEN)
		}

		c.Set(constant.CTX_CUST_ID, token_data.IDCustomer)

		return next(c)
	}
}

func (h *middlewareHTTP) MiddlewareAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if len(c.Request().Header.Get("Authorization")) == 0 {
			return pkg_error.CreateError(c, pkg_error.MISSING_HEADER_AUTHORIZATION)
		}

		t := token.Token
		duration := time.Duration(pkg_util.Atoi(os.Getenv("TOKEN_DURATION"))) * time.Second

		// validate token
		token_data, err := t.VerifyRequest(c.Request(),
			t.WithKeypair(pkg_util.EkstractKeypairCognotif()),
			t.WithDuration(duration))

		if err != nil {
			return pkg_error.CreateError(c, err.Error())
		}

		if !token_data.IsAdmin {
			return pkg_error.CreateError(c, pkg_error.NEED_ADMIN_ACCESS)
		}

		c.Set(constant.CTX_CUST_ID, token_data.IDCustomer)

		return next(c)
	}
}

func generateHashcode() string {
	return "(" + strings.Split(uuid.New().String(), "-")[0] + ")"
}

func GetHashcode(ctx context.Context) string {
	return ctx.Value(constant.Hashcode{}).(string)
}

func CopyContext(ctx context.Context) (context.Context, context.CancelFunc) {
	c := context.Background()
	c = context.WithValue(c, constant.Hashcode{}, ctx.Value(constant.Hashcode{}))
	return context.WithCancel(c)
}
