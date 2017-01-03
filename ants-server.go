package main

import (
  "net/http"
  "log"
  "github.com/gorilla/mux"
  "strconv"
  "image"
  "image/png"
  "image/color"
  "image/draw"
  "bytes"
  wm "./worldmap"
)

var worldmap *wm.WorldMap

func main() {

  worldmap = wm.NewWorldMap("hello world testttttt", 8, 1.9, 0.65)

  router := mux.NewRouter()

  router.HandleFunc("/tiles/{x}/{y}", tileHandler)
  router.PathPrefix("/vendor/").Handler(http.StripPrefix("/vendor/", http.FileServer(http.Dir("node_modules"))))
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

  log.Println("Tile: ", tileX, tileY)

  pixelSize := 2
  tileSize := 256

  m := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))

  var region *wm.Region

  for localTileY := 0; localTileY < tileSize; localTileY += pixelSize {
    for localTileX := 0; localTileX < tileSize; localTileX += pixelSize {

      globalX := (tileX*tileSize)+localTileX
      globalY := (tileY*tileSize)+localTileY

      region = worldmap.GetRegion(globalX, globalY)

      bgColor := color.RGBA{0, 0, 0, 255}
      switch(region.LandscapeType) {
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
