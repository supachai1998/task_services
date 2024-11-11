package interfaces

import (
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/supachai1998/task_services/docs"
	"github.com/supachai1998/task_services/internal/configs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewEchoInterface(config *configs.ServerConfig) *echo.Echo {
	e := echo.New()
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.CORS(),
		middleware.Secure(),
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: func() string {
				return fmt.Sprintf("%s-%d", configs.AppConfig.Server.AppName, time.Now().UnixNano())
			},
		}),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: func(c echo.Context) bool {
				return strings.Contains(c.Request().URL.Path, "swagger")
			},
		}),
	)

	// Initialize custom validator
	e.Validator = NewCustomValidator()

	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}
