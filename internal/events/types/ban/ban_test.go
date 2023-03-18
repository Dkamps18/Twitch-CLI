// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ban

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/test_setup"
)

var fromUser = "1234"
var toUser = "4567"

func TestEventSubBan(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	params := events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          models.TransportWebhook,
		Trigger:            "ban",
		SubscriptionStatus: "enabled",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err, "Error generating body.")

	var body models.BanEventSubResponse

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err, "Error unmarshalling JSON")

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal(fromUser, body.Event.UserID, "Expected from user %v, got %v", r.ToUser, body.Event.UserID)

	// test for unban
	params = events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          models.TransportWebhook,
		Trigger:            "unban",
		SubscriptionStatus: "enabled",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal(fromUser, body.Event.UserID, "Expected  from user %v, got %v", fromUser, body.Event.UserID)
}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "unban",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}

func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("ban")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("unban")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("notban")
	a.Equal(false, r)
}

func TestValidTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTransport(models.TransportWebhook)
	a.Equal(true, r)

	r = Event{}.ValidTransport("noteventsub")
	a.Equal(false, r)
}

func TestGetTopic(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.GetTopic(models.TransportWebhook, "ban")
	a.NotNil(r)
}
