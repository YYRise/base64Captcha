package base64Captcha

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"

	"github.com/golang/freetype/truetype"
)

var (
	AllOperators = []string{"+", "-", "*"}
)

// DriverMath captcha config for captcha math
type DriverMathCustom struct {
	//Height png height in pixel.
	Height int

	// Width Captcha png width in pixel.
	Width int

	//NoiseCount text noise count.
	NoiseCount int

	//ShowLineOptions := OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine .
	ShowLineOptions int

	//BgColor captcha image background color (optional)
	BgColor *color.RGBA

	//fontsStorage font storage (optional)
	fontsStorage FontsStorage

	//Fonts loads by name see fonts.go's comment
	Fonts      []string
	fontsArray []*truetype.Font
	Operators  []string
	MaxNum     int32
}

// DriverMathCustom creates a driver of math
func NewDriverMathCustom(height int, width int, noiseCount int, showLineOptions int, bgColor *color.RGBA,
	fontsStorage FontsStorage, fonts []string, operators []string, maxNum int32) *DriverMathCustom {
	if fontsStorage == nil {
		fontsStorage = DefaultEmbeddedFonts
	}

	tfs := []*truetype.Font{}
	for _, fff := range fonts {
		tf := fontsStorage.LoadFontByName("fonts/" + fff)
		tfs = append(tfs, tf)
	}

	if len(tfs) == 0 {
		tfs = fontsAll
	}

	if len(operators) == 0 {
		operators = AllOperators
	}

	return &DriverMathCustom{
		Height:          height,
		Width:           width,
		NoiseCount:      noiseCount,
		ShowLineOptions: showLineOptions,
		fontsArray:      tfs,
		BgColor:         bgColor,
		Fonts:           fonts,
		Operators:       operators,
		MaxNum:          maxNum,
	}
}

// ConvertFonts loads fonts from names
func (d *DriverMathCustom) ConvertFonts() *DriverMathCustom {
	if d.fontsStorage == nil {
		d.fontsStorage = DefaultEmbeddedFonts
	}

	tfs := []*truetype.Font{}
	for _, fff := range d.Fonts {
		tf := d.fontsStorage.LoadFontByName("fonts/" + fff)
		tfs = append(tfs, tf)
	}
	if len(tfs) == 0 {
		tfs = fontsAll
	}
	d.fontsArray = tfs

	return d
}

// GenerateIdQuestionAnswer creates id,captcha content and answer
func (d *DriverMathCustom) GenerateIdQuestionAnswer() (id, question, answer string) {
	id = RandomId()
	var mathResult int32
	switch d.Operators[rand.Intn(len(d.Operators))] {
	case "+":
		a := rand.Int31n(d.MaxNum)
		b := rand.Int31n(d.MaxNum)
		question = fmt.Sprintf("%d+%d=?", a, b)
		mathResult = a + b
	case "-":
		a := rand.Int31n(d.MaxNum) + 1
		b := rand.Int31n(a)
		question = fmt.Sprintf("%d-%d=?", a, b)
		mathResult = a - b
	case "x":
		a := rand.Int31n(d.MaxNum)
		b := rand.Int31n(d.MaxNum)
		question = fmt.Sprintf("%dx%d=?", a, b)
		mathResult = a * b
	default:
		a := rand.Int31n(d.MaxNum) + rand.Int31n(d.MaxNum)
		b := rand.Int31n(d.MaxNum)

		question = fmt.Sprintf("%d-%d=?", a, b)
		mathResult = a - b

	}
	answer = fmt.Sprintf("%d", mathResult)
	return
}

// DrawCaptcha creates math captcha item
func (d *DriverMathCustom) DrawCaptcha(question string) (item Item, err error) {
	var bgc color.RGBA
	if d.BgColor != nil {
		bgc = *d.BgColor
	} else {
		bgc = RandLightColor()
	}
	itemChar := NewItemChar(d.Width, d.Height, bgc)

	//波浪线 比较丑
	if d.ShowLineOptions&OptionShowHollowLine == OptionShowHollowLine {
		itemChar.drawHollowLine()
	}

	//背景有文字干扰
	if d.NoiseCount > 0 {
		noise := RandText(d.NoiseCount, strings.Repeat(TxtNumbers, d.NoiseCount))
		err = itemChar.drawNoise(noise, fontsAll)
		if err != nil {
			return
		}
	}

	//画 细直线 (n 条)
	if d.ShowLineOptions&OptionShowSlimeLine == OptionShowSlimeLine {
		itemChar.drawSlimLine(3)
	}

	//画 多个小波浪线
	if d.ShowLineOptions&OptionShowSineLine == OptionShowSineLine {
		itemChar.drawSineLine()
	}

	//draw question
	err = itemChar.drawText(question, d.fontsArray)
	if err != nil {
		return
	}
	return itemChar, nil
}
