package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/francescomari/iterm2"
	"github.com/mitchellh/go-wordwrap"
)

func main() {
	rand.Seed(time.Now().Unix())

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var random bool

	flag.BoolVar(&random, "random", false, "Show a random comic")
	flag.Parse()

	if random {
		return showRandomComic()
	}

	switch flag.NArg() {
	case 0:
		return showLatestComic()
	case 1:
		return showComicByNumber(flag.Arg(0))
	default:
		return fmt.Errorf("too many arguments")
	}
}

func showLatestComic() error {
	comic, err := readCurrentComic()
	if err != nil {
		return fmt.Errorf("read current comic: %v", err)
	}

	if err := showComic(comic); err != nil {
		return fmt.Errorf("show comic: %v", err)
	}

	return nil
}

func showRandomComic() error {
	currentComic, err := readCurrentComic()
	if err != nil {
		return fmt.Errorf("read current comic: %v", err)
	}

	var attempts int

	for {
		attempts++

		if attempts > 3 {
			return fmt.Errorf("cannot find a random comic")
		}

		randomComic, err := readComicByNumber(1 + rand.Intn(currentComic.Num))
		if err != nil {
			return fmt.Errorf("read comic by number: %v", err)
		}
		if randomComic == nil {
			continue
		}

		if err := showComic(randomComic); err != nil {
			return fmt.Errorf("show comic: %v", err)
		}

		return nil
	}
}

func showComicByNumber(numArg string) error {
	num, err := strconv.Atoi(numArg)
	if err != nil {
		return fmt.Errorf("not a comic number: %v", err)
	}

	if num < 1 {
		return fmt.Errorf("not a comic number: %v", numArg)
	}

	comic, err := readComicByNumber(num)
	if err != nil {
		return fmt.Errorf("read comic by number: %v", err)
	}
	if comic == nil {
		return fmt.Errorf("comic not found")
	}

	if err := showComic(comic); err != nil {
		return fmt.Errorf("show comic: %v", err)
	}

	return nil
}

func showComic(comic *comic) error {
	var err error

	image, err := readBytes(comic.Image)
	if err != nil {
		return fmt.Errorf("read image: %v", err)
	}

	fmt.Printf("#%d %s\n\n", comic.Num, comic.Title)

	if err := inlineImage(image); err != nil {
		return fmt.Errorf("inline image: %v", err)
	}

	fmt.Printf("\n\n%s\n", wordwrap.WrapString(comic.Alt, 70))

	return nil
}

func inlineImage(image []byte) error {
	switch {
	case os.Getenv("TERM") == "xterm-kitty":
		return inlineImageWithKitty(ensurePNG(image))
	default:
		return inlineImageWithIterm2(image)
	}
}

func inlineImageWithIterm2(image []byte) error {
	if _, err := iterm2.InlineImage(image, iterm2.WithInline(true)); err != nil {
		return err
	}

	return nil
}

func inlineImageWithKitty(image []byte) error {
	const chunkSize = 4096

	for len(image) > 0 {
		var chunk []byte

		if len(image) > chunkSize {
			chunk, image = image[:chunkSize], image[chunkSize:]
		} else {
			chunk, image = image, nil
		}

		fmt.Print("\033_G")
		fmt.Print("a=T,")
		fmt.Print("f=100,")

		if len(image) > 0 {
			fmt.Print("m=1")
		} else {
			fmt.Print("m=0")
		}

		fmt.Print(";")
		fmt.Print(base64.StdEncoding.EncodeToString(chunk))
		fmt.Print("\033\\")
	}

	return nil
}

func ensurePNG(image []byte) []byte {
	decoded, err := jpeg.Decode(bytes.NewReader(image))
	if err != nil {
		return image
	}

	var encoded bytes.Buffer

	if err := png.Encode(&encoded, decoded); err != nil {
		return image
	}

	return encoded.Bytes()
}

type comic struct {
	Title string `json:"title"`
	Image string `json:"img"`
	Alt   string `json:"alt"`
	Num   int    `json:"num"`
}

func readCurrentComic() (*comic, error) {
	comic, err := readComicByURL("https://xkcd.com/info.0.json")
	if err != nil {
		return nil, fmt.Errorf("current comic not found")
	}

	return comic, err
}

func readComicByNumber(num int) (*comic, error) {
	return readComicByURL(fmt.Sprintf("https://xkcd.com/%d/info.0.json", num))
}

func readComicByURL(url string) (_ *comic, e error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("perform request: %v", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil && e == nil {
			e = fmt.Errorf("close body: %v", err)
		}
	}()

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response status code: %d", res.StatusCode)
	}

	var c comic

	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		return nil, fmt.Errorf("invalid response body: %v", err)
	}

	return &c, nil
}

func readBytes(url string) (_ []byte, e error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("perform request: %v", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil && e == nil {
			e = fmt.Errorf("close body: %v", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %v", err)
	}

	return data, nil
}
