package transport

import "api-shapes/store"

type UserReq struct {
	Name string `xml:"name" json:"name"`
}

type UserRes struct {
	ID   string `xml:"id" json:"id"`
	Name string `xml:"name" json:"name"`
}

func (u *UserRes) Bind(user store.User) {
	u.ID = user.ID
	u.Name = user.Name
}
