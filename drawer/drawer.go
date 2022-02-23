package drawer

import (
	"bytes"
	"fmt"
	"gopkg.in/fogleman/gg.v1"
	"image"
	"image/color"
	"io"
	"unicode/utf8"
)

const lineSpacing = 1.3

func pointsSize(phrase string, width int) float64 {
	runeCount := float64(utf8.RuneCount([]byte(phrase)))

	if runeCount < 15 {
		return 45
	} else if runeCount < 10 {
		return 60
	} else {
		return 28
	}
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

	pointSize := pointsSize(phrase, txt.Width())
	fmt.Printf("string: %s / point size: %f\n", phrase, pointSize)
	if err := txt.LoadFontFace(font, pointSize); err != nil {
		return nil, fmt.Errorf("error loading font: %v", err)
	}
	_, h := txt.MeasureMultilineString(phrase, lineSpacing)

	width := float64(txt.Width())
	height := float64(txt.Height()) - h

	offset := 5.0
	if float64(utf8.RuneCount([]byte(phrase))) > 20 {
		offset = 25.0
	}

	txt.SetColor(color.Black)
	txt.DrawStringWrapped(phrase, width/2+3, height-offset+3, 0.5, 0.5, width-20, lineSpacing, gg.AlignCenter)

	txt.SetColor(color.White)
	txt.DrawStringWrapped(phrase, width/2, height-offset, 0.5, 0.5, width-20, lineSpacing, gg.AlignCenter)

	im.DrawImage(txt.Image(), 0, 0)

	buf := bytes.NewBuffer(nil)
	if err := im.EncodePNG(buf); err != nil {
		return nil, fmt.Errorf("error encoding image: %v", err)
	}

	return buf, nil
}
