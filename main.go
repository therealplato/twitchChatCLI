package main

import (
	"crypto/tls"
	"log"
	"strings"
	"time"

	irc "github.com/fluffle/goirc/client"
	termbox "github.com/nsf/termbox-go"
	"github.com/simplyserenity/twitchOAuth"
)

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
		// panic(err)
	}

	scopes := []string{"chat_login", "user_read"}

	token, err := twitchAuth.GetToken(clientID, scopes)

	if err != nil {
		panic(err)
	}

	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	username := getUsername(token, clientID)
	cfg := irc.NewConfig(strings.ToLower(username))
	cfg.SSL = true
	cfg.SSLConfig = &tls.Config{ServerName: "irc.chat.twitch.tv"}
	cfg.Server = twitch
	cfg.Pass = "oauth:" + token

	client.Username = username

	client.MenuOpen = true

	channel := ""

	var c *irc.Conn

	startMenu := newMenuBox("What channel do you want to talk in? (e.g. vinesauce) or press ESC to quit", func() {
		channel = "#" + strings.ToLower(uInput)
		uInput = ""
		c = connect(cfg, channel)
		client.MenuOpen = false
		client.CurrentChatBox = channel
		client.ChatBoxes[channel] = getChatBox(channel, chat)
		client.Initialized = true
	})

	client.CurrentMenu = startMenu

	//redraw_all()
	go func() {
		for {
			time.Sleep(33 * time.Millisecond)
			redraw_all()
		}
	}()

chat_loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break chat_loop
			case termbox.KeyCtrlC:
				if client.Initialized {
					uInput = client.ChatBoxes[client.CurrentChatBox].Lines[len(client.ChatBoxes[client.CurrentChatBox].Lines)-1].Line
				}
			case termbox.KeyCtrlB:
				if client.MenuOpen && client.Initialized {
					client.MenuOpen = false
				}
			case termbox.KeyF9:
				if client.Initialized {
					updateMode()
				}
			case termbox.KeyF2:
				if client.Initialized {
					addChatChannel(c)
				}
			case termbox.KeyF5:
				if client.Initialized {
					switchChannel()
				}
			case termbox.KeyF3:
				if client.Initialized {
					leaveChatChannel(c)
				}
			case termbox.KeyF12:
				client.SidebarActive = !client.SidebarActive
			case termbox.KeyTab:
				if client.Initialized {
					remainingLetters()
				}
			case termbox.KeySpace:
				uInput += " "
			case termbox.KeyEnter:
				if client.MenuOpen {
					client.CurrentMenu.callback()
				} else {
					if len(uInput) > 1 {
						client.FoundUsername = ""
						newMessage(ChatLine{
							username,
							termbox.ColorBlue,
							uInput,
						}, client.CurrentChatBox)
						c.Privmsg(client.CurrentChatBox, uInput)
						uInput = ""
					}
				}
			case termbox.KeyBackspace2:
				if len(uInput) > 0 {
					uInput = uInput[0 : len(uInput)-1]
				}
			case termbox.KeyBackspace:
				if len(uInput) > 0 {
					uInput = uInput[0 : len(uInput)-1]
				}
			default:
				if ev.Ch != 0 {
					switch client.ChatMode {
					case 0:
						uInput += string(ev.Ch)
					case 1:
						uInput += string(ev.Ch) + " "
					}
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		case termbox.EventResize:
			//redraw_all()
		}
		//redraw_all()
	}
}
