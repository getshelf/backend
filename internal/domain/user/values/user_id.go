package values

type UserID string

func NewUserID(raw string) UserID {
	return UserID(raw)
}

func (id UserID) String() string {
	return string(id)
}
