# xkcd

Download and print the latest [xkcd][1] comic. `xkcd` only works in iterm2 or
kitty (with the `icat` kitten).

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

This program inlines the comics using either:

- [this library][2] for iterm2 (Your terminal emulator must support the iterm2
  protocol for this program to work).
- [the icat kitten for kitty][3].

## License

This program is licensed under [MIT](LICENSE).

[1]: https://xkcd.com/
[2]: https://github.com/francescomari/iterm2
[3]: https://sw.kovidgoyal.net/kitty/kittens/icat/
