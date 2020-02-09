//     Copyright (C) 2020, IrineSistiana
//
//     This file is part of mos-chinadns.
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

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// 注意，很长的re匹配会带来严重性能问题
type regExpDomainList struct {
	res []*regexp.Regexp
}

func loadRegExpDomainListFormFile(file string) (rd *regExpDomainList, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return loadRegExpDomainListFormReader(f)
}

func loadRegExpDomainListFormReader(r io.Reader) (rd *regExpDomainList, err error) {
	l := &regExpDomainList{
		res: make([]*regexp.Regexp, 0, 0),
	}

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		//ignore lines begin with # and empty lines
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		re, err := regexp.Compile(line)
		if err != nil {
			return nil, fmt.Errorf("invaild regular expression [%s], err: [%v]", line, err)
		}

		l.res = append(l.res, re)
	}

	return l, nil
}

//a nil list will always return false
func (l *regExpDomainList) match(s string) bool {
	if l == nil {
		return false
	}
	for i := range l.res {
		if l.res[i].MatchString(s) {
			return true
		}
		continue
	}
	return false
}
