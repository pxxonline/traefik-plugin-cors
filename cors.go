package traefik_plugin_cors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pxxonline/traefik-plugin-cors/cors"
)

// Config the plugin configuration.
type Config struct {
	AllowedOrigins     []string `json:"allowedOrigins,omitempty"`
	AllowedMethods     []string `json:"allowedMethods,omitempty"`
	AllowedHeaders     []string `json:"allowedHeaders,omitempty"`
	ExposedHeaders     []string `json:"exposedHeaders,omitempty"`
	AllowCredentials   bool     `json:"allowCredentials,omitempty"`
	OptionsPassthrough bool     `json:"optionsPassthrough,omitempty"`
	MaxAge             int      `json:"maxAge,omitempty"`
	Debug              bool     `json:"debug,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// CorsTraefik a plugin.
type CorsTraefik struct {
	next http.Handler
	name string
	c    *cors.Cors
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// fmt.Println(name)
	// fmt.Println(toJson(config))

	options := cors.Options{
		AllowedOrigins:     config.AllowedOrigins,
		AllowedMethods:     config.AllowedMethods,
		AllowedHeaders:     config.AllowedHeaders,
		ExposedHeaders:     config.ExposedHeaders,
		AllowCredentials:   config.AllowCredentials,
		OptionsPassthrough: config.OptionsPassthrough,
		MaxAge:             config.MaxAge,
		// Debug:              config.Debug,
		// // AllowOriginFunc:        allowOriginFunc,
		// // AllowOriginRequestFunc: allowOriginRequestFunc,
	}

	fmt.Println(toJson(options.AllowedHeaders))
	c := cors.New(options)

	fmt.Println(c)

	return &CorsTraefik{
		name: name,
		next: next,
		c:    c,
	}, nil
}

func (e *CorsTraefik) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// e.c.Log.Printf(req.URL.String())
	e.c.Handler(e.next).ServeHTTP(rw, req)
}

// toJson toJson
func toJson(v interface{}) string {
	if data, err := json.Marshal(v); err == nil {
		return string(data)
	}
	return ""
}

// func allowOriginFunc(origin string) bool {
// 	return true
// }

// // Optional origin validator (with request) function
// func allowOriginRequestFunc(r *http.Request, origin string) bool {
// 	return true
// }
