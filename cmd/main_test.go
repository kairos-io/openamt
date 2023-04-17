package main

import (
	"github.com/stretchr/testify/assert"
	"openamt/pkg/amtrpc"
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
	config := &AMT{
		ServerAddress: "wss://fake",
	}

	err := activateAMT(amtUnavailable, config)

	assert.NoError(t, err)
}

func Test_activateAMTCheckAccessError(t *testing.T) {
	config := &AMT{
		ServerAddress: "wss://fake",
	}

	err := activateAMT(amtAccessError, config)

	assert.Error(t, err)
}

func Test_activateAMTNoConfiguration(t *testing.T) {
	err := activateAMT(amtActive, &AMT{})

	assert.NoError(t, err)
}

func Test_activateAMTApplyError(t *testing.T) {
	config := &AMT{
		ServerAddress: "wss://fake",
	}

	err := activateAMT(amtExecError, config)

	assert.Error(t, err)
}

func Test_activateAMTStandard(t *testing.T) {
	var execCommand string

	config := &AMT{
		ServerAddress: "wss://fake",
		Extra: map[string]string{
			"-foo": "bar",
		},
	}
	amt := amtrpc.AMTRPC{
		MockAccessStatus: func() int { return utils.Success },
		MockExec: func(s string) (string, int) {
			execCommand = s
			return "", utils.Success
		},
	}

	err := activateAMT(amt, config)

	assert.Contains(t, execCommand, "activate")
	assert.Contains(t, execCommand, "-u "+config.ServerAddress)
	assert.Contains(t, execCommand, "-foo bar")
	assert.NoError(t, err)
}
