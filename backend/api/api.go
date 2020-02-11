package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/ajorgensen/go-react-webapp-bootstrap/backend/app"
	"github.com/ajorgensen/go-react-webapp-bootstrap/backend/models"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type API struct {
	App    *app.App
	Config *Config
}

type statusCodeRecorder struct {
	http.ResponseWriter
	http.Hijacker
	StatusCode int
}

func (r *statusCodeRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func New(a *app.App) (api *API, err error) {
	api = &API{App: a}
	api.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}

	return api, nil
}

func test(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte(`{ "message": "bar" }`))
	return nil
}

func (a *API) initFrontend(router *chi.Mux) {
	root := "./frontend/build"
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", a.handler(func(w http.ResponseWriter, r *http.Request) error {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.Error(w, http.StatusText(404), 404)
		} else {
			fs.ServeHTTP(w, r)
		}

		return nil
	}))
}

func (a *API) initBackend(router *chi.Mux) {
	router.Route("/api", func(r chi.Router) {
		r.Get("/bar", a.apiHandler(test))
	})
}

func (a *API) Init(router *chi.Mux) {
	a.initFrontend(router)
	a.initBackend(router)
}

func (a *API) apiHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return a.handler(func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")
		return f(w, r)
	})
}

func (a *API) handler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)

		beginTime := time.Now()

		hijacker, _ := w.(http.Hijacker)
		w = &statusCodeRecorder{
			ResponseWriter: w,
			Hijacker:       hijacker,
		}

		ctx := a.App.NewContext().WithRemoteAddress(a.IPAddressForRequest(r))
		ctx = ctx.WithLogger(ctx.Logger.WithField("request_id", base64.RawURLEncoding.EncodeToString(models.NewId())))

		defer func() {
			statusCode := w.(*statusCodeRecorder).StatusCode
			if statusCode == 0 {
				statusCode = 200
			}
			duration := time.Since(beginTime)

			logger := ctx.Logger.WithFields(logrus.Fields{
				"duration":    duration,
				"status_code": statusCode,
				"remote":      ctx.RemoteAddress,
			})
			logger.Info(r.Method + " " + r.URL.RequestURI())
		}()

		defer func() {
			if r := recover(); r != nil {
				ctx.Logger.Error(fmt.Errorf("%v: %s", r, debug.Stack()))
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		if err := f(w, r); err != nil {
			ctx.Logger.Error(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
}

func (a *API) IPAddressForRequest(r *http.Request) string {
	addr := r.RemoteAddr
	if a.Config.ProxyCount > 0 {
		h := r.Header.Get("X-Forwarded-For")
		if h != "" {
			clients := strings.Split(h, ",")
			if a.Config.ProxyCount > len(clients) {
				addr = clients[0]
			} else {
				addr = clients[len(clients)-a.Config.ProxyCount]
			}
		}
	}
	return strings.Split(strings.TrimSpace(addr), ":")[0]
}
