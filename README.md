# Slack Source

This is a _work in progress_ Slack source using Knative Eventing Sources to 
bridge Slack events via a Slack App/Bot to a Kubernetes Cluster running the 
Knative Eventing stack. 

You will need [Knative](https://github.com/knative/docs/) and 
[ko](https://github.com/google/go-containerregistry/tree/master/cmd/ko)
setup and installed with a running k8s cluster.

## Quick Start

1. Get Slack App API keys in the Slack Admin Console. 
    1. Replace the _base64_ encoded keys in `config/once/secret.yaml`
    1. Then, `kubectl apply -f config/once/secret.yaml`   
2. Install the Slack Source, `ko apply -f config/`

Optionally install one of the demos inside of `./config/demo` or write your own.

## Overview

Slack Source creates a Container Source running a RTM connection to Slack using
the provided API keys. Applying `config/slack.yaml` also creates 2 channels and 
a subscription for the Slack Source to receive messages from within the cluster.

### Incoming events

```
Slack -> Source -> channel/slack-in
``` 

The Slack Source only forwards cloudevents of type slack.Message to `slack-in`.

### Outgoing Events

```
channel/slack-out -> Source -> Slack
``` 
The Slack Source is only looking for CloudEvents with
[Response](./pkg/events/response.go) objects.

```go
type Response struct {
	Channel string
	Text    string
}
```

## Demos

The receiver demo just prints to the pod log.

The echo demo looks for slack.Message events and if the message.Text starts with
`"echo "` then sends a new cloudevent to `slack-out` with the same text with
`"echo "` trimmed off to the same channel as it was sent.
