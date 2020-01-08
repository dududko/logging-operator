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

package fluentd

import (
	"bytes"
	"fmt"
	"html/template"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type fluentdConfig struct {
	LogLevel string
	Monitor  struct {
		Enabled      bool
		Port         int32
		Path         string
		OutputLabels string
	}
}

func generateConfig(input fluentdConfig) (string, error) {
	output := new(bytes.Buffer)
	tmpl, err := template.New("test").Parse(fluentdInputTemplate)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template")
	}
	err = tmpl.Execute(output, input)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}
	outputString := fmt.Sprint(output.String())
	return outputString, nil
}

func (r *Reconciler) secretConfig() (runtime.Object, k8sutil.DesiredState, error) {
	input := fluentdConfig{Monitor: struct {
		Enabled      bool
		Port         int32
		Path         string
		OutputLabels string
	}{},
	}

	if r.Logging.Spec.FluentdSpec.Metrics != nil {
		input.Monitor.Enabled = true
		input.Monitor.Port = r.Logging.Spec.FluentdSpec.Metrics.Port
		input.Monitor.Path = r.Logging.Spec.FluentdSpec.Metrics.Path
		input.Monitor.OutputLabels = r.Logging.Spec.FluentdSpec.Metrics.OutputLabels
	}
	if r.Logging.Spec.FluentdSpec.LogLevel != "" {
		input.LogLevel = r.Logging.Spec.FluentdSpec.LogLevel
	} else {
		input.LogLevel = "info"
	}

	inputConfig, err := generateConfig(input)
	if err != nil {
		return nil, k8sutil.StatePresent, err
	}

	configs := &corev1.Secret{
		ObjectMeta: r.FluentdObjectMeta(SecretConfigName),
		Data: map[string][]byte{
			"fluent.conf":  []byte(fluentdDefaultTemplate),
			"input.conf":   []byte(inputConfig),
			"devnull.conf": []byte(fluentdOutputTemplate),
		},
	}

	configs.Data["fluentlog.conf"] = []byte(fmt.Sprintf(fluentLog, r.Logging.Spec.FluentdSpec.FluentLogDestination))

	return configs, k8sutil.StatePresent, nil
}
