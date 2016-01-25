# Gophette's Great Adventure

Gophette, the only female gopher in the village, was out for a walk in the forest when she heard a strange voice...

Evil Doctor Barney Starsoup is sitting in his cabin, looking at the programming language news groups and he has found out about the nice little language that Gophette so admires.

Doctor Starsoup has a reputation of adding terrible features to perfectly fine languages and hence he seeks to find the secret Gopher Cave and make it his own.

Can you beat Evil Doctor Barney Starsoup in a race to your home and warn the other Gophers about the threat before it is too late?

# Build

Make sure you have the [SDL2 Go wrapper](https://github.com/veandco/go-sdl2) installed.

To install under OS X you can do this:

    brew install sdl2
    brew install sdl2_ttf
    brew install --with-libvorbis sdl2_mixer
    brew install sdl2_image
    go get -v github.com/veandco/go-sdl2/sdl
    go get -v github.com/veandco/go-sdl2/sdl_image
    go get -v github.com/veandco/go-sdl2/sdl_mixer
    go get -v github.com/veandco/go-sdl2/sdl_ttf

You need to have a C compiler installed for that so if you are on Windows you need [MinGW](http://sourceforge.net/projects/mingw/files/) installed as well.

After installing these you should be able to go get it:

    go get github.com/gophergala2016/gophette

and be ready to play.
