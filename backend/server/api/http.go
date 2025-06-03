package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ApiResponseData struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

func ApiError(err string) []byte {
	ret := &ApiResponseData{
		Data:  "",
		Error: err,
	}
	slog.Error("API call failed", "error", err)
	respBuf, _ := json.Marshal(ret)
	return respBuf
}

func ApiResponse(data any, bMarshal bool) []byte {
	var ret any
	if bMarshal {
		ret = &ApiResponseData{
			Data:  data,
			Error: "",
		}
	} else {
		ret = &struct {
			Data  json.RawMessage `json:"data"`
			Error string          `json:"error"`
		}{
			Data:  json.RawMessage(data.([]byte)),
			Error: "",
		}
		// return data.([]byte)
	}
	respBuf, _ := json.Marshal(ret)
	return respBuf
}

func ApiResponseList(data any, total int) []byte {
	retList := &struct {
		Data  any `json:"data"`
		Total int `json:"total"`
	}{
		Data:  data,
		Total: total,
	}
	ret := &ApiResponseData{
		Data:  retList,
		Error: "",
	}
	respBuf, _ := json.Marshal(ret)
	return respBuf
}

func WebResponseAny(w http.ResponseWriter, r *http.Request, data []byte, contentType string, code int) {
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	w.Write(data)
	slog.Info(fmt.Sprintf("%s %s", r.Method, r.URL.Path), "len(data)", len(data), "contentType", contentType)
}

func WebResponseJson(w http.ResponseWriter, r *http.Request, data []byte, code int) {
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

// ====================================== server =========================================
const c_params_parsed = "parsed"

func WebParams(r *http.Request) url.Values {
	ret := make(url.Values, 0)
	if r.Header.Get(c_params_parsed) != "true" {
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			if err := r.ParseMultipartForm(32 << 20); err != nil {
				slog.Error(fmt.Sprintf("%s %s", r.Method, r.URL.Path), "error", err.Error())
				return ret
			}
		} else {
			if err := r.ParseForm(); err != nil {
				slog.Error(fmt.Sprintf("%s %s", r.Method, r.URL.Path), "error", err.Error())
				return ret
			}
		}
	}

	for k := range r.URL.Query() {
		ret.Add(k, r.URL.Query().Get(k))
	}
	for k := range r.Header {
		ret.Add(k, r.Header.Get(k))
	}
	for k := range r.Form {
		if ret.Get(k) == "" {
			ret.Add(k, r.Form.Get(k))
		}
	}
	for k := range r.PostForm {
		if ret.Get(k) == "" {
			ret.Add(k, r.PostForm.Get(k))
		}
	}
	r.Header.Add(c_params_parsed, "true")
	return ret
}

func WebBody[T any](r *http.Request) ([]byte, *T, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}
	t := new(T)
	return body, t, json.Unmarshal(body, t)
}
