package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Attempts int
type Retry int

type Server struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (s *Server) SetAlive(alive bool) {
	s.mux.Lock()
	s.Alive = alive
	s.mux.Unlock()
}

func (s *Server) IsAlive() (alive bool) {
	s.mux.RLock()
	alive = s.Alive
	s.mux.RUnlock()
	return
}

type ServerPool struct {
	servers []*Server
	current atomic.Uint32
}

func (s *ServerPool) AddServer(server *Server) {
	s.servers = append(s.servers, server)
}

func (s *ServerPool) NextIndex() int {
	return int(s.current.Add(1)) & len(s.servers)
}

func (s *ServerPool) MarkServerStatus(serverUrl *url.URL, alive bool) {
	for _, b := range s.servers {
		if b.URL.String() == serverUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}

func (s *ServerPool) GetNext() *Server {
	next := s.NextIndex()
	for i := next; i < len(s.servers)+next; i++ {
		idx := i % len(s.servers)
		if s.servers[idx].IsAlive() {
			if i != next {
				s.current.Store(uint32(idx))
			}
			return s.servers[idx]
		}
	}
	return nil
}

func (s *ServerPool) HealthCheck() {
	for _, b := range s.servers {
		alive := isServerAlive(b.URL)
		b.SetAlive(alive)
		status := "up"
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL, status)
	}
}

func lb(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := serverPool.GetNext()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func isServerAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Server unreachable, error: ", err)
		return false
	}
	defer conn.Close()
	return true
}

func healthCheck(period int) {
	t := time.NewTicker(time.Duration(period) * time.Second)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}
}

func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(Attempts(0)).(int); ok {
		return attempts
	}
	return 1
}

func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry(0)).(int); ok {
		return retry
	}
	return 0
}

var serverPool ServerPool

const maxRetriesNumber = 3

func main() {
	var serverList string
	var port int
	var hCheckPeriod int
	flag.StringVar(&serverList, "servers", "", "List of load balanced servers, use commas to separate ()")
	flag.IntVar(&port, "port", 3030, "Port to serve")
	flag.IntVar(&hCheckPeriod, "hCheckPeriod", 1, "Period for health checking (in seconds)")
	flag.Parse()

	if len(serverList) == 0 {
		log.Fatal("Please provide one or more servers to load balance")
	}

	serverInfo := strings.Split(serverList, ",")
	for _, info := range serverInfo {
		serverUrl, err := url.Parse(info)
		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		proxy.ErrorHandler = func(w http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
			retries := GetRetryFromContext(request)
			if retries < maxRetriesNumber {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), Retry(0), retries+1)
					proxy.ServeHTTP(w, request.WithContext(ctx))
				}
				return
			}

			serverPool.MarkServerStatus(serverUrl, false)

			attempts := GetAttemptsFromContext(request)
			log.Printf("[%s] Attempting retry %d\n", request.RequestURI, attempts)
			ctx := context.WithValue(request.Context(), Attempts(0), attempts+1)
			lb(w, request.WithContext(ctx))
		}

		serverPool.AddServer(&Server{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
		log.Printf("Configured server: %s\n", serverUrl)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(lb),
	}

	go healthCheck(hCheckPeriod)

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
