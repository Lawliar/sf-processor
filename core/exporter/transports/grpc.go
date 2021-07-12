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
	"context"
	"log"
	"strconv"
	"reflect"

	"google.golang.org/grpc"

	pb "github.com/lawliar1/pb"
	"github.com/sysflow-telemetry/sf-apis/go/logger"
	"github.com/sysflow-telemetry/sf-processor/core/exporter/commons"
)

type GrpcProto struct {
	config commons.Config
	conn   *grpc.ClientConn
	stream pb.SysflowGrpc_UploadClient
}

func NewGrpcProto(conf commons.Config) TransportProtocol {
	return &GrpcProto{config: conf}
}

// connect to a remote port
func (s *GrpcProto) Init() (err error) {
	sock := s.config.GRPCHost + ":" + strconv.Itoa(s.config.GRPCPort)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	s.conn, err = grpc.Dial(sock, opts...)
	if err != nil {
		logger.Error.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewSysflowGrpcClient(s.conn)

	ctx := context.Background()
	s.stream, err = client.Upload(ctx)
	if err != nil {
		log.Fatalf("%v.Upload(_) = _, %v", client, err)
	}

	return
}

// write the buffer to the remote port
func (s *GrpcProto) Export(data []commons.EncodedData) (err error) {
	for _, d := range data {
		if entry,ok := d.(*pb.SysflowEntry); ok{
			if err := s.stream.Send(entry); err != nil {
				logger.Error.Println("%v.Send(%v) = %v", s.stream, d, err)
			}
		}else{
			logger.Error.Println("wrong type,",reflect.TypeOf(d).String())
		}
	}
	return
}

func (s *GrpcProto) Register(eps map[commons.Transport]TransportProtocolFactory) {
	eps[commons.GRPCTransport] = NewGrpcProto
}

func (s *GrpcProto) Cleanup() {
	s.conn.Close()
}
