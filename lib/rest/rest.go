package rest

import (
	"net/http"
)

func Init() (r *rest) {

	r = &rest{}
	return
}

func (o *rest) Render(w http.ResponseWriter, r *http.Request) {

}
