# Image cropping

A demo of server-side image cropping with [Go](https://golang.org/) and [cropper.js](https://github.com/fengyuanchen/cropperjs).

## How it works?

In this demo, we're using [cropper.js](https://github.com/fengyuanchen/cropperjs) to crop the image in the browser and send the crop details to the server to process the final result.

### Why do the crop work on server?

[Cropper.js](https://github.com/fengyuanchen/cropperjs) uses the browser's native [canvas.toBlob API](https://github.com/fengyuanchen/cropperjs/issues/534) which means, it will produce a reduced quality image. To mend this, we do the crop work on the server.

## Running

1. Install web app dependencies.

```console
$ yarn // or npm install
```

2. Build the web app.

```console
$ yarn build // or npm run build
```

3. Run the server.

```console
$ go run main.go
```

4. Go to [localhost:3002](http://localhost:3002) and start playing!
