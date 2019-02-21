package main

import (
	"github.com/coredns/coredns/core/dnsserver"
	_ "github.com/coredns/coredns/core/plugin"
	"github.com/coredns/coredns/coremain"

	"github.com/coredns/metadata_edns0/plugin/metadata_edns0"
)

func init() {
	for i, d := range dnsserver.Directives {
		if d == "metadata" {
			dnsserver.Directives = append(dnsserver.Directives, "")
			copy(dnsserver.Directives[i+1:], dnsserver.Directives[i:])
			dnsserver.Directives[i+1] = metadataEdns0.Name()
		}
	}
}

func main() {
	coremain.Run()
}
