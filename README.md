# Repo moved

The new repository location is [this one](https://github.com/gonutz/gophette).

After the contest was over, I created this fork and continued working on it. There is now a Windows version that does not need any external library to create a window and it uses DirectX so no extra DLLs have to be installed on the target Windows machine since all necessary DLLs are packaged since Windows XP. This makes it easy to build on Windows as well, you do not have to install SDL2, you only need Go and MinGW installed, the DirectX wrapper libraries will then build without any dependencies.

Besides the Windows version, bugs were fixed, so your opponent will now run correctly even after the level is reset.

# Gophette's Great Adventure

Gophette, the only female Gopher in the Gocave, was out for a walk in the forest when she heard a strange voice...

Evil Doctor Barney Starsoup is sitting in his cabin, looking at the programming language news groups and he has found out about the nice little language that Gophette so admires.

![Barney Starsoup](https://raw.githubusercontent.com/gophergala2016/gophette/master/screenshots/barney_starsoup.png)
![Gophette](https://raw.githubusercontent.com/gophergala2016/gophette/master/screenshots/gophette.png)

Doctor Starsoup has a reputation of adding terrible features to perfectly fine languages and hence he seeks to find the secret Gocave and make it his own.

Can you beat Evil Doctor Barney Starsoup in a race to your home and warn the other Gophers about the threat before it is too late?

![Race](https://raw.githubusercontent.com/gophergala2016/gophette/master/screenshots/race.png)

Here is a vide of the gameplay:
![Gameplay](https://github.com/gophergala2016/gophette/raw/master/screenshots/gameplay.flv)

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

# About

I created this as a solo project, meaning this is all programmer art (graphics and sound). I have created small games in the past, first in C++ and now in Go.

I hope people enjoy this game and realize that Go is very capable of creating desktops apps.