package rest

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
	xdraw "golang.org/x/image/draw"
)

type Calendar interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type calendar struct{}

func NewCalendar() Calendar {
	return &calendar{}
}

// FIXME: Move to external storage
//
//go:embed img/2023-08.png
var calender202308 []byte

//go:embed img/2023-09.png
var calender202309 []byte

//go:embed img/2023-10.png
var calender202310 []byte

//go:embed img/2023-11.png
var calender202311 []byte

//go:embed img/2023-12.png
var calender202312 []byte

var calendarMap = map[string][]byte{
	"2023-08": calender202308,
	"2023-09": calender202309,
	"2023-10": calender202310,
	"2023-11": calender202311,
	"2023-12": calender202312,
}

// Create : Create calendar images
// FIXME: Move logics from controller (手抜き実装 now)
func (c calendar) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	baseFile, baseFileHeader, err := r.FormFile("base_file")
	if err != nil {
		log.Error(ctx, fmt.Errorf("failed to http.Request.FormFile(): %w", err))
		DoResponse(ctx, w, err, http.StatusInternalServerError)
	}
	defer func() { _ = baseFile.Close() }()

	ext := strings.ToLower(filepath.Ext(baseFileHeader.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		DoResponse(ctx, w, "base file has invalid extension", http.StatusBadRequest)
		return
	}

	baseImg, _, err := image.Decode(baseFile)
	if err != nil {
		log.Error(ctx, fmt.Errorf("failed to image.Decode(): %w", err))
		DoResponse(ctx, w, "failed to decode base file", http.StatusInternalServerError)
		return
	}
	dateImg, err := png.Decode(bytes.NewReader(calendarMap[r.FormValue("target_date")]))
	if err != nil {
		log.Error(ctx, fmt.Errorf("failed to png.Decode(): %w", err))
		DoResponse(ctx, w, "failed to decode date file", http.StatusInternalServerError)
		return
	}

	baseImgRect := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: baseImg.Bounds().Size(),
	}

	baseImgBounds := baseImg.Bounds()
	dateImgBounds := dateImg.Bounds()

	dateImgRectMap := map[string]image.Rectangle{
		"upper_left":  getDateImgRectAtUpperLeft(baseImgBounds),
		"upper_right": getDateImgRectAtUpperRight(baseImgBounds),
		"lower_left":  getDateImgRectAtLowerLeft(baseImgBounds),
		"lower_right": getDateImgRectAtLowerRight(baseImgBounds),
	}

	rect, ok := dateImgRectMap[r.FormValue("date_position")]
	if !ok {
		log.Error(ctx, errors.New("invalid date position"))
		DoResponse(ctx, w, err.Error(), http.StatusBadRequest)
		return
	}

	out := image.NewRGBA(baseImgRect)
	draw.Draw(out, baseImgRect, baseImg, image.Point{X: 0, Y: 0}, draw.Src)
	xdraw.CatmullRom.Scale(out, rect, dateImg, dateImgBounds, draw.Over, nil)

	var output bytes.Buffer
	if err := png.Encode(&output, out); err != nil {
		log.Error(ctx, fmt.Errorf("failed to png.Encode(): %w", err))
		DoResponse(ctx, w, "failed to encode output file", http.StatusInternalServerError)
		return
	}

	DoImageResponse(ctx, w, output.Bytes(), "image/png", http.StatusCreated)
}

func getDateImgRectAtUpperLeft(baseImgBounds image.Rectangle) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{
			X: baseImgBounds.Dx() - baseImgBounds.Dx()/4,
			Y: baseImgBounds.Dy() / 4,
		},
	}
}

func getDateImgRectAtUpperRight(baseImgBounds image.Rectangle) image.Rectangle {
	startPoint := image.Point{
		X: baseImgBounds.Dx() - (baseImgBounds.Dx() - baseImgBounds.Dx()/4),
		Y: 0,
	}
	return image.Rectangle{
		Min: startPoint,
		Max: image.Point{
			X: startPoint.X + baseImgBounds.Dx() - baseImgBounds.Dx()/4,
			Y: startPoint.Y + baseImgBounds.Dy()/4,
		},
	}
}

func getDateImgRectAtLowerLeft(baseImgBounds image.Rectangle) image.Rectangle {
	startPoint := image.Point{
		X: 0,
		Y: baseImgBounds.Dy() - baseImgBounds.Dy()/4,
	}
	return image.Rectangle{
		Min: startPoint,
		Max: image.Point{
			X: startPoint.X + baseImgBounds.Dx() - baseImgBounds.Dx()/4,
			Y: startPoint.Y + baseImgBounds.Dy()/4,
		},
	}
}

func getDateImgRectAtLowerRight(baseImgBounds image.Rectangle) image.Rectangle {
	startPoint := image.Point{
		X: baseImgBounds.Dx() - (baseImgBounds.Dx() - baseImgBounds.Dx()/4),
		Y: baseImgBounds.Dy() - baseImgBounds.Dy()/4,
	}
	return image.Rectangle{
		Min: startPoint,
		Max: image.Point{
			X: startPoint.X + baseImgBounds.Dx() - baseImgBounds.Dx()/4,
			Y: startPoint.Y + baseImgBounds.Dy()/4,
		},
	}
}
