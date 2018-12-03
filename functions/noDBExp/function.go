package noDBExp

import (
	"github.com/makuc/GCFunctions-Private-Go-Modules/modules/mytest"
	"net/http"
)

func BrezBaze(w http.ResponseWriter, r *http.Request) {
	w.Write(mytest.Message())
}
