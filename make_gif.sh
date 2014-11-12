ffmpeg -i in.mov -s 600x400 -pix_fmt rgb8 -f gif - | gifsicle --delay 2 --optimize=3 > out.gif
