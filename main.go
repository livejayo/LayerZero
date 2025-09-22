package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// A struct to hold our metrics
type NodeStatus struct {
	BlockHeight string `json:"blockHeight"`
	ChainID     string `json:"chainID"`
	Syncing     string `json:"syncing"`
}

var rpcURL string

func main() {
	rpcURL = os.Getenv("GOAT_RPC_NODE")
	if rpcURL == "" {
		// As per the exercise, fall back to a public node if the env var isn't set.
		rpcURL = "https://rpc.goat.network"
		log.Printf("GOAT_RPC_NODE not set, using public node: %s", rpcURL)
	} else {
		log.Printf("Monitoring Goat RPC node at: %s", rpcURL)
	}

	http.HandleFunc("/", statusHandler)
	log.Println("Simple monitor starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// statusHandler fetches the data and writes it to the webpage
func statusHandler(w http.ResponseWriter, r *http.Request) {
	status, err := fetchNodeStatus()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch node status: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Goat Node Status:\n\n")
	fmt.Fprintf(w, "Chain ID: %s\n", status.ChainID)
	fmt.Fprintf(w, "Current Block Height: %s\n", status.BlockHeight)
	fmt.Fprintf(w, "Syncing Status: %s\n", status.Syncing)
}

// --- Helper Functions to call the RPC Endpoint ---

type rpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

func fetchNodeStatus() (*NodeStatus, error) {
	// For simplicity, we make 3 separate calls.
	blockNumHex, err := callRPC("eth_blockNumber", nil)
	if err != nil {
		return nil, err
	}

	chainIDHex, err := callRPC("eth_chainId", nil)
	if err != nil {
		return nil, err
	}

	syncingResult, err := callRPC("eth_syncing", nil)
	if err != nil {
		return nil, err
	}

	status := &NodeStatus{
		BlockHeight: hexToString(blockNumHex),
		ChainID:     hexToString(chainIDHex),
		Syncing:     fmt.Sprintf("%v", syncingResult),
	}
	return status, nil
}

func callRPC(method string, params []interface{}) (interface{}, error) {
	reqBody, _ := json.Marshal(rpcRequest{"2.0", method, params, 1})
	resp, err := http.Post(rpcURL, "application/json", strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp struct {
		Result interface{} `json:"result"`
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &rpcResp)

	return rpcResp.Result, nil
}

func hexToString(hexVal interface{}) string {
	hexStr, ok := hexVal.(string)
	if !ok {
		return "N/A"
	}
	var intVal int
	fmt.Sscanf(hexStr, "0x%x", &intVal)
	return fmt.Sprintf("%d", intVal)
}
