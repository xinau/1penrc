package duration

import "time"

type Duration time.Duration

func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var plain string
	if err := unmarshal(&plain); err != nil {
		return err
	}

	t, err := time.ParseDuration(plain)
	if err != nil {
		return err
	}

	*d = Duration(t)
	return nil
}
