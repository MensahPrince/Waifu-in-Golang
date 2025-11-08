package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Define the interface
type Waifu interface {
	Check() error
	Download() error
}

// Base struct
type BaseContent struct {
	link   string
	w      []string // working links returned from API
	IsDone bool
}

// Concrete types that embed BaseContent
type Sfw struct {
	BaseContent
}

type Nsfw struct {
	BaseContent
}

// Implement the Check method for *BaseContent
func (b *BaseContent) Check() error {
	resp, err := http.Get(b.link)
	if err != nil {
		b.IsDone = false
		return fmt.Errorf("content with link \"%s\" could not be reached: %v", b.link, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b.IsDone = false
		return fmt.Errorf("bad status: %d for %s", resp.StatusCode, b.link)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	b.w = []string{}
	for _, v := range data {
		if str, ok := v.(string); ok {
			b.w = append(b.w, str)
			fmt.Println(str)
		}
	}

	b.IsDone = true
	fmt.Printf("%s: hit status 200 OK\n", b.link)
	return nil
}

var downloadCounter int = 0

func (b *BaseContent) Download() error {
	if !b.IsDone {
		return errors.New("content not verified yet; run Check() first")
	}

	if len(b.w) == 0 {
		return errors.New("no valid URLs to download")
	}

	src := b.w[0]
	dst := "./downloads"

	// Create downloads folder if it doesn't exist
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		if err := os.Mkdir(dst, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create downloads folder: %v", err)
		}
	}

	// Increment the counter
	downloadCounter++
	fileName := fmt.Sprintf("%s/waifu_%d.gif", dst, downloadCounter)

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	resp, err := http.Get(src)
	if err != nil {
		return fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	fmt.Printf("✅ Download successful: %s\n", fileName)
	return nil
}


// Function to operate on any Waifu
func GoWaifu(w Waifu) {
	fmt.Println("Starting Waifu Checker...")
	if err := w.Check(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Success! Proceeding to download, hang tight...")
	if err := w.Download(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Download Success!")
}




func main() {
	s := &Sfw{BaseContent{
		link:   "https://api.waifu.pics/sfw/pat",
		w:      []string{},
		IsDone: false,
	}}

	GoWaifu(s)
}
