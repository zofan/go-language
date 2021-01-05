package language

type Language struct {
	Alpha2 string
	Alpha3 string

	Name string

	Users []string
}

func ByAlpha3(v string) *Language {
	for _, c := range List {
		if c.Alpha3 == v {
			return &c
		}
	}

	return nil
}

func ByAlpha2(v string) *Language {
	for _, c := range List {
		if c.Alpha2 == v {
			return &c
		}
	}

	return nil
}
