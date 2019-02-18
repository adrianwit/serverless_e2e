package logger

import (
	"cloud.google.com/go/logging"
	"io/ioutil"
	"net/http"
)

// LogFn logs payload to customized stream
func LogFn(w http.ResponseWriter, r *http.Request) {
	client, logger, err := newLogger()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Close()
	if r.ContentLength == 0 {
		return
	}
	payload, err := ioutil.ReadAll(r.Body)
	stdLogger := logger.StandardLogger(logging.Info)
	stdLogger.Println(string(payload))
}
