package traefikcors

import (
	"context"
	"net/http"
	"time"
)

// Config the plugin configuration.
type Config struct {
	AllowAllOrigins bool

	// AllowMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (GET and POST)
	AllowMethods []string

	// AllowHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	AllowHeaders []string

	AllowOrigins []string
	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposeHeaders []string

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge time.Duration

	// Allows to add origins like http://some-domain/*, https://api.* or http://some.*.subdomain.com
	AllowWildcard bool

	// Allows usage of popular browser extensions schemas
	AllowBrowserExtensions bool

	// Allows usage of WebSocket protocol
	AllowWebSockets bool

	// Allows usage of file:// schema (dangerous!) use it only when you 100% sure it's needed
	AllowFiles bool
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
}

// Cors a plugin.
type Cors struct {
	Name             string
	Next             http.Handler
	Config           *Config
	allowOrigins     []string
	exposeHeaders    []string
	normalHeaders    http.Header
	preflightHeaders http.Header
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Cors{
		Name:             name,
		Next:             next,
		Config:           config,
		allowOrigins:     normalize(config.AllowOrigins),
		normalHeaders:    generateNormalHeaders(config),
		preflightHeaders: generatePreflightHeaders(config),
	}, nil
}

func (e *Cors) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	method := req.Method               //请求方法
	origin := req.Header.Get("Origin") //请求头部

	if len(origin) == 0 {
		// request is not a CORS request
		return
	}

	host := req.Host

	if origin == "http://"+host || origin == "https://"+host {
		return
	}

	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		e.handleNormal(rw, req)
		rw.WriteHeader(http.StatusNoContent)
		return
	} else {
		e.handlePreflight(rw, req)
		e.Next.ServeHTTP(rw, req)
	}
}

func (cors *Cors) handleNormal(rw http.ResponseWriter, req *http.Request) {
	header := rw.Header()
	for key, value := range cors.normalHeaders {
		header[key] = value
	}
}

func (cors *Cors) handlePreflight(rw http.ResponseWriter, req *http.Request) {
	header := rw.Header()
	for key, value := range cors.preflightHeaders {
		header[key] = value
	}
}

func (cors *Cors) validateOrigin(origin string) bool {
	if cors.Config.AllowAllOrigins {
		return true
	}
	for _, value := range cors.allowOrigins {
		if value == origin {
			return true
		}
	}
	// if len(cors.Config.WildcardOrigins) > 0 && cors.validateWildcardOrigin(origin) {
	// 	return true
	// }
	// if cors.Config.AllowOriginFunc != nil {
	// 	return cors.allowOriginFunc(origin)AllowWildcard
	// }
	return false
}

// func (cors *Cors) validateWildcardOrigin(origin string) bool {
// 	for _, w := range cors.wildcardOrigins {
// 		if w[0] == "*" && strings.HasSuffix(origin, w[1]) {
// 			return true
// 		}
// 		if w[1] == "*" && strings.HasPrefix(origin, w[0]) {
// 			return true
// 		}
// 		if strings.HasPrefix(origin, w[0]) && strings.HasSuffix(origin, w[1]) {
// 			return true
// 		}
// 	}

// 	return false
// }
