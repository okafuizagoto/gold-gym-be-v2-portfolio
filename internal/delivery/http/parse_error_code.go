package http

import (
	"gold-gym-be/pkg/response"
	"strings"
)

// ParseErrorCode ...
func ParseErrorCode(err string) response.Response {
	errResp := response.Error{}

	switch {
	case strings.Contains(err, "401"):
		errResp = response.Error{
			Status: false,
			Msg:    "Unauthorized",
			Code:   401,
		}
	case strings.Contains(err, "10001"):
		errResp = response.Error{
			Status: false,
			Msg:    "Failed to fetch data",
			Code:   10001,
		}
	case strings.Contains(err, "10002"):
		errResp = response.Error{
			Status: false,
			Msg:    "Failed to insert data",
			Code:   10001,
		}
	}

	errResp.Msg = errResp.Msg + " | " + err

	return response.Response{
		Error: errResp,
	}
}
