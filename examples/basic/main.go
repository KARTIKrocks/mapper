// Example: basic registry usage with Map and MapSlice.
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

type Address struct {
	Street string
	City   string
}

type UserWithAddress struct {
	ID      int
	Name    string
	Address Address
}

type FlatUserDTO struct {
	ID     int
	Name   string
	Street string
	City   string
}

func init() {
	// Basic mapping
	mapper.Register(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})

	// With custom transform
	mapper.Register(func(u User) string {
		return fmt.Sprintf("%s <%s>", u.Name, u.Email)
	})

	// Nested struct flattening
	mapper.Register(func(u UserWithAddress) FlatUserDTO {
		return FlatUserDTO{
			ID:     u.ID,
			Name:   u.Name,
			Street: u.Address.Street,
			City:   u.Address.City,
		}
	})

	// Bidirectional
	mapper.Register(func(d UserDTO) User {
		return User{ID: d.ID, Name: d.Name}
	})
}

func main() {
	// Single mapping
	user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	dto := mapper.Map[UserDTO](user)
	fmt.Printf("Map:       %+v\n", dto)

	// Custom transform (User -> string)
	display := mapper.Map[string](user)
	fmt.Printf("Transform: %s\n", display)

	// Slice mapping
	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}
	dtos := mapper.MapSlice[User, UserDTO](users)
	fmt.Printf("Slice:     %+v\n", dtos)

	// Nested struct flattening
	rich := UserWithAddress{
		ID:      1,
		Name:    "Alice",
		Address: Address{Street: "123 Main St", City: "Springfield"},
	}
	flat := mapper.Map[FlatUserDTO](rich)
	fmt.Printf("Flatten:   %+v\n", flat)

	// Bidirectional
	back := mapper.Map[User](dto)
	fmt.Printf("Reverse:   %+v\n", back)

	// Error handling
	_, err := mapper.MapErr[int](user) // no int mapping registered
	fmt.Printf("Error:     %v\n", err)

	// Has check
	fmt.Printf("Has User->UserDTO: %v\n", mapper.Has[User, UserDTO]())
	fmt.Printf("Has User->int:     %v\n", mapper.Has[User, int]())

}
