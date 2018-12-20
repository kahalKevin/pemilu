package liliput

import (
	"fmt"
	"strings"
	"time"

	"github.com/discordapp/lilliput"
)

var EncodeOptions = map[string]map[int]int{
	".jpeg": map[int]int{lilliput.JpegQuality: 85},
	".png":  map[int]int{lilliput.PngCompression: 7},
	".webp": map[int]int{lilliput.WebpQuality: 85},
}

func ReduceSize(inputBuf []byte) (reducedImg []byte, err error) {
	var outputWidth int
	var outputHeight int
	stretch := true

	decoder, err := lilliput.NewDecoder(inputBuf)
	if err != nil {
		return
	}
	defer decoder.Close()

	header, err := decoder.Header()
	if err != nil {
		return
	}

	// print some basic info about the image
	fmt.Printf("file type: %s\n", decoder.Description())
	fmt.Printf("%dpx x %dpx\n", header.Width(), header.Height())

	if decoder.Duration() != 0 {
		fmt.Printf("duration: %.2f s\n", float64(decoder.Duration())/float64(time.Second))
	}

	// get ready to resize image,
	// using 8192x8192 maximum resize buffer size
	ops := lilliput.NewImageOps(8192)
	defer ops.Close()

	// create a buffer to store the output image, 8MB in this case
	outputImg := make([]byte, 8*1024*1024)
	outputType := "." + strings.ToLower(decoder.Description())
	if header.Width() > header.Height() {
		fmt.Println("wider")
		outputWidth = 800
		outputHeight = 600
	} else if header.Width() < header.Height() {
		fmt.Println("taller")
		outputWidth = 600
		outputHeight = 800
	} else {
		fmt.Println("square")
		outputWidth = 700
		outputHeight = 700
	}

	resizeMethod := lilliput.ImageOpsFit
	if stretch {
		resizeMethod = lilliput.ImageOpsResize
	}

	opts := &lilliput.ImageOptions{
		FileType:             outputType,
		Width:                outputWidth,
		Height:               outputHeight,
		ResizeMethod:         resizeMethod,
		NormalizeOrientation: true,
		EncodeOptions:        EncodeOptions[outputType],
	}

	// resize and transcode image
	reducedImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		return
	}
	return
}
