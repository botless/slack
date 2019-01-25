package cloudevents

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/knative/pkg/cloudevents"
	"net/http"
	"time"
)

type CloudEventEncoding int

const (
	Binary     CloudEventEncoding = 0
	Structured CloudEventEncoding = 1
)

type Builder struct {
	Source    string
	EventType string
	Encoding  CloudEventEncoding
}

func (b *Builder) Build(target string, data interface{}) (*http.Request, error) {
	if b.Source == "" {
		return nil, fmt.Errorf("Build.Source is empty")
	}
	if b.EventType == "" {
		return nil, fmt.Errorf("Build.EventType is empty")
	}

	ctx := b.cloudEventsContext()

	switch b.Encoding {
	case Binary:
		return cloudevents.Binary.NewRequest(target, data, ctx)
	case Structured:
		return cloudevents.Binary.NewRequest(target, data, ctx)
	default:
		return nil, fmt.Errorf("unsupported encoding: %v", b.Encoding)
	}
}

// Creates a CloudEvent Context for a given heartbeat.
func (b *Builder) cloudEventsContext() cloudevents.EventContext {
	return cloudevents.EventContext{
		CloudEventsVersion: cloudevents.CloudEventsVersion,
		EventType:          b.EventType,
		EventID:            uuid.New().String(),
		Source:             b.Source,
		ContentType:        "application/json",
		EventTime:          time.Now(),
	}
}
