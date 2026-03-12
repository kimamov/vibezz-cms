package domain

// userQuery is the implementation of UserQueryInterface
type userQuery struct {
	*query
}

// NewUserQuery creates a new UserQuery instance
func NewUserQuery() UserQueryInterface {
	return &userQuery{
		query: NewQuery().(*query),
	}
}

// HasRole returns true if role is set
func (q *userQuery) HasRole() bool {
	_, ok := q.properties["role"]
	return ok
}

// Role returns the role
func (q *userQuery) Role() UserRole {
	if val, ok := q.properties["role"].(UserRole); ok {
		return val
	}
	return ""
}

// SetRole sets the role
func (q *userQuery) SetRole(role UserRole) UserQueryInterface {
	q.properties["role"] = role
	return q
}

// HasEmail returns true if email is set
func (q *userQuery) HasEmail() bool {
	_, ok := q.properties["email"]
	return ok
}

// Email returns the email
func (q *userQuery) Email() string {
	if val, ok := q.properties["email"].(string); ok {
		return val
	}
	return ""
}

// SetEmail sets the email
func (q *userQuery) SetEmail(email string) UserQueryInterface {
	q.properties["email"] = email
	return q
}
