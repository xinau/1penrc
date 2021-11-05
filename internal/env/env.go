package env

import (
	"github.com/xinau/1penrc/internal/op"
)

type Supplier interface {
	Retrieve(item Item) (Variables, error)
}

type Item struct {
	Ref     string
	Account string
}

type Variables map[string]string

func (v Variables) Merge(o Variables) Variables {
	vars := make(map[string]string)
	for key, val := range v {
		vars[key] = val
	}

	for key, val := range o {
		if _, ok := vars[key]; !ok {
			vars[key] = val
		}
	}

	return vars
}

type Environment struct {
	Name  string
	Extra Variables

	Items     map[string]Item
	Suppliers map[string]Supplier
}

func (e Environment) Export() (Variables, error) {
	var vars Variables
	for key, item := range e.Items {
		supplier, ok := e.Suppliers[key]
		if !ok {
			continue
		}

		tmp, err := supplier.Retrieve(item)
		if err != nil {
			return nil, err
		}
		vars = vars.Merge(tmp)
	}

	return vars.Merge(e.Extra), nil
}

var clients = make(map[string]*op.Client)

func GetClient(account string) (*op.Client, error) {
	if account == "" {
		return op.NewClient(), nil
	}

	if client, ok := clients[account]; ok && client != nil {
		return client, client.SignIn()
	}

	client, err := op.SignIn(account)
	if err != nil {
		return nil, err
	}

	clients[account] = client
	return client, nil
}
