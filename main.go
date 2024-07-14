package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rafaelcn/zrm/fisy"
)

var (
	blockSize *int    = flag.Int("bs", 1, "the block size used to write data")
	delete    *bool   = flag.Bool("d", false, "delete files after writing")
	recursive *bool   = flag.Bool("r", false, "enable directory recursive walk")
	input     *string = flag.String("i", "", "a list of files or directories")
)

func main() {
	flag.Parse()

	if len(*input) == 0 {
		log.Fatalf("no files or directories provided")
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
				log.Printf("won't recurse into dir %s (enable recursion)", file)
			} else {
				wg := sync.WaitGroup{}

				filepath.Walk(f.Name(), func(path string, info fs.FileInfo, err error) error {
					if !info.IsDir() {
						wg.Add(1)
						go func() {
							fisy.WriteFill(path, []byte{0})
							if *delete {
								fisy.Delete(path)
							}
							wg.Done()
						}()
					}

					return nil
				})

				wg.Wait()
			}
		} else {
			fisy.WriteFill(f.Name(), []byte{0})
			if *delete {
				fisy.Delete(f.Name())
			}
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
