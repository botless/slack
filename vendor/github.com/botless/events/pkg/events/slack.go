package events

import (
	"fmt"
	"log"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
)

const (
	slack_source_channel_template = "https://%s.slack.com/messages/%s/" // domain, channel
	slack_source_domain_template  = "https://%s.slack.com/"             // domain
	slack_type_template           = "botless.slack.%s"                  // type
)

var slackKnownEvents = []string{
	"welcome",
	"message",
	"latency",
}

type slack int

var Slack slack // export Slack
var _ = Slack

func (slack) Type(t string) string {
	if !contains(slackKnownEvents, t) {
		log.Printf("[WARN] unknown slack event type: %q", t)
	}
	return strings.ToLower(fmt.Sprintf(slack_type_template, t))
}

func (slack) SourceForDomain(domain string) types.URIRef {
	source := types.ParseURIRef(fmt.Sprintf(slack_source_domain_template, domain))
	if source == nil {
		return types.URIRef{}
	}
	return *source
}

func (slack) SourceForChannel(domain, channel string) types.URIRef {
	source := types.ParseURIRef(fmt.Sprintf(slack_source_channel_template, domain, channel))
	if source == nil {
		return types.URIRef{}
	}
	return *source
}
