## Notes
WEBRTC - Forked from https://github.com/poi5305/go-yuv2webRTC
GoCV - OpenCV for go https://gocv.io/computer-vision/

## Build and run
```
go mod download
go mod vendor
source setup.sh
go build main.go
./main

Open http://localhost:8000
And press `Start Session`
```

## Docker
```
docker build -t go-opencv-webrtc
docker run -it --rm -p 8000:8000 -p 52000-52100:52000-52100/udp go-opencv-webrtc
```

