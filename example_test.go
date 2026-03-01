package mapper_test

import (
	"fmt"
	"strings"

	"github.com/KARTIKrocks/mapper"
)

func Example() {
	mapper.Register(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	defer mapper.Global().Reset()

	user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	dto := mapper.Map[UserDTO](user)
	fmt.Printf("ID: %d, Name: %s\n", dto.ID, dto.Name)
	// Output: ID: 1, Name: Alice
}

func Example_mapSlice() {
	mapper.Register(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	defer mapper.Global().Reset()

	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	dtos := mapper.MapSlice[User, UserDTO](users)
	for _, dto := range dtos {
		fmt.Printf("%d:%s ", dto.ID, dto.Name)
	}
	fmt.Println()
	// Output: 1:Alice 2:Bob
}

func Example_customTransform() {
	mapper.Register(func(u User) UserDTO {
		return UserDTO{
			ID:   u.ID,
			Name: strings.ToUpper(u.Name),
		}
	})
	defer mapper.Global().Reset()

	dto := mapper.Map[UserDTO](User{ID: 1, Name: "alice"})
	fmt.Println(dto.Name)
	// Output: ALICE
}

func Example_errorHandling() {
	m := mapper.New()
	// No mapping registered — MapErr returns an error instead of panicking.
	_, err := mapper.MapErrFrom[UserDTO](m, User{ID: 1})
	fmt.Println(err)
	// Output: mapper: no mapping registered from mapper_test.User to mapper_test.UserDTO
}

func ExampleNew() {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})

	dto := mapper.MapFrom[UserDTO](m, User{ID: 42, Name: "Bob"})
	fmt.Println(dto.ID, dto.Name)
	// Output: 42 Bob
}

func ExampleFunc() {
	toDTO := mapper.NewFunc(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})

	dto := toDTO.Map(User{ID: 1, Name: "Alice"})
	fmt.Printf("ID: %d, Name: %s\n", dto.ID, dto.Name)
	// Output: ID: 1, Name: Alice
}

func ExampleFunc_mapSlice() {
	toDTO := mapper.NewFunc(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})

	users := []User{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}
	dtos := toDTO.MapSlice(users)
	fmt.Println(len(dtos), dtos[0].Name, dtos[1].Name)
	// Output: 2 A B
}
