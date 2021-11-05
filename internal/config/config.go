package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"

	"github.com/xinau/1penrc/internal/env"
	"github.com/xinau/1penrc/internal/env/static"
)

type Config struct {
	EnvironmentConfigs []EnvironmentConfig `hcl:"env,block"`
}

type EnvironmentConfig struct {
	Name        string        `hcl:"name,label"`
	Extra       env.Variables `hcl:"extra,optional"`
	ItemConfigs []ItemConfig  `hcl:"item,block"`
}

type ItemConfig struct {
	Supplier string   `hcl:"supplier,label"`
	Account string   `hcl:"account,label"`
	Ref     string   `hcl:"ref,label"`
	Body    hcl.Body `hcl:",remain"`
}

func (c *Config) GetEnvironments() ([]env.Environment, error) {
	var envs []env.Environment
	var diags hcl.Diagnostics

	for _, ec := range c.EnvironmentConfigs {
		var tmp env.Environment
		if errs := c.DecodeEnvironment(ec, nil, &tmp); errs.HasErrors() {
			diags = diags.Extend(errs)
			continue
		}

		envs = append(envs, tmp)
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return envs, nil
}

var EnvironmentNotFoundError = errors.New("environment not found")

func (c *Config) FindEnvironment(name string) (env.Environment, error) {
	envs, err := c.GetEnvironments()
	if err != nil {
		return env.Environment{}, err
	}

	for _, e := range envs {
		if e.Name == name {
			return e, nil
		}
	}

	return env.Environment{}, EnvironmentNotFoundError
}

func (c Config) DecodeEnvironment(ec EnvironmentConfig, ctx *hcl.EvalContext, environment *env.Environment) hcl.Diagnostics {
	items := make(map[string]env.Item)
	suppliers := make(map[string]env.Supplier)
	var diags hcl.Diagnostics

	for _, ic := range ec.ItemConfigs {
		var tmp env.Supplier
		if errs := c.DecodeSupplier(ic, ctx, &tmp); errs.HasErrors() {
			diags = diags.Extend(errs)
			continue
		}

		name := fmt.Sprintf("%s/%s", ic.Supplier, ic.Ref)
		items[name] = env.Item{
			Ref:     ic.Ref,
			Account: ic.Account,
		}
		suppliers[name] = tmp
	}

	if diags.HasErrors() {
		return diags
	}

	*environment = env.Environment{
		Name:  ec.Name,
		Extra: ec.Extra,

		Items:     items,
		Suppliers: suppliers,
	}
	return nil
}

func (c Config) DecodeSupplier(ic ItemConfig, ctx *hcl.EvalContext, supplier *env.Supplier) hcl.Diagnostics {
	switch ic.Supplier {
	case "static":
		cfg := static.DefaultSupplierConfig()
		if diags := gohcl.DecodeBody(ic.Body, ctx, &cfg); diags.HasErrors() {
			return diags
		}
		*supplier = static.NewSupplier(cfg)

	default:
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Unknown item supplier",
			Detail:   fmt.Sprintf("Can't recognize item supplier %q", ic.Supplier),
			Subject:  ic.Body.MissingItemRange().Ptr(),
		}}
	}
	return nil
}

func LoadFile(filename string) (Config, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	return Load(filename, src)
}

func Load(filename string, src []byte) (Config, error) {
	var config Config

	err := hclsimple.Decode(filename, src, nil, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
