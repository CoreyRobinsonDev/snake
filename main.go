package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 128
	screenHeight = 128
	playerSize   = 8
	sampleRate   = 48000
	volume       = .05
)

type Stream struct {
	pos  int64
	freq float64
}

func (self Stream) Read(buf []byte) (int, error) {
	const bytesPerSample = 8
	length := sampleRate / self.freq
	n := len(buf) / bytesPerSample * bytesPerSample

	for i := 0; i < n/bytesPerSample; i++ {
		v := math.Float32bits(float32(math.Sin(2 * math.Pi * float64(self.pos/bytesPerSample+int64(i)) / length)))
		buf[8*i] = byte(v)
		buf[8*i+1] = byte(v >> 8)
		buf[8*i+2] = byte(v >> 16)
		buf[8*i+3] = byte(v >> 24)
		buf[8*i+4] = byte(v)
		buf[8*i+5] = byte(v >> 8)
		buf[8*i+6] = byte(v >> 16)
		buf[8*i+7] = byte(v >> 24)
	}
	self.pos += int64(n)
	self.pos %= int64(length) * bytesPerSample

	return n, nil
}

func (self *Stream) Close() error {
	return nil
}

type PlayArea struct {
	Top, Bottom, Left, Right float32
}

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

func (self Direction) String() string {
	switch self {
	case North:
		return "North"
	case South:
		return "South"
	case East:
		return "East"
	case West:
		return "West"
	default:
		return ""
	}
}

type Player struct {
	Body []*BodyPart
}

type BodyPart struct {
	X, Y          float32
	Direction     Direction
	PrevDirection Direction
}

type Food struct {
	X, Y float32
}

type Game struct {
	AudioContext *audio.Context
	Player       *Player
	PlayArea     *PlayArea
	Food         Food
	Delay        int
	Score        int
	IsOver       bool
	HasFood      bool
	KeyPressed   bool
}

