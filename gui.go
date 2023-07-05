

package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
	"github.com/goki/freetype/truetype"
	"os"
)

func init()  {
	fp,err := findfont.Find("./frontLib/arial.ttf")
	if err != nil {
		panic(err)
	}
	fd,err := os.ReadFile(fp)
	if err != nil {
		panic(err)
	}
	_,err  = truetype.Parse(fd)
	if err != nil {
		panic(err)
	}
	os.Setenv("FYNE_FONT",fp)
}
func main() {
	a := app.New()

	w := a.NewWindow("Hello")
	w.Resize(fyne.NewSize(600, 400))

	ymEntry := widget.NewEntry()
	ymEntry.SetPlaceHolder("输入连接地址")
	ymEntry.OnChanged = func(content string) {
		fmt.Println("name:", ymEntry.Text, "entered")
	}
	w.SetContent(container.NewVBox(
		ymEntry,
		widget.NewButton("login", func() {}),
	))


	w.ShowAndRun()
}
