package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/jaipradeesh/iptablesbeat/config"
)

type Iptablesbeat struct {
	done          chan struct{}
	config        config.Config
	client        beat.Client
	lastIndexTime time.Time
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Iptablesbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func (bt *Iptablesbeat) Run(b *beat.Beat) error {
	logp.Info("iptablesbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			err := bt.gather(b.Info.Name)
			if err != nil {
				return nil
			}
			counter++
		}
	}
}

func (bt *Iptablesbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
