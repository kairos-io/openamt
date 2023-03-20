package main

import (
	"encoding/json"
	"github.com/kairos-io/kairos-sdk/bus"
	"github.com/mudler/go-pluggable"
	"github.com/stretchr/testify/assert"
	"provider-amt/pkg/amtrpc"
	"rpc/pkg/utils"
	"testing"
)

var amtUnavailable = amtrpc.AMTRPC{
	MockAccessStatus: func() int { return utils.AmtNotDetected },
	MockExec:         func(s string) (string, int) { return "", utils.Success },
}

var amtAccessError = amtrpc.AMTRPC{
	MockAccessStatus: func() int { return utils.IncorrectPermissions },
	MockExec:         func(s string) (string, int) { return "", utils.Success },
}

var amtActive = amtrpc.AMTRPC{
	MockAccessStatus: func() int { return utils.Success },
	MockExec:         func(s string) (string, int) { return "", utils.Success },
}

var amtExecError = amtrpc.AMTRPC{
	MockAccessStatus: func() int { return utils.Success },
	MockExec:         func(s string) (string, int) { return "", utils.ActivationFailed },
}

func Test_activateAMTUnavailable(t *testing.T) {
	config := Configuration{
		AMT: &AMT{
			ServerAddress: "wss://fake",
		},
	}
	event := encodeConfiguration(config)

	resp := activateAMT(amtUnavailable, event)

	assert.Equal(t, StateUnavailable, resp.State)
	assert.Equal(t, event.Data, resp.Data)
}

func Test_activateAMTCheckAccessError(t *testing.T) {
	config := Configuration{
		AMT: &AMT{
			ServerAddress: "wss://fake",
		},
	}
	event := encodeConfiguration(config)

	resp := activateAMT(amtAccessError, event)

	assert.Equal(t, StateError, resp.State)
	assert.Equal(t, event.Data, resp.Data)
}

func Test_activateAMTNoConfiguration(t *testing.T) {
	config := Configuration{}
	event := encodeConfiguration(config)

	resp := activateAMT(amtActive, event)

	assert.Equal(t, StateSkipped, resp.State)
	assert.Equal(t, event.Data, resp.Data)
}

func Test_activateAMTInvalidEventData(t *testing.T) {
	config := Configuration{
		AMT: &AMT{
			ServerAddress: "wss://fake",
		},
	}
	event := encodeConfiguration(config)
	event.Data = event.Data[1:]

	resp := activateAMT(amtActive, event)

	assert.Equal(t, StateError, resp.State)
	assert.Equal(t, event.Data, resp.Data)
}

func Test_activateAMTInvalidConfiguration(t *testing.T) {
	event := &pluggable.Event{Data: `{"config":"{"}`}

	resp := activateAMT(amtActive, event)

	assert.Equal(t, StateError, resp.State)
	assert.Equal(t, event.Data, resp.Data)
}

func Test_activateAMTApplyError(t *testing.T) {
	config := Configuration{
		AMT: &AMT{
			ServerAddress: "wss://fake",
		},
	}
	event := encodeConfiguration(config)

	resp := activateAMT(amtExecError, event)

	assert.Equal(t, StateError, resp.State)
	assert.Equal(t, event.Data, resp.Data)
}

func Test_activateAMTStandard(t *testing.T) {
	var execCommand string

	config := Configuration{
		AMT: &AMT{
			ServerAddress: "wss://fake",
			Extra: map[string]string{
				"-foo": "bar",
			},
		},
	}
	event := encodeConfiguration(config)
	amt := amtrpc.AMTRPC{
		MockAccessStatus: func() int { return utils.Success },
		MockExec: func(s string) (string, int) {
			execCommand = s
			return "", utils.Success
		},
	}

	resp := activateAMT(amt, event)

	assert.Contains(t, execCommand, "activate")
	assert.Contains(t, execCommand, "-u "+config.AMT.ServerAddress)
	assert.Contains(t, execCommand, "-foo bar")
	assert.Equal(t, StateActive, resp.State)
	assert.Equal(t, event.Data, resp.Data)
	assert.False(t, resp.Errored())
}

func encodeConfiguration(config Configuration) *pluggable.Event {
	inner, _ := json.Marshal(config)
	data, _ := json.Marshal(bus.EventPayload{
		Config: string(inner),
	})

	return &pluggable.Event{
		Name: bus.EventInstall,
		Data: string(data),
	}
}
