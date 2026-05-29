package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// helper function to initialize an HTTP endpoint to accept GET requests
func MakeGetHandler(
	responseGenerator func() map[string]any,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Debug().Str("url", r.URL.Path).
			Str("request_addr", r.RemoteAddr).
			Str("method", r.Method).
			Msg("")

		w.Header().Set("Content-Type", "application/json")
		b, _ := json.MarshalIndent(responseGenerator(), "", "  ")
		fmt.Fprintf(w, string(b))
	}
}

// Helper function to create an HTTP POST endpoint and populate a specific ouptut struct
func MakePostHandler[T any](callback func(T) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var outputStruct T
		if err := json.NewDecoder(r.Body).Decode(&outputStruct); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Debug().Str("url", r.URL.Path).
			Str("request_addr", r.RemoteAddr).
			Str("method", r.Method).
			Msg(fmt.Sprintf("%+v", outputStruct))

		if err := callback(outputStruct); err != nil {
			log.Err(err).Str("url", r.URL.Path).Msg("Callback failed")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]string{"message": "submission successfully received."}
		json.NewEncoder(w).Encode(response)
	}
}
