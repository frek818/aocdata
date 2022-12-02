package aocdata

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"strings"
)

const URL = "https://adventofcode.com/%d/day/%d"

func getHTTPClient() http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	client := http.Client{
		Jar: jar,
	}
	return client
}

func getAOCSessionCookie() string {
	session := os.Getenv("AOC_SESSION")
	if session != "" {
		return session
	}

	return getSessionCookieFromFile()
}

func getSessionCookieFromFile() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	token_file := filepath.Join(homedir, ".config", "aocd", "token")
	dat, err := ioutil.ReadFile(token_file)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Trim(string(dat), "\r\n")
}

func GetInputData(year int, day int) string {
	session_cookie_value := getAOCSessionCookie()
	session_cookie := &http.Cookie{
		Name:   "session",
		Value:  session_cookie_value,
		MaxAge: 300,
	}

	input_url := fmt.Sprintf(URL, year, day) + "/input"
	req, err := http.NewRequest("GET", input_url, nil)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	}
	req.AddCookie(session_cookie)

	client := getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error occured. Error is: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Input data request returned not StausOK")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	return bodyString
}
