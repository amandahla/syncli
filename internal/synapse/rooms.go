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
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"maunium.net/go/mautrix"
	mauevent "maunium.net/go/mautrix/event"
)

type StateResponse struct {
	State []mauevent.Event `json:"state"`
}

type Space struct {
	ID         string
	Name       string
	Members    int
	ChildCount int
	ChildRooms []string
}

func (s Space) Header() []string {
	return []string{"Name", "Members", "Child Count", "Child Rooms"}
}

func (s Space) Row() []interface{} {
	return []interface{}{s.Name, s.Members, s.ChildCount, strings.Join(s.ChildRooms, ",")}
}

const maxConcurrentRequests = 10
const maxConcurrentRequestsTimeout = 5 * time.Minute

func GetSpaces(client SynapseClientInterface, logger *logrus.Logger) ([]Space, error) {
	var spaces []Space
	ctx, cancel := context.WithTimeout(context.Background(), maxConcurrentRequestsTimeout) // Set a timeout for the entire operation to avoid hanging indefinitely in case of issues with the server
	defer cancel()
	payload := []byte(`{"limit": 200, "filter": {"room_types": ["m.space"]}}`)
	output, err := client.Call(ctx, "/_matrix/client/v3/publicRooms", "POST", payload, false)
	if err != nil {
		return spaces, err
	}

	spaces, err = parseSpaces(output)
	if err != nil {
		return spaces, err
	}

	g, ctx := errgroup.WithContext(ctx) // errgroup allows us to wait for all goroutines to finish and captures the first error that occurs
	var mu sync.Mutex
	sem := make(chan struct{}, maxConcurrentRequests)

	logger.WithFields(logrus.Fields{
		"event":                           "fetching_space_details",
		"count":                           len(spaces),
		"max_concurrent_requests":         maxConcurrentRequests,
		"max_concurrent_requests_timeout": maxConcurrentRequestsTimeout,
	}).Debug("Fetching details for spaces")

	for i := range spaces {
		i := i // no need on Go 1.21+ with the new for loop variable scoping, but it's a good practice to avoid bugs in older versions

		g.Go(func() error {
			// Check if another goroutine already failed
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return ctx.Err()
			}

			defer func() { <-sem }()

			logger.WithFields(logrus.Fields{
				"event": "fetching_space_details",
				"space": spaces[i].ID,
			}).Debug("Fetching details for space")
			output, err := client.Call(ctx, "/_synapse/admin/v1/rooms/"+spaces[i].ID+"/state", "GET", nil, false)
			if err != nil {
				return err
			}

			var resp StateResponse
			if err := json.Unmarshal(output, &resp); err != nil {
				return err
			}

			childCount := 0
			childRooms := make([]string, 0)

			for _, event := range resp.State {
				if event.Type.String() == "m.space.child" && event.StateKey != nil {
					childCount++
					childRooms = append(childRooms, *event.StateKey)
				}
			}

			mu.Lock()
			spaces[i].ChildCount += childCount
			spaces[i].ChildRooms = append(spaces[i].ChildRooms, childRooms...)
			mu.Unlock()

			return nil
		})
	}

	// Wait blocks until all calls finish OR one returns an error
	if err := g.Wait(); err != nil {
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"event": "fetched_space_details",
		"count": len(spaces),
	}).Debug("Fetched details for spaces")

	return spaces, nil
}

func parseSpaces(data []byte) ([]Space, error) {
	var result mautrix.RespPublicRooms
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	var spaces []Space
	for _, room := range result.Chunk {
		spaces = append(spaces, Space{
			ID:         room.RoomID.String(),
			Name:       room.Name,
			Members:    room.NumJoinedMembers,
			ChildCount: 0,
			ChildRooms: []string{},
		})
	}
	return spaces, nil
}
