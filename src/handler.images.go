package main

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/fs"
	"net/http"

	"github.com/kennygrant/sanitize"
	"github.com/nfnt/resize"

	//_ "golang.org/x/image/webp"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

//const defaultTo = "resources/DecaFans-big.png"

// todo delete this folder from time to time
const imageDir = "cacheImages"
const originalImage = "original"

// todo, maybe preserve aspect ration

type CacheImageFs string

func (c CacheImageFs) Open(name string) (fs.File, error) {
	path := filepath.Join(string(c), name)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

var cacheFS CacheImageFs

// put 0 in width or height for default / auto-scale
// url can be a path
func getImage(url string, width, height int) (string, error) {
	name := getImageName(url)

	path := filepath.Join(imageDir, name)

	sizeRequested := fmt.Sprintf("%dx%d", width, height)
	if sizeRequested == "0x0" {
		sizeRequested = originalImage
	}

	var img image.Image
	var format string

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", err
		}
		// get new image
		img, format, err = fetchImage(url)
		if err != nil {
			if err == image.ErrFormat {
				log.Println("bad format", url)
			}
			errRem := os.Remove(path)
			if errRem != nil {
				return "", errRem
			}
			return "", err
		}
		imgPath := fmt.Sprintf("%s.%s", filepath.Join(path, originalImage), format)
		err = writeImage(img, imgPath, format)
		if err != nil {
			return "", err
		}

		// dont try to resize gifs
		if format == "gif" {
			sizeRequested = originalImage
		}

		if sizeRequested == originalImage {
			return imgPath, nil
		}
	} else {
		images, err := ioutil.ReadDir(path)
		if err != nil {
			return "", err
		}
		// check to see if size is already available
		for _, f := range images {
			if strings.Contains(f.Name(), sizeRequested) {
				return filepath.Join(path, f.Name()), nil
			}
		}

		// get original to rescale
		for _, f := range images {
			if strings.Contains(f.Name(), originalImage) {
				imgPath := filepath.Join(path, f.Name())
				img, format, err = loadImage(imgPath)
				if err != nil {
					return "", err
				}
				break
			}
		}

		// dont try to resize gifs
		if format == "gif" {
			return getImage(url, 0, 0)
		}
	}

	imgPath := fmt.Sprintf("%s.%s", filepath.Join(path, sizeRequested), format)
	newImg := resizeImage(img, width, height)
	err := writeImage(newImg, imgPath, format)
	if err != nil {
		return "", err
	}
	return imgPath, nil
}

func getImageFromId(id string, width, height int) (string, error) {
	leak, err := getArticleByID(id)
	if err != nil {
		//log.Println("Error getting leak:", id)
		//return defaultTo
		return "", err
	}

	if leak.ImageUrl == "" {
		return "", nil
	}

	path, err := getImage(leak.ImageUrl, width, height)
	if err != nil {
		//log.Println(err)
		//return defaultTo
		return "", err
	}
	return path, nil
}

func resizeImage(original image.Image, targetWidth, targetHeight int) (output image.Image) {
	p := original.Bounds().Size()

	w := boundNumber(targetWidth, 0, p.X)
	h := boundNumber(targetHeight, 0, p.Y)

	return resize.Resize(w, h, original, resize.Bilinear) // resize.NearestNeighbor
}

func boundNumber(number, min, max int) uint {
	return uint(math.Min(math.Max(float64(number), float64(min)), float64(max)))
}

func fetchImage(url string) (img image.Image, format string, err error) {
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		return nil, "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)

	img, format, err = decodeImage(res.Body)
	return img, format, err
}

func getImageName(url string) (name string) {
	s := strings.Split(url, "/")
	name = sanitize.Name(s[len(s)-1])
	return fmt.Sprintf("%d-%s", hashTo32(url), name)
}

func writeImage(img image.Image, path, format string) error {
	b := new(bytes.Buffer)

	switch format {
	case "jpeg":
		err := jpeg.Encode(b, img, &jpeg.Options{
			Quality: 90,
		})
		if err != nil {
			return err
		}

	case "gif":
	case "png":
		err := png.Encode(b, img)
		if err != nil {
			return err
		}
	}

	cacheFile, err := os.Create(path)
	defer func(cacheFile *os.File) {
		err := cacheFile.Close()
		if err != nil {
			log.Println(err)
		}
	}(cacheFile)
	if err != nil {
		return err
	}
	_, err = b.WriteTo(cacheFile)
	if err != nil {
		return err
	}
	return nil
}

func loadImage(path string) (img image.Image, format string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}

	img, format, err = decodeImage(f)
	if err != nil {
		return nil, "", err
	}
	err = f.Close()
	if err != nil {
		return nil, "", err
	}
	return img, format, nil
}

func decodeImage(f io.Reader) (img image.Image, format string, err error) {
	img, format, err = image.Decode(f)
	if err != nil {
		return nil, "", err
	}
	//https://stackoverflow.com/a/54210633 for gifs
	return img, format, err
}
