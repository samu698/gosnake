package screen

type Pos struct {
	X, Y int
}

func NewPos[T, U int | uint](x T, y U) Pos {
	return Pos{X: int(x), Y: int(y)}
}

func (this *Pos) Equal(other Pos) bool {
	return this.X == other.X && this.Y == other.Y
}

func (this *Pos) Add(other Pos) {
	this.X += other.X
	this.Y += other.Y
}

func (this *Pos) Sub(other Pos) {
	this.X -= other.X
	this.Y -= other.Y
}

func (this *Pos) Addv(x, y int) {
	this.X += x
	this.Y += y
}

func min(a, b int) int {
	if a < b { return a }
	return b
}

func max(a, b int) int {
	if a > b { return a }
	return b
}

// Checks if the point is inside the rect constructed from the points p1 and p2
func (this *Pos) IsInside(p1, p2 Pos) bool {
	maxP := Pos{ max(p1.X, p2.X), max(p1.Y, p2.Y) }
	minP := Pos{ min(p1.X, p2.X), min(p1.Y, p2.Y) }

	return this.X >= minP.X && this.X <= maxP.X &&
		this.Y >= minP.Y && this.Y <= maxP.Y
}
