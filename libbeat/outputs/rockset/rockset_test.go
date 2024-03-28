// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//go:build !integration

package rockset

import (
	"context"
	"testing"

	"github.com/elastic/elastic-agent-libs/config"
	"github.com/golang/mock/gomock"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elastic/elastic-agent-libs/mapstr"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/outputs/codec"
	"github.com/elastic/beats/v7/libbeat/outputs/codec/json"
	"github.com/elastic/beats/v7/libbeat/outputs/outest"
)

func TestRocksetOutput(t *testing.T) {
	cfgMap := map[string]interface{}{
		"apikey":     "apikey",
		"apiserver":  "apiserver",
		"workspace":  "ws",
		"collection": "coll",
	}
	tests := []struct {
		title  string
		codec  codec.Codec
		events []beat.Event
		config map[string]interface{}
	}{
		{
			"single event",
			json.New("1.2.3", json.Config{
				Pretty:     false,
				EscapeHTML: false,
			}),
			[]beat.Event{
				{Fields: event("field", "value")},
			},
			cfgMap,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockRocker(ctrl)
	m.EXPECT().
		AddDocuments(gomock.Any(), "ws", "coll", gomock.Any()).
		Return([]openapi.DocumentStatus{
			{Status: openapi.PtrString("SUCCESS")},
		}, nil).
		AnyTimes()

	ctx := context.TODO()
	beatInfo := beat.Info{Beat: "libbeat", Version: "1.2.3"}

	for _, tst := range tests {
		tst := tst
		t.Run(tst.title, func(t *testing.T) {
			cfg := defaultConfig
			conf, err := config.NewConfigFrom(tst.config)
			assert.NoError(t, err)
			err = conf.Unpack(&cfgMap)
			require.NoError(t, err)

			enc, err := codec.CreateEncoder(beatInfo, cfg.Codec)
			require.NoError(t, err)

			groups, err := makeOutputGroup(cfg.Queue, 10, enc, outputs.NewNilObserver(), m, "ws", "coll")
			require.NoError(t, err)

			batch := outest.NewBatch(tst.events...)

			err = groups.Clients[0].Publish(ctx, batch)
			require.NoError(t, err)

			// check batch correctly signalled
			if !assert.Len(t, batch.Signals, 1) {
				return
			}
			assert.Equal(t, outest.BatchACK, batch.Signals[0].Tag)
		})
	}
}

func event(k, v string) mapstr.M {
	return mapstr.M{k: v}
}
