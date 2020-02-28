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
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Sirupsen/logrus"
)

var (
	configPath  = flag.String("c", "", "[path] load config from file")
	genConfigTo = flag.String("gen", "", "[path] generate a config template here")

	dir     = flag.String("dir", "", "[path] change working directory to here")
	verbose = flag.Bool("v", false, "more log")
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
		} else {
			logrus.Print("config template generated")
		}
		return
	}

	//change working dir
	if *dir == "" {
		*dir = filepath.Dir(os.Args[0])
	}
	dir, err := filepath.Abs(*dir)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		logrus.Fatal(err)
	}

	//checking

	if len(*configPath) == 0 {
		logrus.Fatal("need a config file")
	}

	d, err := initDispather(getConfigOrFatal())
	if err != nil {
		logrus.Fatal(err)
	}

	go func() {
		logrus.Print("server started")
		if err := d.ListenAndServe(); err != nil {
			logrus.Fatalf("server exited with err: %v", err)
		} else {
			logrus.Print("server exited")
			os.Exit(0)
		}
	}()

	//wait signals
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
	s := <-osSignals
	logrus.Printf("exiting: signal: %v", s)
	os.Exit(0)
}

func getConfigOrFatal() *Config {
	c, err := loadJSONConfig(*configPath)
	if err != nil {
		logrus.Fatalf("can not load config file, %v", err)
	}
	return c
}
