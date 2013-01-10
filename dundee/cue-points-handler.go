package dundee

import (
	"errors"
	"fmt"
	"github.com/ign/ipl-dundee/dundee/cuepoints"
	"github.com/ign/ipl-dundee/dundee/liveevents"
	"github.com/ign/ipl-dundee/dundee/streams"
	"log"
	"net/http"
)

func CuePointsHandler(w http.ResponseWriter, r *http.Request, c *Config) {

	streamID, cuePointType, err := fetchParams(r)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err)
		log.Println(err)
		return
	}

	cuePoint, err := cuepoints.New(cuePointType, r)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err)
		log.Println(err)
		return
	}

	stream, err := resolveStream(streamID, w, c)
	if err != nil {
		fmt.Fprint(w, err)
		log.Println(err)
		return
	}

	go injectCuePoint(stream, cuePoint, c)

	w.WriteHeader(201)
	fmt.Fprint(w, stream.Name)
}

func injectCuePoint(stream *streams.Stream, cuePoint interface{}, c *Config) {
	liveEvents := liveevents.Gather(c.Elementals)

	liveEvent, err := liveevents.Find(stream, liveEvents)
	if err != nil {
		log.Println(err)
		return
	}

	err = cuepoints.Inject(liveEvent, cuePoint)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Successfully injected cuepoint into event:", liveEvent.Name)
}

func resolveStream(streamID string, w http.ResponseWriter, c *Config) (*streams.Stream, error) {
	streamData, err := streams.Fetch(c.Streams_url)
	if err != nil {
		w.WriteHeader(500)
		return nil, err
	}

	streamList, err := streams.Parse(streamData)
	if err != nil {
		w.WriteHeader(500)
		return nil, err
	}

	stream, err := streams.Find(streamID, streamList)
	if err != nil {
		w.WriteHeader(400)
		return nil, err
	}

	return stream, nil
}

func fetchParams(r *http.Request) (string, string, error) {
	streamID := r.FormValue("streamid")
	if streamID == "" {
		return "", "", errors.New("A streamid must be included.")
	}

	cuePointType := r.FormValue("cue-point-type")
	if cuePointType == "" {
		return "", "", errors.New("A cue-point-type must be included.")
	}

	return streamID, cuePointType, nil
}
