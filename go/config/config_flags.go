package config

const (
	overrideFlag = "override"
	configFlag   = "config"
)

func getFlagUsageMap() map[string]string {
	return map[string]string{
		overrideFlag: "Additive config file to apply on top of default or -config",
		configFlag:   "The path to the host's config file. Overrides all other flags",
	}
}
