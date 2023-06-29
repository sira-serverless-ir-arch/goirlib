package filter

type Filter interface {
	Process(text []string) []string
}
