# xkcd

Download and print the latest [xkcd][1] comic. `xkcd` only works in iterm2 or
kitty.

## Install

```shell
go install github.com/francescomari/xkcd@latest
```

## Run

Show the latest comic:

```shell
xkcd
```

Show a random comic:

```shell
xkcd -random
```

Show a specific comic:

```shell
xkcd 927
```

## Compatibility

This program displays the comics by using [this library][2] for iterm2 or
iterm2-compatible terminals, and by implementing kitty's [terminal graphics
protocol][3] directly.

## License

This program is licensed under [MIT](LICENSE).

[1]: https://xkcd.com/
[2]: https://github.com/francescomari/iterm2
[3]: https://sw.kovidgoyal.net/kitty/graphics-protocol/
