package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func verifyURL(instagramLink string) error {
	// checking for a 200 from server
	resp, err := http.Get(instagramLink)
	if err != nil {
		return errors.New("ERROR: failed to make http request")
	}
	if resp.StatusCode != 200 {
		return errors.New("ERROR: received non-OK http resonse")
	}
	return nil
}

func getPostTitle(instagramLink string) string {
	resp, _ := http.Get(instagramLink)
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			break
		case html.StartTagToken:
			t := z.Token()
			tagname := t.Data
			if tagname == "title" {
				tt := z.Next()
				if tt == html.TextToken {
					fullTitle := (z.Token().Data)
					titleList := strings.Split(fullTitle, ":")
					title := strings.Join(titleList[1:], "")
					// instagram returns some funky characters
					title = strings.Replace(title, "”", "", 1)
					title = strings.Replace(title, "“", "", 1)
					return title
				}
			}
		}
	}
}

func getFileName(s string) (string, error) {
	var fileName string
	fileNameSlice := strings.Split(s, ".")
	if len(fileNameSlice) > 1 {
		fileExtention := fileNameSlice[len(fileNameSlice)-1]
		if strings.Compare(fileExtention, "jpeg") != 0 && strings.Compare(fileExtention, "png") != 0 {
			return "", errors.New("ERROR: only .jpeg and .png file extentions supported")
		}
		fileName = s

	} else {
		fileName = s + ".jpeg"
	}
	fileName = strings.ReplaceAll(fileName, "\n", "")
	return fileName, nil
}

func isPrivateUser(instagramLink string) error {
	title := getPostTitle(instagramLink)
	if title == "" {
		return errors.New("This account is private")
	}
	return nil
}

func getMediaLink(instagramLink string) string {
	var mediaLink string
	if instagramLink[len(instagramLink)-1] == '/' {
		mediaLink = instagramLink + "media"
	} else {
		mediaLink = instagramLink + "/media"
	}
	return mediaLink
}

func writeFile(fileName string, bytes []byte, downloadPath string) error {
	if strings.Compare(downloadPath, "") != 0 {
		// determine if local path or not
		var isAbsolutePath bool = false
		if downloadPath[0] == '/' {
			isAbsolutePath = true
		}
		if isAbsolutePath {
			err := os.MkdirAll(downloadPath, 0770)
			if err != nil {
				return err
			}
			if downloadPath[len(downloadPath)-1] == '/' {
				downloadPath = downloadPath + fileName
			} else {
				downloadPath = downloadPath + "/" + fileName
			}
		} else {
			currentWorkingDirectory, err := os.Getwd()
			if err != nil {
				return err
			}
			err = os.MkdirAll(currentWorkingDirectory+"/"+downloadPath, 0770)
			if err != nil {
				return err
			}
			if downloadPath[len(downloadPath)-1] == '/' {
				downloadPath = currentWorkingDirectory + "/" + downloadPath + fileName
			} else {
				downloadPath = currentWorkingDirectory + "/" + downloadPath + "/" + fileName
			}
		}
		fileName = downloadPath
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return err
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
