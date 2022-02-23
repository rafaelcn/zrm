package fisy

import (
	"log"
	"os"
	"time"
)

func open(file string) *os.File {

	f, err := os.OpenFile(file, os.O_RDWR, 0777)

	if err != nil {
		return nil
	}

	return f
}

func WriteFill(file string, fill []byte) {

	f := open(file)

	if f == nil {
		return
	}
	defer f.Close()

	props, err := os.Stat(file)

	if err != nil {
		return
	}

	var written, size int64

	size = props.Size()
	now := time.Now()

	for written < size {
		// ignore fill for now
		n, err := f.Write([]byte{0})

		if err != nil {
			break
		}

		written += int64(n)
	}

	log.Printf("[%s|%s]\n", file, time.Since(now))
}

func Delete(file string) error {
	return os.Remove(file)
}
