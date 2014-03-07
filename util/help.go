package util

var Usage = `
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
  -port=8080            TCP port for progress bar. (Default: 8080)
  -from=00:00:00        Time of starting frame in hh:mm:ss (Default: 00:00:00)
`
