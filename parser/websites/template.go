package websites

type WebBase struct {
	urlBase  string
	resCache []map[string]string
}

func (web *WebBase) ShowFindResult(filterMap map[string]int, len int) [][]string {
	infos := [][]string{}
	for _, row := range web.resCache {
		info := make([]string, len)
		for key, idx := range filterMap {
			info[idx] = row[key]
		}
		infos = append(infos, info)
	}

	return infos
}

func (web *WebBase) GetMagnets(selectNums []int) []string {
	res := make([]string, len(selectNums))
	for idx, num := range selectNums {
		res[idx] = web.resCache[num]["magnet"]
	}
	return res
}

type WebParser interface {
	Request(args []string) error
	ShowFindResult(map[string]int, int) [][]string
	GetMagnets(selectNums []int) []string
}
