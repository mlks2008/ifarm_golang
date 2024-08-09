package api

import (
	"net/http"
)

type Api struct {
}

//能被子类重写(执行Execute前，会先执行Prepare)
func (this *Api) Prepare(w http.ResponseWriter, r *http.Request) bool {
	return true
}

//能被子类重写(执行Execute后，会执行Finish)
func (this *Api) Finish(w http.ResponseWriter, r *http.Request) {

}
