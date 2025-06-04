package web

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"meme3/global"
	"net/http"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

type eFS struct {
	fs      *embed.FS
	webPath string
}

var wfs *embed.FS
var wpath string

func (p eFS) Open(name string) (http.File, error) {
	if p.fs != nil {
		if name == "/" {
			return http.FS(p.fs).Open(p.webPath)
		}
		if _, err := p.fs.Open(p.webPath + "/" + strings.TrimPrefix(name, "/")); err == nil {
			return http.FS(p.fs).Open(p.webPath + "/" + strings.TrimPrefix(name, "/"))
		}
		return http.FS(p.fs).Open(p.webPath)
	} else {
		if name == "/" {
			return http.Dir(p.webPath).Open(name)
		}
		if f, err := http.Dir(p.webPath).Open(name); err == nil {
			return f, nil
		}
		// return nil, errors.New("permission denied")
		return http.Dir(p.webPath).Open("/index.html")
	}
}

func RedirectHome(w http.ResponseWriter, r *http.Request) {
	if wfs != nil {
		http.ServeFile(w, r, wpath+"/index.html")
	} else {
		http.ServeFile(w, r, wpath+"/index.html")
	}
}

func Init(webPath string, fs *embed.FS, routerInit func(*httprouter.Router)) *httprouter.Router {
	router := httprouter.New()

	wfs = fs
	wpath = webPath
	router.NotFound = http.FileServer(&eFS{fs, webPath})

	if routerInit != nil {
		routerInit(router)
	}

	return router
}

func Run(ctx context.Context, port int, router *httprouter.Router, allowCredentials bool, allowOrigins ...string) {
	n := negroni.New()
	n.Use(newCors(allowCredentials, allowOrigins...))
	n.UseFunc(newGzip)
	// n.Use(newRateLimite())
	n.UseFunc(newAPILog)
	// n.UseFunc(newAuth)

	n.UseHandlerFunc(router.ServeHTTP)
	server := &http.Server{Addr: fmt.Sprintf(":%v", port), Handler: n}

	fmt.Println("")
	slog.Info("run web server successed", "listen", port)
	fmt.Println("")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			slog.Error("run web server failed", "detail", err.Error())
		}
	}()
	<-ctx.Done()
	server.Shutdown(ctx)
}

func newAPILog(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !strings.HasPrefix(r.URL.Path, "/api") {
		next(rw, r)
		return
	}

	ts := time.Now()
	next(rw, r)
	slog.Info(fmt.Sprintf("%s %s", r.Method, r.RequestURI), "time", time.Since(ts).Milliseconds())
}

func newAuth(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !strings.HasPrefix(r.URL.Path, "/auth") {
		next(rw, r)
		return
	}

	strTokens := strings.Split(r.URL.Path, "/")
	authID := strTokens[len(strTokens)-1]
	http.Redirect(rw, r, fmt.Sprintf("%s?authID=%s", global.Config.HostURL, authID), http.StatusFound)

	// if _, err := api.WebPlayerAuth(rw, r, ""); err == nil {
	// 	r.URL.Path = "/"
	// 	next(rw, r)
	// }
}

func newGzip(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if strings.HasPrefix(r.URL.Path, "/api") {
		next(rw, r)
		return
	}
	gzip.Gzip(gzip.DefaultCompression).ServeHTTP(rw, r, next)
}

func newRateLimite() negroni.Handler {
	limiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour, ExpireJobInterval: time.Second})
	limiter.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}).SetMethods([]string{"POST"})
	limiter.SetMessage("You have reached maximum request limit.")
	return tollbooth_negroni.LimitHandler(limiter)
}

func newCors(allowCredentials bool, allowOrigins ...string) negroni.Handler {
	return cors.New(cors.Options{
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodPut, http.MethodOptions},
		// AllowedOrigins:   []string{"*"},
		AllowCredentials:    allowCredentials,
		AllowedOrigins:      []string{"*"},
		AllowedHeaders:      []string{"*"},
		AllowPrivateNetwork: true,
	})
}
