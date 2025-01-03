package exception

import (
	"boilerplate-go/internal/delivery/rest/response"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

func RecoverFromPanicWithMonitor(w http.ResponseWriter, ctx context.Context) {
	if err := recover(); err != nil {
		formattedError := fmt.Errorf("%+v %s", err, string(debug.Stack()))

		logrus.Errorf(fmt.Sprint(formattedError))

		httpResponse := response.BuildErrorResponse(response.GenericServerError)
		response.JSON(w, httpResponse.StatusCode, httpResponse)
	}
}

func HandleError(ctx context.Context, err error) *response.ApiResponse {
	logrus.Errorf(err.Error())
	return response.BuildErrorResponse(response.GenericServerError)
}
