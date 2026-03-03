package brownsugar

type RenderContext struct {
	Width  int
	Height int
}

func RenderWithCondition(condition bool, fn func() string) string {
	if condition {
		return fn()
	}

	return ""
}
