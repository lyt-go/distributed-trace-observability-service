package registration

import "tracing/internal/policy"

func Register(cfg *policy.Config, service string) error {
	if cfg.Validator != nil {
		if err := cfg.Validator.Validate(service); err != nil {
			return err
		}
	}
	cfg.Labels[service] = "registered"
	return nil
}
