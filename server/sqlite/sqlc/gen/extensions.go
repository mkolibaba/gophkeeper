package sqlc

func (l Login) GetUser() string {
	return l.User
}

func (n Note) GetUser() string {
	return n.User
}

func (b Binary) GetUser() string {
	return b.User
}

func (c Card) GetUser() string {
	return c.User
}
