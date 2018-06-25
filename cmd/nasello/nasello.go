// Nasello is a DNS proxy server
package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"github.com/piger/nasello"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	configFile = flag.String("config", "nasello.json", "Configuration file")
	listenAddr = flag.String("listen", "localhost:8053", "Local bind address")
)

func serve(net string, address string) {
	err := dns.ListenAndServe(address, net, nil)
	if err != nil {
		log.Fatalf("Failed to setup the "+net+" server: %s\n", err.Error())
	}
}

func main() {
	flag.Parse()

	configuration := nasello.ReadConfig(*configFile)
	for _, rule := range configuration.Rules {
		// Ensure that each pattern is a FQDN name
		pattern := dns.Fqdn(rule.Match)
		resolver := configuration.Resolvers[rule.Resolver]

		log.Printf("Proxing %s on %v(%s)\n", pattern, strings.Join(resolver.Servers, ", "), resolver.Protocol)
		var addresses []string
		var port int
		if resolver.Port == 0 {
			port = 53
		} else {
			port = resolver.Port
		}
		for _, a := range resolver.Servers {
			addresses = append(addresses, fmt.Sprintf("%s:%d", a, port))
		}

		dns.HandleFunc(pattern, nasello.ServerHandler(addresses, resolver.Protocol))
	}

	go serve("tcp", *listenAddr)
	go serve("udp", *listenAddr)

	log.Printf("Started DNS server on: %s\n", *listenAddr)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	recvSig := <-sig
	log.Printf("Signal (%s) received, stopping\n", recvSig.String())
}
