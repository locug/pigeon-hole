package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// wait for something on availableDirs channel and then move files when one shows up
func (h *hole) mover() {

	for {
		dir := <-h.availableDirs
		file := h.nextFile()
		inFile := path.Join(h.holdDir, file)
		outFile := path.Join(dir, file)
		fmt.Println(inFile, outFile)

		data, err := ioutil.ReadFile(inFile)
		if err != nil {
			// ignore the error and put the file back into the channel
			h.availableDirs <- dir
			continue
		}

		err = ioutil.WriteFile(outFile, data, 0666)
		if err == nil {
			// since there wasn't an error meaning file was written delete original
			os.RemoveAll(inFile)
		}
	}
}
