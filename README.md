# Student Manager System

## Entities

### Student
Student is an entity that has:
- id
- fullName
- groupNumber
- age
- email

### Group (a group of students)
Group includes a bunch of students, it has:
- id
- groupNumber

## Api possibilities

### Student Service

- Add student
- Get student
- Update student data
- Delete student

### Group Service

- Add group
- Get group
- Update group data
- Delete group

### Student Manager Service
It's a structure that provides different operations with Groups and Students
- add Student to Group
- remove student from group
- contains Student Service
- contains Group Service
