package api

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stormsync/database"
)

type RouterConfig struct {
	ROKey  string
	RWKey  string
	DB     *database.Queries
	Logger *slog.Logger
}
type ServerAndDB struct {
	Web    *echo.Echo
	DB     *database.Queries
	Logger *slog.Logger
}

// NewRouter will setup the router and endpoints and
// returns a db connection and endpoint in a ServerAndDB stuct
// that provides DB access to the handlers.
func NewRouter(config RouterConfig) ServerAndDB {
	s := ServerAndDB{
		Web:    nil,
		DB:     config.DB,
		Logger: config.Logger,
	}
	e := echo.New()

	e.GET("/api/v1/report/all", s.GetAllReports)
	e.GET("/api/v1/report/hail", s.GetHailReports)
	e.GET("/api/v1/report/tornado", s.GetTornadoReports)
	e.GET("/api/v1/report/wind", s.GetWindReports)
	e.POST("/api/v1/maint/report", s.AddReport)

	e.Use(middleware.Secure())
	e.Use(middleware.Recover())

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-Api-Key",
		Validator: func(key string, c echo.Context) (bool, error) {
			matchKey := config.ROKey
			if c.Path() == "/api/v1/maint/report" {
				matchKey = config.RWKey
			}
			return key == matchKey, nil
		},
	}))
	s.Web = e
	return s
}
