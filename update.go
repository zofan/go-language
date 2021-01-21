package language

import (
	"fmt"
	"github.com/zofan/go-country"
	"github.com/zofan/go-fwrite"
	"github.com/zofan/go-req"
	"github.com/zofan/go-xmlre"
	"html"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func Update() error {
	var (
		httpClient = req.New(req.DefaultConfig)
		list       = make(map[string]*Language)
	)

	resp := httpClient.Get(`http://www.loc.gov/standards/iso639-2/php/code_list.php`)
	if resp.Error() != nil {
		return resp.Error()
	}

	body := string(resp.ReadAll())
	body = html.UnescapeString(body)

	rowRe := xmlre.Compile(`<td scope="row">(\w+)</td><td>(\w+)</td><td>(\w+)</td>`)

	for _, row := range rowRe.FindAllStringSubmatch(body, -1) {
		l := &Language{
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
		tpl = append(tpl, `	{`)
		tpl = append(tpl, `		Alpha3:    "`+l.Alpha3+`",`)
		tpl = append(tpl, `		Alpha2:    "`+l.Alpha2+`",`)
		tpl = append(tpl, `		Name:      "`+l.Name+`",`)
		tpl = append(tpl, `		Users:     `+fmt.Sprintf(`%#v`, l.Users)+`,`)
		tpl = append(tpl, `		AltNames:  `+fmt.Sprintf(`%#v`, l.AltNames)+`,`)
		tpl = append(tpl, `		Tags:      `+fmt.Sprintf(`%#v`, l.Tags)+`,`)
		tpl = append(tpl, `	},`)
	}

	tpl = append(tpl, `}`)
	tpl = append(tpl, ``)

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)

	return fwrite.WriteRaw(dir+`/language_db.go`, []byte(strings.Join(tpl, "\n")))
}

func updateTags(list map[string]*Language) {
	wordSplitRe := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	wordMap := map[string][]*Language{}

	for _, c := range list {
		name := strings.ToLower(c.Name + ` ` + strings.Join(c.AltNames, ` `))
		words := wordSplitRe.Split(name, -1)
		for _, w := range words {
			if len(w) > 0 {
				wordMap[w] = append(wordMap[w], c)
			}
		}
		c.Tags = []string{}
	}

	for w, cs := range wordMap {
		if len(cs) == 1 {
			cs[0].Tags = append(cs[0].Tags, w)
		}
	}
}
