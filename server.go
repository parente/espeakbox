package main

import (
    "io"
    "os/exec"
    "net/http"
    "log"
    // "encoding/json"
    // "fmt"
)

// http://stackoverflow.com/questions/19292113/not-buffered-http-responsewritter-in-golang
// http://play.golang.org/p/PpbPyXbtEs

type flushWriter struct {
    f http.Flusher
    w io.Writer
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
    n, err = fw.w.Write(p)
    if fw.f != nil {
        fw.f.Flush()
    }
    return
}

func speechHandler(w http.ResponseWriter, r *http.Request) {
    fw := flushWriter{w: w}
    if f, ok := w.(http.Flusher); ok {
        fw.f = f
    }

    values := r.URL.Query()
    text := values.Get("text")

    if len(text) == 0 {
        http.Error(w, "Missing required parameter: text", 400)
        return
    }

    encoding := values.Get("encoding")

    var encode *exec.Cmd
    speak := exec.Command("espeak", "--stdout", text)
    if(len(encoding) == 0 || encoding == "mp3") {
        encode = exec.Command("lame", "-", "-")
        w.Header().Set("Content-Type", "audio/mpeg")
    } else if(encoding == "opus") {
        encode = exec.Command("opusenc", "-", "-")
        w.Header().Set("Content-Type", "audio/ogg")
    } else {
        http.Error(w, "Unknown encoding requested: " + encoding, 400)
        return
    }

    encode.Stdin, _ = speak.StdoutPipe()
    encode.Stdout = &fw

    if err := encode.Start(); err != nil {
        http.Error(w, "Failed to start encoder", 500)
        return
    }
    
    if err := speak.Run(); err != nil {
        http.Error(w, "Failed to run synthesizer", 500)
        return
    }

    if err := encode.Wait(); err != nil {
        http.Error(w, "Failed to finish encoding", 500)
        return
    }
}

func main() {
    http.HandleFunc("/speech", speechHandler)

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Panic("Failed to start server: " + err.Error())
    }
}