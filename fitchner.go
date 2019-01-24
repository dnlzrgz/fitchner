package fitchner

import (
	"fmt"
	"net/http"
)

func Fetch(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("bad request: %v", err)
	}

	if checkBadStatus(resp.StatusCode) {
		return nil, fmt.Errorf("bad status code %v at %v", resp.StatusCode, req.URL)
	}

	return resp, nil
}

func checkBadStatus(s int) bool {
	return s != http.StatusOK
}
