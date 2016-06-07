package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func Test_fileMd5(t *testing.T) {

	tempdir, err := ioutil.TempDir("", "")
	if nil != err {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	file := filepath.Join(tempdir, "def.txt")

	err = ioutil.WriteFile(file, []byte("def\n"), 0755)
	if nil != err {
		t.Fatal(err)
	}

	var hash string
	hash, err = fileMD5(file)
	if nil != err {
		t.Fatal(err)
	}

	assert.Equal(t, "614dd0e977becb4c6f7fa99e64549b12", hash)

}

func Test_FindFilesWithSameMd5Hash(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "")
	if nil != err {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	fileContents := []string{"abc", "def", "abc", "gh", "abc", "def"}
	var filePaths []string

	for index, contents := range fileContents {
		filepath := filepath.Join(tempdir, "same_md5_hash_test_file_"+strconv.Itoa(index)+".txt")
		filePaths = append(filePaths, filepath)

		err = ioutil.WriteFile(filepath, []byte(contents+"\n"), 0755)
		if nil != err {
			t.Fatal(err)
		}
	}

	duplcateFiles, err := findFilesWithSameMd5Hash(filePaths)
	if nil != err {
		t.Fatal(err)
	}

	assert.Len(t, duplcateFiles, 2)

	var twoDuplicateFiles *DuplicateFiles
	var threeDuplicateFiles *DuplicateFiles
	if 2 == len(duplcateFiles[0].Filepaths) {
		twoDuplicateFiles = duplcateFiles[0]
		threeDuplicateFiles = duplcateFiles[1]
	} else {
		twoDuplicateFiles = duplcateFiles[1]
		threeDuplicateFiles = duplcateFiles[0]
	}

	assert.Len(t, twoDuplicateFiles.Filepaths, 2)
	assert.Len(t, threeDuplicateFiles.Filepaths, 3)
	assert.NotEqual(t, twoDuplicateFiles.Md5Hash, threeDuplicateFiles.Md5Hash)
}

func Test_humaniseBytes(t *testing.T) {
	assert.Equal(t, "1.00 KiB", humaniseBytes(int64(1024)))
	assert.Equal(t, "1.50 KiB", humaniseBytes(int64(1540)))
	assert.Equal(t, "5.50 MiB", humaniseBytes(int64(5767168)))
	assert.Equal(t, "40.70 GiB", humaniseBytes(int64(43701292236)))
}
