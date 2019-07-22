package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

// wait for something on availableDirs channel and then move files when one shows up
func (h *hole) mover() {

	for {
		dir := <-h.availableDirs
		// lock other operations while reading the directory
		file := h.nextFile()
		inFile := path.Join(h.holdDir, file)
		outFile := path.Join(dir, file)

		data, err := ioutil.ReadFile(inFile)
		if err != nil {
			// if there was an error then just continue on, the file will be re-read
			continue
		}

		err = ioutil.WriteFile(outFile, data, 0666)
		if err == nil {
			// remove the archive bit on the file
			err := removeArchive(outFile)
			if err != nil {
				log.Printf("error setting archive bit: %s", file)
			}
			// since there wasn't an error meaning file was written delete original
			os.RemoveAll(inFile)

		}
		// unlock so other operations can do stuff
	}
}
