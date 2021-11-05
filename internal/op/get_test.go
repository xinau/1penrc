package op

import (
	"errors"
	"reflect"
	"testing"
)

func TestCommand_GetItemFields(t *testing.T) {
	cmd := &Client{Executable: "./testdata/op-get-item-fields.sh"}

	tests := []struct {
		name   string
		item   string
		fields []string
		want   map[string]string
		err    error
	}{{
		"item fields well formed",
		"item-fields-well-formed",
		[]string{"field"},
		map[string]string{"foo": "bar", "baz": "qux"},
		nil,
	}, {
		"item fields malformed",
		"item-fields-malformed",
		[]string{"field"},
		nil,
		ParsingItemFieldsError,
	}, {
		"zero item fields provided",
		"",
		nil,
		map[string]string{},
		nil,
	}, {
		"single empty item field provided",
		"",
		[]string{""},
		map[string]string{"": ""},
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := cmd.GetItemFields(test.item, test.fields)
			Assertf(t, reflect.DeepEqual(got, test.want), "got %q, expected %q", got, test.want)
			Assertf(t, errors.Is(err, test.err), "got %q, expected %q", got, test.err)
		})
	}
}

func TestCheckMissingItemFields(t *testing.T) {
	tests := []struct{
		name string
		values map[string]string
		fields []string
		want error
	}{{
		"no missing fields",
		map[string]string{"foo":"1", "bar":"2"},
		[]string{"foo", "bar"},
		nil,
	}, {
		"missing single field",
		map[string]string{"foo":"1", "bar":"2"},
		[]string{"foo", "bar", "baz"},
		&MissingItemFieldsError{[]string{"baz"}},
	}, {
		"missing multiple fields",
		map[string]string{"foo":"1"},
		[]string{"foo", "bar", "baz"},
		&MissingItemFieldsError{[]string{"baz", "bar"}},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CheckMissingItemFields(test.values, test.fields)
			Assertf(t, errors.Is(got, test.want), "got %q, expected %q", got, test.want)
		})
	}

}
