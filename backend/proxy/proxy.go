package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

var portsForApps map[string]string = map[string]string{"minio": "9000"}

type ProxyWrapper struct {
	Proxy *httputil.ReverseProxy
}

func (pw *ProxyWrapper) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	splitPath := strings.Split(req.URL.Path, "/")
	if len(splitPath) <= 3 {
		fmt.Fprintf(rw, "Generated at %s", time.Now().UTC().String())
		return
	}
	pw.Proxy.ServeHTTP(rw, req)
}

// The URL is of the format: /:shortcode/:appname/:instancenum

func (pw *ProxyWrapper) DefaultDirector(req *http.Request) {
	splitPath := strings.Split(req.URL.Path, "/")
	shortcode := "invalid"
	appname := "noapp"
	instancenum := "X"
	if len(splitPath) <= 3 {
		req.URL.Scheme = "http"
		req.URL.Host = "google.com" // TODO: Change this to siacdn.com once it's live
		req.URL.Path = "/"
		unsetUserAgent(req)
		return
	}
	shortcode = splitPath[1]
	appname = splitPath[2]
	instancenum = splitPath[3]
	req.URL.Scheme = "http"
	req.URL.Host = fmt.Sprintf("siacdn-%s-%s%s.sia.svc.cluster.local:%s", shortcode, appname, instancenum, portsForApps[appname])
	req.URL.Path = "/" + strings.Join(splitPath[4:], "/")
	unsetUserAgent(req)
}

func New() (http.Handler, error) {
	pw := &ProxyWrapper{}
	pw.Proxy = &httputil.ReverseProxy{
		Director: pw.DefaultDirector,
	}
	return pw, nil
}

func unsetUserAgent(req *http.Request) {
	if _, ok := req.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent so it's not set to default value
		req.Header.Set("User-Agent", "")
	}
}
