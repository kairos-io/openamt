package main

import (
	"encoding/json"
	"fmt"
	"github.com/kairos-io/kairos-sdk/bus"
	"github.com/mudler/go-pluggable"
	"github.com/sirupsen/logrus"
	"os"
	"provider-amt/pkg/amtrpc"
	"reflect"
	"rpc/pkg/utils"
	"strings"
)

const (
	StateActive      = "active"
	StateError       = "error"
	StateUnavailable = "unavailable"
	StateSkipped     = "skipped"
)

type AMT struct {
	DnsSuffixOverride string            `json:"dns_suffix_override,omitempty" flag:"-d"`
	Hostname          string            `json:"hostname,omitempty" flag:"-h"`
	LMSAddress        string            `json:"lms_address,omitempty" flag:"-lmsaddress"`
	LMSPort           string            `json:"lms_port,omitempty" flag:"-lmsport"`
	ProxyAddress      string            `json:"proxy_address,omitempty" flag:"-p"`
	Password          string            `json:"password,omitempty" flag:"-password"`
	Profile           string            `json:"profile,omitempty" flag:"-profile"`
	ServerAddress     string            `json:"server_address,omitempty" flag:"-u"`
	Timeout           string            `json:"timeout,omitempty" flag:"-t"`
	Extra             map[string]string `json:"extra,omitempty"`
}

type Configuration struct {
	AMT *AMT `yaml:"amt,omitempty" json:"amt,omitempty"`
}

func main() {
	err := pluggable.NewPluginFactory(
		pluggable.FactoryPlugin{
			EventType: bus.EventInstall,
			PluginHandler: func(event *pluggable.Event) pluggable.EventResponse {
				return activateAMT(amtrpc.AMTRPC{}, event)
			},
		},
	).Run(pluggable.EventType(os.Args[1]), os.Stdin, os.Stdout)

	if err != nil {
		logrus.Fatal(err)
	}
}

func getConfiguration(event *pluggable.Event) (*AMT, error) {
	var payload bus.EventPayload
	var config Configuration

	// parse the event to get the configuration
	if err := json.Unmarshal([]byte(event.Data), &payload); err != nil {
		return nil, fmt.Errorf("failed to parse event payload  %s", err.Error())
	}

	// parse the configuration to get the amt configuration
	if err := json.Unmarshal([]byte(payload.Config), &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration  %s", err.Error())
	}

	return config.AMT, nil
}

func activateAMT(rpc amtrpc.Interface, event *pluggable.Event) pluggable.EventResponse {
	config, err := getConfiguration(event)
	if err != nil {
		return pluggable.EventResponse{
			State: StateError,
			Data:  event.Data,
			Error: err.Error(),
		}
	}

	// if no amt configuration is given we do nothing
	if config == nil {
		return pluggable.EventResponse{
			State: StateSkipped,
			Data:  event.Data,
		}
	}

	if status := rpc.CheckAccess(); status != utils.Success {
		if status == utils.AmtNotDetected {
			return pluggable.EventResponse{
				State: StateUnavailable,
				Data:  event.Data,
				Logs:  "no intel AMT device detected",
			}
		}

		return pluggable.EventResponse{
			State: StateError,
			Data:  event.Data,
			Error: fmt.Sprintf("failed to access AMT device with status code: %d", status),
		}
	}

	if response, status := rpc.Exec(toCLIFlags("activate", config, config.Extra)); status != utils.Success {
		return pluggable.EventResponse{
			State: StateError,
			Data:  event.Data,
			Error: fmt.Sprintf("failed to activate AMT device with status code: %d", status),
			Logs:  response,
		}
	}

	return pluggable.EventResponse{
		State: StateActive,
		Data:  event.Data,
	}
}

const flagTag = "flag"

func toCLIFlags(command string, config any, extra map[string]string) string {
	t := reflect.TypeOf(config).Elem()
	val := reflect.ValueOf(config).Elem()
	for i := 0; i < t.NumField(); i++ {
		flag := t.Field(i).Tag.Get(flagTag)
		fval := val.Field(i).String()

		if flag == "" || val.Field(i).IsZero() {
			continue
		}

		command = strings.Join([]string{command, flag, fval}, " ")
	}

	if extra != nil {
		for k, v := range extra {
			command = strings.Join([]string{command, k, v}, " ")
		}
	}

	return command
}
