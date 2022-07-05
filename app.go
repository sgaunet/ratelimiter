package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sgaunet/ratelimit"
)

type App struct {
	cfg   ConfigYaml
	rates map[string]*ratelimit.RateLimit
	mu    sync.Mutex
}

func NewApp(cfg ConfigYaml) (*App, error) {
	app := &App{
		cfg:   cfg,
		rates: make(map[string]*ratelimit.RateLimit),
	}
	log.Debugln("           DaemonPort=", cfg.DaemonPort)
	log.Debugln("RateDurationInSeconds=", cfg.RateDurationInSeconds)
	log.Debugln("           RateNumber=", cfg.RateNumber)
	log.Debugln("        TargetService=", cfg.TargetService)

	return app, nil
}

// Serve a reverse proxy for a given url
func (a *App) serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, err := url.Parse(target)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if url.Scheme == "" {
		url.Scheme = "http"
	}
	if url.Host == "" {
		url.Host = "localhost"
	}

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &myTransport{}

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.Host = url.Host
	req.URL.Scheme = url.Scheme

	// req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func (a *App) handleRate(ip string) error {
	var err error
	a.mu.Lock()
	defer a.mu.Unlock()
	_, found := a.rates[ip]
	if !found {
		a.rates[ip], err = ratelimit.New(context.Background(), time.Duration(a.cfg.RateDurationInSeconds)*time.Second, a.cfg.RateNumber)
	}
	return err
}

// Given a request send it to the appropriate url
func (a *App) handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	var httpStatusCode int
	err := a.handleRate(GetIP(req))
	if err != nil {
		log.Errorln(err.Error())
	} else {
		a.mu.Lock()
		isLimitReached := a.rates[GetIP(req)].IsLimitReached()
		a.mu.Unlock()
		if !isLimitReached {
			a.serveReverseProxy(a.cfg.TargetService, res, req)
		} else {
			httpStatusCode = http.StatusTooManyRequests
			http.Error(res, "Too many requests", httpStatusCode)
			formatLog := "req.RemoteAddr=%s req.Host=%s req.URL.Path=%s req.URL.Query()=%s StatusCode=%d\n"
			log.Infof(formatLog, GetIP(req), req.Host, req.URL.Path, req.URL.Query(), httpStatusCode)
		}
	}
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return strings.Split(forwarded, ":")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func (a *App) LaunchWebServer() {
	http.HandleFunc("/", a.handleRequestAndRedirect)
	port := fmt.Sprintf(":%d", a.cfg.DaemonPort)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalln(err.Error())
	}
}
