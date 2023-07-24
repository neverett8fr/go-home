package home

import (
	"fmt"
	"net/http"
	"time"
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

type Home struct {
	Endpoints map[string]HTTPEndpoint `json:"http_endpoints"`
	stopCh    chan struct{}
}

func NewHome() (Home, error) {
	return Home{
		Endpoints: make(map[string]HTTPEndpoint),
		stopCh:    make(chan struct{}),
	}, nil
}

func (h *Home) RegisterHTTPEndpoint(name string, route string, method string, conditions ...func() bool) error {

	h.Endpoints[name] = HTTPEndpoint{
		Route:      route,
		Method:     method,
		Conditions: conditions,
	}

	return nil
}

func (h *Home) AddCondition(name string, conditions ...func() bool) error {

	endpoint, exists := h.Endpoints[name]
	if !exists {
		return fmt.Errorf("err, endpoint '%s' does not exist", name)
	}

	endpoint.Conditions = append(endpoint.Conditions, conditions...)
	h.Endpoints[name] = endpoint

	return nil
}

func (h *Home) StopHandlers() {
	close(h.stopCh)
}

func (h *Home) StartHandlers() {

	ticker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ticker.C:
			for _, endpoint := range h.Endpoints {
				for _, condition := range endpoint.Conditions {
					if condition() {
						// call method
						_, err := CallHTTP(endpoint.Route, endpoint.Method)
						if err != nil {
							fmt.Printf("error calling function, err %v", err)
						}
						break
					}
				}
			}

		case <-h.stopCh:
			return
		}
	}

}

// func main() {
// 	home, _ := NewHome()

// 	home.RegisterHTTPEndpoint("endpoint1", "192.168....", "GET")
// 	home.AddCondition("endpoint1", func() bool {
// 		return true
// 	})

// 	go home.StartHandlers()

// 	time.Sleep(5 * time.Minute)
// 	home.StopHandlers()

// }
