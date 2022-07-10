package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/sgaunet/ratelimit"
)

type App struct {
	http.Handler
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

	go app.removeObsoleteRates()

	return app, nil
}

func (a *App) removeObsoleteRates() {
	tick := time.NewTicker(10 * time.Second)
	for range tick.C {
		a.mu.Lock()
		for ip := range a.rates {
			// if IP have not requested service since 2*a.cfg.RateDurationInSeconds
			if a.rates[ip].GetLastCall().Before(time.Now().Add(time.Duration(-2*a.cfg.RateDurationInSeconds) * time.Second)) {
				// delete the IP from the a.rates map
				delete(a.rates, ip)
			}
		}
		a.mu.Unlock()
	}
}

// Serve a reverse proxy for a given url
func (a *App) serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, err := url.Parse(target)
	if err != nil {
		log.Errorln(err.Error())
		fmt.Fprintf(res, "Internal error (ratelimiter)")
		http.Error(res, "Too many requests", http.StatusInternalServerError)
		return
	}

	log.Debugln("req.URL.Host=", req.URL.Host)
	log.Debugln("req.RequestURI=", req.RequestURI)
	log.Debugln("req.URL.Path=", req.URL.Path)
	log.Debugln("req.URL.RawPath=", req.URL.RawPath)
	log.Debugln("req.URL.RawQuery=", req.URL.RawQuery)
	hostCopy := req.Host // Need to preserve original Host request in the response
	if url.Host == "" {
		log.Debugln("url.Host empty, set to localhost")
		url.Host = "localhost"
	}

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &myTransport{}
	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Debugln("modify response Host=", hostCopy)
		resp.Request.Host = hostCopy
		return nil
	}

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	// req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

	log.Debugln("modified req.URL.Host=", req.URL.Host)
	log.Debugln("not modified req.Host=", req.Host)
	log.Debugln("modified req.URL.Scheme=", req.URL.Scheme)

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
func (a *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
	// X-Forwarded-For: <client>, <proxy1>, <proxy2>
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	firstIP := strings.Split(forwarded, ", ")[0]
	if firstIP != "" {
		return firstIP
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func (a *App) LaunchWebServer() {
	port := fmt.Sprintf(":%d", a.cfg.DaemonPort)
	if err := http.ListenAndServe(port, a); err != nil {
		log.Fatalln(err.Error())
	}
}
