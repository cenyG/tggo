package heroku

import (
"fmt"
"net/http"
	"os"
)

func MockServer() {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is mock web for Heroku.")
	})

	http.ListenAndServe(":"+os.Getenv(`PORT`), nil)
}