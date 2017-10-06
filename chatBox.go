package main

import "github.com/nsf/termbox-go"

//I really should have simplified this design by introducing a
//user struct that would contain their nick and color...

type ChatLine struct {
	Nick string
	NickColor termbox.Attribute
	Line string
}

type ChatBox struct {
	Lines []ChatLine
	UserNames []string
	Name string
}

func NewChatBox(name, chatType string)*ChatBox{
	return &ChatBox{
		Lines: []ChatLine{
			{
				Nick: "",
				NickColor: termbox.ColorDefault,
				Line: "Talking in the " + name + " " + chatType + "! :",
			},
		},
		UserNames: []string{
		},
		Name: name,
	}
}