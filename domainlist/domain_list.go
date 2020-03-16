package domainlist

import (
	"errors"

	"github.com/miekg/dns"
)

// error
var (
	ErrInvalidDomainName = errors.New("invalid doamin name")
)

type List struct {
	s map[[16]byte]struct{}
	m map[[32]byte]struct{}
	l map[[256]byte]struct{}
}

func New() *List {
	return &List{
		s: make(map[[16]byte]struct{}),
		m: make(map[[32]byte]struct{}),
		l: make(map[[256]byte]struct{}),
	}
}

func (l *List) Add(domain string) error {

	if _, ok := dns.IsDomainName(domain); !ok {
		return ErrInvalidDomainName
	}
	fqdn := dns.Fqdn(domain)
	n := len(fqdn)

	switch {
	case n <= 16 && n > 0:
		var b [16]byte
		copy(b[:], fqdn)
		l.s[b] = struct{}{}
	case n <= 32 && n > 16:
		var b [32]byte
		copy(b[:], fqdn)
		l.m[b] = struct{}{}
	case n > 32 && n <= 256:
		var b [256]byte
		copy(b[:], fqdn)
		l.l[b] = struct{}{}
	default:
		return ErrInvalidDomainName
	}

	return nil
}

func (l *List) Has(domain string) bool {
	fqdn := dns.Fqdn(domain)
	if _, ok := dns.IsDomainName(fqdn); !ok {
		return false
	}
	e := dns.Split(fqdn)
	for i := range e {
		p := e[len(e)-1-i]
		if l.has(fqdn[p:]) {
			return true
		}
	}
	return false
}

func (l *List) has(fqdn string) bool {
	n := len(fqdn)
	switch {
	case n <= 16 && n > 0:
		var b [16]byte
		copy(b[:], fqdn)
		_, ok := l.s[b]
		return ok
	case n <= 32 && n > 16:
		var b [32]byte
		copy(b[:], fqdn)
		_, ok := l.m[b]
		return ok
	case n > 32 && n <= 256:
		var b [256]byte
		copy(b[:], fqdn)
		_, ok := l.l[b]
		return ok
	default:
		return false
	}
}

func (l *List) Len() int {
	return len(l.l) + len(l.m) + len(l.s)
}
