package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"time"
)

// get all the files in the hold directory and add them to the holdFiles array
func (h *hole) getFiles(s time.Duration) {
	for {
		files, err := ioutil.ReadDir(h.holdDir)
		if err != nil {
			log.Println(err)
		}
	loop:
		for _, f := range files {
			// TODO: run regex over files to prioritis
			priority := getFilePriority(f.Name())

			// ranging over the current holdFiles with priority to see if file has already been added
			// should probably be a map to simplify
			for _, hf := range h.holdFiles[priority] {
				if f.Name() == hf {
					// matched a current file so jump back to start
					continue loop
				}
			}
			h.holdFiles[priority] = append(h.holdFiles[priority], f.Name())
		}
		time.Sleep(s)
	}

}

// nextFile returns the path of the next file to be processed
func (h *hole) nextFile() string {
	// TODO: this is where we should get the next up priority
	// setting to 1
	priority := h.nextPriority()

	if priority == 101 {
		time.Sleep(1 * time.Second)
		return h.nextFile()
	}

	_, ok := h.holdFiles[priority]

	// if a file does not exist at the given priorty then wait and loop on this function
	if !ok || len(h.holdFiles[priority]) == 0 {
		// should build in some logic which adjusts the sleep time based load or activity
		time.Sleep(1 * time.Second)
		return h.nextFile()
	}

	file := h.holdFiles[priority][0]
	// recreate the array without the file
	h.holdFiles[priority] = h.holdFiles[priority][1:]

	return file
}

// checkOut looks in the out dirs and adds any empty ones to the availableDirs channel
func (h *hole) checkOut(s time.Duration) error {
	for {
		for _, d := range h.outDirs {
			files, err := ioutil.ReadDir(d)
			if err != nil {
				return err
			}
			if len(files) == 0 {
				h.availableDirs <- d
			}
		}
		time.Sleep(s)
	}
}

func getFilePriority(f string) int {
	priority := 100

	high := regexp.MustCompile(`^h_.+`)
	med := regexp.MustCompile(`^g_.+`)

	switch {
	case high.MatchString(f):
		priority = 1
	case med.MatchString(f):
		priority = 2

	}
	return priority
}

// firstPriority returns the lowest priority int
func (h *hole) nextPriority() int {
	// should find the max priority first but defaulting to 100
	p := 101
	for pr := range h.holdFiles {
		if pr < p && len(h.holdFiles[pr]) > 0 {
			p = pr
		}
	}
	return p
}
