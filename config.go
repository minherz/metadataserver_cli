package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/minherz/metadataserver"
)

const undefinedStringFlag = "_undefined_"

type metadataTuple struct {
	path  string
	value string
}

func (v metadataTuple) String() string {
	return fmt.Sprintf("%s=%s", v.path, v.value)
}

// MetadataFlagSlice defines a custom type that implements the flag.Value interface
// to store multiple values of the metadata configuration
type MetadataFlagSlice []metadataTuple

func (v *MetadataFlagSlice) String() string {
	s := []string{}
	for _, e := range *v {
		s = append(s, e.String())
	}
	return strings.Join(s, ", ")
}

func (v *MetadataFlagSlice) Set(value string) error {
	s := strings.Split(value, "=")
	if len(s) != 2 {
		return errors.New("cannot parse metadata tuple: " + value)
	}
	*v = append(*v, metadataTuple{
		path:  s[0],
		value: s[1],
	})
	return nil
}

func ConfigOptions() ([]metadataserver.Option, error) {
	file := flag.String("config-file", undefinedStringFlag, "Path to JSON configuration file")
	address := flag.String("a", undefinedStringFlag, "IP address to serve requests")
	port := flag.Int("p", metadataserver.DefaultPort, "Port to listen for requests")
	endpoint := flag.String("endpoint", undefinedStringFlag, "Endpoint prefix")
	var (
		metadataValues  MetadataFlagSlice
		metadataEnvVars MetadataFlagSlice
	)
	flag.Var(&metadataValues, "", "A metadata path that returns a literal")
	flag.Var(&metadataEnvVars, "", "A metadata path that returns value of an environment variable")

	flag.Parse()

	if *file != undefinedStringFlag &&
		(*address != undefinedStringFlag || *port != 0 || *endpoint != undefinedStringFlag || len(metadataValues) > 0 || len(metadataEnvVars) > 0) {
		return nil, errors.New("'-config-file' flag cannot be use with other flags")
	}
	if *port < 1024 && *port != 80 {
		return nil, errors.New("'-p' flag cannot have negative value or a value in range 0-1024 other than 80")
	}

	handlers := map[string]metadataserver.Metadata{}
	for _, t := range metadataValues {
		if _, ok := handlers[t.path]; ok {
			return nil, errors.New("metadata path can be defined only once. '" + t.path + "' is found defined multiple times")
		}
		v := func() string {
			return t.value
		}
		handlers[t.path] = metadataserver.Metadata(v)
	}
	for _, t := range metadataEnvVars {
		if _, ok := handlers[t.path]; ok {
			return nil, errors.New("metadata path can be defined only once. '" + t.path + "' is found defined multiple times")
		}
		v := func() string {
			return os.Getenv(t.value)
		}
		handlers[t.path] = metadataserver.Metadata(v)
	}

	ops := []metadataserver.Option{}
	if *file != undefinedStringFlag {
		return append(ops, metadataserver.WithConfigFile(*file)), nil
	}
	if *address != undefinedStringFlag {
		ops = append(ops, metadataserver.WithAddress(*address))
	}
	if *port > 0 {
		ops = append(ops, metadataserver.WithPort(*port))
	}
	if *endpoint != undefinedStringFlag {
		ops = append(ops, metadataserver.WithEndpoint(*endpoint))
	}
	if len(handlers) > 0 {
		ops = append(ops, metadataserver.WithHandlers(handlers))
	}
	return ops, nil
}
