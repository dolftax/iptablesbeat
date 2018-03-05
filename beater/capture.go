package beater

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

// Validation regexs
var errParse = errors.New("Cannot parse iptables list information")
var chainNameRe = regexp.MustCompile(`^Chain\s+(\S+)`)
var fieldsHeaderRe = regexp.MustCompile(`^\s*pkts\s+bytes\s+`)
var valuesRe = regexp.MustCompile(`^\s*(\d+)\s+(\d+)\s+.*?\/\*\s*(.+?)\s*\*\/\s*`)

func (bt *Iptablesbeat) chainList(table, chain string) (string, error) {
	iptablePath, err := exec.LookPath("iptables")
	if err != nil {
		return "", err
	}
	var args []string
	name := iptablePath
	if bt.config.Sudo {
		name = "sudo"
		args = append(args, iptablePath)
	}
	iptablesBaseArgs := "-nvL"
	if bt.config.Lock {
		iptablesBaseArgs = "-wnvL 1"
	}
	args = append(args, iptablesBaseArgs, chain, "-t", table, "-x")
	c := exec.Command(name, args...)
	out, err := c.Output()
	return string(out), err
}

func (bt *Iptablesbeat) gather(beatname string) error {
	for _, chain := range bt.config.Chain {
		data, e := bt.chainList(bt.config.Table, chain)
		if e != nil {
			continue
		}

		lines := strings.Split(data, "\n")
		if len(lines) < 3 {
			return errParse
		}

		mchain := chainNameRe.FindStringSubmatch(lines[0])
		if mchain == nil {
			return errParse
		}

		if !fieldsHeaderRe.MatchString(lines[1]) {
			return errParse
		}

		for _, line := range lines[2:] {
			matches := valuesRe.FindStringSubmatch(line)

			if len(matches) != 4 {
				continue
			}

			pkts := matches[1]
			bytes := matches[2]
			comment := matches[3]

			var err error
			respPkts, err := strconv.ParseUint(pkts, 10, 64)
			if err != nil {
				return errParse
			}
			respBytes, err := strconv.ParseUint(bytes, 10, 64)
			if err != nil {
				return errParse
			}

			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"type":    beatname,
					"table":   bt.config.Table,
					"chain":   mchain[1],
					"rule_id": comment,
					"pkts":    respPkts,
					"bytes":   respBytes,
				},
			}
			bt.client.Publish(event)
			logp.Info("Iptables event sent")
		}
	}
	return nil
}
