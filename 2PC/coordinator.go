/*
To implement a Two-Phase Commit (2PC) protocol using MongoDB, Golang, and HTTP REST APIs, we create a coordinator
service and participant services. The coordinator manages the transaction lifecycle, while participants handle 
resource operations atomically.

# Coordinator Service

The coordinator exposes an endpoint to start a transaction, communicates with participants to prepare, and based 
on responses, commits or aborts the transaction.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TransactionRequest struct {
	Participants []Participant `json:"participants"`
}

type Participant struct {
	URL     string      `json:"url"`
	Payload interface{} `json:"payload"`
}

func main() {
	http.HandleFunc("/start-transaction", func(w http.ResponseWriter, r *http.Request) {
		var txReq TransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&txReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Phase 1: Prepare
		allPrepared := true
		for _, p := range txReq.Participants {
			prepared, err := sendPrepareRequest(p.URL, p.Payload)
			if err != nil || !prepared {
				allPrepared = false
				break
			}
		}

		if !allPrepared {
			for _, p := range txReq.Participants {
				sendAbortRequest(p.URL)
			}
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Transaction aborted")
			return
		}

		// Phase 2: Commit
		for _, p := range txReq.Participants {
			if err := sendCommitRequest(p.URL); err != nil {
				fmt.Printf("Commit failed for %s: %v\n", p.URL, err)
			}
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Transaction committed")
	})

	http.ListenAndServe(":8080", nil)
}

func sendPrepareRequest(url string, payload interface{}) (bool, error) {
	jsonPayload, _ := json.Marshal(payload)
	resp, err := http.Post(url+"/prepare", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

func sendCommitRequest(url string) error {
	resp, err := http.Post(url+"/commit", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("commit failed for %s", url)
	}
	return nil
}

func sendAbortRequest(url string) error {
	resp, err := http.Post(url+"/abort", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("abort failed for %s", url)
	}
	return nil
}