package main

import (
	"fmt"
	"github.com/zofan/go-country"
	"github.com/zofan/go-fwrite"
	"github.com/zofan/go-language"
	"github.com/zofan/go-req"
	"github.com/zofan/go-xmlre"
	"html"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func main() {
	fmt.Println(Update())
}

func Update() error {
	var (
		httpClient = req.New(req.DefaultConfig)
		list       = make(map[string]*language.Language)
	)

	resp := httpClient.Get(`http://www.loc.gov/standards/iso639-2/php/code_list.php`)
	if resp.Error() != nil {
		return resp.Error()
	}

	body := string(resp.ReadAll())
	body = html.UnescapeString(body)

	rowRe := xmlre.Compile(`<td scope="row">(\w+)</td><td>(\w+)</td><td>(\w+)</td>`)

	for _, row := range rowRe.FindAllStringSubmatch(body, -1) {
		l := &language.Language{
			Alpha3: strings.ToUpper(strings.TrimSpace(row[1])),
			Alpha2: strings.ToUpper(strings.TrimSpace(row[2])),
			Name:   strings.TrimSpace(row[3]),
		}

		list[l.Alpha3] = l
	}

	// ---

	for _, c := range country.List {
		for _, cl := range c.Languages {
			if _, ok := list[cl]; ok {
				list[cl].Users = append(list[cl].Users, c.Alpha3)
			}
		}
	}

	// ---

	updateTags(list)

	var tpl []string

	tpl = append(tpl, `package language`)
	tpl = append(tpl, ``)
	tpl = append(tpl, `// Updated at: `+time.Now().String())
	tpl = append(tpl, `var List = []Language{`)

	for _, l := range list {
		s := fmt.Sprintf(`%#v`, *l) + `,`
		s = strings.ReplaceAll(s, `language.Language`, ``)
		tpl = append(tpl, s)
	}

	tpl = append(tpl, `}`)
	tpl = append(tpl, ``)

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)

	return fwrite.WriteRaw(dir+`/../db.go`, []byte(strings.Join(tpl, "\n")))
}

func updateTags(list map[string]*language.Language) {
	wordSplitRe := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	wordMap := map[string][]*language.Language{}

	for _, l := range list {
		name := strings.ToLower(l.Name + ` ` + strings.Join(l.AltNames, ` `))
		words := wordSplitRe.Split(name, -1)
		for _, w := range words {
			if len(w) > 0 {
				wordMap[w] = append(wordMap[w], l)
			}
		}
		l.Tags = []string{}
	}

	for w, ls := range wordMap {
		if len(ls) == 1 {
			ls[0].Tags = append(ls[0].Tags, w)
		}
	}
}
