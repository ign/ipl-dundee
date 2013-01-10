package liveevents

import (
	"github.com/ign/ipl-dundee/dundee/elemental"
)

const liveEventsPath = "/live_events"

func Fetch(elementalServer *elemental.ElementalServer) ([]byte, error) {
	req, err := elemental.GenerateRequest("GET", elementalServer, liveEventsPath, nil)
	if err != nil {
		return nil, err
	}

	_, body, err := elemental.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}

	return body, nil
}
