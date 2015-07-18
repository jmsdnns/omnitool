# Omnitool

The idea is simple. Using SSH, you can tell one or more machines to do something. Either tell the pool of machines to run some command or a copy a script to each machine and execute it in parallel.

Omnitool's purpose is to make it easy to manage N machines as easily as we manage one. It follows the [ZOI rule](https://en.wikipedia.org/wiki/Zero_one_infinity_rule) and provides layers above that for the context of managing networks of computers.

It currently builds machines on Digital Ocean using [Boombox](https://github.com/jmsdnns/boombox).

Omnitool believes it is easier to throw machines away than to upgrade them. To build a network of machines is to 1) ask Digital Ocean for N machines, 2) SSH to each in parallel to configure them, 3) provide an easy mechanism for performing commands on all live machines, and 4) terminate them at will.

In a somewhat accurate, single sentence: _Omnitool is tool for managing networks via SSH pools_.

## Installing

Omnitool requires a working Go environment. [Installing Go](https://golang.org/doc/install).

```
$ go get https://github.com/jmsdnns/omnitool
```

## Using Omnitool

Here is an example where 5 machines are created for a load test using [siege](https://www.joedog.org/siege-home/). Omnitool runs the siege command on all five machines in parallel

```
$ omnitool create machines 5 -g 'microarmy'
$ omnitool install siege -g 'microarmy'
$ omnitool run -c 'siege -c 100 -t 60s http://jmsdnns.com/' -g 'microarmy'
$ omnitool terminate microarmy
```

Here is an example setting up 5 machines to run a go service.

```
$ omnitool create machines 10 -g 'api servers'
$ omnitool install git go myserver
```

You're probably wondering, "how do I write that _myserver_ part?"

## Extending Omnitool

To extend Omnitool is to change the tracks you use with boombox. Just change the repo Omnitool uses to one with your changes and those tracks will be available as Omnitool installers.