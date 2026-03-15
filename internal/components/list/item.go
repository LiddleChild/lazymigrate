package list

type Item interface {
	Render(ctx Context) string

	Height() int
}
