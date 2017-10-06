package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	irc "github.com/fluffle/goirc/client"
	"crypto/tls"
	"github.com/simplyserenity/twitchOAuth"
	"net/http"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"strings"
	"math/rand"
	"strconv"
	"time"
)

const twitch = "irc.chat.twitch.tv:443"
const clientID = "dlpf1993tub698zw0ic6jlddt9e893"

const chat = "chat"


var sideBarOptions = []string{
	"Here's a list of nifty tricks",
	"Ctrl+C to copy the last meme",
	"F2: Start a new chat tab",
	"F3: Leave current chat tab",
	"F5: Change chat tabs",
	"F9: Change typing modes",
	"F12: Hide this menu",
}

var client = GetClient()
var uInput = ""
var totalMessages = 0

//TODO add functionality like exiting menus if you don't want them & remove totalMessage counter since it was just for debugging(maybe)

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
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

	startMenu := newMenuBox("What channel do you want to talk in? (e.g. vinesauce)", func(){
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
	go func(){
		for{
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
				uInput = client.ChatBoxes[client.CurrentChatBox].Lines[len(client.ChatBoxes[client.CurrentChatBox].Lines)-1].Line
			case termbox.KeyF9:
				updateMode()
			case termbox.KeyF2:
				addChatChannel(c)
			case termbox.KeyF5:
				switchChannel()
			case termbox.KeyF3:
				leaveChatChannel(c)
			case termbox.KeyF12:
				client.SidebarActive = !client.SidebarActive
			case termbox.KeyTab:
				remainingLetters()
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
				if len(uInput) > 0{
					uInput = uInput[0:len(uInput)-1]
				}
			case termbox.KeyBackspace:
				if len(uInput) > 0 {
					uInput = uInput[0:len(uInput)-1]
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


func redraw_all() {
	if client.Locked { //used to prevent this entire thing for a bit because if we don't we can get some weird nil pointers / seg faults
		return
	}

	w, h := termbox.Size()

	const coldef = termbox.ColorDefault

	termbox.Clear(coldef, coldef)

	pos := 0

	if client.MenuOpen {
		userString := "Input: " + uInput
		tbprint(w / 2 - len(client.CurrentMenu.title) / 2, h / 2, termbox.ColorWhite, coldef, client.CurrentMenu.title)
		tbprint(w / 2 - len("Input: " + uInput) / 2, h / 2 + 2, termbox.ColorWhite, coldef, userString)      //TODO make this more flexible
		termbox.SetCursor(len(userString) + w / 2 - len("Input: " + uInput) / 2, h / 2 + 2)
	} else if client.Initialized {
		//draw the current chat lines
		for i := len(client.ChatBoxes[client.CurrentChatBox].Lines) - 1; i >= 0; i-- {
			//this is a little dense lmao
			tbprint(1, h - pos - 2, termbox.ColorWhite, coldef, client.ChatBoxes[client.CurrentChatBox].Lines[i].Line)
			tbprint(1, h - pos - 2, client.ChatBoxes[client.CurrentChatBox].Lines[i].NickColor, coldef, client.ChatBoxes[client.CurrentChatBox].Lines[i].Nick + ": ")
			tbprint(1 + len(client.ChatBoxes[client.CurrentChatBox].Lines[i].Nick + ": "), h - pos - 2, termbox.ColorWhite, termbox.ColorDefault, client.ChatBoxes[client.CurrentChatBox].Lines[i].Line)
			pos++
		}


		//draw the chat tabs at the top of the screen

		for i := 0; i < w; i++{
			termbox.SetCell(i, 0, ' ', coldef, coldef)
		}

		tabText := "Current channels: "

		tbprint(0, 0, termbox.ColorWhite, coldef, tabText)

		var j = len(tabText) + 1
		for i := 0; i < len(client.ChatBoxList); i++{

			if client.CurrentChatBox == client.ChatBoxList[i] {
				tbprint(j, 0, termbox.ColorBlack, termbox.ColorWhite, client.ChatBoxList[i])
			} else {
				tbprint(j, 0, termbox.ColorWhite, termbox.ColorBlack, client.ChatBoxList[i])
			}

			j += len(client.ChatBoxList[i]) + 1
		}


		//draw the sidebar if it's active

		if client.SidebarActive {
			for i := 0; i < len(sideBarOptions); i++{
				tbprint(w - 30, i, termbox.ColorWhite, coldef, sideBarOptions[i])
			}

			//TODO make chat mode more flexible maybe map it or something
			modeText := "current mode: "
			tbprint(w - 30, len(sideBarOptions), termbox.ColorWhite, coldef, modeText)
			if client.ChatMode == 0 {
				tbprint(w - 30 + len(modeText), len(sideBarOptions), termbox.ColorBlue, coldef, "normal")
			} else {
				tbprint(w - 30 + len(modeText), len(sideBarOptions), termbox.ColorRed, coldef, "fullwidth")
			}
			tbprint(w - 30, len(sideBarOptions) + 1, termbox.ColorWhite, coldef, "Total Messages received: " + strconv.Itoa(totalMessages))
		}


		//draw user input
		userString := client.Username + ": " + uInput
		tbprint(1, h - 1, termbox.ColorBlue, coldef, userString)
		termbox.SetCursor(len(userString) + 1, h - 1)

		//draw @ search menu
		if strings.Contains(uInput, "@") && strings.Compare(atUsername(uInput), client.FoundUsername) != 0{
			client.SearchMenuOpen = true

			possibleNames := Filter(client.ChatBoxes[client.CurrentChatBox].UserNames, func(v string)bool{
				return strings.Contains(v, atUsername(uInput))
			})

			if len(possibleNames) > 0 {
				client.FoundUsername = possibleNames[0]
			}

			for i := 0; i < len(possibleNames); i++ {
				tbprint(1, h-i-2, termbox.ColorWhite, termbox.ColorWhite, "                             ")
				tbprint(1, h-i-2, client.UserColors[possibleNames[i]], termbox.ColorWhite, possibleNames[i])
			}
		} else {
			client.SearchMenuOpen = false
		}
	}

	termbox.Flush()
}

//get the user's twitch username through the twitch API
func getUsername(token, cId string) (username string){
	rclient := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/kraken/user", nil)
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("Client-ID", cId)
	req.Header.Set("Authorization", "OAuth " + token)

	res, err := rclient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	name, err :=jsonparser.GetString(body, "display_name")

	if err != nil {
		panic(err)
	}

	return name
}

func connect(cfg *irc.Config, channel string)*irc.Conn{
	c := irc.Client(cfg)

	c.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(channel)
	})


	c.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
		totalMessages++
		if !Contains(client.ChatBoxes[line.Target()].UserNames, line.Nick){
			client.ChatBoxes[line.Target()].UserNames = append(client.ChatBoxes[line.Target()].UserNames, line.Nick)
		}

		newMessage(ChatLine{
			Nick: line.Nick,
			NickColor: getNickColor(line.Nick),
			Line: line.Args[1],
		}, line.Target())
	})

	if err := c.Connect(); err != nil {
		panic(err)
	}

	return c
}

func addChatChannel(conn *irc.Conn){
	client.MenuOpen = true

	newChannel := ""

	newChannelMenu := newMenuBox("What channel do you want to talk in? (e.g. vinesauce)", func(){
		client.Locked = true //although it shouldn't be an issue here, better safe than sorry
		newChannel = "#" + strings.ToLower(uInput)
		uInput = ""
		conn.Join(newChannel)
		client.MenuOpen = false
		client.CurrentChatBox = newChannel
		client.ChatBoxes[newChannel] = getChatBox(newChannel, chat)
		client.Initialized = true
		client.CurrentChatBoxInt++
		client.Locked = false
		//redraw_all()
	})

	client.CurrentMenu = newChannelMenu
}

func leaveChatChannel(conn *irc.Conn){
	client.Locked = true
	if len(client.ChatBoxList) > 1 {
		debugMessage("trying to leave")
		conn.Raw("LEAVE " + client.CurrentChatBox)
		client.ChatBoxList = append(client.ChatBoxList[:client.CurrentChatBoxInt], client.ChatBoxList[client.CurrentChatBoxInt+1:]...)
		delete(client.ChatBoxes, client.CurrentChatBox)

		if client.CurrentChatBoxInt > 0 {
			client.CurrentChatBoxInt--
		}

		client.CurrentChatBox = client.ChatBoxList[client.CurrentChatBoxInt]
	}
	client.Locked = false
	//redraw_all()
}

func switchChannel(){
	client.Locked = true
	if client.CurrentChatBoxInt == len(client.ChatBoxList) - 1 {
		client.CurrentChatBoxInt = 0
	} else {
		client.CurrentChatBoxInt++
	}

	client.CurrentChatBox = client.ChatBoxList[client.CurrentChatBoxInt]
	client.Locked = false
	//redraw_all()
}

func updateMode(){
	if client.ChatMode > 0 {
		client.ChatMode = 0
	} else {
		client.ChatMode++
	}
}

func getChatBox(name, chatType string)*ChatBox{
	client.ChatBoxList = append(client.ChatBoxList, name)
	return NewChatBox(name, chatType)
}

func debugMessage(text string){
	newMessage(ChatLine{
		Line: text,
		Nick: "",
		NickColor: termbox.ColorWhite,
	}, client.CurrentChatBox)
}

func newMessage(chatLine ChatLine, channel string){
	_, h := termbox.Size()

	if len(client.ChatBoxes[channel].Lines) == h - 2 {
		client.ChatBoxes[channel].Lines = client.ChatBoxes[channel].Lines[1:]
	}

	client.ChatBoxes[channel].Lines = append(client.ChatBoxes[channel].Lines, chatLine)
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func getNickColor(nick string)termbox.Attribute{
	if client.UserColors[nick] == termbox.ColorDefault{
		client.UserColors[nick] = getRandomColor()
	}

	return client.UserColors[nick]
}

func atUsername(parsed string)(username string){
	username = ""
	parsing := false

	for i := 0; i < len(parsed); i++ {
		if parsed[i] == '@' && !parsing{
			parsing = true
		} else if parsed[i] != ' ' && parsing {
			username += string(parsed[i])
		} else if parsed[i] == ' ' && parsing{
			break
		}
	}

	return username
}

func remainingLetters(){
	if !client.SearchMenuOpen {
		return
	}

	newInput := ""

	for i := 0; i < len(uInput); i++ {
		if uInput[i] == '@' {
			newInput += string('@')
			newInput += client.FoundUsername
			break
		} else {
			newInput += string(uInput[i])
		}
	}
	uInput = newInput
}

func getRandomColor()(termbox.Attribute){
	rng := rand.Intn(5)
	switch rng {
	case 0:
		return termbox.ColorBlue
	case 1:
		return termbox.ColorRed
	case 2:
		return termbox.ColorGreen
	case 3:
		return termbox.ColorCyan
	case 4:
		return termbox.ColorMagenta
	case 5:
		return termbox.ColorYellow
	default:
		return termbox.ColorCyan
	}
}