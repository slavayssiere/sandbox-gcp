package laststat

import (
	"fmt"
	"net/http"
)

// LastStat test function
func LastStat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
