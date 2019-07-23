package websites

type WebBase struct {
	urlBase  string
	resCache []map[string]string
}

type WebParser interface {
	Request(args []string) error
	ShowFindResult(map[string]int, int) [][]string
	GetMagnets(selectNums []int) []string
}
