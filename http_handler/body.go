package http_handler

import (
	"encoding/json"
	"net/http"
)

func ParseBody(r *http.Request, input interface{}, logger Logger) error {
	if r.Body == nil {
		logger.DebugError("empty body error", nil)
		return ErrEmptyBody
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.DebugError("body parse error: %s", err)
		return ErrBodyParseFailed
	}

	return nil
}
