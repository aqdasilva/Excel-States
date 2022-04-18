package main

import (
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"image/color"
)

//straight from the ebitenUI demo https://github.com/blizzy78/ebitenui/tree/master/_demo
type listResources struct {
	image        *widget.ScrollContainerImage
	track        *widget.SliderTrackImage
	trackPadding widget.Insets
	handle       *widget.ButtonImage
	handleSize   int
	face         font.Face
	entry        *widget.ListEntryColor
	entryPadding widget.Insets
}

//Based on the Ebiten demo edited to add some defaults and use go:embed
func newListResources() (*listResources, error) {
	idle := loadPNGImageFromEmbedded("list-idle.png")

	disabled := loadPNGImageFromEmbedded("list-disabled.png")

	mask := loadPNGImageFromEmbedded("list-mask.png")

	trackIdle := loadPNGImageFromEmbedded("list-track-idle.png")

	trackDisabled := loadPNGImageFromEmbedded("list-track-disabled.png")

	handleIdle := loadPNGImageFromEmbedded("slider-handle-idle.png")

	handleHover := loadPNGImageFromEmbedded("slider-handle-hover.png")

	return &listResources{
		image: &widget.ScrollContainerImage{
			Idle:     image.NewNineSlice(idle, [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(disabled, [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Mask:     image.NewNineSlice(mask, [3]int{26, 10, 23}, [3]int{26, 10, 26}),
		},

		track: &widget.SliderTrackImage{
			Idle:     image.NewNineSlice(trackIdle, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Hover:    image.NewNineSlice(trackIdle, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(trackDisabled, [3]int{0, 5, 0}, [3]int{25, 12, 25}),
		},

		trackPadding: widget.Insets{
			Top:    5,
			Bottom: 24,
		},

		handle: &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(handleIdle, 0, 5),
			Hover:    image.NewNineSliceSimple(handleHover, 0, 5),
			Pressed:  image.NewNineSliceSimple(handleHover, 0, 5),
			Disabled: image.NewNineSliceSimple(handleIdle, 0, 5),
		},

		handleSize: 5,
		face:       basicfont.Face7x13,

		entry: &widget.ListEntryColor{
			Unselected:         color.RGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
			DisabledUnselected: color.RGBA{R: 0x5a, G: 0x7a, B: 0x91, A: 0x1f},

			Selected:         color.RGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
			DisabledSelected: color.RGBA{R: 0x5a, G: 0x7a, B: 0x91, A: 0x1f},

			SelectedBackground:         color.RGBA{R: 0x4b, G: 0x68, B: 0x7a, A: 0xff},
			DisabledSelectedBackground: color.RGBA{R: 0x2a, G: 0x39, B: 0x44, A: 0xff},
		},

		entryPadding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    2,
			Bottom: 2,
		},
	}, nil
}
