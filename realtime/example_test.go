package realtime_test

import (
	"log"

	"github.com/aws/amazon-ivs-client-go/realtime"
)

// token is an IVS Stage participant token, created with the IVS Real-Time
// CreateParticipantToken API.
var token string

// hooks receives the Stage events this program cares about. Every other hook
// stays nil.
var hooks = realtime.StageHooks{
	OnConnectionStateChanged: func(event realtime.ConnectionStateChangedEvent) {
		if event.Err != nil {
			log.Printf("connection failed: %v", event.Err)
			return
		}
		log.Printf("connection: %s", event.State)
	},
	OnParticipantJoined: func(event realtime.ParticipantJoinedEvent) {
		log.Printf("participant joined: %s", event.Participant.ParticipantID)
	},
}

// audioSubscriber subscribes to every remote participant's audio and publishes
// nothing.
type audioSubscriber struct{}

func (audioSubscriber) SubscriptionTypes(realtime.ParticipantInfo) realtime.StreamType {
	return realtime.StreamTypeAudio
}

func (audioSubscriber) ShouldPublish(realtime.ParticipantInfo) bool {
	return false
}

func (audioSubscriber) StreamsToPublish(realtime.ParticipantInfo) realtime.PublishStreams {
	return realtime.PublishStreams{}
}

// Joining a Stage and receiving its mixed audio.
func Example() {
	stage, err := realtime.NewStage(realtime.StageConfig{Token: token, Strategy: audioSubscriber{}, Hooks: hooks})
	if err != nil {
		log.Fatal(err)
	}
	defer stage.Close()

	err = stage.SetPCMCallback(func(buf realtime.PCMBuffer) {
		// Called on an audio thread, so return promptly. buf is owned by this
		// function and may be retained.
		log.Printf("%d frames at %d Hz", buf.NumFrames, buf.SamplesPerSec)
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := stage.Join(); err != nil {
		log.Fatal(err)
	}
}

// subscribeStrategy subscribes to audio from hosts and publishes nothing.
type subscribeStrategy struct{}

func (subscribeStrategy) SubscriptionTypes(participant realtime.ParticipantInfo) realtime.StreamType {
	if participant.Attributes["role"] == "host" {
		return realtime.StreamTypeAudio
	}
	return realtime.StreamTypeNone
}

func (subscribeStrategy) ShouldPublish(realtime.ParticipantInfo) bool {
	return false
}

func (subscribeStrategy) StreamsToPublish(realtime.ParticipantInfo) realtime.PublishStreams {
	return realtime.PublishStreams{}
}

// Answering the strategy decisions yourself.
func ExampleStageStrategy() {
	stage, err := realtime.NewStage(realtime.StageConfig{Token: token, Strategy: subscribeStrategy{}, Hooks: hooks})
	if err != nil {
		log.Fatal(err)
	}
	defer stage.Close()

	if err := stage.Join(); err != nil {
		log.Fatal(err)
	}
}
