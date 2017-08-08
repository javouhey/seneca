## Seneca 

[![GoDoc](https://godoc.org/github.com/javouhey/seneca?status.png)](https://godoc.org/github.com/javouhey/seneca)

Creates animated GIFs from videos.

<img src="logo.png" width="289" height="309" alt="seneca animated gif logo"/>

## Dependencies

* [Go](http://golang.org/) >= 1.2
* [ffmpeg](http://www.ffmpeg.org/) 2.1.4

## Usage

```bash
$ seneca -h

Usage:
  seneca -video-infile <path>
  seneca -h
  seneca -version

Options:
  -h                    Show this screen.
  -version              Show version.
  -dry-run              Show what would be done without
                        real invocations.
  -vv                   More verbose output
  -video-infile=<path>  Path (relative/full) to your mp4/flv/mov etc..
  -from=00:00:00        Starting frame offset in hh:mm:ss
                        (Default: 00:00:00)
  -length=<duration>    Duration to capture (Default: 3s) 
                        E.g. 2m35s, 1h2m15s

Codec Options:
  -scale width:height   Scale dimensions of input video (Optional)
                        constraint: width & height must be even integers
                        e.g. 300:_  calc height to maintain aspect ratio
                             _:250  calc width to maintain aspect ratio.

  -fps=<value>          frames per second. (Default: 25)
                        Range [1, 30]

Progress Reporting Options:
  -port=8080            TCP port for progress bar. (Default: 8080)

Animated GIF Options:
  -speed=<value>        Slow down / speed up animation(Default: placebo)
                        e.g veryfast, faster, placebo, slower, veryslow

  -repeat=<count>   **  Number of times to loop. (Default: loop forever)
  -delay=<seconds>  **  Seconds to pause before repeating animation
  -optimize         **  Attempts to reduce size of generated GIF.
  -upload           **  Uploads to imgur.com

Exit status:
  0  if OK,
  1  if invalid cli arguments (e.g. unable to read supplied video file),
126  if execution of ffmpeg failed,
127  if ffmpeg & ffprobe are not found on $PATH.


DEVELOPMENT STATUS:
  Options tagged with ** are not implemented yet.
```

## Installation

### Ubuntu

On my Ubuntu development machine:

1. Install `Go`
2. Install `ffmpeg` 
3. Compile `seneca`

```bash
$ git clone --recursive https://github.com/javouhey/seneca.git
$ cd seneca
$ ./build
```

You can now run the binary. It is advisable to add it to your $PATH.

```bash
$ ./bin/youtube-dl --version
2014.02.06

$ bin/youtube-dl --restrict-filenames -f mp4 -v 
       https://www.youtube.com/watch?v=4_pS46YRMIQ

$ ./bin/seneca -video-infile=Die_Antwoord_-_Zef_Side.-4_pS46YRMIQ.mp4 
      -from=00:01:34 -fps=3

 5 9 Completed
 Completed
 Completed

Your animated GIF is ready at location:
  /tmp/seneca/1502155863/Die_Antwoord_-_Zef_Side.gif
```

### Windows 7/8+

* Static binaries for `ffmpeg` are available from the [download page](http://www.ffmpeg.org/download.html)
* There is a `build.bat` for windows OS. The batch script will run the `mklink` program to create a Directory symbolic link and you need to run that with Administrator's privilege.

```bash
> build
```

## Sample

```bash
$ seneca -video-infile=./goproplane.mp4 -scale 280:_
         -fps 18 -from 00:00:39 -length 9s -speed=slower
```
![animated gif](http://i.imgur.com/4VdXgx3.gif)

## License

* Code is released under Apache license. See [LICENSE][license] file.
* The license for code under the `vendor` subdirectory remains under the purview of their respective owners.
* The [logo](http://commons.wikimedia.org/wiki/File:Nuremberg_chronicles_f_105r_1.png) above is from the public domain.


[license]: https://github.com/javouhey/seneca/blob/master/LICENSE
