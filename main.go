package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// TODO:
// [x] ファイルを引数で受取る
// [x] 行を解析してURLだけ抽出する
// [ ] getpocket.com に登録する(レートリミットに注意)
// [ ] ビルドスクリプト(Makefile)を作る
// [ ] 除外リストを作ってメンテできるようにする
// [ ] APIトークンを環境変数から読み取れるようにする
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Accepts a filename as an argument.
// Structure of the document file:
// <URL>\s\|\s"Description" or empty line.
//
// for example:
// https://github.com/japan-clojurians/curriculum | japan-clojurians/curriculum

func main() {
	var (
		f = flag.String("f", "", "File name with URL listed.")
	)
	flag.Parse()

	// There has to be a flag.
	if flag.NFlag() == 0 {
		fmt.Print("we need flag option. please specify filename.")
		os.Exit(1)
	}

	fp, err := os.Open(*f)
	handleError(err)
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	re := regexp.MustCompile("^http(.*)://(.*)")
	googleDriveURLRE := regexp.MustCompile(`^https://drive.google.com/(.*)`)
	googleDocslURLRE := regexp.MustCompile(`^https://docs.google.com/(.*)`)
	gmailURLRE := regexp.MustCompile(`^https://mail.google.com/`)
	googleCalendarURLRE := regexp.MustCompile(`https://calendar.google.com/(.*)`)
	googleSearchURLRE := regexp.MustCompile(`^https://www.google.com/search(.*)`)
	googleMapURLRE := regexp.MustCompile(`^https://www.google.co.jp/maps(.*)`)
	urls := []string{}

	// Create a slice of the URL
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line and take out only the URL part.
		url_part := strings.Split(line, " |")[0]
		// Extract only lines beginning with http:// or https://
		fs := re.FindString(url_part)
		if len(fs) != 0 {
			// Exclude Google services.
			is_googleSearchURL := googleSearchURLRE.MatchString(fs)
			is_googleDriveURL := googleDriveURLRE.MatchString(fs)
			is_googleDocsURL := googleDocslURLRE.MatchString(fs)
			is_gmailURL := gmailURLRE.MatchString(fs)
			is_googleCalendarURL := googleCalendarURLRE.MatchString(fs)
			is_googleMapURL := googleMapURLRE.MatchString(fs)
			if is_googleSearchURL || is_googleDriveURL || is_googleDocsURL || is_gmailURL || is_googleCalendarURL || is_googleMapURL {
				continue
			}
			// Adding a URL to a slice.
			urls = append(urls, fs)
		}
	}

	// If there are zero elements in the slice, i.e., there are none to register, then the program is terminated.
	if len(urls) == 0 {
		fmt.Println("No more than one URL was found.")
		os.Exit(1)
	}

	// TODO: urls を使って getpocket.com に登録する。
	// TODO: タグを今日の日付とする。
	for _, u := range urls {
		fmt.Println(u)
	}
}
