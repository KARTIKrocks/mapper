// Example: Func[S, D] for zero-reflection, fully type-safe mapping.
package main

import (
	"fmt"
	"strings"

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

// Declare typed mappers as package-level variables.
// These are zero-reflection and fully compile-time safe.
var (
	toDTO = mapper.NewFunc(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})

	toModel = mapper.NewFunc(func(d UserDTO) User {
		return User{ID: d.ID, Name: d.Name}
	})

	toUpper = mapper.NewFunc(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: strings.ToUpper(u.Name)}
	})
)

func main() {
	user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}

	// Single mapping
	dto := toDTO.Map(user)
	fmt.Printf("To DTO:    %+v\n", dto)

	// Back to model
	model := toModel.Map(dto)
	fmt.Printf("To Model:  %+v\n", model)

	// With transform
	upper := toUpper.Map(user)
	fmt.Printf("Uppercase: %+v\n", upper)

	// Slice mapping
	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}
	dtos := toDTO.MapSlice(users)
	fmt.Printf("Slice:     %+v\n", dtos)
}
