package home

import (
	"fmt"
	"net/http"
)

func CallHTTP(route string, method string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, route, nil) // no body atm
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request, err %v", err)
	}

	// no header atm

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling HTTP request, err %v", err)
	}

	return resp, nil
}

type HTTPEndpoint struct {
	Route  string `json:"route"`
	Method string `json:"method"`

	Conditions []func() bool // tbd
}
