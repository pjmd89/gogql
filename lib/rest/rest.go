package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pjmd89/gogql/lib/resolvers"
)

func Init() (r *Rest) {
	r = &Rest{}
	r.objectTypes = make(map[string]ObjectType)
	return
}

func (o *Rest) RestRender(w http.ResponseWriter, r *http.Request, sessionID string) (isErr bool) {
	isErr = false
	if len(o.objectTypes) > 0 {
		for k, v := range o.objectTypes {
			post := map[string][]string{}
			if err := r.ParseForm(); err == nil {
				post = r.PostForm
			}
			if k == r.URL.Path {
				resolverInfo := resolvers.ResolverInfo{
					Operation: "rest",
					Resolver:  v.Alias,
					SessionID: sessionID,
					RestInfo: &resolvers.RestInfo{
						Path:      r.URL.Path,
						PathSplit: strings.Split(r.URL.Path, "/"),
						GET:       r.URL.Query(),
						POST:      post,
						Writer:    w,
					},
				}
				resolverInfo.RestInfo.SetHTTPRequest(r)
				response, restError := v.ObjectType.Resolver(resolverInfo)
				if restError != nil {
					isErr = true
				}
				for header, value := range resolverInfo.RestInfo.GetHeaders() {
					w.Header().Set(header, value)
				}
				if response != nil {
					fmt.Fprint(w, response)
				}
				break
			}
		}
	}
	return
}
func (o *Rest) ObjectType(url, alias string, object resolvers.ObjectTypeInterface) {
	o.objectTypes[url] = ObjectType{Alias: alias, ObjectType: object}
}
