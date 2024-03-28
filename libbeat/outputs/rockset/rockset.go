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

package rockset

import (
	"context"
	"encoding/json"
	"time"

	"github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/openapi"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/outputs/codec"
	"github.com/elastic/beats/v7/libbeat/publisher"
)

type output struct {
	log        *logp.Logger
	observer   outputs.Observer
	codec      codec.Codec
	rockset    Rocker
	workspace  string
	collection string
}

const name = "rockset"

func init() {
	outputs.RegisterType(name, makeRocksetOutput)
}

// Run 'go generate' to create mocks that are used in tests.
//go:generate go install github.com/golang/mock/mockgen@v1.6.0
//go:generate mockgen -source=rockset.go -destination=mock_interface_test.go -package rockset

type Rocker interface {
	GetOrganization(ctx context.Context) (openapi.Organization, error)
	AddDocuments(ctx context.Context, workspace, collection string, documents []interface{}) ([]openapi.DocumentStatus, error)
}

func makeRocksetOutput(
	_ outputs.IndexManager,
	beat beat.Info,
	observer outputs.Observer,
	configuration *config.C,
) (outputs.Group, error) {
	cfg := defaultConfig
	err := configuration.Unpack(&cfg)
	if err != nil {
		return outputs.Fail(err)
	}

	enc, err := codec.CreateEncoder(beat, cfg.Codec)
	if err != nil {
		return outputs.Fail(err)
	}

	// TODO add other config options?
	rs, err := rockset.NewClient(rockset.WithAPIKey(cfg.APIKey), rockset.WithAPIServer(cfg.APIServer))
	if err != nil {
		return outputs.Fail(err)
	}

	return makeOutputGroup(cfg.Queue, cfg.BatchSize, enc, observer, rs, cfg.Workspace, cfg.Collection)
}

func makeOutputGroup(
	queue config.Namespace,
	batchSize int,
	encoder codec.Codec,
	observer outputs.Observer,
	rs Rocker,
	workspace, collection string,
) (outputs.Group, error) {
	return outputs.Success(queue, batchSize, 0, &output{
		log:        logp.NewLogger(name),
		codec:      encoder,
		observer:   observer,
		rockset:    rs,
		workspace:  workspace,
		collection: collection,
	})
}

func (o *output) String() string { return name }
func (o *output) Close() error   { return nil }

func (o *output) Connect() error {
	// TODO should the timeout be configurable?
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	org, err := o.rockset.GetOrganization(ctx)
	cancel()
	if err != nil {
		return err
	}
	o.log.Infof("connected to Rockset organization %s", org.GetDisplayName())

	return nil
}

func (o *output) Publish(ctx context.Context, batch publisher.Batch) error {
	events := batch.Events()
	o.observer.NewBatch(len(batch.Events()))

	// TODO rockset caps each write to 1000 documents, so we should enforce it here too
	docs := make([]interface{}, 0, len(events))

	dropped := 0
	for i, e := range events {
		// TODO should we use the index as the collection?
		serializedEvent, err := o.codec.Encode("index", &e.Content)
		if err != nil {
			dropped++
			o.log.Errorf("failed to serialize event %d: %v", i, err)
			continue
		}

		// generic way to describe a json object
		var doc map[string]interface{}

		err = json.Unmarshal(serializedEvent, &doc)
		if err != nil {
			dropped++
			o.log.Errorf("failed to unmarshal serialized event %d: %v", i, err)
			continue
		}

		docs = append(docs, doc)
	}

	t0 := time.Now()
	response, err := o.rockset.AddDocuments(ctx, o.workspace, o.collection, docs)
	if err != nil {
		o.observer.WriteError(err)
		return err
	}
	o.observer.ReportLatency(time.Since(t0))

	var nonAddedDocs int
	for i, doc := range response {
		if doc.GetStatus() != "ADDED" {
			o.log.Debugf("failed event %d: %s", i, events[i].Content.String())
			nonAddedDocs++
		}
	}
	if nonAddedDocs > 0 {
		o.log.Warnf("failed to add %d documents", nonAddedDocs)
		dropped += nonAddedDocs
	}
	o.log.Infof("added %d documents", len(docs)-nonAddedDocs)

	o.observer.Dropped(dropped)
	batch.ACK()
	o.observer.Acked(len(events) - dropped)

	return nil
}
