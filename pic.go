package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
	"sync"

	"golang.org/x/mobile/event/size"

	"golang.org/x/mobile/event/paint"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
)

const chunk int = 500

type uploadEvent struct{}

type share struct {
	mu sync.Mutex
}

var (
	delta  = 3 //int
	buf    *image.NRGBA
	fin    *image.NRGBA
	tex    screen.Texture
	shared share
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func pixelAverage(buf *image.NRGBA, x, y int) (uint8, uint8, uint8, uint8) {
	var tr, tg, tb, ta uint16
	var count uint8
	for i := y - delta; i <= y+delta; i++ {
		for j := x - delta; j <= x+delta; j++ {
			pix := buf.NRGBAAt(j, i)
			tr = tr + uint16(pix.R)
			tg = tg + uint16(pix.G)
			tb = tb + uint16(pix.B)
			ta = ta + uint16(pix.A)
			count++
		}
	}

	// pix := buf.NRGBAAt(x, y)
	// fmt.Println(pix, tr, tg, tb, count)
	return uint8(tr / uint16(count)), uint8(tg / uint16(count)), uint8(tb / uint16(count)), uint8(ta / uint16(count))
}

func change(start, finish int, buf *image.NRGBA, fin *image.NRGBA, wg *sync.WaitGroup, pic image.Rectangle) {
	//modify image
	//height
	for y := start; y < finish && y < pic.Dy(); y++ {
		//width
		for x := 0; x < pic.Dx(); x++ {
			r, g, b, a := pixelAverage(buf, x, y)
			fin.SetNRGBA(x, y, color.NRGBA{
				R: r,
				G: g,
				B: b,
				A: a,
			})
		}
	}
	wg.Done()
}

// Used to copy an NRGBA image into the RGBA buffer
func copy2buf(b *image.RGBA, src *image.NRGBA, pic image.Rectangle) {
	for y := 0; y < pic.Dy(); y++ {
		for x := 0; x < pic.Dx(); x++ {
			pix := src.NRGBAAt(x, y)
			b.SetRGBA(x, y, color.RGBA{
				R: pix.R,
				G: pix.G,
				B: pix.B,
				A: pix.A,
			})
		}
	}
}

// Used to copy the finished pass back into the buffer
func copyRestart(buf *image.NRGBA, fin *image.NRGBA, pic image.Rectangle) {
	for y := 0; y < pic.Dy(); y++ {
		for x := 0; x < pic.Dx(); x++ {
			pix := fin.NRGBAAt(x, y)
			buf.SetNRGBA(x, y, pix)
			// fmt.Println(pix)
		}
	}
}

func blurLogic(pic image.Rectangle, q screen.EventDeque) {
	var wg sync.WaitGroup
	var chunk = 1000
	passes, err := strconv.Atoi(os.Args[2])
	check(err)
	pieces := fin.Bounds().Dy()/chunk + 1

	println("Main waiting.")
	for a := 0; a < passes; a++ {
		//spawning goroutines
		for i := 0; i < pieces; i++ {
			wg.Add(1)
			start := i * chunk
			finish := start + chunk
			go change(start, finish, buf, fin, &wg, pic)
		}
		wg.Wait()

		// shared.mu.Lock()
		fmt.Println("Finished pass: ", a+1)
		copyRestart(buf, fin, pic)
		q.Send(uploadEvent{})
		// shared.mu.Unlock()
	}
}

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("usage: ./pic <image> <passes>")
		log.Fatal()
	}

	//reading in cat pic
	fmt.Println("Reading in image.")
	reader, err := os.Open(args[1])
	check(err)
	old, _, err := image.Decode(reader)
	check(err)

	//drawing cat pic into new image
	fmt.Println("Making duplicate to manipulate.")
	buf = image.NewNRGBA(old.Bounds())
	fin = image.NewNRGBA(old.Bounds())
	draw.Draw(buf, buf.Bounds(), old, image.Point{0, 0}, draw.Src)

	//pic.Dy() = height, pic.Dx() = width
	picSize := fin.Bounds()

	var up, right int
	if picSize.Dy() < 1000 && picSize.Dx() < 1000 {
		up = picSize.Dy()
		right = picSize.Dx()
	} else if picSize.Dy() > picSize.Dx() {
		right = picSize.Dx() / picSize.Dy() * 1000
		up = 1000
	} else {
		up = 1000
		right = picSize.Dx() / picSize.Dy() * 1000
	}

	//Drawing
	driver.Main(func(s screen.Screen) {
		var sz size.Event
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Height: up,
			Width:  right,
			Title:  args[1],
		})
		check(err)

		b, err := s.NewBuffer(image.Point{picSize.Dx(), picSize.Dy()})
		check(err)
		tex, err := s.NewTexture(image.Point{picSize.Dx(), picSize.Dy()})
		check(err)

		go blurLogic(picSize, w)

		for {

			publish := false

			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}

			case key.Event:
				if e.Code == key.CodeEscape {
					fmt.Println("Outputting image.")
					f, err := os.Create("image.png")
					check(err)
					err = png.Encode(f, fin)
					check(err)
					fmt.Println("Done.")

					return
				}
			case size.Event:
				fmt.Println("hit size.Event")
				sz = e

			case paint.Event:
				publish = true

			case uploadEvent:
				publish = true
				copy2buf(b.RGBA(), fin, picSize)
				tex.Upload(image.Point{}, b, b.Bounds())

			case error:
				log.Print(e)
			}

			if publish {
				w.Scale(sz.Bounds(), tex, tex.Bounds(), screen.Src, nil)
				w.Publish()
			}
		}
	})

	fmt.Println("Done.")
}
