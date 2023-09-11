package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	S "github.com/udon-code-sudios/vaidya-signal-service/service"
)

func UselessHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	fmt.Println("Request received for URI:", r.RequestURI, "and method:", r.Method)

	replyBody := `{"message": "Hello world. I am alive."}`

	fmt.Println("[ INFO ] Replying with:", replyBody)
	fmt.Fprintf(w, replyBody)
}

func GetVaidyaSignalsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// log request invocation
	fmt.Println("[ INFO ] Request received for URI:", r.RequestURI, "and method:", r.Method)

	// return wrong method if not GET
	if r.Method != "GET" {
		fmt.Println("[ INFO ] Method", r.Method, "is not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	keys, ok := r.URL.Query()["tickers"]

	if !ok || len(keys[0]) < 1 {
		fmt.Println("[ INFO ] Url Param 'tickers' is missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// convert tickers from query param (comma-separated list) to array of strings
	tickers := strings.Split(keys[0], ",")

	// get signals for tickers
	signals := S.FindAllVaidyaSignalsForTickers(tickers)

	// return signals as JSON
	signalsJSON, _ := json.Marshal(signals)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(signalsJSON)
}
