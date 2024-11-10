package game

import (
	"time"
)

const (
	SHOULD_RENDER_FRAME = true
	SHOULD_STREAM_FRAME = true
	FPS                 = 10
	FRAME_WIDTH         = 1024
	FRAME_HEIGHT        = 768
	CHUNK_SIZE          = 48
	ROWS                = FRAME_HEIGHT / CHUNK_SIZE
	COLS                = FRAME_WIDTH / CHUNK_SIZE
)

type SnakeLoop struct {
	state          *Snake
	commandChannel chan string
	stateCh        chan *Snake
	closeSignal    chan bool
	frameTicker    *time.Ticker
}

type SnakeLoopInit struct {
	CommandChannel chan string
	SnakeChannel   chan *Snake
	CloseSignal    chan bool
}

func NewSnakeLoop(options *SnakeLoopInit) *SnakeLoop {
	return &SnakeLoop{
		state:          NewSnake(ROWS, COLS),
		commandChannel: options.CommandChannel,
		stateCh:        options.SnakeChannel,
		closeSignal:    options.CloseSignal,
		frameTicker:    time.NewTicker(2 * time.Second / FPS),
	}
}

func (gl *SnakeLoop) Start() {
	defer gl.frameTicker.Stop()

	for {
		select {
		case command := <-gl.commandChannel:
			gl.handleCommand(command)

		case <-gl.frameTicker.C:
			gl.updateSnake(nil)
		}
	}
}

func (gl *SnakeLoop) handleCommand(command string) {
	gl.updateSnake(&command)
}

func (gl *SnakeLoop) updateSnake(command *string) {
	gameOver := !gl.state.HandleCommand(command)
	if gameOver {
		gl.closeSignal <- true
	}
	if command == nil {
		gl.stateCh <- gl.state
	}
}
