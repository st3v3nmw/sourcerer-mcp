package testdata

// User struct for testing method parsing
type User struct {
	ID   int
	Name string
}

// GetName is a value receiver method
func (u User) GetName() string {
	return u.Name
}

// SetName is a pointer receiver method
func (u *User) SetName(name string) {
	u.Name = name
}

// Service struct for testing multiple types with methods
type Service struct {
	users []User
}

// AddUser adds a user to the service
func (s *Service) AddUser(user User) {
	s.users = append(s.users, user)
}

// FindUser finds a user by ID
func (s *Service) FindUser(id int) *User {
	for i := range s.users {
		if s.users[i].ID == id {
			return &s.users[i]
		}
	}

	return nil
}

// Count returns the number of users
func (s Service) Count() int {
	return len(s.users)
}

// Generic type with methods
type Repository[T any] struct {
	items []T
}

// Add adds an item to the repository
func (r *Repository[T]) Add(item T) {
	r.items = append(r.items, item)
}

// Get retrieves an item by index
func (r Repository[T]) Get(index int) T {
	return r.items[index]
}

// Size returns the number of items
func (r Repository[T]) Size() int {
	return len(r.items)
}

// Methods with same name on different receiver types
type ServiceA struct{}
type ServiceB struct{}

func (s ServiceA) Helper() string {
	return "service A helper"
}

func (s ServiceB) Helper() string {
	return "service B helper"
}
