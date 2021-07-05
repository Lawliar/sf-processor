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

package sysflow

import (
	"bytes"
	"net"

	"github.com/linkedin/goavro"
	"github.com/sysflow-telemetry/sf-apis/go/converter"
	"github.com/sysflow-telemetry/sf-apis/go/logger"
	"github.com/sysflow-telemetry/sf-apis/go/plugins"
)

const (
	tcpDriverName = "tcp"
)

const (
	tcpBuffSize   = 1024
	tcpOOBuffSize = 1024
)

// TcpDriver represents a tcp sysflow datasource
type TcpDriver struct {
	pipeline plugins.SFPipeline
	conn     net.Conn
}

// create a new tcp driver
func NewTcpDriver() plugins.SFDriver {
	return &TcpDriver{}
}

// GetName returns the driver name.
func (s *TcpDriver) GetName() string {
	return tcpDriverName
}

// Register registers driver to plugin cache
func (s *TcpDriver) Register(pc plugins.SFPluginCache) {
	pc.AddDriver(tcpDriverName, NewTcpDriver)
}

// Init initializes the driver
func (s *TcpDriver) Init(pipeline plugins.SFPipeline) error {
	s.pipeline = pipeline
	return nil
}

func (s *TcpDriver) Run(path string, running *bool) error {
	channel := s.pipeline.GetRootChannel()
	sfChannel := channel.(*plugins.SFChannel)

	records := sfChannel.In

	l, err := net.Listen("tcp", ":"+path)
	if err != nil {
		logger.Error.Println("Cannot listen to port ", path, err)
		return err
	}
	defer l.Close()

	sfobjcvter := converter.NewSFObjectConverter()

	for *running {
		buf := make([]byte, tcpBuffSize)
		reader := bytes.NewReader(buf)
		sreader, err := goavro.NewOCFReader(reader)
		if err != nil {
			logger.Error.Println("Reader error: ", err)
			return err
		}

		s.conn, err = l.Accept()
		if err != nil {
			logger.Error.Println("Tcp accept error: ", err)
			break
		}
		for sreader.Scan() {
			if !*running {
				break
			}
			datum, err := sreader.Read()
			if err != nil {
				logger.Error.Println("Datum reading error: ", err)
				break
			}
			records <- sfobjcvter.ConvertToSysFlow(datum)
		}
		s.conn.Close()
		if !*running {
			break
		}

	}
	logger.Trace.Println("Closing main channel")
	close(records)
	s.pipeline.Wait()
	return nil
}

// Cleanup tears down the driver resources.
func (s *TcpDriver) Cleanup() {
	logger.Trace.Println("Exiting ", tcpDriverName)
	if s.conn != nil {
		s.conn.Close()
	}
}
