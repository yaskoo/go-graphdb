package graphdb

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	PathReport       = "/rest/report"
	PathReportStatus = "/rest/report/status"

	StateNone       = "NONE"
	StateInProgress = "IN_PROGRESS"
	StateReady      = "READY"
	StateError      = "ERROR"
)

type ReportStatus struct {
	State      string
	Time       time.Time
	ErrMessage string
}

type ReportClient struct {
	client *Client
}

func (r *ReportClient) Generate(ctx context.Context, config ...RequestConfig) (bool, error) {
	var accepted bool
	return accepted, r.client.post(ctx, PathReport, nil, func(resp *http.Response) error {
		if resp.StatusCode == http.StatusInternalServerError {
			return errors.New("report: internal server error")
		}

		accepted = true
		return nil
	}, config...)
}

func (r *ReportClient) Status(ctx context.Context, config ...RequestConfig) (ReportStatus, error) {
	var status ReportStatus
	return status, r.client.get(ctx, PathReportStatus, func(resp *http.Response) error {
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		parts := strings.Split(string(all), "|")
		for i, part := range parts {
			switch i {
			case 0:
				status.State = part
			case 1:
				millis, _ := strconv.ParseInt(part, 10, 64)
				status.Time = time.Unix(0, millis)
			case 2:
				status.ErrMessage = part
			}
		}
		return nil
	}, config...)
}

func (r *ReportClient) Download(ctx context.Context, consumer func(filename string, r io.Reader) error, config ...RequestConfig) error {
	return r.client.get(ctx, PathReport, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			all, err := io.ReadAll(resp.Body)
			if err != nil {
				return errors.Join(errors.New("report"), errors.New(string(all)), err)
			}
		}

		filename, err := extractFilename(resp.Header.Get("content-disposition"))
		if err != nil || filename == "" {
			filename = "report.zip"
		}
		return consumer(filename, resp.Body)
	}, config...)
}
