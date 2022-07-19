package system

import (
	"os"
	"path"
	"runtime"
	"runtime/pprof"

	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type profileSystem struct {
	Weapon *component.Weapon

	player     gohan.Entity
	cpuProfile *os.File
}

func NewProfileSystem(player gohan.Entity) *profileSystem {
	return &profileSystem{
		player: player,
	}
}

func (s *profileSystem) Update(_ gohan.Entity) error {
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if s.cpuProfile == nil {
			runtime.SetCPUProfileRate(1000)

			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			s.cpuProfile, err = os.Create(path.Join(homeDir, "bbbb.prof"))
			if err != nil {
				return err
			}

			err = pprof.StartCPUProfile(s.cpuProfile)
			if err != nil {
				return err
			}
		} else {
			pprof.StopCPUProfile()

			s.cpuProfile.Close()
			s.cpuProfile = nil

		}
	}
	return nil
}

func (s *profileSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrUnregister
}
