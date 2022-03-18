package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testDir  string = "test_upload"
	testFile string = "test-file-3.txt"
	testData string = "test data for a text file"
)

func TestUpload(t *testing.T) {
	err := createTestFolder()
	if err != nil {
		t.Fatal(err)
	}

	// add file to request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", testFile)
	io.Copy(part, strings.NewReader(testData))
	writer.Close()

	// add body to request
	r := httptest.NewRequest(http.MethodPost, "/upload", body)
	r.Header.Add("Content-Type", writer.FormDataContentType())

	w := &httptest.ResponseRecorder{}
	h := New(testDir)
	h.uploadFunc(w, r)

	// expect file existence
	if _, err := os.Stat(filepath.Join(testDir, testFile)); errors.Is(err, os.ErrNotExist) {
		t.Error("file was not added to work-folder")
	}
	os.RemoveAll(testDir)
}

func TestList(t *testing.T) {
	expectedData := []File{
		{
			Name:        "test-file-1.txt",
			Extension:   ".txt",
			SizeInBytes: 25,
			FileLink:    "/download/test-file-1.txt",
		},
		{
			Name:        "test-file-2.txt",
			Extension:   ".txt",
			SizeInBytes: 25,
			FileLink:    "/download/test-file-2.txt",
		},
	}

	err := createTestFolder()
	if err != nil {
		t.Fatal(err)
	}

	h := New(testDir)
	r := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()
	h.listFilesFunc(w, r)

	outData := make([]File, 0)
	fileData := &File{}
	sc := bufio.NewScanner(w.Result().Body)
	for sc.Scan() {
		json.Unmarshal([]byte(sc.Text()), fileData) // GET the line string
		outData = append(outData, *fileData)
	}

	if len(outData) != len(expectedData) {
		t.Errorf("wrong number of files detected: expected %d, got %d", len(expectedData), len(outData))
	}

	for _, gotData := range outData {
		found := false
		for _, expectData := range expectedData {
			if gotData == expectData {
				found = true
			}
		}
		if !found {
			t.Errorf("got extra element: %+v", gotData)
		}
	}

	for _, expectData := range expectedData {
		found := false
		for _, gotData := range outData {
			if gotData == expectData {
				found = true
			}
		}
		if !found {
			t.Errorf("expected element was not found: %+v", expectData)
		}
	}
	os.RemoveAll(testDir)
}

func createTestFolder() error {
	_, err := os.Stat(testDir)
	if os.IsNotExist(err) {
		err := os.Mkdir(testDir, 0755)
		if err != nil {
			log.Println("cannot create folder", err)
			return err
		}
		f, err := os.Create(filepath.Join(testDir, "test-file-1.txt"))
		if err != nil {
			log.Println("cannot create file 1", err)
			return err
		}
		io.Copy(f, strings.NewReader(testData))
		f.Close()

		f, err = os.Create(filepath.Join(testDir, "test-file-2.txt"))
		if err != nil {
			log.Println("cannot create file 2", err)
			return err
		}
		io.Copy(f, strings.NewReader(testData))
		f.Close()
	}
	return nil
}
