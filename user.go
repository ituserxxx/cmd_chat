package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Con net.Conn
}

func NewUser(name string,con net.Conn)*User {
	return &User{
		Name:name,
		Con:con,
	}
}

func (u *User)Send(msg string)  {
	_,err := u.Con.Write([]byte(msg))
	if err != nil {
		fmt.Println("send msg fail",err.Error())
	}
}