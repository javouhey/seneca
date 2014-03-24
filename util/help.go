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
  -video-infile=<path>  Path (relative/full) to your mp4/flv/mov etc.. video 
  -from=00:00:00        Time of starting frame in hh:mm:ss (Default: 00:00:00)

Codec Options:
  -scale width:height   Scale dimensions of input video (Optional)
                        constraint: width & height must be even integers.
                        e.g. 300:_  height calculated to maintain aspect ratio.
                             _:250  width calculated to maintain aspect ratio.

Progress Reporting Options:
  -port=8080            TCP port for progress bar. (Default: 8080)


AUTHOR:
  Written by Gavin Bong 
  Report bugs to https://github.com/javouhey/seneca

COPYRIGHT:
  Licensed under the Apache License, Version 2.0


`
