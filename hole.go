package main

import (
	"log"
	"regexp"
	"sync"
	"time"
)

type hole struct {
	holdDir string   // where the files live while waiting for an output directory
	outDirs []string // the directorys where the files go
	inDirs  []string // where the files get put by external program
	mutex   sync.Mutex
	// holdFiles should probably be map of channels to allow for adding and removing concurrently
	holdFiles       map[int][]string // array of the files, the mapped int is the priority
	availableDirs   chan string
	priorities      []priority
	defaultPriority int
}

type priority struct {
	level int
	regex *regexp.Regexp
}

func main() {
	// TODO: need a "new" function and these to be passed by variable
	var h hole
	err := h.make("hole.ini")
	if err != nil {
		log.Panicln(err)
	}
	// make the holdFiles map
	h.holdFiles = make(map[int][]string)
	// make the availableDirs chan with length of number of possible out dirs
	h.availableDirs = make(chan string, len(h.outDirs))

	// get the files, needs to loop after a delay
	go h.getFiles(time.Duration(1 * time.Second))

	// launch the mover
	go h.mover()

	// h.getFiles()

	// launch the out folder check
	err = h.checkOut(time.Duration(1 * time.Second))
	if err != nil {
		log.Panicln(err)
	}
}
