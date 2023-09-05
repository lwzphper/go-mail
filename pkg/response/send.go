package response

import (
	"encoding/json"
	"fmt"
	"github.com/lwzphper/go-mall/pkg/logger"
	"net/http"
)

type Response struct {
	Page           int         `json:"page,omitempty"`
	PageSize       int         `json:"page_size,omitempty"`
	Code           int         `json:"code"`
	Msg            string      `json:"msg"`
	Data           interface{} `json:"data"`
	httpStatusCode int
	headers        map[string]string
}

type RespOption func(*Response)

func WithPage(page int) RespOption {
	return func(r *Response) {
		r.Page = page
	}
}

func WithPageSize(pageSize int) RespOption {
	return func(r *Response) {
		r.PageSize = pageSize
	}
}

func WithData(data interface{}) RespOption {
	return func(r *Response) {
		r.Data = data
	}
}

func WithMsg(msg string) RespOption {
	return func(r *Response) {
		r.Msg = msg
	}
}

func WithCode(code int) RespOption {
	return func(r *Response) {
		r.Code = code
	}
}

func WithHeaders(headers map[string]string) RespOption {
	return func(r *Response) {
		r.headers = headers
	}
}

func WithAuthHeader(token string) RespOption {
	return func(r *Response) {
		r.headers["Authorization"] = token
	}
}

func WithHttpStatusCode(code int) RespOption {
	return func(r *Response) {
		r.httpStatusCode = code
	}
}

func Success(w http.ResponseWriter, data interface{}, options ...RespOption) {
	SendResponse(w, data, CodeSuccess, options...)
}

func PageSuccess(w http.ResponseWriter, data interface{}, page, pageSize int) {
	SendResponse(w, data, CodeSuccess, WithPage(page), WithPageSize(pageSize))
}

func NotFoundError(w http.ResponseWriter, msg string) {
	SendResponse(w, nil, CodeNotFound, WithMsg(msg), WithHttpStatusCode(http.StatusNotFound))
}

func FormValidError(w http.ResponseWriter, msg string) {
	SendResponse(w, nil, CodeInvalidParams, WithMsg(msg), WithHttpStatusCode(http.StatusBadRequest))
}

func InternalError(w http.ResponseWriter, options ...RespOption) {
	options = append(options, WithHttpStatusCode(http.StatusInternalServerError))
	SendResponse(w, nil, CodeError, options...)
}

func Failed(w http.ResponseWriter, code int, options ...RespOption) {
	SendResponse(w, nil, code, options...)
}

func SendResponse(w http.ResponseWriter, data interface{}, code int, options ...RespOption) {

	resp := &Response{
		Data:           data,
		Code:           code,
		httpStatusCode: http.StatusOK,
		headers:        map[string]string{},
	}

	for _, option := range options {
		option(resp)
	}

	if len(resp.Msg) == 0 {
		resp.Msg = GetMsg(resp.Code)
	}

	respByt, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf(`{"code":"%d", "msg": "encoding to json error, %s"}`, CodeError, err)
		_, err = w.Write([]byte(errMsg))
		if err != nil {
			logger.Error("send response error: " + err.Error()) // 错误默认输出到终端
			return
		}
		return
	}

	// set response header
	w.Header().Set("Content-Type", "application/json")
	for key, val := range resp.headers {
		w.Header().Set(key, val)
	}
	w.WriteHeader(resp.httpStatusCode)

	_, err = w.Write(respByt)
	if err != nil {
		logger.Error("send response error: " + err.Error()) // 错误默认输出到终端
		return
	}
}
