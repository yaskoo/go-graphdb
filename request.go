package graphdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type Part struct {
	Key      string
	Type     string
	Filename string
	Value    io.Reader
}

func (p *Part) ContentType() string {
	if p.Type == "" {
		return "application/octet-stream"
	}
	return p.Type
}

type RequestConfig func(req *http.Request)

type ResponseHandler func(resp *http.Response) error

func Query(key, val string) RequestConfig {
	return func(req *http.Request) {
		q := req.URL.Query()
		q.Add(key, val)
		req.URL.RawQuery = q.Encode()
	}
}

func Header(key, val string) RequestConfig {
	return func(req *http.Request) {
		req.Header.Set(key, val)
	}
}

func MultipartFormData(parts ...Part) RequestConfig {
	return func(req *http.Request) {
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)

		go func() {
			defer pw.Close()

			for _, p := range parts {
				partHeader := textproto.MIMEHeader{}
				if p.Filename != "" {
					partHeader.Set("content-disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, p.Key, p.Filename))
				} else {
					partHeader.Set("content-disposition", fmt.Sprintf(`form-data; name="%s";`, p.Key))
				}
				partHeader.Set("content-type", p.ContentType())

				part, err := writer.CreatePart(partHeader)
				if err != nil {
					_ = pw.CloseWithError(err)
					return
				}

				_, err = io.Copy(part, p.Value)
				if err != nil {
					_ = pw.CloseWithError(err)
					return
				}
			}

			_ = writer.Close()
		}()

		req.Header.Set("content-type", writer.FormDataContentType())
		req.Body = pr
	}
}

func JsonBody(v any) RequestConfig {
	return func(req *http.Request) {
		req.Header.Set("content-type", "application/json")
		data, err := json.Marshal(v)
		if err != nil {
			return
		}
		req.Body = io.NopCloser(bytes.NewBuffer(data))
	}
}
