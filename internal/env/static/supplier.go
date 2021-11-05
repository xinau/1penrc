package static

import (
	"github.com/xinau/1penrc/internal/env"
	"github.com/xinau/1penrc/internal/op"
)

type SupplierConfig struct {
	Mapping map[string]string `hcl:"mapping,attr"`
}

func DefaultSupplierConfig() SupplierConfig {
	return SupplierConfig{}
}

type Supplier struct {
	config SupplierConfig
}

func NewSupplier(config SupplierConfig) Supplier {
	return Supplier{
		config: config,
	}
}

func (i Supplier) Retrieve(item env.Item) (env.Variables, error) {
	client, err := env.GetClient(item.Account)
	if err != nil {
		return nil, err
	}

	var fields []string
	revs := make(map[string]string)
	for key, field := range i.config.Mapping {
		fields = append(fields, field)
		revs[field] = key
	}

	values, err := client.GetItemFields(item.Ref, fields)
	if err != nil {
		return nil, err
	}

	if err := op.CheckMissingItemFields(values, fields); err != nil {
		return nil, err
	}

	vars := make(map[string]string)
	for field, val := range values {
		if key, ok := revs[field]; ok {
			vars[key] = val
		}
	}

	return vars, nil
}
