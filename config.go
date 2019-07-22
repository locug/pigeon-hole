package main

import (
	"log"
	"regexp"
	"strconv"
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

	h.prioritize(cfg)
	return nil
}

func splitDirs(dirs string) []string {
	// ioutil.ReadDir does not like quotation marks so just remove them
	dirs = strings.Replace(dirs, "\"", "", -1)
	return strings.Split(dirs, ",")
}

func (h *hole) prioritize(cfg *ini.File) {

	for _, r := range cfg.Section("PRIORITY").Keys() {
		l, err := strconv.Atoi(r.Name()[1:])
		if err != nil {
			log.Panicf("priority key did not convert to int key: %s err: %v", r.Name(), err)
		}
		p := priority{
			level: l,
			regex: regexp.MustCompile(`^` + r.String() + `.+`),
		}
		h.priorities = append(h.priorities, p)

	}

}
