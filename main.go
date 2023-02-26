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

// 引数にファイルを受け取る
// ファイルの構造:
// https://github.com/japan-clojurians/curriculum | japan-clojurians/curriculum
// <URL>\s\|\s"Description" か 空行 の2種類で構成されている
func main() {
	fmt.Println("Hello")

	var (
		f = flag.String("f", "", "File name with URL listed.")
	)
	flag.Parse()

	// There has to be a flag.
	if flag.NFlag() == 0 {
		fmt.Print("we need flag option. please specify filename.")
		os.Exit(1)
	}

	if &f != nil {
		fmt.Printf("param -f: %s\n", *f)
	}

	fp, err := os.Open(*f)
	handleError(err)
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	re := regexp.MustCompile("^http(.*)://(.*)")
	urls := []string{}

	// 登録対象とするURLのスライスを作る
	for scanner.Scan() {
		line := scanner.Text()
		seperated_line := strings.Split(line, " |")
		// 空行を間引く
		// http:// もしくは https:// から始まる行だけ抽出する
		// 行を分割してURL部分だけ抽出する
		fs := re.FindString(seperated_line[0])
		if len(fs) != 0 {
			// https://www.google.com/search?q= にマッチしていたなら除外する
			urls = append(urls, fs)
		}
	}
	// urlsの要素が0個なら終了する
	if len(urls) == 0 {
		fmt.Println("No more than one URL was found.")
		os.Exit(1)
	}

	// urls を使って getpocket.com に登録する。
	// タグを今日の日付とする。
	for _, u := range urls {
		fmt.Println(u)
	}
}
