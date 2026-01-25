package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, map[string]string{"error": message})
}

func GetIDFromPath(r *http.Request, prefix string) (int, error) {
	idStr := strings.TrimPrefix(r.URL.Path, prefix)
	return strconv.Atoi(idStr)
}
