package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":8080"
	}

	h := &http.Server{Addr: addr, Handler: &server{}}

	go func() {
		log.Printf("Listening on http://0.0.0.0%s\n", addr)

		if err := h.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Println("\nShutting down the server...")
	h.Shutdown(context.Background())
	log.Println("Server gracefully stopped")
}

type server struct{}

type Response struct {
	LocalIPs       []string `json:"local_ips"`
	RemoteAddr     string   `json:"remote_addr"`
	ServerHostname string   `json:"server_hostname"`
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ifaces, _ := net.Interfaces()
	ips := []string{}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}

			if ip.IsLoopback() || ip.IsInterfaceLocalMulticast() ||
				ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() ||
				ip.IsMulticast() {
				continue
			}

			ips = append(ips, ip.String())
		}
	}

	hostname, _ := os.Hostname()
	w.Header().Set("Content-Type", "application/json")
	resp, _ := json.MarshalIndent(Response{
		LocalIPs:       ips,
		RemoteAddr:     r.RemoteAddr,
		ServerHostname: hostname,
	}, "", "    ")
	w.Write(resp)

	log.Printf("Got query from %s\n", r.RemoteAddr)
}
