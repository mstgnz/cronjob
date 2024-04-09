package models

// users should be able to define cronjob for the urls they approve.
type UserUrl struct {
	UserId int    `json:"user_id" validate:"required"`
	Url    string `json:"url" validate:"required"`
}

func (uu *UserUrl) GetUserUrl(userId int) *UserUrl {
	return uu
}

func (uu *UserUrl) CreateUserUrl(userId int, url string) *UserUrl {
	return uu
}

func (uu *UserUrl) UpdateUserUrl(id, userId int, url string) *UserUrl {
	return uu
}

func (uu *UserUrl) DeleteUserUrl(id, userId int) *UserUrl {
	return uu
}
