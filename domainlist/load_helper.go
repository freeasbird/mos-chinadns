package domainlist

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func LoadFormFile(file string) (*List, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return LoadFormReader(f)
}

func LoadFormReader(r io.Reader) (*List, error) {
	l := New()

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		//ignore lines begin with # and empty lines
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		err := l.Add(line)
		if err != nil {
			return nil, fmt.Errorf("invaild domain [%s], err: [%v]", line, err)
		}
	}

	return l, nil
}
