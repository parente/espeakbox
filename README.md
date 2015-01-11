# espeakbox

[progrium/busybox](https://github.com/progrium/busybox) + [espeak](http://espeak.sourceforge.net/) + [lame](http://lame.sourceforge.net/) / [opus](http://www.opus-codec.org/) + a golang http server

=

a text-to-speech server in a 16.43 MB Docker image.

## Status

Late night hacking complete. It works in the easy cases:

```
http://192.168.59.103:8080/speech?text=go%20to%20sleep%20pete
```

Stay tuned for code cleanup, a build on Docker Hub, and a real, documented API.

## License

MIT