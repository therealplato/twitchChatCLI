package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/davecgh/go-spew/spew"
	irc "github.com/fluffle/goirc/client"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const twitch = "irc.chat.twitch.tv:443"

var clientID = os.Getenv("TWITCH_TOKEN")

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
//add name highlighting when people use @username in their messages
//add emote mode maybe?

func redraw_all() {
	spew.Dump(client.BoxLines())
	if client.Locked { //used to prevent this entire thing for a bit because if we don't we can get some weird nil pointers / seg faults
		return
	}

	/*
		w, h := termbox.Size()

		const coldef = termbox.ColorDefault

		termbox.Clear(coldef, coldef)

		pos := 0

		if client.MenuOpen {
			userString := "Input: " + uInput
			tbprint(w/2-len(client.CurrentMenu.title)/2, h/2, termbox.ColorWhite, coldef, client.CurrentMenu.title)
			tbprint(w/2-len("Input: "+uInput)/2, h/2+2, termbox.ColorWhite, coldef, userString) //TODO make this more flexible
			termbox.SetCursor(len(userString)+w/2-len("Input: "+uInput)/2, h/2+2)
		} else if client.Initialized {
			termbox.Flush()
			return
		}
		//draw the current chat lines
		// for i := len(client.Boxes[client.Box].Lines)
		// for i := len(client.BoxLines())
		for i := len(client.ChatBoxes[client.CurrentChatBox].Lines) - 1; i >= 0; i-- {
			tbprint(1, h-pos-2, client.ChatBoxes[client.CurrentChatBox].Lines[i].NickColor, coldef, client.ChatBoxes[client.CurrentChatBox].Lines[i].Nick+": ")

			if strings.Contains(client.ChatBoxes[client.CurrentChatBox].Lines[i].Line, "@") {

				if strings.ToLower(atUsername(client.ChatBoxes[client.CurrentChatBox].Lines[i].Line)) == strings.ToLower(client.Username) {
					tbprint(1+len(client.ChatBoxes[client.CurrentChatBox].Lines[i].Nick+": "), h-pos-2, termbox.ColorBlack, termbox.ColorWhite, client.ChatBoxes[client.CurrentChatBox].Lines[i].Line)
				} else {
					userName := strings.ToLower(atUsername(client.ChatBoxes[client.CurrentChatBox].Lines[i].Line))

					userColor := client.UserColors[userName]

					if userColor == termbox.ColorDefault {
						userColor = getRandomColor()
						client.UserColors[userName] = userColor
					}

					parsing := false

					for j := 0; j < len(client.ChatBoxes[client.CurrentChatBox].Lines[i].Line); j++ {

						if parsing {

							if client.ChatBoxes[client.CurrentChatBox].Lines[i].Line[j] == ' ' {
								parsing = false
							}

							tbprint(1+len(client.ChatBoxes[client.CurrentChatBox].Lines[i].Nick+": ")+j, h-pos-2, userColor, coldef, string(client.ChatBoxes[client.CurrentChatBox].Lines[i].Line[j]))

						} else if client.ChatBoxes[client.CurrentChatBox].Lines[i].Line[j] == '@' {
							parsing = true
							tbprint(1+len(client.ChatBoxes[client.CurrentChatBox].Lines[i].Nick+": ")+j, h-pos-2, userColor, coldef, string(client.ChatBoxes[client.CurrentChatBox].Lines[i].Line[j]))
						} else {
							tbprint(1+len(client.ChatBoxes[client.CurrentChatBox].Lines[i].Nick+": ")+j, h-pos-2, termbox.ColorWhite, coldef, string(client.ChatBoxes[client.CurrentChatBox].Lines[i].Line[j]))
						}

					}

				}

			} else {
				tbprint(1+len(client.ChatBoxes[client.CurrentChatBox].Lines[i].Nick+": "), h-pos-2, termbox.ColorWhite, termbox.ColorDefault, client.ChatBoxes[client.CurrentChatBox].Lines[i].Line)
				// tbprint(1+len(client.BoxLines()[i].Nick+": "), h-pos-2, termbox.ColorWhite, termbox.ColorDefault, client.BoxLines()[i].Line)
				// for i := len(client.BoxLines())
			}

			pos++
		}

		//draw the chat tabs at the top of the screen

		for i := 0; i < w; i++ {
			termbox.SetCell(i, 0, ' ', coldef, coldef)
		}

		tabText := "Current channels: "

		tbprint(0, 0, termbox.ColorWhite, coldef, tabText)

		var j = len(tabText) + 1
		for i := 0; i < len(client.ChatBoxList); i++ {

			if client.CurrentChatBox == client.ChatBoxList[i] {
				tbprint(j, 0, termbox.ColorBlack, termbox.ColorWhite, client.ChatBoxList[i])
			} else {
				tbprint(j, 0, termbox.ColorWhite, termbox.ColorBlack, client.ChatBoxList[i])
			}

			j += len(client.ChatBoxList[i]) + 1
		}

		//draw the sidebar if it's active

		if client.SidebarActive {
			for i := 0; i < len(sideBarOptions); i++ {
				tbprint(w-30, i, termbox.ColorWhite, coldef, sideBarOptions[i])
			}

			//TODO make chat mode more flexible maybe map it or something
			modeText := "current mode: "
			tbprint(w-30, len(sideBarOptions), termbox.ColorWhite, coldef, modeText)
			if client.ChatMode == 0 {
				tbprint(w-30+len(modeText), len(sideBarOptions), termbox.ColorBlue, coldef, "normal")
			} else {
				tbprint(w-30+len(modeText), len(sideBarOptions), termbox.ColorRed, coldef, "fullwidth")
			}
			tbprint(w-30, len(sideBarOptions)+1, termbox.ColorWhite, coldef, "Total Messages received: "+strconv.Itoa(totalMessages))
		}

		//draw user input
		userString := client.Username + ": " + uInput
		tbprint(1, h-1, termbox.ColorBlue, coldef, userString)
		termbox.SetCursor(len(userString)+1, h-1)

		//draw @ search menu
		if strings.Contains(uInput, "@") && strings.Compare(atUsername(uInput), client.FoundUsername) != 0 {
			client.SearchMenuOpen = true

			possibleNames := Filter(client.ChatBoxes[client.CurrentChatBox].UserNames, func(v string) bool {
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
	*/

}

