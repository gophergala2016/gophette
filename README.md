# Gophette's Great Adventure

Gophette, the only female gopher in the village, is in trouble.

Evil Doctor Barney Starsoup has found out about the nice little village and seeks to compromise they idyll.

Beat the Evil Doctor and get back home earlier to warn the other gophers about the upcoming modernization of their home...

# Build

Make sure you have the [SDL2 Go wrapper](https://github.com/veandco/go-sdl2) installed.

To install under OS X you can do this:

    brew install sdl2
    brew install sdl2_ttf
    brew install --with-libvorbis sdl2_mixer
    brew install sdl2_image
    go get -v github.com/veandco/go-sdl2/...

You need to have a C compiler installed for that so if you are on Windows you need [MinGW](http://sourceforge.net/projects/mingw/files/) installed as well.

After installing these you should be able to go get it:

    go get github.com/gophergala2016/gophette

and be ready to play.