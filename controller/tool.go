package controller

import (
	"bytes"
	"fmt"
	"gowebtool/common"
	"gowebtool/services"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"mime"
	"os"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"

	"golang.org/x/image/font"

	"github.com/golang/freetype"

	"github.com/go-playground/validator/v10"

	//_ "image/jpeg"
	_ "image/png"

	"github.com/gin-gonic/gin"
)

type ToolController struct {
	BaseController
}

func (t *ToolController) SendEmail(c *gin.Context) {
	e := &services.EmailServices{}
	e.Ctx = c
	var data services.EmailData
	err := c.ShouldBind(&data)
	if err != nil {
		for _, filedError := range err.(validator.ValidationErrors) {
			t.Fail(c, &Fail{ErrorInfo: filedError.Translate(common.MyTran)})
			return
		}
	}
	data.Body = strings.Trim(data.Body, " ")
	if data.Body == "" {
		t.Fail(c, &Fail{ErrorInfo: "body不能为空"})
		return
	}
	err, result := e.Send(data)
	if err != nil {
		t.Fail(c, &Fail{ErrorInfo: err.Error()})
		return
	}
	t.Success(c, &Success{Data: result})
}

func (t *ToolController) WaterMark(c *gin.Context) {
	img, _ := os.Open("static/a.jpg")
	defer img.Close()
	rawImg, _, err := image.Decode(img)
	if err != nil {
		fmt.Println(err)
		return
	}

	wimg, _ := os.Open("static/w.png")
	defer wimg.Close()
	wrawImg, _, err := image.Decode(wimg)
	if err != nil {
		fmt.Println(err)
		return
	}

	maskImg := image.NewAlpha(wrawImg.Bounds())
	for x := 0; x < wrawImg.Bounds().Dx(); x++ {
		for y := 0; y < wrawImg.Bounds().Dy(); y++ {
			maskImg.Set(x, y, color.Alpha{A: 200})
			//r, _, _, _ := wrawImg.At(x, y).RGBA()
			//maskImg.Set(x, y, color.Alpha{A: uint8(256 - r)})
		}
	}

	//fmt.Println(wrawImg.Bounds().Dx(), wrawImg.Bounds().Dy())
	//fmt.Println(wimageType)
	//
	//fmt.Println(rawImg.ColorModel())
	//fmt.Println(rawImg.Bounds().Dx(), rawImg.Bounds().Dy())

	offset := image.Pt(rawImg.Bounds().Dx()-wrawImg.Bounds().Dx(), rawImg.Bounds().Dy()-wrawImg.Bounds().Dy())
	newImage := image.NewRGBA(rawImg.Bounds())

	draw.Draw(newImage, rawImg.Bounds(), rawImg, image.Point{}, draw.Src)
	draw.DrawMask(newImage, wrawImg.Bounds().Add(offset), wrawImg, image.Point{}, maskImg, image.Point{}, draw.Over)
	//draw.Draw(newImage, wrawImg.Bounds().Add(offset), maskImg, image.Point{}, draw.Over)

	// 设置字体
	fontFile, err := ioutil.ReadFile("static/font1716.ttf")
	if err != nil {
		fmt.Println(err)
		return
	}
	fonts, err := freetype.ParseFont(fontFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	f := freetype.NewContext()
	f.SetDPI(72)
	f.SetFont(fonts)
	f.SetFontSize(65)
	f.SetClip(newImage.Bounds())
	f.SetDst(newImage)
	f.SetHinting(font.HintingFull)

	col := colorful.FastHappyColor()
	fmt.Println(col.Hex())
	s := image.NewUniform(col)
	f.SetSrc(s)
	fpt := freetype.Pt(newImage.Bounds().Dx()-10*70, newImage.Bounds().Dy()-35)
	_, err = f.DrawString("TEST WATER MARK", fpt)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := bytes.NewBuffer(nil)
	err = jpeg.Encode(buf, newImage, &jpeg.Options{Quality: 100})
	if err != nil {
		c.String(200, err.Error())
		return
	}
	mimeType := mime.TypeByExtension(".jpg")
	c.Header("Content-Length", strconv.Itoa(len(buf.Bytes())))
	c.Data(200, mimeType, buf.Bytes())
}

func CombineColor() color.Color {
	// 合并颜色, 透明色
	col := color.RGBA{
		R: 7,
		G: 40,
		B: 187,
		A: 255,
	}
	r, g, b, a := col.RGBA()
	r2, g2, b2, a2 := color.Alpha{A: 128}.RGBA()
	col.R = uint8((r + r2) >> 9) // div by 2 followed by">> 8"  is">> 9"
	col.G = uint8((g + g2) >> 9)
	col.B = uint8((b + b2) >> 9)
	col.A = uint8((a + a2) >> 9)
	//fmt.Println(col.A, col.G, col.B, col.A)
	return col
}
