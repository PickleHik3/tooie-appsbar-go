package config

// Config represents the launcher configuration.
type Config struct {
	Grid  GridConfig  `yaml:"grid"`
	Style StyleConfig `yaml:"style"`
	Apps  []AppConfig `yaml:"apps"`
}

// GridConfig defines the grid layout.
type GridConfig struct {
	Rows    int `yaml:"rows"`
	Columns int `yaml:"columns"`
}

// StyleConfig defines visual styling options.
type StyleConfig struct {
	Border  bool `yaml:"border"`
	Padding int  `yaml:"padding"`
}

// AppConfig defines a single app entry.
type AppConfig struct {
	Name     string `yaml:"name"`
	Icon     string `yaml:"icon"`
	Package  string `yaml:"package"`
	Activity string `yaml:"activity,omitempty"`
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() Config {
	return Config{
		Grid: GridConfig{
			Rows:    1,
			Columns: 5,
		},
		Style: StyleConfig{
			Border:  true,
			Padding: 1,
		},
		Apps: []AppConfig{},
	}
}
