package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// get all the files in the hold directory and add them to the holdFiles array
func (h *hole) getFiles(s time.Duration) {
	for {
		// lock other operations while reading the directory
		h.mutex.Lock()
		files, err := ioutil.ReadDir(h.holdDir)
		if err != nil {
			log.Println(err)
		}
		_, files = eligibleFiles(files)
	loop:
		for _, f := range files {
			// TODO: run regex over files to prioritis
			priority := h.getFilePriority(f.Name())

			// ranging over the current holdFiles with priority to see if file has already been added
			// should probably be a map to simplify
			for _, hf := range h.holdFiles[priority] {
				if f.Name() == hf {
					// matched a current file so jump back to start
					continue loop
				}
			}
			fmt.Printf("Adding File: %s with priority %d\n", f.Name(), priority)
			h.holdFiles[priority] = append(h.holdFiles[priority], f.Name())
		}
		// unlock so other operations can do stuff
		h.mutex.Unlock()
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
		// lock other operations while reading the directory
		h.mutex.Lock()
		for _, d := range h.outDirs {

			files, err := ioutil.ReadDir(d)
			if err != nil {
				return err
			}
			// change from len to eligible files, this will eventually also look for the archive bit on windows systems
			l, _ := eligibleFiles(files)
			if l == 0 {
				h.availableDirs <- d
				// sleep here to give time to make the directory used
				// time.Sleep(100 * time.Millisecond)
			}
		}
		// unlock so other operations can do stuff
		h.mutex.Unlock()
		time.Sleep(s)
	}
}

// eligibleFiles returns the length of eligible files in a folder
func eligibleFiles(files []os.FileInfo) (length int, outFiles []os.FileInfo) {
	// needs to check for archive bit on windows
	for _, file := range files {
		if file.Name()[0:1] != "." {
			length++
			outFiles = append(outFiles, file)
		}
	}
	return
}

func (h *hole) getFilePriority(f string) int {
	priority := 100
	for _, p := range h.priorities {
		if p.regex.MatchString(f) {
			priority = p.level
		}
	}

	// testregex := regexp.MustCompile(`^Z.+`)
	// if testregex.MatchString(f) {

	// }
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
