package pkg_error

import (
	"embed"
	"encoding/json"

	"github.com/labstack/echo/v4"
)

//go:embed error_properties.json
var error_properties embed.FS
var error_data map[string]ErrorData

type ErrorData struct {
	Body   any `json:"body"`
	Status int `json:"status"`
}

func constructErrorData() {
	fileName := "error_properties.json"
	error_properties, _ := error_properties.ReadFile(fileName)
	json.Unmarshal(error_properties, &error_data)
}

func CreateError(c echo.Context, pkg_error string) error {
	if error_data == nil {
		constructErrorData()
	}

	body := error_data[pkg_error].Body
	if body == nil {
		body = map[string]string{
			"message": pkg_error,
		}
	}

	return c.JSON(error_data[pkg_error].Status, body)
}
