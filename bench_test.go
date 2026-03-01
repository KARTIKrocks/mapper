package mapper_test

import (
	"runtime"
	"testing"

	"github.com/KARTIKrocks/mapper"
)

// --- Single mapping benchmarks ---

func BenchmarkManualMap(b *testing.B) {
	user := User{ID: 1, Name: "Alice"}
	var result UserDTO
	for b.Loop() {
		result = UserDTO{ID: user.ID, Name: user.Name}
	}
	runtime.KeepAlive(result)
}

func BenchmarkFuncMap(b *testing.B) {
	f := mapper.NewFunc(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	user := User{ID: 1, Name: "Alice"}
	var result UserDTO
	for b.Loop() {
		result = f.Map(user)
	}
	runtime.KeepAlive(result)
}

func BenchmarkRegistryMap(b *testing.B) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	user := User{ID: 1, Name: "Alice"}
	var result UserDTO
	for b.Loop() {
		result = mapper.MapFrom[UserDTO](m, user)
	}
	runtime.KeepAlive(result)
}

// --- Slice mapping benchmarks (100 elements) ---

func BenchmarkManualMapSlice100(b *testing.B) {
	users := make([]User, 100)
	for i := range users {
		users[i] = User{ID: i, Name: "User"}
	}
	var result []UserDTO
	for b.Loop() {
		dtos := make([]UserDTO, len(users))
		for i, u := range users {
			dtos[i] = UserDTO{ID: u.ID, Name: u.Name}
		}
		result = dtos
	}
	runtime.KeepAlive(result)
}

func BenchmarkFuncMapSlice100(b *testing.B) {
	f := mapper.NewFunc(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	users := make([]User, 100)
	for i := range users {
		users[i] = User{ID: i, Name: "User"}
	}
	var result []UserDTO
	for b.Loop() {
		result = f.MapSlice(users)
	}
	runtime.KeepAlive(result)
}

func BenchmarkRegistryMapSlice100(b *testing.B) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	users := make([]User, 100)
	for i := range users {
		users[i] = User{ID: i, Name: "User"}
	}
	var result []UserDTO
	for b.Loop() {
		result = mapper.MapSliceFrom[User, UserDTO](m, users)
	}
	runtime.KeepAlive(result)
}

// --- Parallel benchmarks ---

func BenchmarkRegistryMapParallel(b *testing.B) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	user := User{ID: 1, Name: "Alice"}
	b.RunParallel(func(pb *testing.PB) {
		var result UserDTO
		for pb.Next() {
			result = mapper.MapFrom[UserDTO](m, user)
		}
		runtime.KeepAlive(result)
	})
}

func BenchmarkFuncMapParallel(b *testing.B) {
	f := mapper.NewFunc(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	user := User{ID: 1, Name: "Alice"}
	b.RunParallel(func(pb *testing.PB) {
		var result UserDTO
		for pb.Next() {
			result = f.Map(user)
		}
		runtime.KeepAlive(result)
	})
}
