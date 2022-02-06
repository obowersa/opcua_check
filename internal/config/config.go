package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envPrefix = "OPCUA_CHECK"
)

type Config struct {
	Endpoint      string
	BaseVariables map[string][]string
	Variables     []string
	Timeout       int
	ConfigFile    string
	Output        string
}

func NewConfig() *Config {

	return &Config{}
}

func (c *Config) InitializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigFile(c.ConfigFile)
	v.SetEnvPrefix(envPrefix)

	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return fmt.Errorf("error loading config: err: %w", err)
		}
	}

	// Unmarshalls our BaseVariables from config
	t := Config{}
	_ = v.Unmarshal(&t)
	c.BaseVariables = t.BaseVariables

	v.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, v)
	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			_ = v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
