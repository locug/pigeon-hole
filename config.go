package main

import (
	"strings"

	"github.com/go-ini/ini"
)

// make creates a "hole" from supplied ini file
func (h *hole) make(filename string) error {
	cfg, err := ini.Load(filename)
	if err != nil {
		return err
	}
	// single directory don't split
	h.holdDir = cfg.Section("DIRECTORIES").Key("HOLD").String()
	// split the in direcotories into an array to be used
	// in direcotry not actually used yet
	h.inDirs = splitDirs(cfg.Section("DIRECTORIES").Key("IN").String())
	// split the out directoryies into an array to be used
	h.outDirs = splitDirs(cfg.Section("DIRECTORIES").Key("OUT").String())
	return nil
}

func splitDirs(dirs string) []string {
	// ioutil.ReadDir does not like quotation marks so just remove them
	dirs = strings.Replace(dirs, "\"", "", -1)
	return strings.Split(dirs, ",")
}