//get the user's twitch username through the twitch API
func getUsername(token, cId string) (username string) {
	rclient := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/kraken/user", nil)
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("Client-ID", cId)
	req.Header.Set("Authorization", "OAuth "+token)

	res, err := rclient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	name, err := jsonparser.GetString(body, "display_name")

	if err != nil {
		panic(err)
	}

	return name
}

func connect(cfg *irc.Config, channel string) *irc.Conn {
	c := irc.Client(cfg)

	c.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(channel)
	})

	c.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
		totalMessages++
		if !Contains(client.ChatBoxes[line.Target()].UserNames, line.Nick) {
			client.ChatBoxes[line.Target()].UserNames = append(client.ChatBoxes[line.Target()].UserNames, line.Nick)
		}

		w, _ := termbox.Size()
		splitMessage := make([]string, 0)
		newString := ""
		for i := 0; i < len(line.Args[1]); i++ {
			chr := string(line.Args[1][i])

			if len(newString)+1 <= int(float64(w)*0.8) || chr != " " {
				newString += chr
			} else {
				splitMessage = append(splitMessage, newString)
				newString = ""
			}
		}
		splitMessage = append(splitMessage, newString)

		for i := 0; i < len(splitMessage); i++ {
			newMessage(ChatLine{
				Nick:      line.Nick,
				NickColor: getNickColor(line.Nick),
				Line:      splitMessage[i],
			}, line.Target())
		}
	})

	if err := c.Connect(); err != nil {
		panic(err)
	}

	return c
}

func addChatChannel(conn *irc.Conn) {
	client.MenuOpen = true

	newChannel := ""

	newChannelMenu := newMenuBox("What channel do you want to talk in? (e.g. vinesauce) or press CRTL+B to go back", func() {
		client.Locked = true //although it shouldn't be an issue here, better safe than sorry
		newChannel = "#" + strings.ToLower(uInput)
		uInput = ""
		conn.Join(newChannel)
		client.MenuOpen = false
		client.CurrentChatBox = newChannel
		client.ChatBoxes[newChannel] = getChatBox(newChannel, chat)
		client.Initialized = true
		client.CurrentChatBoxInt = len(client.ChatBoxList) - 1
		client.Locked = false
		//redraw_all()
	})

	client.CurrentMenu = newChannelMenu
}

func leaveChatChannel(conn *irc.Conn) {
	client.Locked = true
	if len(client.ChatBoxList) > 1 {
		conn.Part(client.CurrentChatBox)
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

func switchChannel() {
	client.Locked = true
	if client.CurrentChatBoxInt == len(client.ChatBoxList)-1 {
		client.CurrentChatBoxInt = 0
	} else {
		client.CurrentChatBoxInt++
	}

	client.CurrentChatBox = client.ChatBoxList[client.CurrentChatBoxInt]
	client.Locked = false
	//redraw_all()
}

func updateMode() {
	if client.ChatMode > 0 {
		client.ChatMode = 0
	} else {
		client.ChatMode++
	}
}

func getChatBox(name, chatType string) *ChatBox {
	client.ChatBoxList = append(client.ChatBoxList, name)
	return NewChatBox(name, chatType)
}

func debugMessage(text string) {
	newMessage(ChatLine{
		Line:      text,
		Nick:      "",
		NickColor: termbox.ColorWhite,
	}, client.CurrentChatBox)
}

func newMessage(chatLine ChatLine, channel string) {
	_, h := termbox.Size()

	if len(client.ChatBoxes[channel].Lines) == h-2 {
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

func getNickColor(nick string) termbox.Attribute {
	if client.UserColors[nick] == termbox.ColorDefault {
		client.UserColors[nick] = getRandomColor()
	}

	return client.UserColors[nick]
}

func atUsername(parsed string) (username string) {
	username = ""
	parsing := false

	for i := 0; i < len(parsed); i++ {
		if parsed[i] == '@' && !parsing {
			parsing = true
		} else if parsed[i] != ' ' && parsing {
			username += string(parsed[i])
		} else if parsed[i] == ' ' && parsing {
			break
		}
	}

	return username
}

func remainingLetters() {
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

func getRandomColor() termbox.Attribute {
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
