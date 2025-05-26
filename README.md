<!-- markdownlint-disable MD033 -->
# Metadata Server CLI

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/minherz/metadataserver_cli)](https://github.com/minherz/metadataserver_cli/releases)
[![Build](https://github.com/minherz/metadataserver_cli/actions/workflows/go.yaml/badge.svg)](https://github.com/minherz/metadataserver_cli/actions/workflows/go.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/minherz/metadataserver_cli)](https://goreportcard.com/report/github.com/minherz/metadataserver_cli)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/minherz/metadataserver_cli)](https://pkg.go.dev/github.com/minherz/metadataserver_cli)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/minherz/metadataserver_cli?label=go-version)
![Repo license](https://badgen.net/badge/license/Apache%202.0/blue)

This is a command line interface for the [metadataserver](https://github.com/minherz/metadataserver) that allows to run a metadata server emulator from a command line.
You can also run it as a sidecar container in your [local K8s or Docker environments](http://about).

## Download instructions

You can download CLI from the [releases page](https://github.com/minherz/metadataserver/releases) by expanding "Assets" and picking a binary for Win/Arm/x86.
You can also use a container image at us-docker.pkg.dev/minherz/examples/metadataserver. You can use a tag matching the release version to pull an images of the metadata server for that release.

## Using CLI

CLI runs an instance of the metadata server and gracefully terminates on receiving SIGINT or a user pressing <kbd>Ctrl</kbd> + <kbd>C</kbd>.
CLI allows to customize its launch and metadata configuration using flags. A flag's format supports both `=` and space delimetered syntax.
The following command will launch a server listening for requests using local host interface at port 8080:

```console
metadataserver_cli -a=0.0.0.0 -p 8080
```

### Full list of flags

| Flag | Default | Description |
| ---- | ------- | ----------- |
| `-a` | `169.254.169.254` | IP address of the server. |
| `-p` | `80` | Port number at which the server listens. |
| `-endpoint` | `computeMetadata/v1` | Endpoint prefix for all serving requests. Heading and trailing slashes can be omitted. |
| `-config-file` | _None_ | Path to JSON configuration file. It cannot be used with other flags. See [file format](https://github.com/minherz/metadataserver#custom-configuration) for JSON schema. |
| `-metadata` | _None_ | A pair in the form `"path/to/metadata"=VALUE` where `VALUE` is a literal to be returned when a request is sent to the '{endpoint}/path/to/metadata'. This flag can be repreated to define multiple metadata paths. |
| `-metadata-env` | _None_ | A pair in the form `"path/to/metadata"=ENV_NAME` where `ENV_NAME` is a name of the environment variable that stores the value of the metadata. The value is read on each request is sent to the '{endpoint}/path/to/metadata'. This flag can be repreated to define multiple metadata paths. |
