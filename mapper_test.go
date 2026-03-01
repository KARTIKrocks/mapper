package mapper_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/KARTIKrocks/mapper"
)

// --- Test types ---

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

type UserWithAddressDTO struct {
	ID     int
	Name   string
	Street string
	City   string
}

// --- Register + Map ---

func TestRegisterAndMap(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})

	user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	dto := mapper.MapFrom[UserDTO](m, user)

	if dto.ID != 1 {
		t.Errorf("ID: got %d, want 1", dto.ID)
	}
	if dto.Name != "Alice" {
		t.Errorf("Name: got %q, want %q", dto.Name, "Alice")
	}
}

func TestMapErr(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})

	dto, err := mapper.MapErrFrom[UserDTO](m, User{ID: 1, Name: "Bob"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.Name != "Bob" {
		t.Errorf("Name: got %q, want %q", dto.Name, "Bob")
	}
}

func TestMapErrNotRegistered(t *testing.T) {
	m := mapper.New()
	_, err := mapper.MapErrFrom[UserDTO](m, User{})
	if err == nil {
		t.Fatal("expected error for unregistered mapping")
	}
	if !strings.Contains(err.Error(), "no mapping registered") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestMapPanicsWhenNotRegistered(t *testing.T) {
	m := mapper.New()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unregistered mapping")
		}
	}()
	mapper.MapFrom[UserDTO](m, User{})
}

func TestMapErrNilSource(t *testing.T) {
	m := mapper.New()
	_, err := mapper.MapErrFrom[UserDTO](m, nil)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// --- MapSlice ---

func TestMapSlice(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})

	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	dtos := mapper.MapSliceFrom[User, UserDTO](m, users)

	if len(dtos) != 2 {
		t.Fatalf("len: got %d, want 2", len(dtos))
	}
	if dtos[0].Name != "Alice" || dtos[1].Name != "Bob" {
		t.Errorf("unexpected: %+v", dtos)
	}
}

func TestMapSliceErr(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})

	dtos, err := mapper.MapSliceErrFrom[User, UserDTO](m, []User{{ID: 1, Name: "A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dtos) != 1 || dtos[0].Name != "A" {
		t.Errorf("unexpected: %+v", dtos)
	}
}

func TestMapSliceErrNotRegistered(t *testing.T) {
	m := mapper.New()
	_, err := mapper.MapSliceErrFrom[User, UserDTO](m, []User{{ID: 1}})
	if err == nil {
		t.Fatal("expected error for unregistered mapping")
	}
}

func TestMapSliceNil(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	result := mapper.MapSliceFrom[User, UserDTO](m, nil)
	if result != nil {
		t.Errorf("got %v, want nil", result)
	}
}

func TestMapSliceEmpty(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	result := mapper.MapSliceFrom[User, UserDTO](m, []User{})
	if len(result) != 0 {
		t.Errorf("got %v, want empty", result)
	}
}

// --- Nested structs ---

func TestNestedStructMapping(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u UserWithAddress) UserWithAddressDTO {
		return UserWithAddressDTO{
			ID:     u.ID,
			Name:   u.Name,
			Street: u.Address.Street,
			City:   u.Address.City,
		}
	})

	user := UserWithAddress{
		ID:      1,
		Name:    "Alice",
		Address: Address{Street: "123 Main St", City: "Springfield"},
	}
	dto := mapper.MapFrom[UserWithAddressDTO](m, user)

	if dto.Street != "123 Main St" {
		t.Errorf("Street: got %q, want %q", dto.Street, "123 Main St")
	}
	if dto.City != "Springfield" {
		t.Errorf("City: got %q, want %q", dto.City, "Springfield")
	}
}

// --- Custom transforms ---

func TestCustomTransform(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: strings.ToUpper(u.Name)}
	})

	dto := mapper.MapFrom[UserDTO](m, User{ID: 1, Name: "alice"})
	if dto.Name != "ALICE" {
		t.Errorf("Name: got %q, want %q", dto.Name, "ALICE")
	}
}

// --- Bidirectional ---

func TestBidirectionalMapping(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	mapper.RegisterTo(m, func(d UserDTO) User {
		return User{ID: d.ID, Name: d.Name}
	})

	user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	dto := mapper.MapFrom[UserDTO](m, user)
	back := mapper.MapFrom[User](m, dto)

	if back.ID != 1 || back.Name != "Alice" {
		t.Errorf("round-trip failed: %+v", back)
	}
	if back.Email != "" {
		t.Errorf("Email should be zero value, got %q", back.Email)
	}
}

// --- Has ---

func TestHas(t *testing.T) {
	m := mapper.New()
	if mapper.HasIn[User, UserDTO](m) {
		t.Error("expected false before registration")
	}
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	if !mapper.HasIn[User, UserDTO](m) {
		t.Error("expected true after registration")
	}
}

// --- Reset ---

func TestReset(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	m.Reset()
	if mapper.HasIn[User, UserDTO](m) {
		t.Error("expected false after reset")
	}
}

// --- Overwrite ---

func TestRegisterOverwrites(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID * 10, Name: "overwritten"}
	})

	dto := mapper.MapFrom[UserDTO](m, User{ID: 1, Name: "Alice"})
	if dto.ID != 10 || dto.Name != "overwritten" {
		t.Errorf("overwrite failed: %+v", dto)
	}
}

// --- Pointer types ---

func TestPointerTypes(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u *User) *UserDTO {
		return &UserDTO{ID: u.ID, Name: u.Name}
	})

	user := &User{ID: 1, Name: "Alice"}
	dto := mapper.MapFrom[*UserDTO](m, user)
	if dto.ID != 1 || dto.Name != "Alice" {
		t.Errorf("pointer mapping failed: %+v", dto)
	}
}

