# City Limits
[![Donate](https://img.shields.io/liberapay/receives/rocketnine.space.svg?logo=liberapay)](https://liberapay.com/rocketnine.space)

[City-building simulation](https://en.wikipedia.org/wiki/City-building_game) video game

This game was created for the [Wasted Resources Game Jam](https://itch.io/jam/wastedresources).

## Play

### Browser

[**Play in your browser**](https://rocketnine.itch.io/citylimits?secret=citylimits)

### Compile

**Note:** You will need to install the dependencies listed for [your platform](https://github.com/hajimehoshi/ebiten/blob/main/README.md#platforms).

Run the following command to build a `citylimits` executable:

`go install code.rocketnine.space/tslocum/citylimits@latest`

Run `~/go/bin/citylimits` to play.

## Support

Please share issues and suggestions [here](https://code.rocketnine.space/tslocum/citylimits/issues).

## Credits

- [Trevor Slocum](https://rocketnine.space) - Game design and programming
- [node punk](https://soundcloud.com/solve_x) - Music

## Dependencies

- [ebiten](https://github.com/hajimehoshi/ebiten) - Game engine
- [go-tiled](https://github.com/lafriks/go-tiled) - Tiled map file (.TMX) parser
- [go-astar](https://github.com/beefsack/go-astar) - Pathfinding library
