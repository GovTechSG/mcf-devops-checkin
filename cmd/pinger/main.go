package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

const DefaultInterface = "0.0.0.0"
const DefaultPingTimeout = 10 * time.Second
const DefaultPort = 8000
const DefaultTargetHost = "localhost"
const DefaultTargetProto = "http"
const DefaultTimeout = 5 * time.Second
const ExitCodeSuccess = 0
const ExitCodeInitFailed = 1
const ExitCodeMainFailed = 2

var config Config
var errorLogger = log.New(os.Stdout, "  error|", log.LstdFlags)
var serverLogger = log.New(os.Stdout, " server|", log.LstdFlags)
var serviceLogger = log.New(os.Stdout, "service|", log.LstdFlags)
var readiness = map[string]bool{
	"target_up": false,
}

func init() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
			os.Exit(ExitCodeInitFailed)
		}
	}()
	c := viper.New()
	c.SetDefault("interface", DefaultInterface)
	c.SetDefault("ping_timeout", DefaultPingTimeout)
	c.SetDefault("port", DefaultPort)
	c.SetDefault("target_proto", DefaultTargetProto)
	c.SetDefault("target_host", DefaultTargetHost)
	c.SetDefault("target_port", DefaultPort)
	c.AutomaticEnv()

	config = Config{
		Interface:   c.GetString("interface"),
		PingTimeout: c.GetDuration("ping_timeout"),
		Port:        uint16(c.GetUint("port")),
		TargetProto: c.GetString("target_proto"),
		TargetHost:  c.GetString("target_host"),
		TargetPort:  uint16(c.GetUint("target_port")),
		TargetPath:  c.GetString("target_path"),
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
			os.Exit(ExitCodeMainFailed)
		}
	}()

	// setup
	serviceLogger.Printf("initialising service...\n")
	addr := fmt.Sprintf("%s:%v", config.Interface, config.Port)
	mux := createMux()
	server := createServer(addr, mux)
	tick := time.Tick(time.Second)
	done := make(chan bool, 1)
	ossig := make(chan os.Signal, 1)
	signal.Notify(ossig, syscall.SIGTERM, syscall.SIGINT)

	// listen
	serviceLogger.Printf("attempting to listen on '%s'...\n", server.Addr)
	go server.ListenAndServe()
	for {
		select {
		// pinger
		case <-tick:
			client := http.Client{
				Timeout: DefaultPingTimeout,
			}
			req, err := http.NewRequest("GET", config.getTargetURL(), nil)
			handleError(err, errorLogger)
			response, err := client.Do(req)
			handleError(err, errorLogger)
			serviceLogger.Printf("> %s -> '%s'\n", config.getTargetURL(), response.Status)
		// handle os termination
		case sig := <-ossig:
			serviceLogger.Printf("received termination signal '%v', shutting down server now\n", sig)
			done <- true
		// handle graceful shutdown
		case <-done:
			server.Close()
			serviceLogger.Println("exiting now...")
			os.Exit(ExitCodeSuccess)
		}
	}
}

func createMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("alive"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if readiness["target_up"] {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ready"))
		} else {
			w.WriteHeader(http.StatusTeapot)
			w.Write([]byte("not ready"))
		}
	})
	return mux
}

func createServer(addr string, handler http.Handler) http.Server {
	return http.Server{
		Addr:              addr,
		Handler:           requestLoggerMiddleware(handler),
		ReadTimeout:       DefaultTimeout,
		ReadHeaderTimeout: DefaultTimeout,
		WriteTimeout:      DefaultTimeout,
		ErrorLog:          errorLogger,
	}
}

func handleError(err error, log *log.Logger) {
	if err != nil {
		log.Println(err)
	}
}

// requestLoggerMiddleware is a middleware that
func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		serverLogger.Printf("< %s <- %s | %s %s %s \n", r.Host, r.RemoteAddr, strings.ToUpper(r.Proto), r.Method, r.URL.Path)
	})
}

// Config provides an interface for the configuration we will be using
// in this service
type Config struct {
	Interface   string        `json:"interface"`
	PingTimeout time.Duration `json:"port"`
	Port        uint16        `json:"port"`
	TargetHost  string        `json:"target_host"`
	TargetPath  string        `json:"target_path"`
	TargetPort  uint16        `json:"target_port"`
	TargetProto string        `json:"target_proto"`
}

// getTargetURL retrieves the exact URL which we can use to ping the
// target server
func (c *Config) getTargetURL() string {
	return fmt.Sprintf("%s://%s:%v/%s", c.TargetProto, c.TargetHost, c.TargetPort, c.TargetPath)
}
