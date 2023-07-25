package home

import (
	"fmt"
	"log"
	"time"
)

type Endpoint interface {
	ConditionMet() bool
	ConditionAdd(conditions ...func() bool) error

	FunctionDo() error
}

type Home struct {
	Endpoints map[string]Endpoint `json:"http_endpoints"`
	stopCh    chan struct{}
}

func NewHome() (Home, error) {
	return Home{
		Endpoints: make(map[string]Endpoint),
		stopCh:    make(chan struct{}),
	}, nil
}

func (h *Home) RegisterHTTPEndpoint(name string, route string, method string, conditions ...func() bool) error {

	h.Endpoints[name] = &HTTPEndpoint{
		Route:      route,
		Method:     method,
		Conditions: conditions,
	}

	return nil
}

func (h *Home) AddCondition(name string, conditions ...func() bool) error {

	_, exists := h.Endpoints[name]
	if !exists {
		return fmt.Errorf("err, endpoint '%s' does not exist", name)
	}

	return h.Endpoints[name].ConditionAdd(conditions...)
}

func (h *Home) StopHandlers() {
	close(h.stopCh)
}

func (h *Home) StartHandlers() {

	ticker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ticker.C:
			for name, endpoint := range h.Endpoints {
				if endpoint.ConditionMet() {
					// call method
					err := endpoint.FunctionDo()
					if err != nil {
						log.Printf("error calling function, err %v", err)
					}
					log.Printf("function with name %v called", name)
					break
				}

			}

		case <-h.stopCh:
			log.Printf("stop called")
			return
		}
	}

}
