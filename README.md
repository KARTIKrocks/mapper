# mapper

[![Go Reference](https://pkg.go.dev/badge/github.com/KARTIKrocks/mapper.svg)](https://pkg.go.dev/github.com/KARTIKrocks/mapper)
[![Go Report Card](https://goreportcard.com/badge/github.com/KARTIKrocks/mapper)](https://goreportcard.com/report/github.com/KARTIKrocks/mapper)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Fast, type-safe struct mapping for Go using generics. No reflection on struct fields, no code generation, no magic — just explicit mapping functions with a clean API.

## Why mapper?

Every layered Go project (clean architecture, hex architecture, microservices) ends up with hundreds of lines of repetitive DTO-to-model conversion code:

```go
func ToUserDTO(u User) UserDTO {
    return UserDTO{
        ID:    u.ID,
        Name:  u.Name,
        Email: u.Email,
    }
}
```

Multiply that by 40+ structs across 3 layers, and it becomes a maintenance burden. Existing solutions rely heavily on reflection, which sacrifices compile-time safety and performance.

**mapper** solves this with two approaches:

| Approach | Reflection | Allocs | Use case |
|---|---|---|---|
| `Func[S, D]` | Zero | 0 | Hot paths, maximum performance |
| Registry (`Map`, `Register`) | Minimal (type lookup only) | 1 | Convenience, dynamic dispatch |

## Install

```bash
go get github.com/KARTIKrocks/mapper
```

Requires **Go 1.22+** (for generics and `reflect.TypeFor`).

## Quick Start

### Registry API

Register mapping functions once, invoke them anywhere by type:

```go
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

func init() {
    mapper.Register(func(u User) UserDTO {
        return UserDTO{ID: u.ID, Name: u.Name}
    })
}

func main() {
    user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}

    // Single mapping
    dto := mapper.Map[UserDTO](user)
    fmt.Println(dto) // {1 Alice}

    // Slice mapping
    users := []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
    dtos := mapper.MapSlice[User, UserDTO](users)
    fmt.Println(dtos) // [{1 Alice} {2 Bob}]
}
```

### Func API (zero reflection)

For maximum performance and full compile-time safety:

```go
var toDTO = mapper.NewFunc(func(u User) UserDTO {
    return UserDTO{ID: u.ID, Name: u.Name}
})

dto := toDTO.Map(user)          // single
dtos := toDTO.MapSlice(users)   // slice
```

`Func[S, D]` is a thin wrapper around your function — it compiles down to a direct function call with **zero allocations**.

## API Reference

### Global Functions (use the default mapper)

```go
// Register a mapping function from S to D.
mapper.Register(func(s S) D { ... })

// Map a single value. Panics if not registered.
d := mapper.Map[D](src)

// Map a single value. Returns error if not registered.
d, err := mapper.MapErr[D](src)

// Map a slice. Panics if not registered.
ds := mapper.MapSlice[S, D](srcSlice)

// Map a slice. Returns error if not registered.
ds, err := mapper.MapSliceErr[S, D](srcSlice)

// Check if a mapping exists.
ok := mapper.Has[S, D]()

// Clear all global mappings.
mapper.ResetGlobal()
```

### Instance Mapper (for DI and testing)

```go
m := mapper.New()

mapper.RegisterTo(m, func(s S) D { ... })
d := mapper.MapFrom[D](m, src)
d, err := mapper.MapErrFrom[D](m, src)
ds := mapper.MapSliceFrom[S, D](m, srcSlice)
ds, err := mapper.MapSliceErrFrom[S, D](m, srcSlice)
ok := mapper.HasIn[S, D](m)
m.Reset()
```

### Func Type (zero reflection)

```go
f := mapper.NewFunc(func(s S) D { ... })
d := f.Map(src)
ds := f.MapSlice(srcSlice)
```

## Examples

### Nested Structs

```go
mapper.Register(func(u UserWithAddress) UserDTO {
    return UserDTO{
        ID:     u.ID,
        Name:   u.Name,
        Street: u.Address.Street,
        City:   u.Address.City,
    }
})
```

### Custom Transforms

```go
mapper.Register(func(u User) UserDTO {
    return UserDTO{
        ID:   u.ID,
        Name: strings.ToUpper(u.Name),
    }
})
```

### Bidirectional Mapping

```go
mapper.Register(func(u User) UserDTO {
    return UserDTO{ID: u.ID, Name: u.Name}
})
mapper.Register(func(d UserDTO) User {
    return User{ID: d.ID, Name: d.Name}
})

dto := mapper.Map[UserDTO](user)
back := mapper.Map[User](dto)
```

### Pointer Types

```go
mapper.Register(func(u *User) *UserDTO {
    return &UserDTO{ID: u.ID, Name: u.Name}
})

dto := mapper.Map[*UserDTO](&user)
```

### Error Handling

```go
dto, err := mapper.MapErr[UserDTO](user)
if err != nil {
    // no mapping registered from User to UserDTO
    log.Fatal(err)
}
```

### Instance Mapper for Testing

```go
func TestUserService(t *testing.T) {
    m := mapper.New()
    mapper.RegisterTo(m, func(u User) UserDTO {
        return UserDTO{ID: u.ID, Name: u.Name}
    })
    // use m in your service under test — isolated from global state
}
```

## Benchmarks

Run on Intel i5-11400H @ 2.70GHz:

```
BenchmarkManualMap-12                  1000000000    0.24 ns/op    0 B/op    0 allocs/op
BenchmarkFuncMap-12                     461350041    2.69 ns/op    0 B/op    0 allocs/op
BenchmarkRegistryMap-12                  15748872   83.58 ns/op   24 B/op    1 allocs/op

BenchmarkManualMapSlice100-12             1000000    1088 ns/op   2688 B/op   1 allocs/op
BenchmarkFuncMapSlice100-12               1000000    1310 ns/op   2688 B/op   1 allocs/op
BenchmarkRegistryMapSlice100-12            123028    9186 ns/op   9888 B/op  201 allocs/op

BenchmarkFuncMapParallel-12           1000000000    0.60 ns/op      0 B/op   0 allocs/op
BenchmarkRegistryMapParallel-12         12804081   92.36 ns/op     72 B/op   2 allocs/op
```

**Key takeaways:**
- `Func[S, D]` is ~11x slower than manual but has **zero allocations** — practically free
- Registry is ~350x slower than manual but still only **84ns per mapping** — fast enough for any real workload
- Both scale linearly under parallel load with no lock contention

## Design Philosophy

- **No struct field reflection.** Your mapping function is the source of truth — if a field is renamed, the compiler tells you.
- **No magic.** Every mapping is an explicit function you write and can read.
- **Small API.** ~10 functions total. Nothing to configure.
- **Two tiers.** `Func[S,D]` for hot paths, registry for convenience. Pick what fits.

## License

[MIT](LICENSE)
