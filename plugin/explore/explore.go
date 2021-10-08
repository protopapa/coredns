package explore

import (
	"context"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"net"
)

type Explore struct{}

var restricted = []string{"example.com.", "protopapa.com.", "protopapa.social."}
var log = clog.NewWithPlugin("explore")

// ServeDNS is the actual handler function, and, unless it fully handles the request by itself, it should call the next handler in the chain
func (e Explore) ServeDNS(ctx context.Context, writer dns.ResponseWriter, msg *dns.Msg) (int, error) {
	state := request.Request{W: writer, Req: msg}

	zone := plugin.Zones(restricted).Matches(state.Name())
	if zone != "" {
		m := new(dns.Msg)
		m.SetRcode(msg, dns.RcodeRefused)
		writer.WriteMsg(m)
		return dns.RcodeSuccess, nil
	}

	m := new(dns.Msg)
	m.SetReply(msg)
	hdr := dns.RR_Header{Name: msg.Question[0].Name, Ttl: 604800, Class: dns.ClassINET, Rrtype: dns.TypeA}
	m.Answer = []dns.RR{&dns.A{Hdr: hdr, A: net.ParseIP("127.0.0.1").To4()}}
	writer.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

func (e Explore) Name() string {
	return "explore"
}
