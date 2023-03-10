package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// TODO:
// [x] ファイルを引数で受取る
// [x] 行を解析してURLだけ抽出する
// [x] getpocket.com に登録する(レートリミットに注意)
// [ ] リファクタリングする
// [ ] ビルドスクリプト(Makefile)を作る
// [ ] 除外リストを作ってメンテできるようにする
// [ ] OAuthトークンのリフレッシュを実装する
// [ ] APIトークンを環境変数から読み取れるようにする
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// url := "https://getpocket.com/api/v3/add"
// リクエストJSONのフォーマット
// {"url":"http:\/\/pocket.co\/s8Kga",
// "title":"iTeaching: The New Pedagogy (How the iPad is Inspiring Better Ways of
// Teaching)",
// "time":1346976937,
// "consumer_key":"1234-abcd1234abcd1234abcd1234",
// "access_token":"5678defg-5678-defg-5678-defg56"}
func doPost(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	method := "POST"
	url := "https://getpocket.com/v3/add"
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("X-Acept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("got %d; X-Error=[%s]", resp.StatusCode, resp.Header.Get("X-Error"))
	}

	// fmt.Printf("%v", json.NewDecoder(resp.Body).Decode(resp))
	fmt.Println(resp.Body)
	fmt.Println(resp.StatusCode)
	return err
}

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
		urlPart := strings.Split(line, " |")[0]
		// Extract only lines beginning with http:// or https://
		fs := re.FindString(urlPart)
		if len(fs) != 0 {
			// Exclude Google services.
			isGoogleSearchURL := googleSearchURLRE.MatchString(fs)
			isGoogleDriveURL := googleDriveURLRE.MatchString(fs)
			isGoogleDocsURL := googleDocslURLRE.MatchString(fs)
			isGmailURL := gmailURLRE.MatchString(fs)
			isGoogleCalendarURL := googleCalendarURLRE.MatchString(fs)
			isGoogleMapURL := googleMapURLRE.MatchString(fs)
			if isGoogleSearchURL || isGoogleDriveURL || isGoogleDocsURL || isGmailURL || isGoogleCalendarURL || isGoogleMapURL {
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
	// for _, u := range urls {
	// 	fmt.Println(u)
	// }
	yahoo := "https://www.yahoo.co.jp"
	tag := "TestHoge"
	// consumer_key :=
	// access_token :=
	r := RegisterURL{
		yahoo,
		tag,
		"106320-f11816ace4ac49e05d72fc6",
		"822886b3-cc45-68e7-d33e-75abd0",
	}
	err = doPost(r)
	if err != nil {
		log.Fatal(err)
	}
}

// {"url":"http:\/\/pocket.co\/s8Kga",
// "title":"iTeaching: The New Pedagogy (How the iPad is Inspiring Better Ways of
// Teaching)",
// "time":1346976937,
// "consumer_key":"1234-abcd1234abcd1234abcd1234",
// "access_token":"5678defg-5678-defg-5678-defg56"
type RegisterURL struct {
	URL         string `json:"url"`
	Tags        string `json:"tags"`
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}
