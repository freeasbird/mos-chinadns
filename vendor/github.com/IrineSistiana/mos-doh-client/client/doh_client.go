// Copyright (c) 2020 IrineSistiana
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package client

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/miekg/dns"

	"github.com/valyala/fasthttp"
)

type DohClient struct {
	url string

	httpClient *fasthttp.Client
	bufPool    *sync.Pool

	timeout time.Duration
}

//NewServer returns a DoHServer
func NewClient(url, addr string, sv bool, maxSize int, timeout time.Duration) *DohClient {
	tlsConf := &tls.Config{InsecureSkipVerify: sv}

	c := &DohClient{
		url: url,
		bufPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, maxSize)
			}},
		httpClient: &fasthttp.Client{
			TLSConfig:           tlsConf,
			MaxResponseBodySize: maxSize,
		},
		timeout: timeout,
	}

	if len(addr) != 0 {
		c.httpClient.Dial = func(_ string) (net.Conn, error) {
			return net.Dial("tcp", addr)
		}
	}
	return c
}

func (c *DohClient) Exchange(q *dns.Msg) (*dns.Msg, error) {
	buf := c.bufPool.Get().([]byte)
	defer c.bufPool.Put(buf)

	wireMsg, err := q.PackBuffer(buf)
	if err != nil {
		return nil, err
	}

	//Note: It is forbidden copying Request instances. Create new instances and use CopyTo instead.
	//Request instance MUST NOT be used from concurrently running goroutines.

	req := fasthttp.AcquireRequest()
	url := c.url + "?dns=" + base64.RawURLEncoding.EncodeToString(wireMsg)
	req.SetRequestURI(url)
	req.Header.Set("accept", "application/dns-message")
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := c.httpClient.DoTimeout(req, resp, c.timeout); err != nil {
		return nil, err
	}

	statusCode := resp.StatusCode()
	if statusCode != fasthttp.StatusOK {
		return nil, fmt.Errorf("HTTP status codes [%d]", statusCode)
	}

	r := new(dns.Msg)
	err = r.Unpack(resp.Body())
	if err != nil {
		return nil, err
	}

	return r, nil
}

//ServeDNS impliment the interface
func (c *DohClient) ServeDNS(w dns.ResponseWriter, q *dns.Msg) {

	r, err := c.Exchange(q)
	if err != nil {
		logrus.Errorf("client exchange: %v", err)
	}

	if r != nil {
		w.WriteMsg(r)
	}

	//We do not need to call Close() here
	//defer w.Close()
}
