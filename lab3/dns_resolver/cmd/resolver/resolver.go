package main

import (
	"dns_resolver/internal"
	"flag"
	"fmt"
	"log"
)

func main() {
	dns, address := getArgs()
	//dns := "8.8.8.8:53"
	//address := "example.com"
	fmt.Printf("You entered dns: %v, address: %v\n", dns, address)

	client, err := internal.NewClient(internal.ClientConfig{
		Address: dns,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	_, err = client.LookupAddr(address, internal.QTypeA)
	_, err = client.LookupAddr(address, internal.QTypeAAAA)
}

func getArgs() (string, string) {
	var dns string // dns server address
	flag.StringVar(&dns, "dns", "8.8.8.8:53", "address of dns server")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("You should specify address for resolving\n Example: go run resolver.go example.com")
	}
	address := args[0] // address for resolving

	return dns, address
}
