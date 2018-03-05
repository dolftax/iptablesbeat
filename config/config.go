// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	Sudo   bool          `config:"sudo"`
	Lock   bool          `config:"lock"`
	Table  string        `config:"table"`
	Chain  []string      `config:"chain"`
}

var DefaultConfig = Config{
	Period: 5 * time.Second,
	Sudo:   true,
	Table:  "filter",
	Chain:  []string{"INPUT"},
}
