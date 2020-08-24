# gocast
gocast is a program for downloading podcasts from acast.com via rss feed.

## Build
```
go get -u -v github.com/philippdrebes/gocast
cd github.com/philippdrebes/gocast
make build
```

## Install
```
make install
```

## Usage
```
Hello Gocast!
usage: gocast [-h|--help] -n|--name "<value>" [-o|--output "<value>"]
              [-l|--list] [-i|--index <integer>] [--all] [--latest]

              gocast is a program for downloading podcasts from acast.com via
              rss feed.

Arguments:

  -h  --help    Print help information
  -n  --name    Podcast name e.g. letstalkaboutcarsyo for
                http://rss.acast.com/letstalkaboutcarsyo
  -o  --output  Specifies where downloaded episodes will be saved to
  -l  --list    List all episodes
  -i  --index   Download a single episode via index. Run the 'list' command in
                order to get the index of the desired episode. Default: -1
      --all     Download all episodes
      --latest  Download the latest episode


```
