package main

import (
	"fmt"
	"log"
	"os"
	"net"
)

func main() {
	args := os.Args;
	if len(args) != 2 {
		fmt.Println("usage nslookup <name>");
		return;
	}

	host := args[1];
	ips, err := net.LookupIP(host);
	if err != nil {
		log.Fatal(err);
		return;
	}
	for _, v := range ips {
		fmt.Println(v);
	}

	return;
}
