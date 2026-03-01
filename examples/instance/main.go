// Example: instance-based Mapper for dependency injection and testing.
package main

import (
	"fmt"

	"github.com/KARTIKrocks/mapper"
)

type User struct {
	ID    int
	Name  string
	Email string
}

type UserDTO struct {
	ID   int
	Name string
}

// UserService demonstrates using an instance mapper for DI.
type UserService struct {
	mapper *mapper.Mapper
}

func NewUserService(m *mapper.Mapper) *UserService {
	return &UserService{mapper: m}
}

func (s *UserService) GetDTO(user User) UserDTO {
	return mapper.MapFrom[UserDTO](s.mapper, user)
}

func (s *UserService) GetDTOs(users []User) []UserDTO {
	return mapper.MapSliceFrom[User, UserDTO](s.mapper, users)
}

func main() {
	// Create an isolated mapper instance
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})

	// Inject into service
	svc := NewUserService(m)

	user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	dto := svc.GetDTO(user)
	fmt.Printf("Single: %+v\n", dto)

	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	dtos := svc.GetDTOs(users)
	fmt.Printf("Slice:  %+v\n", dtos)

	// Reset for testing scenarios
	m.Reset()
	fmt.Printf("After reset, has mapping: %v\n", mapper.HasIn[User, UserDTO](m))
}
