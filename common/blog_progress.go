package common

type BlogProgress struct {
	BlogType  int
	ViewCount int
}

func (bp BlogProgress) GetViewCount() int {
	return bp.ViewCount
}
func (bp *BlogProgress) SetViewCount(val int) {
	bp.ViewCount = val
}
