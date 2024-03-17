package transport

import "api-shapes/store"

type UserReq struct {
	Name string `json:"name"`
}

type UserRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (u *UserRes) Bind(user store.User) {
	u.ID = user.ID
	u.Name = user.Name
}
