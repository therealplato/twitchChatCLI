package main

import (
	"sync"

	"github.com/nsf/termbox-go"
)

type Client struct {
	Lock                                                         *sync.Mutex
	ChatBoxes                                                    map[string]*ChatBox
	ChatBoxList                                                  []string
	UserColors                                                   map[string]termbox.Attribute
	MenuOpen, Initialized, SidebarActive, Locked, SearchMenuOpen bool
	CurrentMenu                                                  MenuBox
	ChatMode, CurrentChatBoxInt                                  int
	Username, FoundUsername, CurrentChatBox                      string
}

func (c *Client) BoxLines() []ChatLine {
	b, ok := client.ChatBoxes[client.CurrentChatBox]
	if !ok {
		return nil
	}
	ll := b.Lines
	return ll
}

func GetClient() Client {
	return Client{
		ChatBoxes:         make(map[string]*ChatBox),
		ChatBoxList:       make([]string, 0),
		UserColors:        make(map[string]termbox.Attribute),
		MenuOpen:          false,
		Initialized:       false,
		SidebarActive:     true,
		Lock:              &sync.Mutex{},
		Locked:            false,
		SearchMenuOpen:    false,
		CurrentMenu:       MenuBox{},
		ChatMode:          0,  // what does zero mean?
		CurrentChatBoxInt: -1, // what does zero mean?
		CurrentChatBox:    "",
		Username:          "",
		FoundUsername:     "",
	}
}
