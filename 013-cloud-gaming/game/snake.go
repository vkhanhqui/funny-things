package game

import (
	"math/rand"
)

type Direction string
type PositionValue int

type Position struct {
	X int
	Y int
}

func (p *Position) Equal(other Position) bool {
	return p.X == other.X && p.Y == other.Y
}

const (
	UP    Direction = "UP"
	DOWN  Direction = "DOWN"
	LEFT  Direction = "LEFT"
	RIGHT Direction = "RIGHT"

	EMPTY PositionValue = 0
	HEAD  PositionValue = 1
	BODY  PositionValue = 2
	FOOD  PositionValue = 3
)

type Snake struct {
	snakeDirection Direction
	snakeSize      int
	snake          []Position
	matrix         [][]PositionValue
	rows           int
	cols           int
	foodPosition   Position
}

func NewSnake(rows, cols int) *Snake {
	matrix := make([][]PositionValue, rows)
	for i := range matrix {
		matrix[i] = make([]PositionValue, cols)
	}

	headPos := Position{X: cols / 2, Y: rows / 2}
	snake := []Position{headPos}
	state := &Snake{
		snakeDirection: RIGHT,
		snakeSize:      3,
		matrix:         matrix,
		rows:           rows,
		cols:           cols,
	}

	state.generateFood()
	for i := 1; i < state.snakeSize; i++ {
		bodyPos := Position{Y: headPos.Y, X: headPos.X - i}
		snake = append(snake, bodyPos)
		state.setAt(bodyPos, BODY)
	}

	state.snake = snake
	state.setAt(headPos, HEAD)
	return state
}

func (g *Snake) GetMatrix() [][]PositionValue {
	return g.matrix
}

func (g *Snake) setAt(pos Position, value PositionValue) {
	g.matrix[pos.Y][pos.X] = value
}

func (g *Snake) at(pos Position) PositionValue {
	return g.matrix[pos.Y][pos.X]
}

func (g *Snake) LastSnakePart() Position {
	return g.snake[len(g.snake)-1]
}

func (g *Snake) IsHeadingUp() bool {
	return g.snakeDirection == UP
}

func (g *Snake) IsHeadingDown() bool {
	return g.snakeDirection == DOWN
}

func (g *Snake) IsHeadingLeft() bool {
	return g.snakeDirection == LEFT
}

func (g *Snake) IsHeadingRight() bool {
	return g.snakeDirection == RIGHT
}

func (g *Snake) generateFood() {
	randomX := rand.Intn(g.cols)
	randomY := rand.Intn(g.rows)

	foodP := Position{X: randomX, Y: randomY}
	if g.at(foodP) != EMPTY {
		g.generateFood()
	}

	g.foodPosition = foodP
	g.setAt(g.foodPosition, FOOD)
}

func (g *Snake) HandleCommand(command *string) bool {
	return g.updateSnake(command)
}

func (g *Snake) updateSnake(command *string) bool {
	if command != nil && g.validateNewDirection(Direction(*command)) {
		g.snakeDirection = Direction(*command)
	}

	collided := g.move()
	return !collided
}

func (g *Snake) validateNewDirection(newDirection Direction) bool {
	if g.IsHeadingUp() && newDirection == DOWN {
		return false
	}
	if g.IsHeadingDown() && newDirection == UP {
		return false
	}
	if g.IsHeadingLeft() && newDirection == RIGHT {
		return false
	}
	if g.IsHeadingRight() && newDirection == LEFT {
		return false
	}
	return true
}

func (g *Snake) move() bool {
	headPos := g.snake[0]
	g.setAt(headPos, BODY)

	switch g.snakeDirection {
	case UP:
		headPos.Y = (headPos.Y - 1 + g.rows) % g.rows
	case DOWN:
		headPos.Y = (headPos.Y + 1) % g.rows
	case LEFT:
		headPos.X = (headPos.X - 1 + g.cols) % g.cols
	case RIGHT:
		headPos.X = (headPos.X + 1) % g.cols
	}

	if g.at(headPos) == BODY {
		return true
	}

	if headPos.Equal(g.foodPosition) {
		g.snakeSize++
		g.generateFood()
	}

	g.setAt(headPos, HEAD)
	g.setAt(g.LastSnakePart(), EMPTY)
	g.snake = append([]Position{headPos}, g.snake[:g.snakeSize-1]...)
	return false
}
