package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func httpHandler(u *url.URL) http.HandlerFunc {
	var reverseProxy = httputil.NewSingleHostReverseProxy(u)
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("http: %s\n", r.URL)
		reverseProxy.ServeHTTP(w, r)
	}
}

func websocketHandler(target string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d, err := net.Dial("tcp", target)
		if err != nil {
			http.Error(w, "Error contacting backend server.", 500)
			log.Printf("Error dialing websocket backend %s: %v", target, err)
			return
		}
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Not a hijacker?", 500)
			return
		}
		nc, _, err := hj.Hijack()
		if err != nil {
			log.Printf("Hijack error: %v", err)
			return
		}
		defer nc.Close()
		defer d.Close()

		err = r.Write(d)
		if err != nil {
			log.Printf("Error copying request to target: %v", err)
			return
		}

		errc := make(chan error, 2)
		cp := func(dst io.Writer, src io.Reader) {
			_, err := io.Copy(dst, src)
			errc <- err
		}
		go cp(d, nc)
		go cp(nc, d)
		<-errc
	})
}

func StartProxy(ipaddr, port string) {
	var httpBackend *url.URL
	var err error

	httpBackend, err = url.Parse(fmt.Sprintf("http://%s:%s", ipaddr, port))
	if err != nil {
		return
	}

	wsBackend := fmt.Sprintf("%s:%s", ipaddr, port)

	http.HandleFunc("/api/kernels/", websocketHandler(wsBackend))
	http.HandleFunc("/terminals/websocket/", websocketHandler(wsBackend))
	http.HandleFunc("/", httpHandler(httpBackend))
	http.ListenAndServe(":3000", nil)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: <ip> <port>")
		os.Exit(0)
	}

	ip := os.Args[1]
	port := os.Args[2]

	fmt.Printf("Proxy http/ws on port 3000 to %s:%s\n", ip, port)

	StartProxy(ip, port)
}
