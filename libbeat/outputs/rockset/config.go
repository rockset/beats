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
	"fmt"

	"github.com/elastic/elastic-agent-libs/config"
	"github.com/prometheus/tsdb/errors"

	"github.com/elastic/beats/v7/libbeat/outputs/codec"
)

type Config struct {
	Codec codec.Config `config:"codec"`

	Workspace  string `config:"workspace"`
	Collection string `config:"collection"`
	APIKey     string `config:"apikey"`
	APIServer  string `config:"apiserver"`

	BatchSize int
	Queue     config.Namespace `config:"queue"`
}

var defaultConfig = Config{}

func (c *Config) Validate() error {
	var err errors.MultiError

	if c.Collection == "" {
		err.Add(fmt.Errorf("must specify Rockset colletion"))
	}
	if c.Workspace == "" {
		err.Add(fmt.Errorf("must specify Rockset workspace"))
	}
	if c.APIKey == "" {
		err.Add(fmt.Errorf("must specify Rockset API key"))
	}
	if c.APIServer == "" {
		err.Add(fmt.Errorf("must specify Rockset API server"))
	}

	return err.Err()
}
