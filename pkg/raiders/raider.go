package raiders

import "sort"

type Raider struct {
	Id     uint
	Name   string
	Points uint
}

func SortSlice(rs []Raider) {
	sort.Slice(rs, func(i, j int) bool {
		if rs[i].Points < rs[j].Points {
			return true
		} else if rs[i].Points > rs[j].Points {
			return false
		}
		return rs[i].Name[0] < rs[j].Name[0]
	})
}
