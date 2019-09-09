module github.com/coredns/metadata_edns0

go 1.12

require (
	github.com/caddyserver/caddy v1.0.3
	github.com/coredns/coredns v1.6.3
	github.com/miekg/dns v1.1.16
)

replace golang.org/x/net v0.0.0-20190813000000-74dc4d7220e7 => golang.org/x/net v0.0.0-20190813141303-74dc4d7220e7
