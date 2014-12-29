package fileutil

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestAppendFile(t *testing.T) {
	f, err := ioutil.TempFile("/tmp", "fastfoodtest")
	if err != nil {
		t.Errorf("unexpected error creating tempfile %v", err)
	}

	// Get the filename and close the file
	fileName := f.Name()
	f.Close()
	defer os.Remove(fileName)

	testString := "this is my test string"
	err = AppendFile(fileName, testString)
	if err != nil {
		t.Errorf("did not expect error appending to file %s: %v", fileName, err)
	}

	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Errorf("did not expect error reading file %s: %v", fileName, err)
	}

	if !strings.Contains(string(fileContents), testString) {
		t.Errorf("file contents did not contain %s", testString)
	}
}

func TestFileExist(t *testing.T) {
	f, err := ioutil.TempFile("/tmp", "fastfoodtest")
	if err != nil {
		t.Errorf("unexpected error creating tempfile %v", err)
	}

	// Get the filename and close the file
	fileName := f.Name()
	defer f.Close()
	defer os.Remove(fileName)

	if !FileExist(fileName) {
		t.Errorf("expected FileExist(%s) to return true")
	}
}

func TestCopy(t *testing.T) {
	f, err := ioutil.TempFile("/tmp", "fastfoodtest")
	if err != nil {
		t.Errorf("unexpected error creating tempfile %v", err)
	}

	fileName := f.Name()
	testString := "this is a test string added by go"
	defer f.Close()
	defer os.Remove(fileName)

	_, err = f.WriteString(testString)
	if err != nil {
		t.Errorf("unexpected error writing to file %s: %v", fileName, err)
	}

	destFile := "/tmp/fastfoodtest_copy_dest"
	if err := Copy(fileName, destFile); err != nil {
		t.Errorf("unexpected error copying file %s to %s: %v", fileName, destFile, err)
	}

	defer os.Remove(destFile)
	destContents, err := ioutil.ReadFile(destFile)
	if err != nil {
		t.Errorf("error reading destination file %s: %v", destFile, err)
	}

	if !strings.Contains(string(destContents), testString) {
		t.Errorf("file contents of %s did not contain %s", destFile, testString)
	}
}
