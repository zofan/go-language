package language

import "strings"

type Language struct {
	Alpha2 string
	Alpha3 string

	Name     string
	AltNames []string
	Tags     []string

	Users []string
}

func Get(v string) *Language {
	for _, l := range List {
		if l.Alpha3 == v || l.Alpha2 == v {
			return &l
		}
	}

	return nil
}

func ByName(v string) *Language {
	v = strings.ToLower(v)

	for _, l := range List {
		if strings.ToLower(l.Name) == v {
			return &l
		}

		for _, n := range l.AltNames {
			if strings.ToLower(n) == v {
				return &l
			}
		}

		for _, n := range l.Tags {
			if strings.ToLower(n) == v {
				return &l
			}
		}
	}

	return nil
}
