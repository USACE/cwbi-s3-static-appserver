package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Based on Virtual Hosting Example at https://echo.labstack.com/cookbook/subdomains/
type (
	Config struct {
		Domain string `envconfig:"domain"`
	}

	Host struct {
		Echo *echo.Echo
	}
)

func main() {

	// Hosts
	hosts := map[string]*Host{}

	// Config
	var cfg Config
	if err := envconfig.Process("appserver", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	// HOME
	// ====
	home := echo.New()
	home.Use(middleware.Recover())
	home.Static("/", "/data/home")
	hosts[cfg.Domain] = &Host{home}

	// CUMULUS
	// =======
	cumulus := echo.New()
	cumulus.Use(middleware.Recover())
	cumulus.Static("/", "/data/cumulus")
	hosts[fmt.Sprintf("cumulus.%s", cfg.Domain)] = &Host{cumulus}

	// MIDAS
	// =====
	midas := echo.New()
	midas.Use(middleware.Recover())
	midas.Static("/", "/data/midas")
	hosts[fmt.Sprintf("midas.%s", cfg.Domain)] = &Host{midas}

	// WATER
	// =====
	water := echo.New()
	water.Use(middleware.Recover())
	water.Static("/", "/data/water")
	hosts[fmt.Sprintf("water.%s", cfg.Domain)] = &Host{water}

	// Server
	e := echo.New()
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]

		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}
		return
	})
	e.Logger.Fatal(e.Start(":80"))
}
