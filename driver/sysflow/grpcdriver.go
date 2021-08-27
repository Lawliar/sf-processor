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
	"bufio"
	"bytes"
	"net"

	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/sysflow-telemetry/sf-apis/go/logger"
	"github.com/sysflow-telemetry/sf-apis/go/plugins"
	"github.com/sysflow-telemetry/sf-apis/go/sfgo"
)

const (
	grpcDriverName = "grpc"
)


// TcpDriver represents a tcp sysflow datasource
type GrpcDriver struct {
	pb.UnimplementedSysflowGrpcServer
	pipeline plugins.SFPipeline
}

func (s* GrpcDriver) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	for {
		sfData, err := stream.Recv()
	}
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

	sFlow := sfgo.NewSysFlow()
	deser, err := compiler.CompileSchemaBytes([]byte(sFlow.Schema()), []byte(sFlow.Schema()))
	if err != nil {
		logger.Error.Println("Compilation error: ", err)
		return err
	}
	for *running {
		buf := make([]byte, tcpBuffSize)
		reader := bytes.NewReader(buf)
		s.conn, err = l.Accept()
		if err != nil {
			logger.Error.Println("Tcp accept error: ", err)
			break
		}
		for *running {
			sFlow = sfgo.NewSysFlow()
			_, err = bufio.NewReader(s.conn).Read(buf[:])
			if err != nil {
				logger.Error.Println("TCP read error: ", err)
				return err
			}
			reader.Reset(buf)
			err = vm.Eval(reader, deser, sFlow)
			if err != nil {
				logger.Error.Println("Deserialization error: ", err)
			}
			records <- sFlow
		}
		s.conn.Close()
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
