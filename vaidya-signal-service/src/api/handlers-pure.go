package api

import (
	"encoding/json"

	S "github.com/udon-code-sudios/vaidya-signal-service/service"

	"fmt"
	"net/http"
)

func UselessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received for URI:", r.RequestURI, "and method:", r.Method)
	fmt.Fprintf(w, `{"message": "hello world."}`)
	S.ExportedFunction()
}

func GetVaidyaSignalsHandler(w http.ResponseWriter, r *http.Request) {
	// log request invocation
	fmt.Println("[ INFO ] Request received for URI:", r.RequestURI, "and method:", r.Method)

	// return wrong method if not GET
	if r.Method != "GET" {
		fmt.Println("[ INFO ] Method", r.Method, "is not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	keys, ok := r.URL.Query()["ticker"]

	if !ok || len(keys[0]) < 1 {
		fmt.Println("[ INFO ] Url Param 'ticker' is missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	signals := S.GetHistoricalVaidyaSignals(keys[0])

	signalsJSON, _ := json.Marshal(signals)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(signalsJSON)
}
