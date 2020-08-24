package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var instagramLink string

func downloadPost(title string, downloadPath string) error {
	fileName, err := getFileName(title)
	if err != nil {
		return err
	}

	mediaLink := getMediaLink(instagramLink) // getting link which returns post media
	resp, err := http.Get(mediaLink)
	if err != nil {
		return err
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	err = writeFile(fileName, bytes, downloadPath)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// determine name of file
	var title string
	userProvidedTitle := flag.String("filename", "", "user-provided text to name the file")
	helpFlag := flag.Bool("h", false, "")
	flag.Parse()

	if *helpFlag {
		fmt.Println("usage: instagramDownloader [-fileName fileName] <instagram link> [download path]")
		return
	}
	args := flag.Args()

	if len(args) == 1 && args[0] == "help" {
		fmt.Println("usage: instagramDownloader [-fileName fileName] <instagram link> [download path]")
	} else if len(args) == 0 || len(args) > 2 {
		fmt.Println("ERORR: Please pass a link to a post you'd like to download")
	} else {

		// grabbing link to content from command line
		instagramLink = args[0]
		downloadPath := "."
		if len(args) == 2 { // user passed a download path
			downloadPath = args[1]
		}

		err := isPrivateUser(instagramLink)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = verifyURL(instagramLink)
		if err != nil {
			fmt.Println(err)
			return
		}

		if *userProvidedTitle == "" {
			title = getPostTitle(instagramLink)
		} else {
			title = *userProvidedTitle
		}

		err = downloadPost(title, downloadPath)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Download complete")
		}
	}
}
