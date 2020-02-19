package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"curbflow/server"

	"github.com/spf13/cobra"
)

// Cmd returns client command
func Cmd() {
	//
	// var value int

	var createEvent = &cobra.Command{
		Use:   "create [increment or decrement] [value]",
		Short: "Post an event to the event server",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			typ := strings.ToUpper(args[0])
			if typ != server.Increment && typ != server.Decrement {
				fmt.Fprintf(os.Stderr, "Type must be either %s or %s", server.Increment, server.Decrement)
				return
			}
			i, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "error converting value to int: %s", err)
			}
			url := "http://localhost:8080/event"
			e := server.Event{Type: typ, Value: i}
			j, _ := json.Marshal(e)
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(j))
			if err != nil {
				fmt.Fprintf(os.Stderr, "err: %s", err)
				return
			}
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "err: %s", err)
				return
			}

			fmt.Printf("succesfully posted %s", b)

		},
	}

	var getAll = &cobra.Command{
		Args:  cobra.MinimumNArgs(0),
		Use:   "get all events",
		Short: "get all events",
		Run: func(cmd *cobra.Command, args []string) {
			url := "http://localhost:8080/events"
			resp, err := http.Get(url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "err: %s", err)
			}

			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "err: %s", err)
				return
			}

			fmt.Fprintf(os.Stdout, "status code: %d\n", resp.StatusCode)
			fmt.Fprintf(os.Stdout, "list of events: %s\n", b)
		},
	}

	var value = &cobra.Command{
		Args:  cobra.MinimumNArgs(0),
		Use:   "value",
		Short: "retrieve the value of the system",
		Run: func(cmd *cobra.Command, args []string) {
			url := "http://localhost:8080/value"
			if len(args[0]) > 0 {
				url = fmt.Sprintf("%s/%s", url, args[0])
			}

			resp, err := http.Get(url)
			fmt.Fprintf(os.Stdout, "\nstatus code: %d\n", resp.StatusCode)
			if err != nil {
				fmt.Fprintf(os.Stderr, "err: %s", err)
			}
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "err: %s", err)
				return
			}

			fmt.Fprintf(os.Stdout, "value: %s\n", b)
		},
	}

	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(createEvent, getAll, value)
	rootCmd.Execute()
}
