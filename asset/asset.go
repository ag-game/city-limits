package asset

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	_ "image/png"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const sampleRate = 44100

//go:embed image map sound
var FS embed.FS

var (
	ImgBlank       = ebiten.NewImage(1, 1)
	ImgWhiteSquare = ebiten.NewImage(64, 64)
	ImgBlackSquare = ebiten.NewImage(64, 64)
	ImgHelp        = LoadImage("image/help.png")
	ImgPower       = LoadImage("image/power.png")
)

var (
	SoundMusic1 *audio.Player
	SoundMusic2 *audio.Player
	SoundMusic3 *audio.Player

	SoundSelect   *audio.Player
	SoundBulldoze *audio.Player

	SoundPop1 *audio.Player
	SoundPop2 *audio.Player
	SoundPop3 *audio.Player
	SoundPop4 *audio.Player
	SoundPop5 *audio.Player

	SoundExplosion1 *audio.Player
	SoundExplosion2 *audio.Player
)

func init() {
	ImgWhiteSquare.Fill(color.White)
	ImgBlackSquare.Fill(color.Black)
}

func LoadSounds(ctx *audio.Context) {
	SoundMusic1 = LoadOGG(ctx, "sound/we_will_build_it.ogg", false)
	SoundMusic1.SetVolume(0.6)

	SoundMusic2 = LoadOGG(ctx, "sound/please_recycle.ogg", false)
	SoundMusic2.SetVolume(0.1)

	SoundMusic3 = LoadOGG(ctx, "sound/the_world_is_a_landfill_and_its_your_fault_you_fucking_son_of_a_bitch_god_damn.ogg", false)
	SoundMusic3.SetVolume(0.4)

	SoundSelect = LoadWAV(ctx, "sound/select/select.wav")
	SoundSelect.SetVolume(0.6)

	SoundBulldoze = LoadOGG(ctx, "sound/bulldozer/bulldozer.ogg", true)
	SoundBulldoze.SetVolume(0.6)

	const popVolume = 0.15
	SoundPop1 = LoadWAV(ctx, "sound/pop/pop1.wav")
	SoundPop2 = LoadWAV(ctx, "sound/pop/pop2.wav")
	SoundPop3 = LoadWAV(ctx, "sound/pop/pop3.wav")
	SoundPop4 = LoadWAV(ctx, "sound/pop/pop4.wav")
	SoundPop5 = LoadWAV(ctx, "sound/pop/pop5.wav")
	SoundPop1.SetVolume(popVolume)
	SoundPop2.SetVolume(popVolume)
	SoundPop3.SetVolume(popVolume)
	SoundPop4.SetVolume(popVolume)
	SoundPop5.SetVolume(popVolume)

	const explosionVolume = 0.1
	SoundExplosion1 = LoadOGG(ctx, "sound/explosion/explosion1.ogg", false)
	SoundExplosion2 = LoadOGG(ctx, "sound/explosion/explosion2.ogg", false)
	SoundExplosion1.SetVolume(explosionVolume)
	SoundExplosion2.SetVolume(explosionVolume)
}

func LoadImage(p string) *ebiten.Image {
	f, err := FS.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	baseImg, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(baseImg)
}

func LoadBytes(p string) []byte {
	b, err := FS.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return b
}

func LoadWAV(context *audio.Context, p string) *audio.Player {
	f, err := FS.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	stream, err := wav.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		panic(err)
	}

	player, err := context.NewPlayer(stream)
	if err != nil {
		panic(err)
	}

	// Workaround to prevent delays when playing for the first time.
	player.SetVolume(0)
	player.Play()
	player.Pause()
	player.Rewind()
	player.SetVolume(1)

	return player
}

func LoadOGG(context *audio.Context, p string, loop bool) *audio.Player {
	b := LoadBytes(p)

	stream, err := vorbis.DecodeWithSampleRate(sampleRate, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	var s io.Reader
	if loop {
		s = audio.NewInfiniteLoop(stream, stream.Length())
	} else {
		s = stream
	}

	player, err := context.NewPlayer(s)
	if err != nil {
		panic(err)
	}

	// Workaround to prevent delays when playing for the first time.
	player.SetVolume(0)
	player.Play()
	player.Pause()
	player.Rewind()
	player.SetVolume(1)

	return player
}
