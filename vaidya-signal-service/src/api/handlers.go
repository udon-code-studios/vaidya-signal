package api

import (
	S "github.com/udon-code-sudios/vaidya-signal-service/service"

	"fmt"
	"net/http"
)

func UselessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received for URI:", r.RequestURI, "and method:", r.Method)
	fmt.Fprintf(w, `{"message": "hello world."}`)
	S.ExportedFunction()
}
