package main

import (
	"path"
	"time"
)

type hole struct {
	holdDir string   // where the files live while waiting for an output directory
	outDirs []string // the directorys where the files go
	inDir   string   // where the files get put by external program
	// holdFiles should probably be map of channels to allow for adding and removing concurrently
	holdFiles     map[int][]string // array of the files, the mapped int is the priority
	availableDirs chan string
}

func main() {
	// need a "new" function and these to be passed by variable
	var h hole
	h.holdDir = "folders/HOLD"
	// TODO: Make it so either provide an array of out folders or regext for matching to that folder
	h.outDirs = []string{"folders/XF999980", "folders/XF999981", "folders/XF999982", "folders/XF999983"}

	h.holdFiles = make(map[int][]string)
	// make the availableDirs chan with length of number of possible out dirs
	h.availableDirs = make(chan string, len(h.outDirs))

	// get the files, needs to loop after a delay
	go h.getFiles(time.Duration(1 * time.Second))

	// launch the mover
	go h.mover()

	// h.getFiles()

	// launch the out folder check
	h.checkOut(time.Duration(1 * time.Second))
}

func (h *hole) hold(f string) string {
	return path.Join(h.holdDir, f)
}
func (h *hole) in(f string) string {
	return path.Join(h.inDir, f)
}

// pathOut takes the i of the array for outDirs to return the folder
func (h *hole) out(i int, f string) string {

	return ""
}
