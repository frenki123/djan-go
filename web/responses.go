package web

import (
	"fmt"
	"net/http"
)

type Response struct {
	MediaType  string
	Charset    string
	StatusCode int
	Content    []byte
	Headers    map[string]string
}

type htmlError struct {
	info      string
	origError error
}

func (err htmlError) Error() string {
	return fmt.Sprintf("Error on %s. Original error msg is: %s", err.info, err.origError)
}

func HTMLResponse(w http.ResponseWriter, content []byte) error {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(content)
	if err != nil {
		return htmlError{
			info:      "content writer error in html response",
			origError: err,
		}
	}
	return nil
}
