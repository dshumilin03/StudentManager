package response

import (
	"StudentManager/internal/domain"
	"net/http"
)

type Response struct {
	Status   string           `json:"status,omitempty"`
	Error    string           `json:"error,omitempty"`
	Student  *domain.Student  `json:"student,omitempty"`
	Students []domain.Student `json:"students,omitempty"`
}

func OK() Response {
	return Response{
		Status: http.StatusText(http.StatusOK),
	}
}

func FoundAll(students []domain.Student) Response {
	return Response{
		Status:   http.StatusText(http.StatusFound),
		Students: students,
	}
}

func Created(student domain.Student) Response {
	return Response{
		Status:  http.StatusText(http.StatusCreated),
		Student: &student,
	}
}

func Error(msg string) Response {
	return Response{
		Status: http.StatusText(http.StatusInternalServerError),
		Error:  msg,
	}
}
