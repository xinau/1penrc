package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/xinau/1penrc/internal/config"
)

var (
	MissingEnvironmentArgError = errors.New("missing required environment argument")
)

type Command struct {
	fs *flag.FlagSet

	ConfigFile  string
	Environment string
}

func NewCommand() *Command {
	cmd := &Command{
		fs: flag.NewFlagSet("", flag.ContinueOnError),
	}
	cmd.fs.Usage = func() {
		fmt.Fprintf(cmd.fs.Output(), "Usage of %s [flags] environment\n", os.Args[0])
		cmd.fs.PrintDefaults()
	}
	cmd.fs.StringVar(&cmd.ConfigFile, "config", "", "Location of config file defining environments")

	return cmd
}

func (c *Command) GetConfigFile() (string, error) {
	if c.ConfigFile != "" {
		return c.ConfigFile, nil
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "1penrc.hcl"), nil
}

func (c *Command) Parse(args []string) error {
	err := c.fs.Parse(args)
	if err != nil {
		return fmt.Errorf("Parsing flags: %s", err)
	}

	args = c.fs.Args()
	if len(args) != 1 {
		return fmt.Errorf("Parsing arguments: %s", MissingEnvironmentArgError)
	}
	c.Environment = args[0]

	return nil
}

func (c *Command) Run() error {
	file, err := c.GetConfigFile()
	if err != nil {
		return fmt.Errorf("locating configuration file: %s", err)
	}

	cfg, err := config.LoadFile(file)
	if err != nil {
		return fmt.Errorf("loading configuration file: %s", err)
	}

	env, err := cfg.FindEnvironment(c.Environment)
	if err != nil {
		return fmt.Errorf("finding environment: %s", err)
	}

	vars, err := env.Export()
	if err != nil {
		return fmt.Errorf("getting variables for environment: %s", err)
	}

	var tmpl = template.Must(template.New("").Parse(`
{{- range $key, $val := . }}
export {{ $key }}='{{ $val }}'
{{- end }}
`))

	if err := tmpl.Execute(os.Stdout, vars); err != nil {
		return fmt.Errorf("rendering source output: %s", err)
	}

	return nil
}

func main() {
	cmd := NewCommand()
	if err := cmd.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(cmd.fs.Output(), "[ERROR] %s\n", err)
		cmd.fs.Usage()
		os.Exit(1)
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}
}
