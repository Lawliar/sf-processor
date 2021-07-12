//
// Copyright (C) 2020 IBM Corporation.
//
// Authors:
// Andreas Schade <san@zurich.ibm.com>
// Frederico Araujo <frederico.araujo@ibm.com>
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
package encoders

import (
	"fmt"
	pb "github.com/lawliar1/pb"
	"github.com/sysflow-telemetry/sf-apis/go/sfgo"
	"github.com/sysflow-telemetry/sf-processor/core/exporter/commons"
	"github.com/sysflow-telemetry/sf-processor/core/policyengine/engine"
)

type PBEncoder struct {
	config commons.Config
	//jsonencoder JSONEncoder
	batch []commons.EncodedData
}

func NewPBEncoder(config commons.Config) Encoder {
	return &PBEncoder{
		config: config,
		batch:  make([]commons.EncodedData, 0, config.EventBuffer)}
}

// Register registers the encoder to the codecs cache.
func (t *PBEncoder) Register(codecs map[commons.Format]EncoderFactory) {
	codecs[commons.PBFormat] = NewPBEncoder
}

// Encodes telemetry records into an ECS representation.
func (t *PBEncoder) Encode(recs []*engine.Record) ([]commons.EncodedData, error) {
	t.batch = t.batch[:0]
	for _, rec := range recs {
		grpc := t.encode(rec)
		t.batch = append(t.batch, grpc)
	}
	return t.batch, nil
}

// Encodes a telemetry record into an ECS representation.
func (t *PBEncoder) encode(rec *engine.Record) *pb.SysflowEntry {
	pb := &pb.SysflowEntry{}
	for i := uint32(0); i < uint32(sfgo.INT_ARRAY_SIZE); i++ {
		opFlag := rec.GetInt(sfgo.Attribute(i), sfgo.SYSFLOW_SRC)
		pb.Ints = append(pb.Ints, opFlag)
	}
	for i := uint32(0); i < uint32(sfgo.STR_ARRAY_SIZE); i++ {
		str := rec.GetStr(sfgo.Attribute(i), sfgo.SYSFLOW_SRC)
		pb.Strs = append(pb.Strs, str)
	}
	return pb
}

// Cleanup cleans up resources.
func (t *PBEncoder) Cleanup() {}
