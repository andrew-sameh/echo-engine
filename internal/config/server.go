package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/andrew-sameh/echo-engine/pkg/slice"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerConfig struct {
	Host       string
	Port       string
	Env        string
	Validator  echo.Validator
	Binder     echo.Binder
	CORSConfig middleware.CORSConfig
}
type AppValidator struct {
	validate *validator.Validate
}

// Implement the bind method to verify the request's struct for parameter validation
type BinderWithValidation struct{}

func (a *AppValidator) Validate(i interface{}) error {
	return a.validate.Struct(i)

}
func GetEchoLogConfig(cfg *Config) middleware.LoggerConfig {
	echoLogConf := middleware.DefaultLoggerConfig
	echoLogConf.CustomTimeFormat = time.RFC3339
	// echoLogConf.Format = fmt.Sprintln(`{"level":"info","source":"echo","id":"${id}","mt":"${method}","uri":"${uri}","st":${status},"e":"${error}","lc":"${latency_human}","ts":"${time_custom}"}`)
	return echoLogConf
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Host:      os.Getenv("HOST"),
		Port:      os.Getenv("PORT"),
		Env:       os.Getenv("ENV"),
		Validator: ValidatorInit(),
		Binder:    &BinderWithValidation{},
		// Validator:  &AppValidator{validate: validator.New()},
		CORSConfig: middleware.DefaultCORSConfig,
	}
}

func ValidatorInit() echo.Validator {
	v := validator.New()

	v.RegisterValidation("json", func(fl validator.FieldLevel) bool {
		var js json.RawMessage
		return json.Unmarshal([]byte(fl.Field().String()), &js) == nil
	})

	v.RegisterValidation("in", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if slice.ContainsString(strings.Split(fl.Param(), ";"), value) || value == "" {
			return true
		}

		return false
	})

	return &AppValidator{validate: v}
}

func (BinderWithValidation) Bind(i interface{}, ctx echo.Context) error {
	binder := &echo.DefaultBinder{}

	if err := binder.Bind(i, ctx); err != nil {
		return errors.New(err.(*echo.HTTPError).Message.(string))
	}

	if err := ctx.Validate(i); err != nil {
		// Validate only provides verification function for struct.
		// When the requested data type is not struct,
		// the variable should be considered legal after the bind succeeds.
		if reflect.TypeOf(i).Kind() != reflect.Struct {
			return nil
		}

		var buf bytes.Buffer
		if ferrs, ok := err.(validator.ValidationErrors); ok {
			for _, ferr := range ferrs {
				buf.WriteString("Validation failed on ")
				buf.WriteString(ferr.Tag())
				buf.WriteString(" for ")
				buf.WriteString(ferr.StructField())
				buf.WriteString("\n")
			}

			return errors.New(buf.String())
		}

		return err
	}

	return nil
}
