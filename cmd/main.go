package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"openamt/pkg/amtrpc"
	"os"
	"reflect"
	"rpc/pkg/utils"
	"strings"
)

type AMT struct {
	DnsSuffixOverride string            `json:"dns_suffix_override,omitempty" yaml:"dns_suffix_override,omitempty" flag:"-d"`
	Hostname          string            `json:"hostname,omitempty" yaml:"hostname,omitempty" flag:"-h"`
	LMSAddress        string            `json:"lms_address,omitempty" yaml:"lms_address,omitempty" flag:"-lmsaddress"`
	LMSPort           string            `json:"lms_port,omitempty" yaml:"lms_port,omitempty" flag:"-lmsport"`
	ProxyAddress      string            `json:"proxy_address,omitempty" yaml:"proxy_address,omitempty" flag:"-p"`
	Password          string            `json:"password,omitempty" yaml:"password,omitempty" flag:"-password"`
	Profile           string            `json:"profile,omitempty" yaml:"profile,omitempty" flag:"-profile"`
	ServerAddress     string            `json:"server_address,omitempty" yaml:"server_address,omitempty" flag:"-u"`
	Timeout           string            `json:"timeout,omitempty" yaml:"timeout,omitempty" flag:"-t"`
	Extra             map[string]string `json:"extra,omitempty" yaml:"extra,omitempty"`
}

func main() {
	var config AMT

	if err := yaml.NewDecoder(io.TeeReader(os.Stdin, os.Stderr)).Decode(&config); err != nil {
		logrus.Fatal("failed to parse configuration: ", err)
	}

	if err := activateAMT(amtrpc.AMTRPC{}, &config); err != nil {
		logrus.Fatal()
	}
}

func activateAMT(rpc amtrpc.Interface, config *AMT) error {
	if status := rpc.CheckAccess(); status != utils.Success {
		if status == utils.AmtNotDetected {
			logrus.Info("amt device could not be detected, skipping configuration")
			return nil
		}

		return fmt.Errorf("failed to access AMT device with status code: %d", status)
	}

	if response, status := rpc.Exec(toCLIFlags("activate", config, config.Extra)); status != utils.Success {
		return fmt.Errorf("failed to access AMT device with status code: %d %s", status, response)
	}

	return nil
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
