package main

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"
)

func Test_fileMd5(t *testing.T) {

	tempdir, err := ioutil.TempDir("", "")
	if nil != err {
		t.Fatal(err)
	}

	file := filepath.Join(tempdir, "def.txt")

	err = ioutil.WriteFile(file, []byte("def\n"), 0755)
	if nil != err {
		t.Fatal(err)
	}

	var hash string
	hash, err = fileMD5(file)

	if "614dd0e977becb4c6f7fa99e64549b12" != hash {
		t.Errorf("hash should be 614dd0e977becb4c6f7fa99e64549b12 but was %s", hash)
	}

}

func Test_FindFilesWithSameMd5Hash(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "")
	if nil != err {
		t.Fatal(err)
	}

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

	if 2 != len(duplcateFiles) {
		t.Errorf("There should be 2 sets of same md5 hash files, but only found %d\n", len(duplcateFiles))
	}

	var twoDuplicateFiles DuplicateFiles
	var threeDuplicateFiles DuplicateFiles
	if 2 == len(duplcateFiles[0].Filepaths) {
		twoDuplicateFiles = *duplcateFiles[0]
		threeDuplicateFiles = *duplcateFiles[1]
	} else {
		twoDuplicateFiles = *duplcateFiles[1]
		threeDuplicateFiles = *duplcateFiles[0]
	}

	if 2 != len(twoDuplicateFiles.Filepaths) {
		t.Error("should be 2 duplicate items here")
	}

	if 3 != len(threeDuplicateFiles.Filepaths) {
		t.Error("should be 3 duplicate items here")
	}

	if twoDuplicateFiles.Md5Hash == threeDuplicateFiles.Md5Hash {
		t.Error("twoDuplicateFiles should have a different md5 hash to threeDuplicateFiles")
	}

}
