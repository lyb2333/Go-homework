package domain

import "time"

type User struct {
	Id       string
	Email    string
	Password string
	Nickname string
	Phone    string
	Birthday string
	Aboutme  string

	// UTC 0 的时区
	Ctime time.Time

	//Addr Address
}

//type Address struct {
//	Province string
//	Region   string
//}

//func (u User) ValidateEmail() bool {
// 在这里用正则表达式校验
//return u.Email
//}
