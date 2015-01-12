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

## API

Request:

```
GET /speech?text=<utterance>
            [&pitch=<0,99; default 50>]
            [&rate=<80,450; default 175 wpm>]
            [&voice=<name; default en>]
            [&encoding=<mp3|opus; default mp3>]
```

Response:

* `audio/mpeg` (mp3) or `audio/ogg` (opus) on success
* `text/plain` with a status code >= 400 on error

Request:

```
GET /voices
```

Response:

* `application/json` encoding an object with a `names` list
* `text/plain` with a status code >= 400 on error

## FAQ

*Why is this not an automated build on Docker Hub?*

I didn't spend time trying to get the Go toolchain into the image which is the only place the automated build can execute commands.

## License

MIT