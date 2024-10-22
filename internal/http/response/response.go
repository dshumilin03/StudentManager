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

func FoundAllStudents(students []domain.Student) Response {
	return Response{
		Students: students,
	}
}

func FoundAllGroups(groups []domain.Group) Response {
	return Response{
		Groups: groups,
	}
}

// TODO Reason to get all to one response with student?
func FoundStudent(student domain.Student) Response {
	return Response{
		Student: &student,
	}
}

func FoundGroup(group domain.Group) Response {
	return Response{
		Group: &group,
	}
}

// TODO make it abstract
func StudentCreated(student domain.Student) Response {
	return Response{
		Student: &student,
	}
}

func CreatedGroup(group domain.Group) Response {
	return Response{
		Group: &group,
	}
}

func StudentUpdated(student domain.Student) Response {
	return Response{
		Student: &student,
	}
}

func UpdatedGroup(group domain.Group) Response {
	return Response{
		Group: &group,
	}
}

func Error(msg string) Response {
	return Response{
		Error: msg,
	}
}
