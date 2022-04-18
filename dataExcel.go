package main

import (
	_ "database/sql"
	"embed"
	"fmt"
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xuri/excelize/v2"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"image/png"
	"log"
)

//go:embed graphics/*
var EmbeddedAssets embed.FS

var counter = 0
var stateApp GuiApp
var textWidget *widget.Text

func main() {
	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("State Population Change 2020-2021")

	stateApp = GuiApp{AppUI: MakeUIWindow()}

	excelDataSheet := openExcelGetStates()
	defer excelDataSheet.Close()
	loadStates()
	defer excelDataSheet.Close()

	err := ebiten.RunGame(&stateApp)
	if err != nil {
		log.Fatalln("error with gui/window", err)
	}
}

func openExcelGetStates() *excelize.File {
	excelFile, err := excelize.OpenFile("countyPopChange2020-2021.xlsx")
	if err != nil {
		log.Fatalln(err)
	}
	all_rows, err := excelFile.GetRows("co-est2021-alldata")
	if err != nil {
		log.Fatalln(err)
	}
	// https://github.com/qax-os/excelize/blob/master/col.go
	for number, row := range all_rows {
		if number < 0 {
			continue
		}
		if len(row) <= 1 {
			continue
		}
		if row[5] == row[6] {
			continue
		}
		fmt.Println(row[5], "\t: ", row[10], "\t:", row[11])
	}
	return excelFile
}

//GUI stuff below

func (g GuiApp) Update() error {
	//TODO finish me
	g.AppUI.Update()
	return nil
}

func (g GuiApp) Draw(screen *ebiten.Image) {
	//TODO finish me
	g.AppUI.Draw(screen)
}

func (g GuiApp) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

type GuiApp struct {
	AppUI *ebitenui.UI
}

func loadImageNineSlice(path string, centerWidth int, centerHeight int) (*image.NineSlice, error) {
	i := loadPNGImageFromEmbedded(path)

	w, h := i.Size()
	return image.NewNineSlice(i,
			[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
			[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight}),
		nil
}

func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := EmbeddedAssets.ReadDir("graphics")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := EmbeddedAssets.Open("graphics/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}

func MakeUIWindow() (GUIhandler *ebitenui.UI) {
	background := image.NewNineSliceColor(color.Gray16{})
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(background))
	textInfo := widget.TextOptions{}.Text("This is our first Window", basicfont.Face7x13, color.White)

	idle, err := loadImageNineSlice("button-idle.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	hover, err := loadImageNineSlice("button-hover.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	pressed, err := loadImageNineSlice("button-pressed.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	disabled, err := loadImageNineSlice("button-disabled.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	buttonImage := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}
	button := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Press Me", basicfont.Face7x13, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),
		// ... click handler, etc. ...
		widget.ButtonOpts.ClickedHandler(percentageStateChangeButton),
	)
	rootContainer.AddChild(button)
	resources, err := newListResources()
	if err != nil {
		log.Println(err)
	}
	allStates := loadStates()
	dataofPoP := make([]interface{}, len(allStates))
	for position, state := range allStates {
		dataofPoP[position] = state
	}

	//add excel data here, and button click arguement

	listWidget := widget.NewList(
		widget.ListOpts.Entries(dataofPoP),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			fullStateName := "%s %s %s"
			fmt.Sprintf(fullStateName, e.(States).StateName, e.(States).Pop20, e.(States).Pop21)
			return e.(States).StateName
		}),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(resources.image)),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(resources.track, resources.handle),
			widget.SliderOpts.HandleSize(resources.handleSize),
			widget.SliderOpts.TrackPadding(resources.trackPadding)),
		widget.ListOpts.EntryColor(resources.entry),
		widget.ListOpts.EntryFontFace(resources.face),
		widget.ListOpts.EntryTextPadding(resources.entryPadding),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			//do something when a list item changes
			//args.Entry
		}))
	rootContainer.AddChild(listWidget)
	textWidget = widget.NewText(textInfo)
	rootContainer.AddChild(textWidget)

	GUIhandler = &ebitenui.UI{Container: rootContainer}
	return GUIhandler
}

func percentageStateChang(old, new int) (delta float64) {
	diff := float64(new - old)
	delta = (diff / float64(old)) * 100
	return
}

//When the user selects a list item, display the percentage change in population for the state in a text label
func percentageStateChangeButton(args *widget.ButtonClickedEventArgs) {
	excelFile, err := excelize.OpenFile("countyPopChange2020-2021.xlsx")
	if err != nil {
		log.Fatalln(err)
	}
	statesPopRows, err := excelFile.GetRows("co-est2021-alldata")
	if err != nil {
		log.Fatalln("statePopRows error", err)
	}
	for number, row := range statesPopRows {
		if number < 0 {
			continue
		}
		if len(row) <= 1 {
			continue
		}
		if row[5] == row[6] {
			continue
		}
		percentageStatePop, _ := (row[10] / row[11]) * 100 //i knwo it has to be similar to this but could not figure out how to implement it with adding the strucs from the row
		message := fmt.Sprintf("You have pressed this state button %d times", percentageStateChang(percentageStatePop))
		textWidget.Label = message
	}
}

//loads states from excel sheet into GUI
func loadStates() []States {
	excelFile, err := excelize.OpenFile("countyPopChange2020-2021.xlsx")
	if err != nil {
		log.Fatalln(err)
	}
	sliceofStates := make([]States, 50, 51)
	statesPopRows, err := excelFile.GetRows("co-est2021-alldata")
	if err != nil {
		log.Fatalln("statePopRows error", err)
	}
	for number, row := range statesPopRows {
		if number < 0 {
			continue
		}
		if len(row) <= 1 {
			continue
		}
		if row[5] == row[6] {
			continue
		}
		fmt.Scanf(row[5], "\t: ", row[10], "\t:", row[11])

		location := 0
		currentState := States{}
		currentState = States{
			row[5],
			row[10],
			row[11],
			row[12],
			row[12],
		}
		sliceofStates[location] = currentState
		location++
	}
	return sliceofStates
}
