# X2

An (eventually encrypted) photo album for home network use.

## Getting Started

The server should run on any Mac, Linux, or Windows machine. To run it execute `./x2` after compiling. Run just `x2` if you're on Windows.

To compile the program simply type `go build`.

Now Point your browser to `http://localhost:8002` to start capturing images.

## More Information

I'm no security expert and this software is experimental. I suggest you keep the server on an internal home network that is not accessible from the internet. The point of this tool is to keep your photos out of your phone's own photo album and to encrypt the photos when they're at rest.

The initial version of this tool is not yet encrypted. I'll add that as soon as I've tested the basics.

I wanted a mobile photo album that was separate from my phone's own photo library and that I could control myself. I wanted it to be portable to any mobile phone or tablet in my household so I created a web based app that could be run on any machine on the local network.

The server is a single binary coded in Go.

The images are stored as base64 encoded records inside of an sqlite database.

The font-end is regular html with a bit of htmx for the javascript parts. I have spent _no_ time making this tool pretty.

You can view the album by pointing your browser to `http://localhost:8002/album`, which is also linked from the home page.

You can view raw images by pointing your browser to `http://localhost:8002/image/{id}` where `{id}` is replaced by the image id that you uploaded.

Although this works fine on desktop it was primarily designed to work on mobile phones using camera uploads.

## Todo

- Allow a key to be passed via http basic auth
- Encrypt the image before storage
- Decrypt the image before display
- Redirect to create if the image db doesn't exist
- Maybe use an alternate "browse" button
