package domainlist

import (
	"testing"
)

func Test_DomainList(t *testing.T) {
	l := New()
	l.Add("cn")
	l.Add("a.com")
	l.Add("b.com")
	l.Add("abc.com")
	l.Add("123456789012345678901234567890.com")

	assertTrue(l.Has("a.cn"))
	assertTrue(l.Has("a.b.cn"))

	assertTrue(l.Has("a.com"))
	assertTrue(l.Has("b.com"))
	assertTrue(!l.Has("c.com"))
	assertTrue(!l.Has("a.c.com"))
	assertTrue(l.Has("123456789012345678901234567890.com"))

	assertTrue(l.Has("abc.abc.com"))

	assertTrue(l.Add("") == ErrInvalidDomainName)
	assertTrue(l.Add(string(make([]byte, 256))) == ErrInvalidDomainName)
}

func assertTrue(b bool) {
	if !b {
		panic("assert failed")
	}
}
