package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
)

var allTestEvents = []Event{
	Event{Type: Increment, Value: 5},
	Event{Type: Increment, Value: 5},
}

func TestCalcSystemValue(t *testing.T) {
	expected := 10
	allEvents = allTestEvents
	v := calcSystemValue(-1)
	if expected != v {
		t.FailNow()
	}
}

func TestAll(t *testing.T) {
	go Start()
	allEvents = allTestEvents
	url := fmt.Sprintf("%s:%s/events", BaseURL, Port)

	r, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
	}

	// if err != nil {
	// 	fmt.Printf("%#+v\n\n", err)
	// 	t.FailNow()
	// }

	e := []Event{}
	json.NewDecoder(r.Body).Decode(&e)
	if len(e) != len(allTestEvents) {
		fmt.Print("the number of expected events does not match")
		t.FailNow()
	}

	fmt.Printf("%#+v\n\n", e)
}

func TestCreateEvent(t *testing.T) {

	e := Event{
		Type:  "INCREMENT",
		Value: 4,
	}

	url := fmt.Sprintf("%s:%s/event", BaseURL, Port)
	j, _ := json.Marshal(e)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(j))
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		t.FailNow()
	}

	var res Event

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		t.FailNow()
	}

	if e.Value != res.Value {
		fmt.Printf("incorrect event created: %#+v\n\n", e)
		t.FailNow()
	}

	fmt.Printf("res %#+v\n\n", res)

}

func TestAllCreateEvent(t *testing.T) {
	url := fmt.Sprintf("%s:%s/value", BaseURL, Port)

	r, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
	}

	var v1 int
	json.NewDecoder(r.Body).Decode(&v1)
	fmt.Printf("\n first value %d\n", v1)

	e1 := Event{
		Type:  "INCREMENT",
		Value: 4,
	}
	e2 := Event{
		Type:  "DECREMENT",
		Value: 3,
	}
	e3 := Event{
		Type:  "MULTIPLY",
		Value: 4,
	}
	createEvent(t, e1)
	createEvent(t, e2)
	createEvent(t, e3)

	r2, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
	}

	var v int
	json.NewDecoder(r2.Body).Decode(&v)
	fmt.Printf("\n last value %d\n", v)

	if (v1+1)*4 != v {
		t.FailNow()
	}

}

func createEvent(t *testing.T, e Event) error {
	url := fmt.Sprintf("%s:%s/event", BaseURL, Port)
	j, _ := json.Marshal(e)

	_, err := http.Post(url, "application/json", bytes.NewBuffer(j))
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		t.FailNow()
	}
	return nil
}
