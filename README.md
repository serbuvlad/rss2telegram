rss2telegram
============

Motivation
----------

rss2telegram is a simple utility intended to turn telegram into an RSS(/atom/json) feed reader.
Traditional feed readers rely on in-feed item descriptions to preview content. The problem
it that there is no agreement on what these should contain. Some websites include the entire
article in HTML form, whereas others only include a plain text sentance, making the experience
awkward all-around.

On the other hand, website embeds looks much better and are gererally standerdized. Which is why
scrolling trough a news feed on Telegram is much more enjoyable than scrolling trough one on
your grandpa's RSS reader.

To this end, rss2telegram is a simple bot that follows a list of RSS feeds and posts the new
items as plain links on a specified Telegram chat. There is currently no support for multiple chats
or feed grouping; if you'd like to see that, let me know.

Obtaining
---------

To obtain the program simply run

```
go get github.com/serbuvlad/rss2telegram/cmd/rss2telegram
```

which will obtain and compile the code and place the executable in `$(go env GOPATH)/bin`, which
you can add to your `$PATH` or move elsewhere.

Running
-------

After you write a config file (see below), just run the program as

```
rss2telegram
```

The program does not fork itself and, as such, corresponds to a `Type=exec` SystemD service.

The command line options are

```
  -appname string
        appname to use for XDG directories (default "rss2telegram")
  -config string
        YAML configuation file (default $XDG_CONFIG_HOME/appname/config.yaml)
  -db string
        sqlite3 database to store persistant data in (default $XDG_DATA_HOME/appname/data.db)
  -t duration
        time in minutes to wait between polling feeds (default 1m0s)
```

Configuration
-------------

The program is configured by a simple YAML file which specifies the list of feeds and telegram information.
It looks like this:

```yaml
---
feeds:
  - link: https://xkcd.com/rss.xml
  - link: https://www.reddit.com/.rss
    top: 5 # optional - truncate the feed to it's first 5 items before processing
telegram:
  token: 1234567890:abcedfghijklmopq
  chatid: -1234567890
```

For people new to Telegram botting: to get a token talk to `@BotFather`; to get a chatid
forward a message from it to `@getidsbot`; to send messages you must add your bot as a
group member or channel admin.
