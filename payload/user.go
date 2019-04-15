package payload

import (
	"encoding/json"
	"net/http"

	"github.com/yansal/http-transactions/httpx"
)

type User struct{ Email string }

func NewUser(r *http.Request) (*User, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, httpx.Error{Code: http.StatusUnsupportedMediaType}
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var p User
	err := dec.Decode(&p)
	return &p, err
}
