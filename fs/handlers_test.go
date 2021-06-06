package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadHandleFunc(t *testing.T) {

	// open file to send to server
	file, err := os.Open("testfile.tst")
	if err != nil {
		fmt.Println(err)
		t.Fatal("error opening 'testfile.tst'")
	}
	defer file.Close()
	// write file as a part of form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("error creating multipart form, %v", err)
	}
	writer.Close()
	// create /upload-request
	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	// create ResponseRecorder
	rr := httptest.NewRecorder()
	// call tested function
	uploadHandleFunc(rr, req)
	// check response status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `testfile.tst`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func TestListHandleFunc(t *testing.T) {

	// create /upload-request
	body := &bytes.Buffer{}
	req, _ := http.NewRequest(http.MethodGet, "/list?ext=.tst", body)
	// create ResponseRecorder
	rr := httptest.NewRecorder()
	// call tested function
	listHandleFunc(rr, req)
	// check response status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `testfile.tst`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
