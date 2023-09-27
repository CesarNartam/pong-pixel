package main

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

const (
	width        = 800
	height       = 600
	gameDuration = time.Minute
)

var (
	score1          = 0
	score2          = 0
	player1         = pixel.Rect{Min: pixel.V(10, height/2-50), Max: pixel.V(20, height/2+50)}
	player2         = pixel.Rect{Min: pixel.V(width-20, height/2-50), Max: pixel.V(width-10, height/2+50)}
	ball            = pixel.V(width/2, height/2)
	ballSpeed       = pixel.V(500, 500)
	player1Speed    = 700.0
	player2Speed    = 700.0
	gameStartTime   = time.Now()
	last            = time.Now()
	gameOver        = false
	gameTimeElapsed time.Duration
)

var (
	player1MoveCh = make(chan float64)
	player2MoveCh = make(chan float64)
	gameOverCh    = make(chan struct{})
)

func movePlayer1(win *pixelgl.Window) {
	for !gameOver {
		if win.Pressed(pixelgl.KeyW) && player1.Max.Y < height {
			player1MoveCh <- player1Speed
		}
		if win.Pressed(pixelgl.KeyS) && player1.Min.Y > 0 {
			player1MoveCh <- -player1Speed
		}
	}
}

func movePlayer2(win *pixelgl.Window) {
	for !gameOver {
		if win.Pressed(pixelgl.KeyUp) && player2.Max.Y < height {
			player2MoveCh <- player2Speed
		}
		if win.Pressed(pixelgl.KeyDown) && player2.Min.Y > 0 {
			player2MoveCh <- -player2Speed
		}
	}
}

func updateBall() {
	for !gameOver {
		dt := time.Since(last).Seconds()
		last = time.Now()
		ball = ball.Add(ballSpeed.Scaled(dt))

		if ball.X < player1.Max.X && ball.Y >= player1.Min.Y && ball.Y <= player1.Max.Y {
			ballSpeed.X = -ballSpeed.X
			relativePos := (ball.Y - player1.Min.Y) / (player1.Max.Y - player1.Min.Y)
			ballSpeed.Y = (relativePos - 0.5) * 800
		}

		if ball.X > player2.Min.X && ball.Y >= player2.Min.Y && ball.Y <= player2.Max.Y {
			ballSpeed.X = -ballSpeed.X
			relativePos := (ball.Y - player2.Min.Y) / (player2.Max.Y - player2.Min.Y)
			ballSpeed.Y = (relativePos - 0.5) * 800
		}

		if ball.Y < 0 || ball.Y > height {
			ballSpeed.Y = -ballSpeed.Y
		}

		time.Sleep(time.Millisecond * 10)
	}
}

func updateScore() {
	for !gameOver {
		// Verificar si la pelota salió de la pantalla a la izquierda
		if ball.X < 0 {
			// Punto para el jugador 2
			score2++
			resetGame()
		}

		// Verificar si la pelota salió de la pantalla a la derecha
		if ball.X > width {
			// Punto para el jugador 1
			score1++
			resetGame()
		}

		// Actualizar el temporizador de tiempo transcurrido
		gameTimeElapsed = time.Since(gameStartTime)

		// Esperar un tiempo antes de la próxima actualización
		time.Sleep(time.Millisecond * 10)
	}
}

func resetGame() {
	// Restablecer la posición de la pelota y el temporizador
	ball = pixel.V(width/2, height/2)
	gameStartTime = time.Now()
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pong",
		Bounds: pixel.R(0, 0, width, height),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	// Iniciar goroutines para el movimiento de jugadores y actualización de la pelota.
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		movePlayer1(win)
	}()
	go func() {
		defer wg.Done()
		movePlayer2(win)
	}()
	go updateBall()
	go updateScore()

	for !win.Closed() {
		// Lógica de comunicación con las goroutines
		select {
		case player1Speed := <-player1MoveCh:
			// Actualiza la posición del jugador 1
			dt := time.Since(last).Seconds()
			last = time.Now()
			player1.Min.Y += player1Speed * dt
			player1.Max.Y += player1Speed * dt
		case player2Speed := <-player2MoveCh:
			// Actualiza la posición del jugador 2
			dt := time.Since(last).Seconds()
			last = time.Now()
			player2.Min.Y += player2Speed * dt
			player2.Max.Y += player2Speed * dt
		case <-gameOverCh:
			// El juego ha terminado
			gameOver = true
		default:
			// No se recibieron acciones, continuar
		}

		win.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})

		imd := imdraw.New(nil)
		imd.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

		imd.Push(player1.Min, player1.Max)
		imd.Rectangle(0)

		imd.Push(player2.Min, player2.Max)
		imd.Rectangle(0)

		imd.Push(ball)
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
		fmt.Fprintf(scoreText1, "Jugador 1: %d", score1)
		scoreText1.Draw(win, pixel.IM.Scaled(scoreText1.Orig, 2))

		scoreText2 := text.New(pixel.V(10, height-80), basicAtlas)
		fmt.Fprintf(scoreText2, "Jugador 2: %d", score2)
		scoreText2.Draw(win, pixel.IM.Scaled(scoreText2.Orig, 2))

		imd.Draw(win)

		if gameOver {
			finalScoreText := text.New(pixel.V(width/2-120, height/2), basicAtlas)
			fmt.Fprintf(finalScoreText, "Marcador final\n Jugador 1: %d\n Jugador 2: %d", score1, score2)
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

	// Cerrar goroutines antes de salir
	close(player1MoveCh)
	close(player2MoveCh)
	wg.Wait()
}

func main() {
	pixelgl.Run(run)
}
