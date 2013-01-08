package dundee

import (
	"dundee/config"
	"dundee/cuepoints"
	"dundee/liveevents"
	"dundee/streams"
	"errors"
	"fmt"
	"net/http"
)

var conf = config.Get()

func CuePointsHandler(w http.ResponseWriter, r *http.Request) {

	streamID := r.FormValue("streamid")
	if streamID == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, errors.New("You must include a valid streamid."))
		return
	}

	cuePoint, err := cuepoints.New(r)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err)
		return
	}

	stream, err := streams.Find(streamID, w)
	if err != nil {
		fmt.Fprint(w, err)
	}

	//Beyond this point the client doesn't care - return 201
	w.WriteHeader(201)
	fmt.Fprint(w, franchise)

	go func() {

		liveEvents, err := liveevents.Retrieve(conf.Elementals)
		if err != nil {
			fmt.Println(err)
			return
		}

		eventID, elemental, err := liveevents.Find(stream, liveEvents)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = cuepoints.Inject(eventID, elemental, cuePoint)
		if err != nil {
			fmt.Println(err)
			return
		}

	}()
}