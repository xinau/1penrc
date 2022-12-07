package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/xinau/1penrc/internal/config"
	"github.com/xinau/1penrc/internal/op"
	"github.com/xinau/1penrc/internal/provider"
	"github.com/xinau/1penrc/internal/provider/awssts"
	"github.com/xinau/1penrc/internal/provider/secret"
	"github.com/xinau/1penrc/internal/provider/value"
)

var (
	configF = flag.String(
		"config",
		filepath.Join(MustUserConfigDir(), "1penrc", "config.yaml"),
		"configuration file to load",
	)
)

var (
	EnvironmentNotFoundError = errors.New("couldn't find environment")
)

var (
	ShellVariableExportTemplate = template.Must(template.New("").Parse(`
{{- range $key, $val := . -}}
export {{ $key }}='{{ $val }}'
{{ end -}}
`))
)

func MustUserConfigDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return dir
}

func FindEnvironmentConfigByName(name string, cfgs []*config.EnvironmentConfig) (*config.EnvironmentConfig, error) {
	for _, cfg := range cfgs {
		if cfg.Name == name {
			return cfg, nil
		}
	}
	return nil, fmt.Errorf("%w with name %q", EnvironmentNotFoundError, name)
}

func GetVariablesFromEnvironmentConfig(client *op.Client, cfg *config.EnvironmentConfig) (provider.Variables, error) {
	vars := make(provider.Variables)
	for _, pcfg := range cfg.AWSSTSConfigs {
		other, err := awssts.GetVariables(client, pcfg)
		if err != nil {
			return nil, err
		}
		vars.Merge(other)
	}
	for _, pcfg := range cfg.SecretConfigs {
		other, err := secret.GetVariables(client, pcfg)
		if err != nil {
			return nil, err
		}
		vars.Merge(other)
	}
	for _, pcfg := range cfg.ValueConfigs {
		other, err := value.GetVariables(client, pcfg)
		if err != nil {
			return nil, err
		}
		vars.Merge(other)
	}
	return vars, nil
}

var (
	signinF = flag.Bool("signin", false, "Sign in to a 1Password account.")
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("fatal: no environment argument provided")
	}

	cfg, err := config.LoadFile(*configF)
	if err != nil {
		log.Fatalf("fatal: loading configuration %q: %s", *configF, err)
	}

	client := op.NewClient(cfg.ClientConfig)

	ecfg, err := FindEnvironmentConfigByName(flag.Args()[0], cfg.EnvironmentConfigs)
	if err != nil {
		log.Fatalf("fatal: %s", err)
	}

	if *signinF {
		if err := client.SignIn(ecfg.Account); err != nil {
			log.Fatalf("fatal: signing into %s: %s", ecfg.Account, err)
		}
	}

	vars, err := GetVariablesFromEnvironmentConfig(client, ecfg)
	if err != nil {
		log.Fatalf("fatal: %s", err)
	}

	if err := ShellVariableExportTemplate.Execute(os.Stdout, vars); err != nil {
		log.Fatalf("fatal: %s", err)
	}
}
