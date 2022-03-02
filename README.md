# xkcd

Download and print the latest [xkcd][1] comic.

## Install

```
go install github.com/francescomari/xkcd@latest
```

## Run

Show the latest comic:

```
xkcd
```

Show a random comic:

```
xkcd -random
```

Show a specific comic:

```
xkcd 927
```

## Compatibility

This program inlines the comics using [this library][2] for iterm2. Your
terminal emulator must support the iterm2 protocol for this program to work.

## License

This library is licensed under [MIT](LICENSE).

[1]: https://xkcd.com/
[2]: https://github.com/francescomari/iterm2
