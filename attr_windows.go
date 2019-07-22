package main

import (
	"fmt"
	"syscall"
)

func isArchive(filename string) (bool, error) {
	pointer, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return false, err
	}
	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return false, err
	}
	// fmt.Println(attributes & syscall.FILE_ATTRIBUTE_ARCHIVE)
	return attributes&syscall.FILE_ATTRIBUTE_ARCHIVE != 0, nil

	// windows.GetFileAttributes()
}

func isHidden(filename string) (bool, error) {
	pointer, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return false, err
	}
	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return false, err
	}
	// fmt.Println(filename, attributes & syscall.FILE_ATTRIBUTE_HIDDEN)
	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil

	// windows.GetFileAttributes()
}

func removeArchive(filename string) error {
	fmt.Println("Removing the archive bit")
	pointer, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}
	return syscall.SetFileAttributes(pointer, syscall.FILE_ATTRIBUTE_NORMAL)
}
