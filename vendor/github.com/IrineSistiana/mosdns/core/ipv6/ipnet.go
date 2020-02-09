//     Copyright (C) 2018 - 2019, IrineSistiana
//
//     This file is part of mosdns.
//
//     mosdns is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     mosdns is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ipv6

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"strings"
)

const (
	//32 or 64
	intSize = 32 << (^uint(0) >> 63)

	//2 or 4
	ipSize = 128 / intSize

	maxUint = ^uint(0)
)

//IPv6 represents a ipv6 addr
type IPv6 [ipSize]uint

//mask is ipv6 IP network mask
type mask [ipSize]uint

//Net represents a ip network
type Net struct {
	ip   IPv6
	mask mask
}

var (
	//ErrParseCIDR raised by ParseCIDR()
	ErrParseCIDR = errors.New("error CIDR format")
)

//NewNet returns a new IPNet, mask should be an ipv6 mask,
//which means you should +96 if you have an ipv4 mask.
//ip must be a valid ipv6 address, or ErrNotIPv6 will return.
func NewNet(ipv6 IPv6, mask uint64) Net {
	return Net{
		ip:   ipv6,
		mask: cidrMask(mask),
	}
}

//Contains reports whether the ipnet includes ip.
func (net Net) Contains(ip IPv6) bool {
	start := net.ip
	mask := net.mask
	for i := 0; i < ipSize; i++ {
		if ip[i]&mask[i] == start[i]&mask[i] {
			continue
		}
		return false
	}
	return true
}

var (
	//ErrNotIPv6 raised by Conv()
	ErrNotIPv6 = errors.New("given ip is not a valid ipv6 address")
)

//Conv converts ip to type IPv6.
//ip must be a valid ipv6 address, or ErrNotIPv6 will return.
func Conv(ip net.IP) (IPv6, error) {
	if ip = ip.To16(); ip == nil {
		return IPv6{}, ErrNotIPv6
	}
	intIP := IPv6{}
	switch intSize {
	case 32:
		for i := 0; i < ipSize; i++ { //0 to 3
			s := i * 4
			intIP[i] = uint(binary.BigEndian.Uint32(ip[s : s+4]))
		}
	case 64:
		for i := 0; i < ipSize; i++ { //0 to 1
			s := i * 8
			intIP[i] = uint(binary.BigEndian.Uint64(ip[s : s+8]))
		}
	}

	return intIP, nil
}

//ParseCIDR parses s as a CIDR notation IP address and prefix length.
//As defined in RFC 4632 and RFC 4291.
func ParseCIDR(s string) (Net, error) {

	sub := strings.SplitN(s, "/", 2)
	if len(sub) == 2 { //has "/"
		addrStr, maskStr := sub[0], sub[1]

		//ip
		ip := net.ParseIP(addrStr)
		ipv6, err := Conv(ip)
		if err != nil {
			return Net{}, err
		}

		//mask
		maskLen, err := strconv.ParseUint(maskStr, 10, 0)
		if err != nil {
			return Net{}, err
		}

		//if string is a ipv4 addr, add 96
		if ip.To4() != nil {
			maskLen = maskLen + 96
		}

		return NewNet(ipv6, maskLen), nil
	}

	ipv6, err := Conv(net.ParseIP(s))
	if err != nil {
		return Net{}, err
	}

	return NewNet(ipv6, 128), nil
}

func cidrMask(n uint64) (m mask) {
	for i := uint(0); i < ipSize; i++ {
		if n != 0 {
			m[i] = ^(maxUint >> n)
		} else {
			m[i] = 0
		}

		if n > intSize {
			n = n - intSize
		} else {
			n = 0
		}
	}

	return m
}
