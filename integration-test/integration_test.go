package integration_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/Eun/go-hit"
)

const (
	// Attempts connection
	host       = "127.0.0.1:8080"
	healthPath = "http://" + host + "/health"
	attempts   = 20

	// HTTP REST
	basePath = "http://" + host + "/v1"
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

// HTTP POST: /ohlc/data.
func TestHTTPPostData(t *testing.T) {
	url := "http://" + host + "/v1/ohlc/data"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("sample.csv")
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("csv", filepath.Base("sample.csv"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	Test(t,
		Description("Do PostData Success"),
		Post(basePath+"/ohlc/data"),
		Request().Set(req),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`"status":`),
	)
}

// HTTP GET: /ohlc/data.
func TestHTTPGetData(t *testing.T) {
	Test(t,
		Description("Do GetData Success"),
		Get(basePath+"/ohlc/data"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`"unix":`),
	)
}
