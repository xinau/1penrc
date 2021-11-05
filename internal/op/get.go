package op

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ParsingItemFieldsError = errors.New("parsing item fields")
)

func (c *Client) GetItemFields(item string, fields []string) (map[string]string, error) {
	// skip 1password lookup for zero fields or a single empty field
	if len(fields) < 1 {
		return map[string]string{}, nil
	}
	if len(fields) == 1 && fields[0] == "" {
		return map[string]string{"": ""}, nil
	}

	data, err := c.Exec([]string{
		"get", "item", item,
		"--fields", strings.Join(fields, ","),
		"--format", "json",
	})
	if err != nil {
		return nil, err
	}

	var values map[string]string
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("%w: %s", ParsingItemFieldsError, err)
	}
	return values, nil
}

func GetItemFields(item string, fields []string) (map[string]string, error) {
	return NewClient().GetItemFields(item, fields)
}

type MissingItemFieldsError struct {
	Fields []string
}

func (e *MissingItemFieldsError) Error() string {
	return fmt.Sprintf("missing item fields %q", strings.Join(e.Fields, ","))
}

func (e *MissingItemFieldsError) Is(target error) bool {
	val, ok := target.(*MissingItemFieldsError)
	if !ok {
		return false
	}

	if len(e.Fields) != len(val.Fields) {
		return false
	}

	lookup := make(map[string]struct{})
	for _, field := range e.Fields {
		lookup[field] = struct{}{}
	}

	for _, field := range val.Fields {
		if _, ok := lookup[field]; !ok {
			return false
		}
	}
	return true
}

func CheckMissingItemFields(values map[string]string, fields []string) error {
	var missing []string
	for _, field := range fields {
		if _, ok := values[field]; !ok {
			missing = append(missing, field)
		}
	}

	if len(missing) > 0 {
		return &MissingItemFieldsError{
			Fields: missing,
		}
	}
	return nil
}
