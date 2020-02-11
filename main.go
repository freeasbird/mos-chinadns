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
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Sirupsen/logrus"
)

var (
	configPath  = flag.String("c", "", "path to config file")
	genConfigTo = flag.String("gen", "", "path to generate a config template")

	dir     = flag.String("dir", "", "change working directory")
	verbose = flag.Bool("v", false, "more log")

	bindAddr               = flag.String("b", "", "bind address")
	localServer            = flag.String("local", "", "local dns server address")
	remoteServer           = flag.String("remote", "", "remote dns server address")
	useTCP                 = flag.String("use-tcp", "", "one of 'l','r', or 'l_r'. Means local only, remote only or local and remote")
	localAllowedIPList     = flag.String("local-allowed-ip", "", "a file that contains a list of IPs that should be allowed by local server")
	localBlockedIPList     = flag.String("local-blocked-ip", "", "a file that contains a list of IPs that should be blocked by local server")
	localBlockedDomainList = flag.String("local-blocked-domain", "", "a file that contains a list of regexp that should be blocked by local server")
	remoteECSSubnet        = flag.String("remote-ecs-subnet", "", "a CIDR notation, EDNS client subnet")
)

func main() {
	flag.Parse()

	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
	}

	//gen config
	if len(*genConfigTo) != 0 {
		err := genJSONConfig(*genConfigTo)
		if err != nil {
			logrus.Errorf("can not generate config template, %v", err)
			return
		}
		logrus.Print("config template generated")
		return
	}

	//change working dir
	if *dir == "" {
		*dir = filepath.Dir(os.Args[0])
	}
	dir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		log.Fatal(err)
	}

	d, err := initDispather(buildConfig())
	if err != nil {
		logrus.Fatal(err)
	}

	go d.ListenAndServe()

	//wait signals
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
	s := <-osSignals
	logrus.Printf("exiting: signal: %v", s)
	os.Exit(0)
}

func buildConfig() *Config {
	c := new(Config)
	if len(*configPath) != 0 {
		c2, err := loadJSONConfig(*configPath)
		if err != nil {
			logrus.Fatalf("can not load config file, %v", err)
		}
		c = c2
	}

	setIfNotNil(&c.BindAddr, bindAddr)
	setIfNotNil(&c.LocalServer, localServer)
	setIfNotNil(&c.RemoteServer, remoteServer)
	setIfNotNil(&c.LocalAllowedIPList, localAllowedIPList)
	setIfNotNil(&c.LocalBlockedIPList, localBlockedIPList)
	setIfNotNil(&c.LocalBlockedDomainList, localBlockedDomainList)
	setIfNotNil(&c.UseTCP, useTCP)
	setIfNotNil(&c.RemoteECSSubnet, remoteECSSubnet)

	return c
}

func setIfNotNil(dst *string, src *string) {
	if src == nil || len(*src) == 0 {
		return
	}
	*dst = *src
}
