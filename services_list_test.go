package pool

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/gateway-fm/prover-pool-lib/service"
)

func TestServicesListShuffle(t *testing.T) {
	numServices := 20
	numTries := 50000
	numServicesToPick := 200 // how many we pick for one try

	threshold := 0.1
	duration, _ := time.ParseDuration("5s")

	var services []service.IService
	for i := 1; i < numServices; i++ {
		srv := newHealthyService(fmt.Sprintf("https://%dgateway.fm", i))
		services = append(services, srv)
	}

	var selectedServices []string
	for i := 1; i < numTries; i++ {
		srvList := NewServicesList("name", &ServicesListOpts{
			TryUpTries:     5,
			TryUpInterval:  duration,
			ChecksInterval: duration,
		})

		for _, srv := range services {
			srvList.Add(srv)
		}
		srvList.Shuffle()

		for j := 1; j < numServicesToPick; j++ {
			selectedServices = append(selectedServices, srvList.Next().ID())
		}
	}

	selectionFrequency := make(map[string]int)
	for _, srvID := range selectedServices {
		selectionFrequency[srvID] = selectionFrequency[srvID] + 1
	}

	expected := selectionFrequency[services[0].ID()]
	for _, count := range selectionFrequency {
		if math.Abs(float64(count-expected))/float64(expected) > threshold {
			t.Errorf("selection frequencies are different more than a threshold")
		}
	}
}

func TestServicesListTryUp(t *testing.T) {
	list := NewServicesList("testServicesList", &ServicesListOpts{
		TryUpTries:     5,
		TryUpInterval:  1 * time.Second,
		ChecksInterval: 1 * time.Second,
	})

	srv := &healthyService{0, &service.BaseService{}}

	list.Add(srv)
	list.FromHealthyToJail(srv.ID())

	if list.CountAll() != 1 {
		t.Errorf("unexpected num of services in list")
	}

	list.TryUpService(srv, 0)

	if list.Next() == nil {
		t.Errorf("unexpected no healthy services")
	}
}
