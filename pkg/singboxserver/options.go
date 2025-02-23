package singboxserver

import (
	"net/netip"

	"github.com/awryme/unchained/appconfig"
	"github.com/awryme/unchained/pkg/protocols"
	"github.com/sagernet/sing-box/option"
	dns "github.com/sagernet/sing-dns"
)

func makeLogOptions(cfg appconfig.Config) *option.LogOptions {
	return &option.LogOptions{
		Disabled:     false,
		Timestamp:    true,
		DisableColor: true,

		Level: cfg.LogLevel,
	}
}

func makeDnsOptions(cfg appconfig.Config) *option.DNSOptions {
	server := option.DNSServerOptions{
		Tag:     "dns-out",
		Address: cfg.DNS,
	}
	if cfg.DNSIPv4Only {
		server.Strategy = option.DomainStrategy(dns.DomainStrategyUseIPv4)
	}
	return &option.DNSOptions{
		Servers: []option.DNSServerOptions{server},
	}
}

func makeInbountOptions(cfg appconfig.Config, addr netip.Addr) (option.Inbound, error) {
	switch cfg.Proto {
	case protocols.Trojan:
		return makeTrojanInboundOptions(cfg, addr), nil
	case protocols.Vless:
		return makeVlessInboundOptions(cfg, addr), nil
	}

	return option.Inbound{}, protocols.Invalid(cfg.Proto)
}
