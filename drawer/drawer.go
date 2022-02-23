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

func pointsSize(phrase string, width int) float64 {
	runeCount := float64(utf8.RuneCount([]byte(phrase)))
	if runeCount < 15 {
		runeCount = 15
	}

	return 20 + float64(width)/runeCount
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
	_, h := txt.MeasureString(phrase)

	txt.SetColor(color.Black)
	txt.DrawStringAnchored(phrase, float64(txt.Width()/2)+5, float64(txt.Height())-h-20+5, 0.5, 0.5)

	txt.SetColor(color.White)
	txt.DrawStringAnchored(phrase, float64(txt.Width()/2), float64(txt.Height())-h-20, 0.5, 0.5)

	im.DrawImage(txt.Image(), 0, 0)

	buf := bytes.NewBuffer(nil)
	if err := im.EncodePNG(buf); err != nil {
		return nil, fmt.Errorf("error encoding image: %v", err)
	}

	return buf, nil
}
