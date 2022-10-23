package color

type Color string

const (
	Green Color = "positive"
	Red   Color = "negative"
	Blue  Color = "primary"
	White Color = "default"
)

func (c Color) String() string {
	return string(c)
}
