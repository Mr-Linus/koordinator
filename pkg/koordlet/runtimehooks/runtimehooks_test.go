/*
 Copyright 2022 The Koordinator Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package runtimehooks

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_runtimeHook_Run(t *testing.T) {
	tmpDir := t.TempDir()
	type fields struct {
		config *Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			// grpcurl -plaintext localhost:9318 runtime.v1alpha1.RuntimeHookService/PreRunPodSandboxHook
			name: "run as tcp server",
			fields: fields{
				config: &Config{
					RuntimeHooksNetwork: "tcp",
					RuntimeHooksAddr:    ":0",
				},
			},
			wantErr: false,
		},
		{
			// grpcurl -plaintext -unix $file runtime.v1alpha1.RuntimeHookService/PreRunPodSandboxHook
			name: "run as unix socket",
			fields: fields{
				config: &Config{
					RuntimeHooksNetwork: "unix",
					RuntimeHooksAddr:    path.Join(tmpDir, "kooordlet.sock"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRuntimeHook(tt.fields.config)
			assert.NoError(t, err)
			stop := make(chan struct{})
			go func() { stop <- struct{}{} }()
			err = r.Run(stop)
			assert.NoError(t, err)
		})
	}
}
