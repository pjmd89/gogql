package rest

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/pjmd89/gogql/lib/resolvers"
)

func (o *rest) RestRender(w http.ResponseWriter, r *http.Request, sessionID string) (isErr bool) {
	isErr = false
	if len(o.objectTypes) > 0 {
		for k, v := range o.objectTypes {
			post := map[string][]string{}
			matchPath, _ := regexp.MatchString(`^`+k, r.URL.Path)
			if err := r.ParseForm(); err == nil {
				post = r.PostForm
			}
			if matchPath {
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
func (o *rest) ObjectType(url, alias string, object resolvers.ObjectTypeInterface) {
	o.objectTypes[url] = ObjectType{Alias: alias, ObjectType: object}
}
