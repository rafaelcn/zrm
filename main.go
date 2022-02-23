package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	recursive *bool   = flag.Bool("recursive", false, "")
	input     *string = flag.String("files", "", "")
)

func main() {
	flag.Parse()

	if input == nil {
		log.Printf("no files or directories provided")
		return
	}

	files := strings.Split(*input, " ")

	for _, file := range files {
		f, err := os.Open(file)

		if err != nil {
			log.Printf("failed to open %s, reason %v", file, err)
			close(f)
			continue
		}

		stat, err := f.Stat()

		if err != nil {
			log.Printf("failed to stat file %s, reason %v", file, err)
			close(f)
			continue
		}

		if stat.IsDir() {
			if !*recursive {
				log.Printf("won't recurse into dir %s (enable recursion)",
					file)
			} else {
				wg := sync.WaitGroup{}

				filepath.Walk(f.Name(), func(path string, info fs.FileInfo, err error) error {
					if !info.IsDir() {
						log.Printf("%v %v", info.Name(), info.Size())

						wg.Add(1)
						go func() {
							write(path, stat)
							wg.Done()
						}()
					}

					return nil
				})

				wg.Wait()
			}
		} else {
			write(f.Name(), stat)
		}

		close(f)
	}

}

func close(f *os.File) {
	err := f.Close()

	if err != nil {
		log.Printf("failed to close file %s, reason %v", f.Name(), err)
	}
}

func write(path string, stat os.FileInfo) {
	now := time.Now()

	f, err := os.OpenFile(path, os.O_RDWR, 0777)

	if err != nil {
		log.Printf("failed to open file %s, reason %v", path, err)
		return
	}
	defer f.Close()

	written, s := 0, stat.Size()

	for s > 0 {
		n, werr := f.Write([]byte{0})
		written += n

		if werr != nil {
			log.Printf("failed to write zeroes in file %s, reason %v",
				f.Name(), werr)
		}

		s -= int64(n)
	}

	ellapsed := time.Since(now)
	fmt.Printf("took %v to write %d bytes\n", ellapsed, written)
}
