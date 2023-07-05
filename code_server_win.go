package main

import (
	"cmd_chat/server"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	server.NewServer( serverGui())
}
func serverGui() (string,string)  {
	var ym,port string
	myApp := app.New()
	myWindow := myApp.NewWindow("cmd chat setting")

	ymEntry := widget.NewEntry()
	portEntry := widget.NewEntry()

	form := widget.NewForm(
		&widget.FormItem{Text: "ip", Widget: ymEntry},
		&widget.FormItem{Text: "port", Widget: portEntry},
	)

	form.OnSubmit = func() {
		ym = ymEntry.Text
		port= portEntry.Text
		myWindow.Close()
	}
	form.OnCancel = func() { }
	myWindow.SetContent(form)
	myWindow.Resize(fyne.NewSize(350, 200))
	myWindow.ShowAndRun()
	fmt.Println("ym:",ym, "port:", port,"name:")
	return ym,port
}