# pigeon-hole

Move files from one folder to others for processing.

I have a program that dumps files into a folder, they can be split between input servers and processors but there still ends up being slowdowns due to one processor being overloaded from its inputs while another is free and open not doing anything. The type of files being sent to the processors vary in importance, some should be processed quickly and others can wait. The normal processing works by copying all the files over, then process, it would have an order so most important would be processed quickly however it would not copy files again until all original files had been processed. This created slowdowns in processing the more important files due to hangups on the slower less important files. Pigeon Hole strives to alleviate this by taking an input folder and constantly moving to a "hole" which then drops into the processors one by one. This way the more important files are kept in front and will get processed in a timely fashion even if there are many non important files to be processed.

## Config
Hard coded right now. Will be INI based.

## TODO
Make a class of file which will only be put in 1 processor at a time