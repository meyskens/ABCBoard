ABCBoard
========

ABCBoard is a soundboar written in Go and React for use in theater productions at ABCTheater. It uses React to render the frontend but sound gets dispatched to the Go application via the webview controller. This last is done to ensure the audio always starts and starts reliably with nearly no delay. This is a must in theater productions where the HTML5 Audio API is not predictable enough.

## Building
I only tested this on Linux but it should work on Mac and Windows too. For build instructions I suggest looking at [webview](https://github.com/zserge/webview) and the `build.sh` script.