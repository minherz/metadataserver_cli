package main

import (
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	metadataTupleCmpOtp = cmpopts.EquateComparable(metadataTuple{})
)

func TestMetadataFlagSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  MetadataFlagSlice
	}{
		{
			name:  "single_flag",
			input: []string{"path1=value1"},
			want: []metadataTuple{
				{path: "path1", value: "value1"},
			},
		},
		{
			name:  "multiple_flags",
			input: []string{"path1=value1", "path2=value2"},
			want: []metadataTuple{
				{path: "path1", value: "value1"},
				{path: "path2", value: "value2"},
			},
		},
		{
			name:  "valid_values",
			input: []string{"this/is/valid/path1=value1", "/this/is/valid/path2=value2", "this/is/valid/path3/=value3"},
			want: []metadataTuple{
				{path: "this/is/valid/path1", value: "value1"},
				{path: "/this/is/valid/path2", value: "value2"},
				{path: "this/is/valid/path3/", value: "value3"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got MetadataFlagSlice
			for _, s := range test.input {
				err := got.Set(s)
				if err != nil {
					t.Errorf("MetadataFlagSlice parsing returned error:\n%v", err)
				}
			}
			if diff := cmp.Diff(test.want, got, metadataTupleCmpOtp); diff != "" {
				t.Errorf("MetadataFlagSlice parsing mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMetadataInvalidFlagSlice(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{
			name:  "no_error",
			input: "key1=value1",
			want:  nil,
		},
		{
			name:  "key_with_spaces",
			input: "key with spaces=value",
			want:  MetadataFlagSliceKeyError("key with spaces"),
		},
		{
			name:  "key_with_trailing_space",
			input: "key/ends/with/space =value",
			want:  MetadataFlagSliceKeyError("key/ends/with/space "),
		},
		{
			name:  "key_with_escape_chars",
			input: "key/with/%e0/scapes=value",
			want:  MetadataFlagSliceKeyError("key/with/%e0/scapes"),
		},
		{
			name:  "key_with_colon",
			input: "key/with:colon=value",
			want:  MetadataFlagSliceKeyError("key/with:colon"),
		},
		{
			name:  "key_with_double_slash",
			input: "key/with//doubleslash=value",
			want:  MetadataFlagSliceKeyError("key/with//doubleslash"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var s MetadataFlagSlice
			got := s.Set(test.input)
			if got != nil && test.want != nil && got.Error() != test.want.Error() {
				t.Errorf("MetadataFlagSlice validation mismatch (want:%v got:%v", test.want, got)
			}
		})
	}
}

func TestInvalidConfigs(t *testing.T) {
	tests := []struct {
		name      string
		arguments []string
	}{
		{
			name:      "negative_port",
			arguments: []string{"test -p=-35", "-p=-35"},
		},
		{
			name:      "port_under_1024_range",
			arguments: []string{"test -p=35", "-p=35"},
		},
		{
			name: "config_with_file_and_endpoint",
			arguments: []string{
				"test -config-file=path/to/config_file.json -endpoint=metadata/v2",
				"-config-file=path/to/config_file.json",
				"-endpoint=metadata/v2",
			},
		},
		{
			name: "config_with_file_and_addr",
			arguments: []string{
				"test -config-file=path/to/config_file.json -a=69.39.69.39",
				"-config-file=path/to/config_file.json",
				"-a=69.39.69.39",
			},
		},
		{
			name: "config_with_file_and_port",
			arguments: []string{
				"test -config-file=path/to/config_file.json -p=4455",
				"-config-file=path/to/config_file.json",
				"-p=4455",
			},
		},
		{
			name: "config_with_file_and_metadata1",
			arguments: []string{
				"test -config-file=path/to/config_file.json -metadata=path1=value1",
				"-config-file=path/to/config_file.json",
				"-metadata=path1=value1",
			},
		},
		{
			name: "config_with_file_and_metadata2",
			arguments: []string{
				"test -config-file=path/to/config_file.json -metadata-env=path1=PROJECT_ID",
				"-config-file=path/to/config_file.json",
				"-metadata-env=path1=PROJECT_ID",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Args = test.arguments
			flag.CommandLine = flag.NewFlagSet(test.arguments[0], flag.ExitOnError)
			_, err := ConfigOptions()
			if err == nil {
				t.Errorf("error is missing")
			}
		})
	}
}
