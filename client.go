package main

import "github.com/nsf/termbox-go"

type Client struct {
	ChatBoxes map[string] *ChatBox
	ChatBoxList []string
	UserColors map[string] termbox.Attribute
	MenuOpen, Initialized, SidebarActive, Locked, SearchMenuOpen bool
	CurrentMenu MenuBox
	ChatMode, CurrentChatBoxInt int
	Username, FoundUsername, CurrentChatBox string
}

func GetClient()Client{
	return Client {
		ChatBoxes: make(map[string]*ChatBox),
		ChatBoxList: make([]string, 0),
		UserColors: make(map[string] termbox.Attribute),
		MenuOpen: false,
		Initialized: false,
		SidebarActive: true,
		Locked: false,
		SearchMenuOpen: false,
		CurrentMenu: MenuBox{},
		ChatMode: 0,
		CurrentChatBoxInt: -1,
		CurrentChatBox: "",
		Username: "",
		FoundUsername: "",
	}
}