// --- Nil function panics ---

func TestRegisterNilFuncPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil function")
		}
	}()
	m := mapper.New()
	mapper.RegisterTo[User, UserDTO](m, nil)
}

func TestNewFuncNilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil function")
		}
	}()
	mapper.NewFunc[User, UserDTO](nil)
}

// --- Func ---

func TestFunc(t *testing.T) {
	f := mapper.NewFunc(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	dto := f.Map(User{ID: 1, Name: "Alice"})
	if dto.ID != 1 || dto.Name != "Alice" {
		t.Errorf("Func.Map failed: %+v", dto)
	}
}

func TestFuncMapSlice(t *testing.T) {
	f := mapper.NewFunc(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	users := []User{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}
	dtos := f.MapSlice(users)

	if len(dtos) != 2 {
		t.Fatalf("len: got %d, want 2", len(dtos))
	}
	if dtos[0].Name != "A" || dtos[1].Name != "B" {
		t.Errorf("unexpected: %+v", dtos)
	}
}

func TestFuncMapSliceNil(t *testing.T) {
	f := mapper.NewFunc(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	if f.MapSlice(nil) != nil {
		t.Error("expected nil for nil input")
	}
}

func TestFuncMapSliceEmpty(t *testing.T) {
	f := mapper.NewFunc(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	result := f.MapSlice([]User{})
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

// --- Global API ---

func TestGlobalAPI(t *testing.T) {
	mapper.Register(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	defer mapper.Global().Reset()

	if !mapper.Has[User, UserDTO]() {
		t.Fatal("expected Has to return true")
	}

	dto := mapper.Map[UserDTO](User{ID: 42, Name: "Global"})
	if dto.ID != 42 || dto.Name != "Global" {
		t.Errorf("global Map failed: %+v", dto)
	}

	dtos := mapper.MapSlice[User, UserDTO]([]User{{ID: 1, Name: "A"}})
	if len(dtos) != 1 || dtos[0].Name != "A" {
		t.Errorf("global MapSlice failed: %+v", dtos)
	}
}

func TestResetGlobal(t *testing.T) {
	mapper.Register(func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	mapper.ResetGlobal()
	if mapper.Has[User, UserDTO]() {
		t.Error("expected false after ResetGlobal")
	}
}

// --- Interface type behavior ---

type Stringer interface {
	String() string
}

type myString struct {
	val string
}

func (s myString) String() string { return s.val }

func TestInterfaceTypeLookupUsesConcreteType(t *testing.T) {
	// Register with concrete type — lookup uses concrete type at runtime.
	m := mapper.New()
	mapper.RegisterTo(m, func(s myString) string {
		return s.String()
	})

	result := mapper.MapFrom[string](m, myString{val: "hello"})
	if result != "hello" {
		t.Errorf("got %q, want %q", result, "hello")
	}
}

func TestInterfaceTypeMismatchReturnsError(t *testing.T) {
	// Registering with an interface type and calling Map with a concrete type
	// will NOT match, because reflect.TypeOf returns the concrete type.
	m := mapper.New()
	mapper.RegisterTo[Stringer, string](m, func(s Stringer) string {
		return s.String()
	})

	_, err := mapper.MapErrFrom[string](m, myString{val: "hello"})
	if err == nil {
		t.Fatal("expected error: interface registration should not match concrete type at runtime")
	}
}

func TestPointerValueTypeMismatchReturnsError(t *testing.T) {
	// Register func(User) but call with *User — should fail, not silently miss.
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})

	_, err := mapper.MapErrFrom[UserDTO](m, &User{ID: 1, Name: "Alice"})
	if err == nil {
		t.Fatal("expected error: *User should not match User registration")
	}
}

// --- Empty slice returns non-nil empty slice ---

func TestMapSliceEmptyReturnsNonNil(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})
	result := mapper.MapSliceFrom[User, UserDTO](m, []User{})
	if result == nil {
		t.Error("expected non-nil empty slice, got nil")
	}
	if len(result) != 0 {
		t.Errorf("expected len 0, got %d", len(result))
	}
}

// --- Panic value type ---

func TestMapFromPanicsWithString(t *testing.T) {
	m := mapper.New()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		// Verify panic value is a string (from fmt.Sprintf), not an error interface
		if _, ok := r.(string); !ok {
			t.Errorf("expected string panic, got %T: %v", r, r)
		}
	}()
	mapper.MapFrom[UserDTO](m, User{})
}

func TestMapFromPanicsOnNilWithString(t *testing.T) {
	m := mapper.New()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		if _, ok := r.(string); !ok {
			t.Errorf("expected string panic, got %T: %v", r, r)
		}
	}()
	mapper.MapFrom[UserDTO](m, nil)
}

// --- Concurrency ---

func TestConcurrentReads(t *testing.T) {
	m := mapper.New()
	mapper.RegisterTo(m, func(u User) UserDTO {
		return UserDTO{ID: u.ID, Name: u.Name}
	})

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			dto := mapper.MapFrom[UserDTO](m, User{ID: id, Name: "test"})
			if dto.ID != id {
				t.Errorf("ID: got %d, want %d", dto.ID, id)
			}
		}(i)
	}
	wg.Wait()
}

func TestConcurrentReadWrite(t *testing.T) {
	m := mapper.New()

	var wg sync.WaitGroup

	// Writers
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mapper.RegisterTo(m, func(u User) UserDTO {
				return UserDTO{ID: u.ID, Name: u.Name}
			})
		}()
	}

	// Readers
	for i := range 100 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// May or may not find a mapping depending on timing — just shouldn't panic.
			_, _ = mapper.MapErrFrom[UserDTO](m, User{ID: id, Name: "test"})
		}(i)
	}

	wg.Wait()
}
