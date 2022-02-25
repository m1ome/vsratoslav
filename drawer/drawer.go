package drawer

import (
	"bytes"
	"fmt"
	"gopkg.in/fogleman/gg.v1"
	"image"
	"image/color"
	"io"
	"strings"
	"unicode/utf8"
)

const lineSpacing = 1.3
const basePointSize = 30

func pointsSize(phrase string, width int, height int) float64 {
	runeCount := float64(utf8.RuneCount([]byte(phrase)))
	biggerSide := width
	if height > width {
		biggerSide = height
	}
	basePointSize := float64(basePoints(biggerSide))

	if runeCount > 20 {
		return basePointSize * 0.8
	}

	return basePointSize
}

func basePoints(width int) int {
	return basePointSize + (width/640)*32
}

func DrawText(reader io.Reader, font, phrase string) (io.Reader, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	im := gg.NewContextForImage(img)

	txt := gg.NewContext(im.Width(), im.Height())
	txt.Clear()
	txt.SetColor(color.White)

	pointSize := pointsSize(phrase, txt.Width(), txt.Height())
	if err := txt.LoadFontFace(font, pointSize); err != nil {
		return nil, fmt.Errorf("error loading font: %v", err)
	}
	lines := txt.WordWrap(phrase, float64(txt.Width())-20)
	_, textHeight := txt.MeasureMultilineString(strings.Join(lines, "\n"), lineSpacing)
	offset := textHeight + 10
	if len(lines) > 0 {
		offset = offset - pointSize*float64(len(lines)-1)*0.5
	}

	x := float64(txt.Width())
	y := float64(txt.Height()) - offset

	txt.SetColor(color.Black)
	txt.DrawStringWrapped(phrase, x/2+2, y+2, 0.5, 0.5, x-20, lineSpacing, gg.AlignCenter)

	txt.SetColor(color.White)
	txt.DrawStringWrapped(phrase, x/2, y, 0.5, 0.5, x-20, lineSpacing, gg.AlignCenter)

	im.DrawImage(txt.Image(), 0, 0)

	buf := bytes.NewBuffer(nil)
	if err := im.EncodePNG(buf); err != nil {
		return nil, fmt.Errorf("error encoding image: %v", err)
	}

	return buf, nil
}