func (self *Game) Update() error {
	if self.IsOver {
		return nil
	}
	if !ebiten.IsKeyPressed(ebiten.KeyA) &&
		!ebiten.IsKeyPressed(ebiten.KeyS) &&
		!ebiten.IsKeyPressed(ebiten.KeyD) &&
		!ebiten.IsKeyPressed(ebiten.KeyW) {
		self.KeyPressed = false
	}
	if !self.HasFood {
		// C4
		audioPlayer, err := self.AudioContext.NewPlayer(&Stream{freq: 261.6256})
		if err != nil {
			log.Fatal(err)
		}
		audioPlayer.SetVolume(volume * 1.3)
		audioPlayer.Play()
		time.Sleep(time.Millisecond * 80)
		audioPlayer.Close()
		time.Sleep(time.Millisecond * 80)
		// C5
		audioPlayer2, err := self.AudioContext.NewPlayer(&Stream{freq: 523.2511})
		if err != nil {
			log.Fatal(err)
		}
		audioPlayer2.SetVolume(volume * 1.3)
		audioPlayer2.Play()
		time.Sleep(time.Millisecond * 80)
		audioPlayer2.Close()
		var cx, cy int
		for {
		loop:
			cx = rand.Intn(int(self.PlayArea.Right - self.PlayArea.Left))
			cy = rand.Intn(int(self.PlayArea.Bottom - self.PlayArea.Top))
			for _, part := range self.Player.Body {
				if part.X == float32(cx) && part.Y == float32(cy) {
					goto loop
				}
			}
			break
		}
		self.Food.X = float32(math.Floor(float64(float32(cx)+self.PlayArea.Left)/8) * 8)
		self.Food.Y = float32(math.Floor(float64(float32(cy)+self.PlayArea.Top)/8) * 8)
		self.HasFood = true
	}
	if math.Abs(float64(self.Food.X-self.Player.Body[0].X)) < playerSize &&
		math.Abs(float64(self.Food.Y-self.Player.Body[0].Y)) < playerSize {
		self.HasFood = false
		i := len(self.Player.Body) - 1
		x := self.Player.Body[i].X
		y := self.Player.Body[i].Y
		if self.Player.Body[i].Direction == North {
			y += playerSize
		}
		if self.Player.Body[i].Direction == East {
			x -= playerSize
		}
		if self.Player.Body[i].Direction == West {
			x += playerSize
		}
		if self.Player.Body[i].Direction == South {
			y -= playerSize
		}
		self.Player.Body = append(self.Player.Body, &BodyPart{
			X: x, Y: y,
			Direction: self.Player.Body[i].Direction,
		})
		self.Score++
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) &&
		(self.Player.Body[0].Direction != East && self.Player.Body[0].Direction != East) &&
		!self.KeyPressed {
		self.Player.Body[0].Direction = West
		self.KeyPressed = true
		// A4
		audioPlayer, err := self.AudioContext.NewPlayer(&Stream{freq: 440})
		if err != nil {
			log.Fatal(err)
		}
		audioPlayer.SetVolume(volume)
		audioPlayer.Play()
		time.Sleep(time.Millisecond * 80)
		audioPlayer.Close()
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) &&
		(self.Player.Body[0].Direction != South && self.Player.Body[0].Direction != North) &&
		!self.KeyPressed {
		self.Player.Body[0].Direction = South
		self.KeyPressed = true
		// G4
		audioPlayer, err := self.AudioContext.NewPlayer(&Stream{freq: 391.9954})
		if err != nil {
			log.Fatal(err)
		}
		audioPlayer.SetVolume(volume)
		audioPlayer.Play()
		time.Sleep(time.Millisecond * 80)
		audioPlayer.Close()
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) &&
		(self.Player.Body[0].Direction != East && self.Player.Body[0].Direction != West) &&
		!self.KeyPressed {
		self.Player.Body[0].Direction = East
		self.KeyPressed = true
		// F#4
		audioPlayer, err := self.AudioContext.NewPlayer(&Stream{freq: 369.9944})
		if err != nil {
			log.Fatal(err)
		}
		audioPlayer.SetVolume(volume)
		audioPlayer.Play()
		time.Sleep(time.Millisecond * 80)
		audioPlayer.Close()
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) &&
		(self.Player.Body[0].Direction != South && self.Player.Body[0].Direction != North) &&
		!self.KeyPressed {
		self.Player.Body[0].Direction = North
		self.KeyPressed = true
		// E4
		audioPlayer, err := self.AudioContext.NewPlayer(&Stream{freq: 329.6276})
		if err != nil {
			log.Fatal(err)
		}
		audioPlayer.SetVolume(volume)
		audioPlayer.Play()
		time.Sleep(time.Millisecond * 80)
		audioPlayer.Close()
	}

	self.Delay++
	if self.Delay%(playerSize*2) == 0 {
		for i := range self.Player.Body {
			self.Player.Body[i].PrevDirection = self.Player.Body[i].Direction
			if i > 0 {
				self.Player.Body[i].Direction = self.Player.Body[i-1].PrevDirection
				if self.Player.Body[0].X == self.Player.Body[i].X &&
					self.Player.Body[0].Y == self.Player.Body[i].Y {
					self.IsOver = true
					return nil
				}
			}

			if self.Player.Body[i].PrevDirection == North {
				self.Player.Body[i].Y -= playerSize
			}
			if self.Player.Body[i].PrevDirection == East {
				self.Player.Body[i].X += playerSize
			}
			if self.Player.Body[i].PrevDirection == West {
				self.Player.Body[i].X -= playerSize
			}
			if self.Player.Body[i].PrevDirection == South {
				self.Player.Body[i].Y += playerSize
			}

			if self.Player.Body[i].X < self.PlayArea.Left {
				self.Player.Body[i].X = self.PlayArea.Right - playerSize
			}
			if self.Player.Body[i].X+playerSize > self.PlayArea.Right {
				self.Player.Body[i].X = self.PlayArea.Left
			}
			if self.Player.Body[i].Y < self.PlayArea.Top {
				self.Player.Body[i].Y = self.PlayArea.Bottom - playerSize
			}
			if self.Player.Body[i].Y+playerSize > self.PlayArea.Bottom {
				self.Player.Body[i].Y = self.PlayArea.Top
			}
		}
	}

	return nil
}

func (self *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(0, 0)
	textOp.ColorScale.ScaleWithColor(color.Black)
	source, _ := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	text.Draw(screen, fmt.Sprintf("score: %d", self.Score), &text.GoTextFace{
		Source: source,
		Size:   8,
	}, textOp)

	vector.DrawFilledRect(
		screen,
		self.PlayArea.Left,
		self.PlayArea.Top,
		self.PlayArea.Right-8,
		self.PlayArea.Bottom-8,
		color.Black,
		false,
	)
	vector.DrawFilledCircle(screen, self.Food.X+4, self.Food.Y+4, 2, color.White, false)
	for _, part := range self.Player.Body {
		vector.DrawFilledRect(
			screen,
			part.X,
			part.Y,
			playerSize,
			playerSize,
			color.White,
			false,
		)
	}
	if self.IsOver {
		textO := &text.DrawOptions{}
		textO.GeoM.Translate(32, 48)
		textO.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, "game over", &text.GoTextFace{
			Source: source,
			Size:   12,
		}, textO)
	}
}

func (self *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowTitle("Game")
	ebiten.SetWindowSize(1920, 1080)

	game := &Game{
		AudioContext: audio.NewContext(sampleRate),
		Player: &Player{
			Body: []*BodyPart{&BodyPart{
				X: 64,
				Y: 64,
			}},
		},
		PlayArea: &PlayArea{
			Left:   8,
			Top:    8,
			Right:  screenWidth - 8,
			Bottom: screenHeight - 8,
		},
		HasFood: false,
		IsOver:  false,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
