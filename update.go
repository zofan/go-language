package language

import (
	"fmt"
	"github.com/zofan/go-country"
	"github.com/zofan/go-fwrite"
	"github.com/zofan/go-req"
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
	body = strings.ReplaceAll(body, `&nbsp;`, ` `)

	rowRe := regexp.MustCompile(`(?s)<td scope="row">(\w+)</td>\s*<td>(\w+)</td>\s*<td>(\w+)</td>`)

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
		tpl = append(tpl, `	},`)
	}

	tpl = append(tpl, `}`)
	tpl = append(tpl, ``)

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)

	return fwrite.WriteRaw(dir+`/language_db.go`, []byte(strings.Join(tpl, "\n")))
}
