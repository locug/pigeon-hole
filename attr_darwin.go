package main

import "path"

func isArchive(filename string) (bool, error) {
	return false, nil
}

func isHidden(filename string) (bool, error) {
	if path.Base(filename)[0:1] == "." {
		return true, nil
	}
	return false, nil
}

func removeArchive(filename string) error {
	return nil
}
