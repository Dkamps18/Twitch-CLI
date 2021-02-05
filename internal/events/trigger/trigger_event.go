// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"fmt"
	"time"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/events/types"
	"github.com/twitchdev/twitch-cli/internal/util"
)

// TriggerParameters defines the parameters used to emit an event.
type TriggerParameters struct {
	Event          string
	Transport      string
	IsAnonymous    bool
	FromUser       string
	ToUser         string
	GiftUser       string
	Status         string
	ItemId         string
	Cost           int64
	ForwardAddress string
	Secret         string
	Verbose        bool
	Count          int
}

type TriggerResponse struct {
	ID        string
	JSON      []byte
	FromUser  string
	ToUser    string
	Timestamp string
}

// Fire emits an event using the TriggerParameters defined above.
func Fire(p TriggerParameters) (string, error) {
	var resp events.MockEventResponse
	var err error

	if p.ToUser == "" {
		p.ToUser = util.RandomUserID()
	}

	if p.FromUser == "" {
		p.FromUser = util.RandomUserID()
	}

	eventParamaters := events.MockEventParameters{
		ID:           util.RandomGUID(),
		Trigger:      p.Event,
		Transport:    p.Transport,
		FromUserID:   p.FromUser,
		FromUserName: "testFromUser",
		ToUserID:     p.ToUser,
		ToUserName:   "testBroadcaster",
		IsAnonymous:  p.IsAnonymous,
		Cost:         p.Cost,
		Status:       p.Status,
	}

	e, err := types.GetByTriggerAndTransport(p.Event, p.Transport)
	if err != nil {
		return "", err
	}

	resp, err = e.GenerateEvent(eventParamaters)
	if err != nil {
		return "", err
	}

	err = util.InsertIntoDB(util.EventCacheParameters{
		ID:        resp.ID,
		Event:     p.Event,
		JSON:      string(resp.JSON),
		FromUser:  resp.FromUser,
		ToUser:    resp.ToUser,
		Transport: p.Transport,
		Timestamp: util.GetTimestamp().Format(time.RFC3339Nano),
	})
	if err != nil {
		return "", err
	}

	if p.ForwardAddress != "" {
		resp, err := ForwardEvent(ForwardParamters{
			ID:             resp.ID,
			Transport:      p.Transport,
			JSON:           resp.JSON,
			Secret:         p.Secret,
			ForwardAddress: p.ForwardAddress,
			Event:          p.Event,
		})
		defer resp.Body.Close()

		if err != nil {
			return "", err
		}

		println(fmt.Sprintf(`[%v] Request Sent`, resp.StatusCode))
	}

	return string(resp.JSON), nil
}
