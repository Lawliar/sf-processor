//
// Copyright (C) 2020 IBM Corporation.
//
// Authors:
// Frederico Araujo <frederico.araujo@ibm.com>
// Teryl Taylor <terylt@ibm.com>
// Changming Liu <Changming.Liu@ibm.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package commons

import (
	"strconv"
)

// Configuration keys.
const (
	TCPHostIpConfigKey string = "grpc.ip"
	TCPPortConfigKey   string = "grpc.port"
)

// ESConfig holds Elastic specific configuration.
type GRPCConfig struct {
	TCPHost string
	TCPPort int
}

func CreateTcpConfig(bc Config, conf map[string]interface{}) (c GRPCConfig, err error) {
	c = GRPCConfig{}

	// parse config map
	if v, ok := conf[TCPHostIpConfigKey].(string); ok {
		c.TCPHost = v
	}
	if v, ok := conf[TCPPortConfigKey].(string); ok {
		c.TCPPort, err = strconv.Atoi(v)
		if err != nil {
			return c, err
		}
	}
	return
}
