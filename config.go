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
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config is config
type Config struct {
	BindAddr                    string `json:"bind_addr"`
	LocalServer                 string `json:"local_server"`
	LocalServerBlockUnusualType bool   `json:"local_server_block_unusual_type"`
	RemoteServer                string `json:"remote_server"`
	RemoteServerURL             string `json:"remote_server_url"`
	RemoteServerSkipVerify      bool   `json:"remote_server_skip_verify"`
	RemoteServerDelayStart      int    `json:"remote_server_delay_start"`

	LocalAllowedIPList     string `json:"local_allowed_ip_list"`
	LocalBlockedIPList     string `json:"local_blocked_ip_list"`
	LocalForcedDomainList  string `json:"local_forced_domain_list"`
	LocalBlockedDomainList string `json:"local_blocked_domain_list"`
	RemoteECSSubnet        string `json:"remote_ecs_subnet"`
}

func loadJSONConfig(configFile string) (*Config, error) {
	c := new(Config)
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}

	return c, nil
}

func genJSONConfig(configFile string) error {
	c := new(Config)

	f, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}
