package response

import (
	"StudentManager/internal/domain"
)

type Response struct {
	Error    string           `json:"error,omitempty"`
	Student  *domain.Student  `json:"student,omitempty"`
	Students []domain.Student `json:"students,omitempty"`
	Groups   []domain.Group   `json:"groups,omitempty"`
	Group    *domain.Group    `json:"group,omitempty"`
}

func StudentResponse(student domain.Student) Response {
	return Response{
		Student: &student,
	}
}

func StudentsResponse(students []domain.Student) Response {
	return Response{
		Students: students,
	}
}

func GroupsResponse(groups []domain.Group) Response {
	return Response{
		Groups: groups,
	}
}

func GroupResponse(group domain.Group) Response {
	return Response{
		Group: &group,
	}
}

func Error(msg string) Response {
	return Response{
		Error: msg,
	}
}
