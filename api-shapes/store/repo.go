package store

func List() []User {
	return users
}

func Retrieve(id string) User {
	for _, u := range users {
		if u.ID == id {
			return u
		}
	}
	return User{}
}

func Create(user User) User {
	users = append(users, user)
	return user
}

func Update(user User) User {
	for i, u := range users {
		if u.ID == user.ID {
			users[i] = user
			return users[i]
		}
	}
	return User{}
}

func Delete(id string) {
	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
		}
	}
}
