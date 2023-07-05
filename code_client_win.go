package main

import (
	User "cmd_chat/client"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)


func main() {
	User.NewUserClient(clientGui())
}
func clientGui() (string,string,string)  {
	var ym,port,name string
	myApp := app.New()
	myWindow := myApp.NewWindow("cmd chat login")
	defer myWindow.Close()
	ymEntry := widget.NewEntry()
	portEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()

	form := widget.NewForm(
		&widget.FormItem{Text: "ip", Widget: ymEntry},
		&widget.FormItem{Text: "port", Widget: portEntry},
		&widget.FormItem{Text: "name", Widget: nameEntry},
	)

	form.OnSubmit = func() {
		ym = ymEntry.Text
		port= portEntry.Text
		name= nameEntry.Text

	}
	form.OnCancel = func() { }
	myWindow.SetContent(form)
	myWindow.Resize(fyne.NewSize(350, 200))
	myWindow.ShowAndRun()
	fmt.Println("ym:",ym, "port:", port,"name:",name)
	return ym,port,name
}