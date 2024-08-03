# tinyface
a tiny go- face based face recongnition lib

## Requirements
tinyface require go-face to compile. go-face need to have [dlib](http://dlib.net/) (>= 19.10) and libjpeg development packages installed.

Ubuntu 18.10+, Debian sid
Latest versions of Ubuntu and Debian provide suitable dlib package so just run:
### Arch
``` bash
sudo pacman -S dlib blas atlas lapack libjpeg-turbo
```
### Ubuntu
``` bash
sudo apt-get install libdlib-dev libblas-dev libatlas-base-dev liblapack-dev libjpeg-turbo8-dev
```
### Debian
```bash
sudo apt-get install libdlib-dev libblas-dev libatlas-base-dev liblapack-dev libjpeg62-turbo-dev
```
### macOS
> Make sure you have [Homebrew](https://brew.sh/) installed.
```bash
brew install dlib
```
### Windows
 > Make sure you have [MSYS2](https://www.msys2.org/) installed.
1. Run MSYS2 MSYS shell from Start menu
2. Run `pacman -Syu` and if it asks you to close the shell do that
3. Run `pacman -Syu` again
4. Run `pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-dlib`
5. If you already have Go and Git installed and available in PATH uncomment `set MSYS2_PATH_TYPE=inherit` line in msys2_shell.cmd located in MSYS2 installation folder,Otherwise run `pacman -S mingw-w64-x86_64-go git`
6. Run MSYS2 MinGW 64-bit shell from Start menu to compile and use go-face
### Other systems
Try to install dlib/libjpeg with package manager of your distribution or [compile from sources](http://dlib.net/compile.html). Note that go-face won't work with old packages of dlib such as libdlib18. Alternatively create issue with the name of your system and someone might help you with the installation process.


## reference
+ https://github.com/Kagami/go-face
+ https://github.com/leandroveronezi/go-recognizer
+ https://hackernoon.com/face-recognition-with-go-676a555b8a7e

## Models
Currently shape_predictor_5_face_landmarks.dat, mmod_human_face_detector.dat and dlib_face_recognition_resnet_model_v1.dat are required. You may download them from [dlib-models](https://github.com/davisking/dlib-models) repo:
```bash
mkdir models && cd models
wget https://github.com/davisking/dlib-models/raw/master/shape_predictor_5_face_landmarks.dat.bz2
bunzip2 shape_predictor_5_face_landmarks.dat.bz2
wget https://github.com/davisking/dlib-models/raw/master/dlib_face_recognition_resnet_model_v1.dat.bz2
bunzip2 dlib_face_recognition_resnet_model_v1.dat.bz2
wget https://github.com/davisking/dlib-models/raw/master/mmod_human_face_detector.dat.bz2
bunzip2 mmod_human_face_detector.dat.bz2
```
you can also use [go-face 's testdata](https://github.com/Kagami/go-face-testdata) repo 
