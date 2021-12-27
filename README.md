# X2

An encrypted photo album for home network use.

The initial version of this tool is not yet encrypted. I'll add that as soon as I have the basics working for me.

I wanted a mobile photo album that was separate from my phone's own photo library and that I could control completely. I wanted it to be portable to any mobile phone or tablet in my household so I created a web based app that could be run on any machine on the local network.

The server is a single binary coded in Go.

The images are stored as base64 encoded records inside of an sqlite database.

The font-end is regular html with a bit of htmx for the javascript bits.

To compile the program simply type `go build`.

The server should run on any Mac, Linux, or Windows machine. To run it execute `./x2` after compiling.

Before you can upload files you need to create the blank image DB. To do that, start the server and point your browser to the `create` endpoint at `http://localhost:8002/create`.

Now that you have an image DB, you can upload images on the apps main page. Point your browser to `http://localhost:8002` to get started capturing images.

You can view the entire albim by pointing your browser to `http://localhost:8002/album`.

You can view images by pointing your browser to `http://localhost:8002/image/{id}` where `{id}` is replaced by the image id that you uploaded.

Although this works fine on desktop it was primarily designed to work on mobile phones using camera uploads.

This is called X2 because it's my second unnamed experiment. Read more about why I create unnamed projects on my blog.

## Todo

- Allow a key to be passed somehow
- Encrypt the image before storage
- Decrypt the image before display
- Redirect after image upload
