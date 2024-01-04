package help

import (
	"encoding/json"
	"net/http"
)

func HelpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		helpData := map[string]interface{}{
			"/habits": map[string]interface{}{
				"DELETE": map[string]string{
					"Description": "Not supported yet, will be used for deleting a habit.",
				},
				"GET": map[string]interface{}{
					"Description": "Retrieve all habits or a specific habit.",
					"Details":     "Use an empty body {} for all habits or {\"id\": \"int\"} for a specific habit.",
				},
				"POST": map[string]interface{}{
					"Description": "Create a new habit.",
					"Details":     "Provide details in the body: {\"name\": \"habit name\", \"frequency\": [\"Monday\", \"Tuesday\", ... up to \"Sunday\"]}.",
				},
				"PUT": map[string]string{
					"Description": "Not supported yet, will be used for updating a habit.",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(helpData)
	}
}
