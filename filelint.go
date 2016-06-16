package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
)

// group files by length
func main() {
	// prepare flags
	verbose := flag.Bool("v", false, "verbose (print developer logs?)")
	flag.Parse()

	var dir string
	if 0 < len(flag.Args()) {
		dir = flag.Arg(0)
	} else {
		dir = "."
	}

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	rootdir, err := expandUser(dir)
	if nil != err {
		fmt.Printf("Error obtaining user's home directory")
		os.Exit(1)
	}

	fmt.Printf("Looking in %s\n\n", rootdir)

	lengthCountMap := calculateLengthCountMap(rootdir)

	var possibleBytesSaving int64
	var wg sync.WaitGroup
	for filesize, files := range lengthCountMap {
		if len(files) > 1 {
			wg.Add(1)
			log.Printf("%d files of the same count (%d bytes) detected: %v\nChecking with md5 hashes...\n", len(files), filesize, files)
			duplicateFiles, err := findFilesWithSameMd5Hash(files)
			if nil != err {
				log.Printf("Error finding duplicate files: %s\n", err)
				continue
			}
			for _, duplicateFile := range duplicateFiles {
				fmt.Printf("\nDUPLICATE files (md5). Hash: %s. Paths:\n", duplicateFile.Md5Hash)
				for _, path := range duplicateFile.Filepaths {
					fmt.Println(path)
				}
			}
			possibleBytesSaving = possibleBytesSaving + (filesize * int64(len(files)-1))
			wg.Done()
		}
	}
	wg.Wait()

	fmt.Printf("\nAmount of space taken up by duplicate files: %s\n", humaniseBytes(possibleBytesSaving))

}

func humaniseBytes(bytes int64) string {
	units := []string{"PiB", "TiB", "GiB", "MiB", "KiB", "B"}
	for index, unit := range units {
		oneUnitInBytes := powInt64(int64(1024), len(units)-(index+2))
		//log.Printf("unit in bytes %s: %d Index %d\n", unit, oneUnitInBytes, len(units)-index)
		if bytes >= oneUnitInBytes {
			amountOfWholeUnits := bytes / oneUnitInBytes
			remainder := bytes % oneUnitInBytes
			decimal := float32(remainder) / float32(oneUnitInBytes)
			return fmt.Sprintf("%.2f %s", float32(amountOfWholeUnits)+decimal, unit)
		}
	}
	return fmt.Sprintf("%d %s", bytes, "B")
}

func powInt64(base int64, exponent int) int64 {
	total := base
	for i := 0; i < exponent; i++ {
		total = total * base
	}
	return total
}

func calculateLengthCountMap(rootdir string) map[int64][]string {

	lengthCountMap := make(map[int64][]string)

	filepath.Walk(rootdir, func(path string, fileinfo os.FileInfo, err error) error {
		if fileinfo.IsDir() {
			return nil
		}
		lengthCountMap[fileinfo.Size()] = append(lengthCountMap[fileinfo.Size()], path)
		return nil
	})

	return lengthCountMap
}

type DuplicateFiles struct {
	Md5Hash   string
	Filepaths []string
}

type md5Map map[string][]string

func findFilesWithSameMd5Hash(paths []string) ([]*DuplicateFiles, error) {
	md5MapHashes := make(md5Map)
	var duplicateFiles []*DuplicateFiles

	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(len(paths))
	for _, path := range paths {
		go func(path string) error {
			hashBytes, err := fileMD5(path)
			if nil != err {
				return err
			}
			hash := string(hashBytes)

			mutex.Lock()
			md5MapHashes[hash] = append(md5MapHashes[hash], path)
			mutex.Unlock()

			wg.Done()
			return nil
		}(path)
	}
	wg.Wait()

	for hash, paths := range md5MapHashes {
		if len(paths) > 1 {
			duplicateFiles = append(duplicateFiles, &DuplicateFiles{
				Md5Hash:   hash,
				Filepaths: paths,
			})

		}
	}
	return duplicateFiles, nil

}

func fileMD5(filePath string) (string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// expand for unix
func expandUser(path string) (string, error) {
	if !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	u, err := user.Current()
	if nil != err {
		return "", err
	}

	return strings.Replace(path, "~", u.HomeDir, 1), nil
}
