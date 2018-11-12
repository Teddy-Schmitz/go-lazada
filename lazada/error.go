package lazada

import (
	"fmt"
	"net/http"
	"strings"
)

// Error response is used to return as much data as possible to the calling application to help with dealing with any API issues.
type ErrorResponse struct {
	Response *http.Response

	Code      string          `json:"code"`
	Type      string          `json:"type"`
	Message   string          `json:"message"`
	RequestID string          `json:"request_id"`
	Detail    []*ErrorDetails `json:"detail,omitempty"`
}

type ErrorDetails struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (r *ErrorResponse) Error() string {

	if len(r.Detail) > 0 {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%v %v: %v %v %v \n",
			r.Response.Request.Method, r.Response.Request.URL, r.Code, r.Type, r.Message))

		for _, v := range r.Detail {
			sb.WriteString(fmt.Sprintf("Field: %s, Message: %s \n", v.Field, v.Message))
		}

		return sb.String()
	}

	return fmt.Sprintf("%v %v: %v %v %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Code, r.Type, r.Message)
}
