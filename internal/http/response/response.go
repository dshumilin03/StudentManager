package response

import (
	"StudentManager/internal/domain"
)

type Response struct {
	Error    string           `json:"error,omitempty"`
	Student  *domain.Student  `json:"student,omitempty"`
	Students []domain.Student `json:"students,omitempty"`
}

func FoundAllStudents(students []domain.Student) Response {
	return Response{
		Students: students,
	}
}

// TODO Reason to get all to one response with student?
func FoundStudent(student domain.Student) Response {
	return Response{
		Student: &student,
	}
}

func Created(student domain.Student) Response {
	return Response{
		Student: &student,
	}
}

func Updated(student domain.Student) Response {
	return Response{
		Student: &student,
	}
}

func Error(msg string) Response {
	return Response{
		Error: msg,
	}
}
