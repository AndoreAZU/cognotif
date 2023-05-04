package validator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator *validator.Validate
}

func (x *Validator) Validate(i interface{}) error {
	return x.Validator.Struct(i)
}

func New() *Validator {
	v := &Validator{
		Validator: validator.New(),
	}

	v.Validator.RegisterValidation("alphawithspace", v.ValidateStringWithSpace)
	v.Validator.RegisterValidation("yyyy-mm-dd", v.ValidateFieldBirthDate)
	v.Validator.RegisterValidation("emailformat", v.ValidateEmail)
	v.Validator.RegisterValidation("place", v.ValidatePlace)
	v.Validator.RegisterValidation("phonenumber", v.ValidatePhoneNumber)

	// register tag json name
	v.Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v
}

func (x *Validator) ValidateStringWithSpace(fl validator.FieldLevel) bool {
	return regexp.MustCompile("^[a-zA-Z .']*$").MatchString(fl.Field().String())
}

func (x *Validator) ValidateFieldBirthDate(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$`).MatchString(fl.Field().String())
}

func (x *Validator) ValidateEmail(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,20})+$`).MatchString(fl.Field().String())
}

func (x *Validator) ValidatePlace(fl validator.FieldLevel) bool {
	return regexp.MustCompile("^[a-zA-Z&.' ]+$").MatchString(fl.Field().String())
}

func (x *Validator) ValidatePhoneNumber(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^628\d{8,11}$`).MatchString(fl.Field().String())
}
