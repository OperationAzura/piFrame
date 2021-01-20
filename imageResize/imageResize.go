package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
)

var err error

//filter will hold the image resizing data, as well as the methods to use it.
var filter = gift.New()

// ImageCTX holds the resize data and the image itself
type ImageCTX struct {
	Name      string
	SrcImg    image.Image
	Img       image.Image
	FileName  string
	Path      string
	Subfolder string
	ImgBuff   *bytes.Buffer
}

func main() {
	fmt.Println("GO imageResize start")

	//var path = `./`
	//var subfolder = `./`
	//var name = `./`

	// Directory we want to get all files from.
	directory := "/home/pi/share/japan/"

	// Open the directory.
	outputDirRead, _ := os.Open(directory)

	// Call Readdir to get all files.
	outputDirFiles, _ := outputDirRead.Readdir(0)

	// Loop over files.
	for _, oFile := range outputDirFiles {
		fmt.Println("attempting file: ",oFile.Name())
		if !oFile.IsDir() {
			// Get name of file.
			oFile.Name()
			src, err := imaging.Open(directory + oFile.Name()) //open the image as an image
			if err != nil {
				fmt.Println("error opening image", err)
			}

			iCTX := BuildImageCTX(src)

			imgJPG := BuildImages(iCTX)

			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, imgJPG, nil)
			if err != nil{
				fmt.Println("error encoding jpeg: ",err)
			}
			err = ioutil.WriteFile(directory+"new/"+oFile.Name(), buf.Bytes(), 0644)
			if err != nil {
				fmt.Println("Error writting jpeg: ", err)
			}
		}
	}
	fmt.Println("GO imageResize done")
}

//BuildImageCTX adds information needed to resize the image
func BuildImageCTX(srcImg image.Image) ImageCTX {
	i := ImageCTX{}
	i.SrcImg = srcImg //source image
	//ImgBuff used in sendToS3()
	i.ImgBuff = new(bytes.Buffer)

	return i
}

//ReadyFlags prepares and parses the flags.
func ReadyFlags() {
	//flag.StringVar(&srcBucketFlag, `b`, ``, `Sets the source s3 bucket`)
	//flag.StringVar(&srcNameFlag, `n`, ``, `Set the source file name`)
	//flag.Parse()
}

//GetImage downloads an image from a url
func GetImage(url, fileName string) image.Image {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	srcImg, err := imaging.Decode(response.Body)
	if err != nil {
		fmt.Println("error decoding response image: ", err)
	}
	return srcImg
}

//SetResize Sets the resize filter
func SetResize(width, height int) {
	filter.Add(gift.Resize(width, height, gift.LanczosResampling))
}

//SetUnsharpMask add the unsharp mask filter
func SetUnsharpMask(sigma, amount, threshold float32) {
	filter.Add(gift.UnsharpMask(sigma, amount, threshold))
}

//ApplyFilters applys the set filters
func ApplyFilters(srcImg image.Image) image.Image {
	img := image.NewRGBA(filter.Bounds(srcImg.Bounds()))
	filter.Draw(img, srcImg)
	return img
}

//BuildImages takes the source image and creates images of the desired sizes.
// Returns the filename, sku, and position in the S3 bucket
func BuildImages(i ImageCTX) image.Image {

	width := 600

	// Reset variables that are used multiple times
	i.ImgBuff.Reset()
	filter = gift.New()

	// Resize, apply unsharp mask, and encode to jpg
	SetResize(width, 0)
	SetUnsharpMask(0.1, 1, 0)
	i.Img = ApplyFilters(i.SrcImg)

	// By default, use GoLang's built-in jpeg encoding
	fmt.Println("Using jpeg.")
	err = jpeg.Encode(i.ImgBuff, i.Img, &jpeg.Options{Quality: 95})
	if err != nil {
		fmt.Println("failed to create buffer", err)
	}
	return i.Img
}

