package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/kelseyhightower/envconfig"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Based on Virtual Hosting Example at https://echo.labstack.com/cookbook/subdomains/
type (
	Config struct {
		Domain          string `envconfig:"domain"`
		SubdomainPrefix string `envconfig:"subdomain_prefix"`
	}

	Host struct {
		Echo *echo.Echo
	}
)

func main() {

	// Rewrite Middleware for Client-Side Routing
	rewriteMiddleware := middleware.RewriteWithConfig(middleware.RewriteConfig{
		RegexRules: map[*regexp.Regexp]string{
			// Real files on-disk with file extension (.js, .css, etc.) paths unmodified
			regexp.MustCompile("^(\\/.+\\.(css|html|ico|js|json|map|png|txt)$)"): "/$1",
			// All other paths return index.html (Client-side Routing)
			regexp.MustCompile("^\\/[a-zA-Z0-9\\/\\-\\_]+$"): "/index.html",
		},
	},
	)

	// Hosts
	hosts := map[string]*Host{}

	// Config
	var cfg Config
	if err := envconfig.Process("appserver", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	// DEFAULT
	// =======
	// Routes for appserver itself. Implemented to support health checks, etc.
	// Note: Router is not included in hosts map
	d := echo.New()
	d.Use(middleware.Recover(), middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	d.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	// HOME
	// ====
	home := echo.New()
	home.Pre(rewriteMiddleware)
	home.Use(middleware.Recover(), middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	home.Static("/", "/data/home")
	hosts[fmt.Sprintf("%shome.%s", cfg.SubdomainPrefix, cfg.Domain)] = &Host{home}

	// COLUMBIA RIVER HYDROLOGY
	// ========================
	crb := echo.New()
	crb.Pre(rewriteMiddleware)
	crb.Use(middleware.Recover(), middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	crb.Static("/", "/data/crb-hydrology")
	hosts[fmt.Sprintf("%scrb-hydrology.%s", cfg.SubdomainPrefix, cfg.Domain)] = &Host{crb}

	// CUMULUS
	// =======
	cumulus := echo.New()
	cumulus.Pre(rewriteMiddleware)
	cumulus.Use(middleware.Recover(), middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	cumulus.Static("/", "/data/cumulus")
	hosts[fmt.Sprintf("%scumulus.%s", cfg.SubdomainPrefix, cfg.Domain)] = &Host{cumulus}

	// (TEMPORARY) CUMULUS (TAILWINDUI VERSION)
	cumulus2 := echo.New()
	cumulus2.Pre(rewriteMiddleware)
	cumulus2.Use(middleware.Recover(), middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	cumulus2.Static("/", "/data/cumulus-tailwindui")
	hosts[fmt.Sprintf("%scumulus-tailwindui.%s", cfg.SubdomainPrefix, cfg.Domain)] = &Host{cumulus2}

	// MIDAS
	// =====
	midas := echo.New()
	midas.Pre(rewriteMiddleware)
	midas.Use(middleware.Recover(), middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	midas.Static("/", "/data/midas")
	hosts[fmt.Sprintf("%smidas.%s", cfg.SubdomainPrefix, cfg.Domain)] = &Host{midas}

	// WATER
	// =====
	water := echo.New()
	water.Pre(rewriteMiddleware)
	water.Use(middleware.Recover(), middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	water.Static("/", "/data/water")
	hosts[fmt.Sprintf("%swater.%s", cfg.SubdomainPrefix, cfg.Domain)] = &Host{water}

	// WORKFORCE
	// =========
	workforce := echo.New()
	workforce.Pre(rewriteMiddleware)
	workforce.Use(middleware.Recover(), middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	workforce.Static("/", "/data/workforce")
	hosts[fmt.Sprintf("%sworkforce.%s", cfg.SubdomainPrefix, cfg.Domain)] = &Host{workforce}

	// Pallid Sturgeon
	// ===============
	sturgeon := echo.New()
	sturgeon.Pre(rewriteMiddleware)
	sturgeon.Use(middleware.Recover(), middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	sturgeon.Static("/", "/data/pallid-sturgeon")
	hosts[fmt.Sprintf("%spallid-sturgeon.%s", cfg.SubdomainPrefix, cfg.Domain)] = &Host{sturgeon}

	// Server
	e := echo.New()
	e.Server.Addr = ":8080"
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]

		if host == nil {
			// Default router; Most commonly used for health checks
			// AWS Target Group sends health checks with host header set to the IP of Load Balancer
			// This IP is not in the host map, so request is served with default router "d".
			// See: https://docs.aws.amazon.com/elasticloadbalancing/latest/network/target-group-health-checks.html
			d.ServeHTTP(res, req)
		} else {
			host.Echo.ServeHTTP(res, req)
		}
		return
	})
	e.Logger.Fatal(gracehttp.Serve(e.Server))
}
