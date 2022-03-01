package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/francescomari/iterm2"
	"github.com/mitchellh/go-wordwrap"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	comic, err := readCurrentComic()
	if err != nil {
		return fmt.Errorf("read current comic: %v", err)
	}

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
}

func readCurrentComic() (*comic, error) {
	res, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
		return nil, fmt.Errorf("perform request: %v", err)
	}

	defer res.Body.Close()

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
