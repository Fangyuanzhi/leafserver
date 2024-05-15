package stl

type Pair struct {
	Key  uint64
	Name string
	Val  float32
}

type Top struct {
	PairTop []Pair
	len     int
}

func (top *Top) Init() {
	top.PairTop = make([]Pair, 0)
	top.len = 0
}

func (top *Top) push(el Pair) int32 {
	if top.len == 0 {
		top.len = 1
		top.PairTop = append(top.PairTop, el)
		return 1
	}
	temp := make([]Pair, 0)
	ix := top.len
	for idx, Val := range top.PairTop {
		if Val.Val > el.Val {
			temp = append(temp, Pair{Key: Val.Key, Name: Val.Name, Val: Val.Val})
		} else {
			ix = idx
			break
		}
	}
	temp = append(temp, el)
	if ix != top.len {
		temp = append(temp, top.PairTop[ix:]...)
	}
	top.PairTop = temp
	top.len++
	return int32(ix + 1)
}

func (top *Top) pop(Key uint64) {
	if top.len == 0 {
		return
	}
	for idx, Val := range top.PairTop {
		if Val.Key == Key {
			top.PairTop = append(top.PairTop[:idx], top.PairTop[idx+1:]...)
			top.len--
		}
	}
}

func (top *Top) Update(el Pair) int32 {
	top.pop(el.Key)
	key := top.push(el)
	return key
}
