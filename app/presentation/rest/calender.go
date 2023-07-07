package rest

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
	xdraw "golang.org/x/image/draw"
)

type Calendar interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type calendar struct{}

func NewCalendar() Calendar {
	return &calendar{}
}

// Create : Create calendar images
// FIXME: Move logics from controller (手抜き実装 now)
func (c calendar) Create(w http.ResponseWriter, r *http.Request) {
	baseFile, baseFileHeader, err := r.FormFile("base_file")
	if err != nil {
		app.Error(fmt.Errorf("error in http.Request.FormFile(): %w", err))
		DoResponse(w, err, http.StatusInternalServerError)
	}
	defer func() { _ = baseFile.Close() }()

	filename := strings.Split(baseFileHeader.Filename, ".")[0]
	ext := strings.ToLower(filepath.Ext(baseFileHeader.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		DoResponse(w, "base file has invalid extension", http.StatusBadRequest)
		return
	}

	targetDate := r.FormValue("target_date")
	dateFileName := fmt.Sprintf("./img/%s.png", targetDate)
	dateFile, err := os.Open(dateFileName)
	if err != nil {
		errMsg := fmt.Sprintf("date file is gone (%s)", dateFileName)
		DoResponse(w, errMsg, http.StatusInternalServerError)
		return
	}
	defer func() { _ = dateFile.Close() }()

	baseImg, _, err := image.Decode(baseFile)
	if err != nil {
		DoResponse(w, "decoding base file is failed", http.StatusInternalServerError)
		return
	}
	dateImg, err := png.Decode(dateFile)
	if err != nil {
		DoResponse(w, "decoding date file is failed", http.StatusInternalServerError)
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

	var wg sync.WaitGroup
	for key, rect := range dateImgRectMap {
		rect := rect
		key := key

		wg.Add(1)
		go func() {
			out := image.NewRGBA(baseImgRect)
			draw.Draw(out, baseImgRect, baseImg, image.Point{X: 0, Y: 0}, draw.Src)
			xdraw.CatmullRom.Scale(out, rect, dateImg, dateImgBounds, draw.Over, nil)

			outputFile, err := os.Create(fmt.Sprintf("./img/%s_%s.jpg", filename, key))
			if err != nil {
				DoResponse(w, "creating output file is failed", http.StatusInternalServerError)
				return
			}

			var opt jpeg.Options
			opt.Quality = 100

			if err := jpeg.Encode(outputFile, out, &opt); err != nil {
				DoResponse(w, "encoding output file is failed", http.StatusInternalServerError)
				return
			}

			wg.Done()
		}()
	}
	wg.Wait()

	DoResponse(w, nil, http.StatusCreated)
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
