// internal/transport/draw.go

package transport

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type DrawingTool string

const (
	ToolBrush  DrawingTool = "brush"
	ToolEraser DrawingTool = "eraser"
)

type Stroke struct {
	Points []Point `json:"points"`

	Color     string      `json:"color"`
	Thickness float64     `json:"thickness"`
	Tool      DrawingTool `json:"tool"`
}

type DrawStart struct {
	Point Point `json:"point"`

	Color     string      `json:"color"`
	Thickness float64     `json:"thickness"`
	Tool      DrawingTool `json:"tool"`
}

type DrawMove struct {
	Stroke Stroke `json:"stroke"`
}
