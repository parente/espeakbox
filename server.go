package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
)

// flushWriter holds a flusher and writer.
// http://stackoverflow.com/questions/19292113/not-buffered-http-responsewritter-in-golang
// http://play.golang.org/p/PpbPyXbtEs
type flushWriter struct {
	f http.Flusher
	w io.Writer
}

// Writes and immediately flushes bytes.
func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}

// Builds the command line for espeak. Check the parameters to ensure they are
// valid. Ensure the text parameter exists.
func buildSpeechCmd(values *url.Values, w *http.ResponseWriter) (cmd *exec.Cmd, err error) {
	args := []string{"--stdout"}

	// pitch, 0 to 99, default 50
	pitch := values.Get("pitch")
	if len(pitch) > 0 {
		pitchInt, err := strconv.ParseUint(pitch, 10, 8)
		if err != nil {
			err := errors.New("Invalid value for pitch: not an uint")
			http.Error(*w, err.Error(), 400)
			return nil, err
		} else if pitchInt < 0 || pitchInt > 99 {
			err := errors.New("Invalid value for pitch: range is [0,99]")
			http.Error(*w, err.Error(), 400)
			return nil, err
		}
		args = append(args, "-p")
		args = append(args, pitch)
	}

	// speech, 80 to 450, default 175
	speed := values.Get("speed")
	if len(speed) > 0 {
		speedInt, err := strconv.ParseUint(speed, 10, 0)
		if err != nil {
			err := errors.New("Invalid value for speed: not an uint")
			http.Error(*w, err.Error(), 400)
			return nil, err
		} else if speedInt < 80 || speedInt > 450 {
			err := errors.New("Invalid value for speed: range is [80,450]")
			http.Error(*w, err.Error(), 400)
			return nil, err
		}
		args = append(args, "-s")
		args = append(args, speed)
	}

	// voice
	voice := values.Get("voice")
	if len(voice) > 0 {
		args = append(args, "-v")
		args = append(args, voice)
	}

	// text is the only required parameter
	text := values.Get("text")
	if len(text) == 0 {
		err := errors.New("Missing required parameter: text")
		http.Error(*w, err.Error(), 400)
		return nil, err
	}
	args = append(args, text)

	return exec.Command("espeak", args...), nil
}

// Builds the command line for the encoder. Ensure the encoding is one supported,
// currently mp3 or opus.
func buildEncodeCmd(values *url.Values, w *http.ResponseWriter) (cmd *exec.Cmd, err error) {
	encoding := values.Get("encoding")

	// default to mp3 encoding, allow opus as well, error on all others
	// set the content-type header appropriately
	var encode *exec.Cmd
	if len(encoding) == 0 || encoding == "mp3" {
		encode = exec.Command("lame", "-", "-")
		(*w).Header().Set("Content-Type", "audio/mpeg")
	} else if encoding == "opus" {
		encode = exec.Command("opusenc", "-", "-")
		// use audio/ogg since it seems better supported by all browsers
		// than audio/opus
		(*w).Header().Set("Content-Type", "audio/ogg")
	} else if encoding == "wav" {
		(*w).Header().Set("Content-Type", "audio/wav")
		encode = exec.Command("cat")
	} else {
		err := errors.New("Unknown encoding requested: " + encoding)
		http.Error(*w, err.Error(), 400)
		return nil, err
	}

	return encode, nil
}

// Handles a request for synthesized speech. Build the synthesizer and encoder
// comands. Pipe the first to the second. Pipe the encoded stream out as the
// response. Supports URL arguments:
// text: string to synthesize (required)
// pitch: [0, 99] default: 50
// speed: [80, 450] default: 175
// voice: voice name to use
func speechHandler(w http.ResponseWriter, r *http.Request) {
	fw := flushWriter{w: w}
	if f, ok := w.(http.Flusher); ok {
		fw.f = f
	}

	// build speech and encoding commands
	values := r.URL.Query()
	speak, err := buildSpeechCmd(&values, &w)
	if err != nil {
		return
	}

	encode, err := buildEncodeCmd(&values, &w)
	if err != nil {
		return
	}

	// pipe synthesizer to encoder
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
		http.Error(w, "Failed to finish encoding: "+err.Error(), 500)
		return
	}
}

// Voices holds the list of voice names.
type Voices struct {
	Names []string `json:"names"`
}

// Caches the response to a request for the list of voices which never changes.
var (
	cachedVoicesJSON []byte
)

// Handles a request to get the list of synthesized voices. Uses awk to parse
// them out of the espeak stdout. Caches the immutable response for future
// requests.
func voicesHandler(w http.ResponseWriter, r *http.Request) {
	if cachedVoicesJSON == nil {
		speak := exec.Command("espeak", "--voices")
		awk := exec.Command("awk", "{print $5}")
		awk.Stdin, _ = speak.StdoutPipe()
		var voicesBuf bytes.Buffer
		awk.Stdout = &voicesBuf

		if err := awk.Start(); err != nil {
			http.Error(w, "Failed to parse voices: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := speak.Run(); err != nil {
			http.Error(w, "Failed to fetch voices: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := awk.Wait(); err != nil {
			http.Error(w, "Failed to finish parsing voices: "+err.Error(), http.StatusInternalServerError)
			return
		}

		voices := new(Voices)
		// leave out the header value
		voices.Names = strings.Split(voicesBuf.String(), "\n")[1:]

		js, err := json.Marshal(voices)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cachedVoicesJSON = js
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(cachedVoicesJSON)
}

func main() {
	http.HandleFunc("/speech", speechHandler)
	http.HandleFunc("/voices", voicesHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic("Failed to start server: " + err.Error())
	}
}
