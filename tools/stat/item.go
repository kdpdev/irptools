package stat

func NewItem() Item {
	return Item{
		Ints: map[string]int{},
		Strs: map[string]map[string]int{},
	}
}

type Item struct {
	Ints map[string]int
	Strs map[string]map[string]int
}

func (this *Item) IncInt(key string, value int) {
	this.Ints[key] += value
}

func (this *Item) AddStr(key string, value string) {
	this.addStrN(key, value, 1)
}

func (this *Item) addStrN(key string, value string, n int) {
	strs, ok := this.Strs[key]
	if !ok {
		this.Strs[key] = map[string]int{value: n}
	} else {
		strs[value] += n
	}
}

func (this *Item) AddItem(item Item) {
	for k, v := range item.Ints {
		this.IncInt(k, v)
	}
	for k, v := range item.Strs {
		for kk, vv := range v {
			this.addStrN(k, kk, vv)
		}
	}
}
