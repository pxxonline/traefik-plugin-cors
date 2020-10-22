package traefik_plugin_cors

import (
	"context"
	"net/http"
	"net/http/httptest"

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
	Cors *cors.Cors
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	options := cors.Options{
		AllowedOrigins:     config.AllowedOrigins,
		AllowedMethods:     config.AllowedMethods,
		AllowedHeaders:     config.AllowedHeaders,
		ExposedHeaders:     config.ExposedHeaders,
		AllowCredentials:   config.AllowCredentials,
		OptionsPassthrough: config.OptionsPassthrough,
		MaxAge:             config.MaxAge,
		Debug:              config.Debug,
	}

	return &CorsTraefik{
		name: name,
		next: next,
		Cors: cors.New(options),
	}, nil
}

func (e *CorsTraefik) replacer() http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		recorder := httptest.NewRecorder()

		e.next.ServeHTTP(recorder, req)

		rw.WriteHeader(recorder.Code)
		_, _ = rw.Write(recorder.Body.Bytes())
		for name, values := range recorder.Header() {
			rw.Header().Set(name, values[0])
		}
	})
}

func (e *CorsTraefik) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	e.Cors.ServeHTTP(rw, req, e.replacer())
}
