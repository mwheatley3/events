package producer

import (
	"bytes"
	"curbflow/server"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Start emits an event every X milliseconds
func Start(ms time.Duration) {
	types := []string{server.Increment, server.Decrement}
	timeout := time.Millisecond * ms
	for {
		typ := types[rand.Intn(len(types))]
		val := rand.Intn(5) + 1

		url := fmt.Sprintf("%s:%s/event", server.BaseURL, server.Port)
		e := server.Event{Type: typ, Value: val}
		j, _ := json.Marshal(e)
		_, err := http.Post(url, "application/json", bytes.NewBuffer(j))
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %s", err)
			return
		}

		time.Sleep(timeout)
	}
}
