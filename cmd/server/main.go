package main

import (
	"bytes"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"strconv"

	wm "github.com/Bornholm/fantasy-world/worldmap"
	"github.com/gorilla/mux"
)

var worldmap *wm.WorldMap

var (
	seed       = "fantasy-world"
	octaves    = 8
	lacunarity = 1.9
	gain       = 0.65
)

func init() {
	flag.StringVar(&seed, "seed", seed, "the world seed")
	flag.IntVar(&octaves, "octaves", octaves, "brownian motion octaves")
	flag.Float64Var(&lacunarity, "lacunarity", lacunarity, "brownian motion lacunarity")
	flag.Float64Var(&gain, "gain", gain, "brownian motion gain")
}

func main() {

	flag.Parse()

	worldmap = wm.NewWorldMap(seed, octaves, lacunarity, gain)

	router := mux.NewRouter()

	router.HandleFunc("/tiles/{x}/{y}", tileHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Fatal(http.ListenAndServe(":3333", router))

}

func tileHandler(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	tileX, err := strconv.Atoi(vars["x"])
	if err != nil {
		panic(err)
	}

	tileY, err := strconv.Atoi(vars["y"])
	if err != nil {
		panic(err)
	}

	// log.Println("Tile: ", tileX, tileY)

	pixelSize := 1
	tileSize := 256

	m := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))

	var region *wm.Region

	for localTileY := 0; localTileY < tileSize; localTileY += pixelSize {
		for localTileX := 0; localTileX < tileSize; localTileX += pixelSize {

			globalX := (tileX * tileSize) + localTileX
			globalY := (tileY * tileSize) + localTileY

			region = worldmap.GetRegion(globalX, globalY)

			bgColor := color.RGBA{0, 0, 0, 255}
			switch region.LandscapeType {
			case "deepwater":
				bgColor = color.RGBA{8, 104, 172, 255}
			case "water":
				bgColor = color.RGBA{67, 162, 202, 255}
			case "plain":
				bgColor = color.RGBA{186, 228, 188, 255}
			case "snow":
				bgColor = color.RGBA{240, 249, 232, 255}
			}

			rect := image.Rect(localTileX, localTileY, localTileX+pixelSize, localTileY+pixelSize)
			draw.Draw(m, rect, &image.Uniform{bgColor}, image.ZP, draw.Src)

		}
	}

	var img image.Image = m

	writeImage(res, &img)

}

func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, *img); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}
