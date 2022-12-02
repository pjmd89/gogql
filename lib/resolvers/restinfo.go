package resolvers

import "net/http"

func (o *RestInfo) SetHTTPRequest(r *http.Request) {
	o.r = r
}
func (o *RestInfo) SetHeader(header, value string) {
	if len(o.headers) == 0 {
		o.headers = make(map[string]string)
	}
	o.headers[header] = value
}
func (o *RestInfo) GetHeaders() map[string]string {
	return o.headers
}
