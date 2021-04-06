package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

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
			regexp.MustCompile("^\\/[a-zA-Z0-9\\/\\-]+$"): "/index.html",
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
	d.Use(middleware.Recover())
	d.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	// HOME
	// ====
	home := echo.New()
	home.Pre(rewriteMiddleware)
	home.Use(middleware.Recover())
	home.Static("/", "/data/home")
	hosts[fmt.Sprintf("%shome.%s", cfg.SubdomainPrefix, cfg.Domain)] = &Host{home}

	// CUMULUS
	// =======
	cumulus := echo.New()
	cumulus.Pre(rewriteMiddleware)
	cumulus.Use(middleware.Recover())
	cumulus.Static("/", "/data/cumulus")
	hosts[fmt.Sprintf("%scumulus.%s", cfg.SubdomainPrefix, cfg.Domain)] = &Host{cumulus}

	// MIDAS
	// =====
	midas := echo.New()
	midas.Pre(rewriteMiddleware)
	midas.Use(middleware.Recover())
	midas.Static("/", "/data/midas")
	hosts[fmt.Sprintf("%smidas.%s", cfg.SubdomainPrefix, cfg.Domain)] = &Host{midas}

	// WATER
	// =====
	water := echo.New()
	water.Pre(rewriteMiddleware)
	water.Use(middleware.Recover())
	water.Static("/", "/data/water")
	hosts[fmt.Sprintf("%swater.%s", cfg.SubdomainPrefix, cfg.Domain)] = &Host{water}

	// Server
	e := echo.New()
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
	e.Logger.Fatal(e.Start(":80"))
}