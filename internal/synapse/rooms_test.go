/*
Copyright © 2026 Amanda Hager Lopes de Andrade Katz amandahla@gmail.com

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
package synapse

import (
	"context"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// MockClient satisfies the ClientInterface
type MockClient struct {
	// We use a map to store "path -> response" so we can handle multiple calls
	Responses map[string][]byte
	Errors    map[string]error
}

func (m *MockClient) Call(ctx context.Context, path string, method string, payload []byte, retry bool) ([]byte, error) {
	if err, ok := m.Errors[path]; ok && err != nil {
		return nil, err
	}
	if resp, ok := m.Responses[path]; ok {
		return resp, nil
	}
	return nil, fmt.Errorf("no mock response for path: %s", path)
}

func TestGetSpaces(t *testing.T) {
	cases := []struct {
		name      string
		responses map[string][]byte
		errors    map[string]error
		wantErr   bool
		wantLen   int
		wantName  string
		wantChild int
	}{
		{
			name: "single space, two children",
			responses: map[string][]byte{
				"/_matrix/client/v3/publicRooms":                     []byte(`{"chunk": [{"room_id": "!room1:xentonix.net", "name": "Ubuntu Community", "num_joined_members": 1500}]}`),
				"/_synapse/admin/v1/rooms/!room1:xentonix.net/state": []byte(`{"state": [{"type": "m.space.child", "state_key": "!child1:matrix.org"}, {"type": "m.space.child", "state_key": "!child2:matrix.org"}]}`),
			},
			wantErr:   false,
			wantLen:   1,
			wantName:  "Ubuntu Community",
			wantChild: 2,
		},
		{
			name:      "publicRooms error",
			responses: map[string][]byte{},
			errors: map[string]error{
				"/_matrix/client/v3/publicRooms": assert.AnError,
			},
			wantErr: true,
			wantLen: 0,
		},
		{
			name: "no spaces",
			responses: map[string][]byte{
				"/_matrix/client/v3/publicRooms": []byte(`{"chunk": []}`),
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "malformed state json",
			responses: map[string][]byte{
				"/_matrix/client/v3/publicRooms":                   []byte(`{"chunk": [{"room_id": "!room2:matrix.org", "name": "Test", "num_joined_members": 1}]}`),
				"/_synapse/admin/v1/rooms/!room2:matrix.org/state": []byte(`{"state": [}`),
			},
			wantErr: true,
			wantLen: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &MockClient{Responses: tc.responses, Errors: tc.errors}
			logger := logrus.New()
			spaces, err := GetSpaces(mock, logger)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantLen, len(spaces))
			if tc.wantLen > 0 {
				assert.Equal(t, tc.wantName, spaces[0].Name)
				assert.Equal(t, tc.wantChild, spaces[0].ChildCount)
			}
		})
	}
}
