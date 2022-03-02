package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
	image, err := readBytes(comic.Image)
	if err != nil {
		return fmt.Errorf("read image: %v", err)
	}

	fmt.Printf("%s\n\n", comic.Title)

	if _, err := iterm2.InlineImage(image, iterm2.WithInline(true)); err != nil {
		return fmt.Errorf("inline image: %v", err)
	}

	fmt.Printf("\n\n%s\n", wordwrap.WrapString(comic.Alt, 70))

	return nil
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

func readComicByURL(url string) (*comic, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("perform request: %v", err)
	}

	defer res.Body.Close()

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

func readBytes(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("perform request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %v", err)
	}

	return data, nil
}
