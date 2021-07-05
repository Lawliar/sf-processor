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
package transports

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"

	"github.com/sysflow-telemetry/sf-apis/go/logger"
	"github.com/sysflow-telemetry/sf-processor/core/exporter/commons"
)

type TcpProto struct {
	config commons.Config
	conn   net.Conn
}

func NewTcpProto(conf commons.Config) TransportProtocol {
	return &TcpProto{config: conf}
}

// connect to a remote port
func (s *TcpProto) Init() (err error) {
	sock := s.config.TCPHost + ":" + strconv.Itoa(s.config.TCPPort)
	s.conn, err = net.Dial("tcp", sock)
	if err != nil {
		logger.Error.Println("Fail to dail to ", sock)
		return err
	}
	return
}

// write the buffer to the remote port
func (s *TcpProto) Export(data []commons.EncodedData) (err error) {
	for _, d := range data {
		if buf, ok := d.([]byte); ok {
			if _, err = s.conn.Write(buf); err != nil {
				return err
			}
		} else if buf, err := json.Marshal(d); err == nil {
			if _, err = s.conn.Write(buf); err != nil {
				return err
			}
		} else {
			return errors.New("Expected byte array or serializable object as export data")
		}
	}
	return
}

// Register the tcp proto object with the exporter
func (s *TcpProto) Register(eps map[commons.Transport]TransportProtocolFactory) {
	eps[commons.TCPTransport] = NewTcpProto
}

func (s *TcpProto) Cleanup() {
	s.conn.Close()
}
