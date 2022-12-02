package rest

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/pjmd89/gogql/lib/resolvers"
)

func Init(serverName string) (r *rest) {
	r = &rest{}
	r.serverName = serverName
	r.objectTypes = make(map[string]ObjectType)
	return
}

func (o *rest) RestRender(w http.ResponseWriter, r *http.Request, sessionID string) {

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
					},
				}
				resolverInfo.RestInfo.SetHTTPRequest(r)
				fmt.Println()
				response, _ := v.ObjectType.Resolver(resolverInfo)
				for header, value := range resolverInfo.RestInfo.GetHeaders() {
					w.Header().Set(header, value)
				}
				fmt.Fprint(w, response)
				break
			}
		}
	}
}
func (o *rest) GetServerName() string {
	return o.serverName
}
func (o *rest) ObjectType(url, alias string, object resolvers.ObjectTypeInterface) {
	o.objectTypes[url] = ObjectType{Alias: alias, ObjectType: object}
}
