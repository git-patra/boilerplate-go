package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	response, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

type ApiResponse struct {
	StatusCode int          `json:"status_code"`
	ErrorCode  *string      `json:"error_code,omitempty"`
	Message    string       `json:"message"`
	Data       *interface{} `json:"data,omitempty"`
	MetaData   *interface{} `json:"metadata,omitempty"`
}

type MetaData struct {
	ClientProperties map[string]interface{} `json:"client_properties,omitempty"`
	Pagination       Pagination             `json:"pagination"`
}

type Pagination struct {
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalData  int `json:"total_data"`
}

func SetPagination(page int, size int, total int) (paging *Pagination) {
	var totalPages = 1
	if total > 0 && total > size {
		mod := total % size
		totalPages = (total - mod) / size
		if mod > 0 {
			totalPages = totalPages + 1
		}
	}

	return &Pagination{
		TotalPages: totalPages,
		Page:       page,
		Size:       size,
		TotalData:  total,
	}
}

func GetMetaData(page int, size int, total int) (paging *MetaData) {
	data := SetPagination(page, size, total)

	return &MetaData{Pagination: *data}
}

func BuildErrorResponse(responseCode Code) *ApiResponse {
	return &ApiResponse{
		StatusCode: responseCode.HttpStatusCode,
		ErrorCode:  &responseCode.ErrorCode,
		Message:    responseCode.Message,
	}
}

func BuildSuccessResponseWithoutData(responseCode Code) *ApiResponse {
	return &ApiResponse{
		StatusCode: responseCode.HttpStatusCode,
		Message:    responseCode.Message,
	}
}

func BuildSuccessResponseWithData(responseCode Code, data interface{}) *ApiResponse {
	return &ApiResponse{
		StatusCode: responseCode.HttpStatusCode,
		Message:    responseCode.Message,
		Data:       &data,
	}
}

func BuildSuccessResponseWithDataAndMetaData(responseCode Code, data interface{}, metadata interface{}) *ApiResponse {
	return &ApiResponse{
		StatusCode: responseCode.HttpStatusCode,
		Message:    responseCode.Message,
		Data:       &data,
		MetaData:   &metadata,
	}
}
