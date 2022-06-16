package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/sgaunet/ratelimit"
)

type App struct {
	cfg ConfigYaml
	r   *ratelimit.RateLimit
}

func NewApp(cfg ConfigYaml) (*App, error) {
	r, err := ratelimit.New(context.Background(), time.Duration(cfg.RateDurationInSeconds)*time.Second, cfg.RateNumber)
	if err != nil {
		return nil, err
	}
	app := &App{
		cfg: cfg,
		r:   r,
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

// Given a request send it to the appropriate url
func (a *App) handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	var httpStatusCode int
	if !a.r.IsLimitReached() {
		a.serveReverseProxy(a.cfg.TargetService, res, req)
	} else {
		httpStatusCode = http.StatusTooManyRequests
		http.Error(res, "Too many requests", httpStatusCode)
		formatLog := "req.RemoteAddr=%s req.Host=%s req.URL.Path=%s req.URL.Query()=%s StatusCode=%d\n"
		log.Infof(formatLog, GetIP(req), req.Host, req.URL.Path, req.URL.Query(), httpStatusCode)
	}
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func (a *App) LaunchWebServer() {
	http.HandleFunc("/", a.handleRequestAndRedirect)
	port := fmt.Sprintf(":%d", a.cfg.DaemonPort)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalln(err.Error())
	}
}
