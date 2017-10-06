package main

type MenuBox struct {
	title string
	callback func()
}

func newMenuBox(title string, callback func())MenuBox{
	return MenuBox{
		title: title,
		callback: callback,
	}
}