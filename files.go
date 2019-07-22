package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

// get all the files in the hold directory and add them to the holdFiles array
func (h *hole) getFiles(s time.Duration) {
	for {
		files, err := ioutil.ReadDir(h.holdDir)
		if err != nil {
			log.Println(err)
		}
		// fmt.Println("Checking for eligible files")
		_, files = h.eligibleFiles(files, h.holdDir)
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
		time.Sleep(s)
	}

}

// nextFile returns the path of the next file to be processed
func (h *hole) nextFile() string {
	// TODO: this is where we should get the next up priority
	// setting to 1
	priority := h.nextPriority()

	// if we get a 0 loop around
	if priority == 0 {
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
			// change from len to eligible files, this will eventually also look for the archive bit on windows systems
			l, _ := h.eligibleFiles(files, d)
			fmt.Println(l)
			if l == 0 {
				h.availableDirs <- d
				// sleep here to give time to make the directory used
				// time.Sleep(100 * time.Millisecond)
			}
		}
		time.Sleep(s)
	}
}

// eligibleFiles returns the length of eligible files in a folder
func (h *hole) eligibleFiles(files []os.FileInfo, dir string) (length int, outFiles []os.FileInfo) {
	// needs to check for archive bit on windows
	for _, file := range files {
		if isEligible(path.Join(dir, file.Name())) {
			length++
			outFiles = append(outFiles, file)
		}
	}
	return
}

func isEligible(filename string) bool {
	a, err := isArchive(filename)
	if err != nil {
		log.Panicf("error checking archive bit: %v", err)
	}

	h, err := isHidden(filename)
	if err != nil {
		log.Panicf("error checking if hidden: %v", err)
	}
	// fmt.Println(filename, a, h)
	if !a && !h {
		return true
	}
	return false
}

func (h *hole) getFilePriority(f string) int {
	priority := h.defaultPriority
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
	// should find the max priority first but defaulting to 1000
	p := 1000
	for pr := range h.holdFiles {
		if pr < p && len(h.holdFiles[pr]) > 0 {
			p = pr
		}
	}
	// reversing the 1000 to 0 until the max priority is figured out and this will be better
	// right now any priority over 1000 will not work
	if p == 1000 {
		p = 0
	}
	return p
}
