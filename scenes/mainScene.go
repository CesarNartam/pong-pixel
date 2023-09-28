package scenes

import (
	"fmt"
	"image/color"
	"pong-inverso-pixel/models"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type MainScene struct{}

func NewMainScene() *MainScene {
	return &MainScene{}
}

const (
	width        = 800
	height       = 600
	gameDuration = time.Minute
)

var (
	player1         = models.NewPlayer(pixel.Rect{Min: pixel.V(10, height/2-50), Max: pixel.V(20, height/2+50)}, 700.00)
	player2         = models.NewPlayer(pixel.Rect{Min: pixel.V(width-20, height/2-50), Max: pixel.V(width-10, height/2+50)}, 700.00)
	ball            = models.NewBall(pixel.V(width/2, height/2), pixel.V(400, 400))
	gameStartTime   = time.Now()
	last            = time.Now()
	gameOver        = false
	gameTimeElapsed time.Duration
)

var (
	gameOverCh = make(chan struct{})
)

func (s *MainScene) updateBall() {
	for !gameOver {
		dt := time.Since(last).Seconds()
		last = time.Now()
		ball.Update(ball.Body.Add(ball.Speed.Scaled(dt)))

		if ball.Body.X < player1.Body.Max.X && ball.Body.Y >= player1.Body.Min.Y && ball.Body.Y <= player1.Body.Max.Y {
			ball.Speed.X = -ball.Speed.X
			relativePos := (ball.Body.Y - player1.Body.Min.Y) / (player1.Body.Max.Y - player1.Body.Min.Y)
			ball.Speed.Y = (relativePos - 0.5) * 800
		}

		if ball.Body.X > player2.Body.Min.X && ball.Body.Y >= player2.Body.Min.Y && ball.Body.Y <= player2.Body.Max.Y {
			ball.Speed.X = -ball.Speed.X
			relativePos := (ball.Body.Y - player2.Body.Min.Y) / (player2.Body.Max.Y - player2.Body.Min.Y)
			ball.Speed.Y = (relativePos - 0.5) * 800
		}

		if ball.Body.Y < 0 || ball.Body.Y > height {
			ball.Speed.Y = -ball.Speed.Y
		}

		time.Sleep(time.Millisecond * 10)
	}
}

func (s *MainScene) updateScore() {
	for !gameOver {
		if ball.Body.X < 0 {
			player2.UpdateScore(player2.Score + 1)
			s.resetGame()
		}

		if ball.Body.X > width {
			player1.UpdateScore(player1.Score + 1)
			s.resetGame()
		}

		gameTimeElapsed = time.Since(gameStartTime)

		time.Sleep(time.Millisecond * 10)
	}
}

func (s *MainScene) resetGame() {
	ball.Update(pixel.V(width/2, height/2))
}

func (s *MainScene) Draw() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pong",
		Bounds: pixel.R(0, 0, width, height),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		player1.Move(win, gameOver, height, 1)
	}()
	go func() {
		defer wg.Done()
		player2.Move(win, gameOver, height, 2)
	}()
	go s.updateBall()
	go s.updateScore()

	for !win.Closed() {
		select {
		case player1Speed := <-player1.MoveCh:
			dt := time.Since(last).Seconds()
			last = time.Now()
			player1.Body.Min.Y += player1Speed * dt
			player1.Body.Max.Y += player1Speed * dt

		case player2Speed := <-player2.MoveCh:
			dt := time.Since(last).Seconds()
			last = time.Now()
			player2.Body.Min.Y += player2Speed * dt
			player2.Body.Max.Y += player2Speed * dt

		case <-gameOverCh:
			gameOver = true
		default:

		}

		win.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})

		imd := imdraw.New(nil)
		imd.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

		imd.Push(player1.Body.Min, player1.Body.Max)
		imd.Rectangle(0)

		imd.Push(player2.Body.Min, player2.Body.Max)
		imd.Rectangle(0)

		imd.Push(ball.Body)
		imd.Circle(10, 0)

		timeElapsed := gameTimeElapsed
		timeRemaining := gameDuration - timeElapsed
		if timeRemaining <= 0 {
			timeRemaining = 0
			gameOver = true
		}
		timeRemainingSeconds := int(timeRemaining.Seconds())

		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		timeText := text.New(pixel.V(10, height-20), basicAtlas)
		fmt.Fprintf(timeText, "Tiempo restante: %02d:%02d", timeRemainingSeconds/60, timeRemainingSeconds%60)
		timeText.Draw(win, pixel.IM.Scaled(timeText.Orig, 2))

		scoreText1 := text.New(pixel.V(10, height-60), basicAtlas)
		fmt.Fprintf(scoreText1, "Jugador 1: %d", player1.Score)
		scoreText1.Draw(win, pixel.IM.Scaled(scoreText1.Orig, 2))

		scoreText2 := text.New(pixel.V(10, height-80), basicAtlas)
		fmt.Fprintf(scoreText2, "Jugador 2: %d", player2.Score)
		scoreText2.Draw(win, pixel.IM.Scaled(scoreText2.Orig, 2))

		imd.Draw(win)

		if gameOver {
			finalScoreText := text.New(pixel.V(width/2-120, height/2), basicAtlas)
			fmt.Fprintf(finalScoreText, "Marcador final\n Jugador 1: %d\n Jugador 2: %d", player1.Score, player2.Score)
			finalScoreText.Draw(win, pixel.IM.Scaled(finalScoreText.Orig, 3))
		}

		win.Update()

		if gameTimeElapsed >= gameDuration {
			gameOver = true
		}

		if win.Closed() {
			break
		}
	}

	close(player1.MoveCh)
	close(player2.MoveCh)
	wg.Wait()
}
