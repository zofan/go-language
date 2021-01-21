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
	for _, c := range List {
		if c.Alpha3 == v || c.Alpha2 == v {
			return &c
		}
	}

	return nil
}

func ByName(v string) *Language {
	v = strings.ToLower(v)

	for _, c := range List {
		if strings.ToLower(c.Name) == v {
			return &c
		}

		for _, n := range c.AltNames {
			if strings.ToLower(n) == v {
				return &c
			}
		}

		for _, n := range c.Tags {
			if strings.ToLower(n) == v {
				return &c
			}
		}
	}

	return nil
}
