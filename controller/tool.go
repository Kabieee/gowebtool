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
	"time"

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

	// 设置图片透明
	maskImg := image.NewAlpha(wrawImg.Bounds())
	for x := 0; x < wrawImg.Bounds().Dx(); x++ {
		for y := 0; y < wrawImg.Bounds().Dy(); y++ {
			maskImg.Set(x, y, color.Alpha{A: 200})
			//r, _, _, _ := wrawImg.At(x, y).RGBA()
			//maskImg.Set(x, y, color.Alpha{A: uint8(256 - r)})
		}
	}

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
	fontsize := 68
	f := freetype.NewContext()
	f.SetDPI(float64(fontsize))
	f.SetFont(fonts)
	f.SetFontSize(float64(fontsize))
	f.SetClip(newImage.Bounds())
	f.SetDst(newImage)
	f.SetHinting(font.HintingFull)

	timeStr := time.Now().Format("Jan Mon 2006-01-02 15:04:05 MST")
	colRandom := colorful.FastHappyColor()
	col := CombineColor(colRandom.RGB255())
	s := image.NewUniform(col)
	f.SetSrc(s)

	fpt := freetype.Pt(newImage.Bounds().Dx()-((fontsize/2)+(fontsize%2))*len(timeStr), newImage.Bounds().Dy()-(fontsize/2))
	_, err = f.DrawString(timeStr, fpt)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := bytes.NewBuffer(nil)
	err = jpeg.Encode(buf, newImage, &jpeg.Options{Quality: 75})
	if err != nil {
		c.String(200, err.Error())
		return
	}
	mimeType := mime.TypeByExtension(".jpg")
	c.Header("Content-Length", strconv.Itoa(len(buf.Bytes())))
	c.Data(200, mimeType, buf.Bytes())
}

func CombineColor(cols ...uint8) color.Color {
	col := color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
	offset := 8

	if len(cols) >= 3 {
		col.R = cols[0]
		col.G = cols[1]
		col.B = cols[2]
		offset = 9
	}
	// 合并颜色, 透明色
	r, g, b, a := col.RGBA()
	r2, g2, b2, a2 := color.Alpha{A: 128}.RGBA()
	col.R = uint8((r + r2) >> offset) // div by 2 followed by">> 8"  is">> 9"
	col.G = uint8((g + g2) >> offset)
	col.B = uint8((b + b2) >> offset)
	col.A = uint8((a + a2) >> offset)
	fmt.Println(col.A, col.G, col.B, col.A)
	return col
}
