// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output

import (
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
)

// +docName:"LogDNA output plugin for Fluentd"
//This plugin has been designed to output logs to LogDNA.
//More info at https://github.com/logdna/fluent-plugin-logdna
//>Example Deployment: [Transport Nginx Access Logs into LogDNA with Logging Operator](../../../docs/example-logdna-nginx.md)
//
// #### Example output configurations
// ```
// spec:
//  logdna:
//    api_key: xxxxxxxxxxxxxxxxxxxxxxxxxxx
//    hostname: logging-operator
//    app: my-app
// ```
type _docLogDNA interface{}

// +kubebuilder:object:generate=true
// +docName:"LogDNA"
// Send your logs to LogDNA
type LogDNAOutput struct {
	// LogDNA Api key
	ApiKey string `json:"api_key"`
	// Hostname
	HostName string `json:"hostname"`
	// Application name
	App string `json:"app,omitempty"`
	// Do not increase past 8m (8MB) or your logs will be rejected by LogDNA server.
	BufferQueueLimit string `json:"buffer_chunk_limit,omitempty" default:"1m"`
}

func (l *LogDNAOutput) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "logdna"
	pluginID := id + "_" + pluginType
	logdna := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        pluginID,
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(l); err != nil {
		return nil, err
	} else {
		logdna.Params = params
	}
	return logdna, nil
}
