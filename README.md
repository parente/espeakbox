# espeakbox

[gliderlabs/alpine](https://github.com/gliderlabs/docker-alpine) + [espeak](http://espeak.sourceforge.net/) + [lame](http://lame.sourceforge.net/) / [opus](http://www.opus-codec.org/) + a golang http server &rarr; a text-to-speech server in a ~17 MB Docker image

## Usage

To run a container from the latest image on Docker Hub:

```
docker run --name espeakbox -d -p 8080:8080 parente/espeakbox
```

## Build

Install Docker. Clone this repository and then run:

```
make build
```

## API

Request:

```
GET /speech?text=<utterance>
            [&pitch=<0,99; default 50>]
            [&speed=<80,450; default 175 wpm>]
            [&voice=<name; default en>]
            [&encoding=<mp3|opus; default mp3>]
```

Response:

* `audio/wav  (wave),  audio/mpeg` (mp3) or `audio/ogg` (opus) on success
* `text/plain` with a status code >= 400 on error

Request:

```
GET /voices
```

Response:

* `application/json` encoding an object with a `names` list
* `text/plain` with a status code >= 400 on error

## FAQ

*How well does it perform?*

My ten-minute effort into benchmarking a single container yields the following on a boot2docker virtual machine with 2 GB RAM and 8 VCPUs hosted on a Macbook with 16 GB RAM, a 2.7 GHz processor, and 4 cores.

```
$ ab -c 100 -n 1000 'http://192.168.59.103:8080/speech?text=the%20cheese%20flies%20at%20midnight'

This is ApacheBench, Version 2.3 <$Revision: 1554214 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 192.168.59.103 (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests


Server Software:
Server Hostname:        192.168.59.103
Server Port:            8080

Document Path:          /speech?text=the%20cheese%20flies%20at%20midnight
Document Length:        7105 bytes

Concurrency Level:      100
Time taken for tests:   13.102 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      7187000 bytes
HTML transferred:       7105000 bytes
Requests per second:    76.33 [#/sec] (mean)
Time per request:       1310.172 [ms] (mean)
Time per request:       13.102 [ms] (mean, across all concurrent requests)
Transfer rate:          535.70 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   16 112.6      3    1038
Processing:   174 1271 689.6   1089    5861
Waiting:      114 1012 697.3    812    5220
Total:        177 1287 698.4   1093    5864

Percentage of the requests served within a certain time (ms)
  50%   1093
  66%   1346
  75%   1555
  80%   1723
  90%   2266
  95%   2751
  98%   3163
  99%   3735
 100%   5864 (longest request)
```

The same test on a SoftLayer virtual instance with 16 GB RAM and 8 VCPUs yields the following:

```
$ ab -c 100 -n 1000 'http://127.0.0.1:8080/speech?text=the%20cheese%20flies%20at%20midnight'

This is ApacheBench, Version 2.3 <$Revision: 1528965 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests


Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /speech?text=the%20cheese%20flies%20at%20midnight
Document Length:        7105 bytes

Concurrency Level:      100
Time taken for tests:   11.329 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      7187000 bytes
HTML transferred:       7105000 bytes
Requests per second:    88.27 [#/sec] (mean)
Time per request:       1132.892 [ms] (mean)
Time per request:       11.329 [ms] (mean, across all concurrent requests)
Transfer rate:          619.53 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   1.5      0       7
Processing:   389 1110 276.9   1068    2494
Waiting:       72  830 226.3    833    1954
Total:        390 1111 276.8   1068    2495

Percentage of the requests served within a certain time (ms)
  50%   1068
  66%   1230
  75%   1301
  80%   1344
  90%   1465
  95%   1606
  98%   1721
  99%   1793
 100%   2495 (longest request)
```

## License

MIT
