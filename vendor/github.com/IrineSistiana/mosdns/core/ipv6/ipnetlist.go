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
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

//NetList is a list of Nets. All Nets will be in ipv6 format, even it's
//ipv4 addr. Cause we use bin search.
type NetList struct {
	elems  []Net
	sorted bool
}

//NewNetList returns a NetList, list can not be nil.
func NewNetList(list []Net) *NetList {
	return &NetList{
		elems: list,
	}
}

//NewNetListFromReader read IP list from a reader, if no valid IP addr was found,
//it will return a empty NetList, NOT nil.
func NewNetListFromReader(reader io.Reader) (*NetList, error) {

	ipNetList := NewNetList(make([]Net, 0))
	s := bufio.NewScanner(reader)

	//count how many lines we have readed.
	lineCounter := 0

	for s.Scan() {
		lineCounter++
		line := strings.TrimSpace(s.Text())

		//ignore lines begin with # and empty lines
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		ipnet, err := ParseCIDR(line)
		if err != nil {
			return nil, fmt.Errorf("invaild CIDR format in line %d", lineCounter)
		}

		ipNetList.Append(ipnet)
	}

	ipNetList.Sort()
	return ipNetList, nil
}

//Append appends new Nets to the list.
func (list *NetList) Append(newNet ...Net) {
	list.elems = append(list.elems, newNet...)
	list.sorted = false
}

//Sort sorts the list, this must be called everytime after
//list was modified.
func (list *NetList) Sort() {
	sort.Sort(list)
	list.sorted = true
}

//implement sort Interface
func (list *NetList) Len() int {
	return len(list.elems)
}

func (list *NetList) Less(i, j int) bool {
	return smallOrEqual(list.elems[i].ip, list.elems[j].ip)
}

func (list *NetList) Swap(i, j int) {
	tmp := list.elems[i]
	list.elems[i] = list.elems[j]
	list.elems[j] = tmp
}

//NewNetListFromFile read IP list from a file, if no valid IP addr was found,
//it will return a empty NetList, NOT nil.
func NewNetListFromFile(file string) (*NetList, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewNetListFromReader(f)
}

//Contains reports whether the list includes given ipv6.
//list must be sorted, or Contains will panic.
func (list *NetList) Contains(ipv6 IPv6) bool {
	if !list.sorted {
		panic("list is not sorted")
	}

	i, j := 0, len(list.elems)
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h

		if smallOrEqual(list.elems[h].ip, ipv6) {
			i = h + 1
		} else {
			j = h
		}
	}

	if i == 0 {
		return false
	}

	return list.elems[i-1].Contains(ipv6)
}

//smallOrEqual IP1 <= IP2 ?
func smallOrEqual(IP1, IP2 IPv6) bool {
	for k := 0; k < ipSize; k++ {
		switch {
		case IP1[k] == IP2[k]:
			continue
		case IP1[k] > IP2[k]:
			return false
		case IP1[k] < IP2[k]:
			return true
		}
	}
	return true
}
