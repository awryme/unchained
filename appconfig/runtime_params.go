package appconfig

type RuntimeParams struct {
	LogLevel string
	DNS      string
	Proto    string
	ID       string
	Tags     []string
}

func setRuntimeParams(cfg *Config, params *RuntimeParams) {
	if params == nil {
		return
	}

	trySet := func(value *string, param string) {
		if param != "" {
			*value = param
		}
	}

	trySet(&cfg.Proto, params.Proto)
	trySet(&cfg.LogLevel, params.LogLevel)
	trySet(&cfg.DNS, params.DNS)
	trySet(&cfg.ID, params.ID)

	if len(params.Tags) > 0 {
		cfg.Tags = params.Tags
	}
}
