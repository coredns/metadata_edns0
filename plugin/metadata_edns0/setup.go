package metadataEdns0

import (
	"fmt"
	"strconv"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/caddyserver/caddy"
)

func init() {
	caddy.RegisterPlugin(Name(), caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	r, err := parse(c)

	if err != nil {
		return plugin.Error(Name(), err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		r.Next = next
		return r
	})

	return nil
}

func parse(c *caddy.Controller) (*metadataEdns0, error) {
	r := New()
	for c.Next() {
		c.RemainingArgs()
		for c.NextBlock() {
			err := r.parseEDNS0(c)
			if err != nil {
				return nil, err
			}
		}
	}
	return r, nil
}

func (m *metadataEdns0) parseEDNS0(c *caddy.Controller) error {
	name := c.Val()
	args := c.RemainingArgs()
	// <label> <definition>
	// <label> <id>
	// <label> <id> <encoded-format> <params of format ...>
	// Valid encoded-format are hex (default), bytes, ip.

	argsLen := len(args)
	if argsLen != 1 && argsLen != 2 && argsLen != 5 {
		return fmt.Errorf("invalid edns0 directive")
	}
	code := args[0]

	dataType := "hex"
	size := "0"
	start := "0"
	end := "0"

	if argsLen > 1 {
		dataType = args[1]
	}

	if argsLen == 5 && dataType == "hex" {
		size = args[2]
		start = args[3]
		end = args[4]
	}

	err := m.addEDNS0Map(code, name, dataType, size, start, end)
	if err != nil {
		return fmt.Errorf("could not add EDNS0 map for %s: %s", name, err)
	}

	return nil
}

func parseEDNS0Map(code, name, dataType, sizeStr, startStr, endStr string) (*edns0Map, error) {
	c, err := strconv.ParseUint(code, 0, 16)
	if err != nil {
		return nil, fmt.Errorf("could not parse EDNS0 code: %s", err)
	}
	size, err := strconv.ParseUint(sizeStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("could not parse EDNS0 data size: %s", err)
	}
	start, err := strconv.ParseUint(startStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("could not parse EDNS0 start index: %s", err)
	}
	end, err := strconv.ParseUint(endStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("could not parse EDNS0 end index: %s", err)
	}
	if end <= start && end != 0 {
		return nil, fmt.Errorf("end index should be > start index (actual %d <= %d)", end, start)
	}
	if end > size && size != 0 {
		return nil, fmt.Errorf("end index should be <= size (actual %d > %d)", end, size)
	}
	ednsType, ok := stringToEDNS0MapType[dataType]
	if !ok {
		return nil, fmt.Errorf("invalid dataType for EDNS0 map: %s", dataType)
	}
	ecode := uint16(c)
	return &edns0Map{name, ecode, ednsType, uint(size), uint(start), uint(end)}, nil
}

func (m *metadataEdns0) addEDNS0Map(code, name, dataType, sizeStr, startStr, endStr string) error {
	p, err := parseEDNS0Map(code, name, dataType, sizeStr, startStr, endStr)
	if err != nil {
		return err
	}
	m.options[p.code] = append(m.options[p.code], p)
	return nil
}
