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

		// Create the file
		fhandle, err := os.Create(outFile)
		// set the archive bit so LOC doesn't process this file
		err = setArchive(outFile)

		err = ioutil.WriteFile(outFile, data, 0666)
		if err != nil {
			log.Println("error writing file: ", err)
		}
		// remove the archive bit on the file
		err = removeArchive(outFile)
		if err != nil {
			log.Printf("error setting archive bit: %s", file)
		} else {
			// since there wasn't an error meaning file was written delete original
			os.RemoveAll(inFile)
			fhandle.Close()
		}
		log.Printf("Moving %s to %s", file, dir)
	}
}
