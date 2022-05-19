# `love-build`: Fast, Easy CLI build tool for [LÖVE](https://love2d.org)

**This is still *very* early in dev and is missing a lot of features. It's also a personal tool I use to build my own projects as-needed, so I can't guarantee it's being actively developed at any given time.**

**If you want something a little more robust/established make sure you check out the [Alternatives](#alternatives) section**

## Quick Start
### How to...
Create a `.love` file in the current directory for project located in `superGame/`
```
$ love-build superGame
```
Create a `.love` file with a specific version number
```
$ love-build --version=0.1.3 superGame 
```
Generate a build for Windows
```
$ love-build -w superGame
```
Generate a build for the Web using [love.js](https://github.com/Davidobot/love.js)
```
$ love-build -b superGame
```
Generate a build for Windows, put it in a separate build directory, give the build a different name, specify the version number, AND delete the `.love` file when done.
```
$ love-build -w -d superGame-builds -o ReallyCoolGame --version=0.1.3 --clean superGame
```
Create a `.love` file and then immediately delete it for some reason?
```
$ love-build --clean superGame
```

## Build Target support
- Windows (64-bit)
- Web (via [love.js](https://github.com/Davidobot/love.js))
    - Runs in compatibility mode by defaul

## LÖVE Version Support
- 11.3
- Probably others too!


## Alternatives
- [`love-release`](https://github.com/MisterDA/love-release) is more robust than this (and supports more platforms), but doesn't handle Web builds
- [`makelove`](https://github.com/pfirsich/makelove)

