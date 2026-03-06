/*
Package config manages the sql-archiver configuration file.

The config file is stored at $HOME/.config/sql-archiver/config.yaml and holds a
list of saved connection profiles under the "connections" key:

	connections:
	  - name: "local"
	    host: "localhost:1433"
	    user: "sa"
	    password: "secret"
*/
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Eventid3/sql-archiver/domain"
	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "yaml"
)

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "sql-archiver"), nil
}

func newViper() (*viper.Viper, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType(configType)
	v.AddConfigPath(dir)
	return v, nil
}

// LoadConnections reads all saved connections from the config file.
// If the config file does not exist or contains no connections, an empty slice
// is returned without an error.
func LoadConnections() ([]domain.Connection, error) {
	v, err := newViper()
	if err != nil {
		return nil, err
	}

	if err := v.ReadInConfig(); err != nil {
		// Config file absent is not an error — just no saved connections yet.
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	type raw struct {
		Name     string `mapstructure:"name"`
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	}

	var rawConns []raw
	if err := v.UnmarshalKey("connections", &rawConns); err != nil {
		return nil, fmt.Errorf("failed to parse connections: %w", err)
	}

	conns := make([]domain.Connection, 0, len(rawConns))
	for _, r := range rawConns {
		conns = append(conns, domain.Connection{
			Name:     r.Name,
			Host:     r.Host,
			User:     r.User,
			Password: r.Password,
		})
	}
	return conns, nil
}

// SaveConnection writes a connection to the config file. If a connection with
// the same Name already exists it is replaced; otherwise it is appended.
func SaveConnection(c domain.Connection) error {
	v, err := newViper()
	if err != nil {
		return err
	}

	// Best-effort read so we preserve existing connections.
	_ = v.ReadInConfig()

	type raw struct {
		Name     string `mapstructure:"name"`
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	}

	var existing []raw
	_ = v.UnmarshalKey("connections", &existing)

	updated := false
	for i, r := range existing {
		if r.Name == c.Name {
			existing[i] = raw{c.Name, c.Host, c.User, c.Password}
			updated = true
			break
		}
	}
	if !updated {
		existing = append(existing, raw{c.Name, c.Host, c.User, c.Password})
	}

	// Convert back to []map[string]interface{} so viper serialises cleanly.
	asSlice := make([]map[string]interface{}, len(existing))
	for i, r := range existing {
		asSlice[i] = map[string]interface{}{
			"name":     r.Name,
			"host":     r.Host,
			"user":     r.User,
			"password": r.Password,
		}
	}
	v.Set("connections", asSlice)

	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	path := filepath.Join(dir, configName+"."+configType)
	if err := v.WriteConfigAs(path); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

// DeleteConnection removes the connection with the given name from the config
// file. If no connection with that name exists the call is a no-op.
func DeleteConnection(name string) error {
	v, err := newViper()
	if err != nil {
		return err
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	type raw struct {
		Name     string `mapstructure:"name"`
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	}

	var existing []raw
	_ = v.UnmarshalKey("connections", &existing)

	filtered := existing[:0]
	for _, r := range existing {
		if r.Name != name {
			filtered = append(filtered, r)
		}
	}

	asSlice := make([]map[string]interface{}, len(filtered))
	for i, r := range filtered {
		asSlice[i] = map[string]interface{}{
			"name":     r.Name,
			"host":     r.Host,
			"user":     r.User,
			"password": r.Password,
		}
	}
	v.Set("connections", asSlice)

	dir, err := configDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, configName+"."+configType)
	if err := v.WriteConfigAs(path); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}
