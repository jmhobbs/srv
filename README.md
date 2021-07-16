# srv

This is a very simple web server meant for quick, local use.  I used to use `caddy`, then when that got a little less "zero config" I started using `php -S 127.0.0.1:5000`.  That's more to type though, and rather than do something crazy like turn it into a shell alias, I wrote this quick Go server.

## Usage

This server is not meant to be a production web server, just for quick little tests and moving files around.

```
usage: srv [options] [directory]

  -access-log string
    	Where to write access logs, default is STDOUT. Pass empty string to disable. (default "-")
  -interface string
    	Network interface to listen on (default "127.0.0.1")
  -p int
    	Port to listen on (default 5000)
  -q	Quiet mode, disable most logging
  -v	Verbose mode, enable debug logging
  -version
    	Show version and exit
```
