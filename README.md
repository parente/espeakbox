# espeakbox

[progrium/busybox](https://github.com/progrium/busybox) + [espeak](http://espeak.sourceforge.net/) + [lame](http://lame.sourceforge.net/) / [opus](http://www.opus-codec.org/) + a golang http server

`=`

a text-to-speech server in a 16.43 MB Docker image.

## Status

Late night hacking complete. It works in the easy cases:

```
http://192.168.59.103:8080/speech?text=go%20to%20sleep%20pete
```

Stay tuned for code cleanup, a build on Docker Hub, and a real, documented API.

## Usage

To run a container from the latest image on Docker Hub:

```
docker run --name espeakbox -d -p 8080:8080 parente/espeakbox
```

To build it yourself, first install Go 1.4 on your platform with cross-compilation support for linux/amd64. For example, on a Mac with homebrew, run:

``` 
brew install go --cross-compile-common
```

Clone this repository and then run:

```
make build
```

## FAQ

*Why is this not an automated build on Docker Hub?*

I didn't spend time trying to get the Go toolchain into the image which is the only place the automated build can execute commands.

## License

MIT