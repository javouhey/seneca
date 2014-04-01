package util

var HelpMessage = `
seneca animated gif generator
      ___          ___          ___          ___          ___          ___     
     /\  \        /\  \        /\__\        /\  \        /\  \        /\  \    
    /::\  \      /::\  \      /::|  |      /::\  \      /::\  \      /::\  \   
   /:/\ \  \    /:/\:\  \    /:|:|  |     /:/\:\  \    /:/\:\  \    /:/\:\  \  
  _\:\~\ \  \  /::\~\:\  \  /:/|:|  |__  /::\~\:\  \  /:/  \:\  \  /::\~\:\  \ 
 /\ \:\ \ \__\/:/\:\ \:\__\/:/ |:| /\__\/:/\:\ \:\__\/:/__/ \:\__\/:/\:\ \:\__\
 \:\ \:\ \/__/\:\~\:\ \/__/\/__|:|/:/  /\:\~\:\ \/__/\:\  \  \/__/\/__\:\/:/  /
  \:\ \:\__\   \:\ \:\__\      |:/:/  /  \:\ \:\__\   \:\  \           \::/  / 
   \:\/:/  /    \:\ \/__/      |::/  /    \:\ \/__/    \:\  \          /:/  /  
    \::/  /      \:\__\        /:/  /      \:\__\       \:\__\        /:/  /   
     \/__/        \/__/        \/__/        \/__/        \/__/        \/__/    
Usage:
  seneca -video-infile <path>
  seneca -h
  seneca -version

Options:
  -h                    Show this screen.
  -version              Show version.
  -dry-run              Show what would be done without real invocations.
  -vv                   More verbose output
  -video-infile=<path>  Path (relative/full) to your mp4/flv/mov etc.. video 
  -from=00:00:00        Starting frame offset in hh:mm:ss (Default: 00:00:00)
  -length=<duration>    Duration to capture (Default: 3s) 
                        E.g. 2m35s, 1h2m15s

Codec Options:
  -scale width:height   Scale dimensions of input video (Optional)
                        constraint: width & height must be even integers.
                        e.g. 300:_  height calculated to maintain aspect ratio.
                             _:250  width calculated to maintain aspect ratio.

  -fps=<value>          frames per second. (Default: 25) 
                        Range [1, 30]

Progress Reporting Options:
  -port=8080            TCP port for progress bar. (Default: 8080)

Animated GIF Options:
  -speed=<value>        Slow down or speed up animation. (Default: placebo)
                        e.g. veryfast, faster, placebo, slower, veryslow

  -repeat=<count>   **  Number of times to loop. (Default: loop forever)
  -delay=<seconds>  **  Seconds to pause before repeating animation
  -optimize         **  Attempts to reduce size of generated GIF
  -upload           **  Uploads to imgur.com

Exit status:
  0  if OK,
  1  if invalid cli arguments (e.g. unable to read supplied video file),
126  if execution of ffmpeg failed,
127  if ffmpeg & ffprobe are not found on $PATH.


DEVELOPMENT STATUS:
  Options tagged with ** are not implemented yet.

AUTHOR:
  Written by Gavin Bong
  Report bugs to https://github.com/javouhey/seneca

COPYRIGHT:
  Licensed under the Apache License, Version 2.0


`

var ShortHelp = `See 'seneca -h' for help`
