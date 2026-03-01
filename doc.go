// Package mapper provides fast, type-safe struct mapping using Go generics.
//
// mapper eliminates repetitive DTO ↔ Model conversion boilerplate while
// maintaining compile-time safety. The actual mapping functions are written
// by you — the library manages registration, lookup, and slice operations.
//
// # Registry API
//
// Register mapping functions and invoke them by type:
//
//	mapper.Register(func(u User) UserDTO {
//	    return UserDTO{ID: u.ID, Name: u.Name}
//	})
//
//	dto := mapper.Map[UserDTO](user)
//	dtos := mapper.MapSlice[User, UserDTO](users)
//
// For dependency injection or testing, create a dedicated Mapper instance:
//
//	m := mapper.New()
//	mapper.RegisterTo(m, func(u User) UserDTO { ... })
//	dto := mapper.MapFrom[UserDTO](m, user)
//
// # Func API
//
// For zero-reflection, fully type-safe mapping with no registry overhead:
//
//	var toDTO = mapper.NewFunc(func(u User) UserDTO {
//	    return UserDTO{ID: u.ID, Name: u.Name}
//	})
//
//	dto := toDTO.Map(user)
//	dtos := toDTO.MapSlice(users)
package mapper
