package utils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func GetAuthenticatedUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(Authenticated).(int)
	return id, ok
}

func GetReqID(ctx context.Context) string {
	return ctx.Value(RidKey).(string)
}

func GetIntUrlParam(r *http.Request, name string) (int, error) {
	s, err := GetUrlParam(r, name)
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		ae := ApiError{
			StatusCode: http.StatusBadRequest,
			Msg:        http.StatusText(http.StatusBadRequest),
		}
		return 0, &ae
	}
	return n, nil
}

func GetUrlParam(r *http.Request, name string) (string, error) {
	value := chi.URLParam(r, name)
	if len(value) == 0 {
		ae := ApiError{
			StatusCode: http.StatusNotFound,
			Msg:        http.StatusText(http.StatusNotFound),
		}
		return value, &ae
	}
	return value, nil
}

func Bind(body io.Reader, out interface{}) error {
	err := json.NewDecoder(body).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

func RenderJson(r *http.Request, w http.ResponseWriter, statusCode int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8 ")
	rid := GetReqID(r.Context())
	if len(rid) > 0 {
		w.Header().Set("X-Request-Id", rid)
	}
	var body []byte
	if res != nil {
		var err error
		body, err = json.Marshal(res)
		if err != nil {
			ae := NewInternalServerError(err.Error())
			statusCode = ae.StatusCode
			body, err = json.Marshal(&ae)
			if err != nil { // this should not happen, but anyway better safe than sorry
				body = []byte(`{"msg": "` + err.Error() + `"}`)
			}
		}
	}
	w.WriteHeader(statusCode)
	if len(body) > 0 {
		w.Write(body)
	}
}
