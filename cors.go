package traefik_plugin_cors

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Config the plugin configuration.
type Config struct {
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

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge time.Duration
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
	Name   string
	Next   http.Handler
	Config *Config
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	fmt.Println(name)
	return &Cors{
		Name: name,
		Next: next,
	}, nil
}

func (e *Cors) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	method := req.Method               //请求方法
	origin := req.Header.Get("Origin") //请求头部
	var headerKeys []string            // 声明请求头keys
	for k := range req.Header {
		if !(k == "Access-Control-Request-Method" || k == "Access-Control-Request-Headers") {
			headerKeys = append(headerKeys, k)
		}
	}
	headerStr := strings.Join(headerKeys, ", ")
	if headerStr != "" {
		headerStr = fmt.Sprintf("Access-Control-Allow-Origin, Access-Control-Allow-Headers, %s", headerStr)
	} else {
		headerStr = "Access-Control-Allow-Origin, Access-Control-Allow-Headers"
	}
	headers := rw.Header()
	if origin == "" {
		headers.Set("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
		headers.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
		headers.Set("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
		headers.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
		headers.Set("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
		headers.Set("Access-Control-Allow-Credentials", "false")                                                                                                                                                  // 设置返回格式是json
	} else {
		headers.Set("Access-Control-Allow-Origin", origin)
		headers.Set("Access-Control-Allow-Headers", headerStr)
		headers.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		headers.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
		headers.Set("Access-Control-Allow-Credentials", "true")
	}

	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	e.Next.ServeHTTP(rw, req)
}
