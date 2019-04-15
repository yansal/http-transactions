package httpx

import "net/http"

type Error struct {
	Err  error
	Code int
}

func (e Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return http.StatusText(e.Code)
}
