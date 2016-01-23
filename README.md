# Omnitool

Omnitool is a tool for using SSH on multiple machines in parallel. It is particularly useful for things where you need a bunch of machines for a short period of time and thus don't want to build something more robust with a tool like [Ansible](https://ansible.com/) or [Puppet](https://puppetlabs.com/).

Omnitool's goal is to let you think in terms of one machine while working with N machines.

## Installing

```
$ go get github.com/jmsdnns/omnitool
$ go install github.com/jmsdnns/omnitool
```

## Using It

Omnitool has help built-in.

Each command line argument can be supplied as an environment variable if you prefer. The variable is listed after the text describing what each argument does

```
$ omnitool -h
NAME:
   omnitool - Simple SSH pools, backed by machine lists

USAGE:
   omnitool [global options] command [command options] [arguments...]

VERSION:
   0.1

COMMANDS:
   run		Runs command on machine group
   scp		Copies file to machine group
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --list, -l 		Path to machine list file [$OMNI_MACHINE_LIST]
   --username, -u 	Username for machine group [$OMNI_USERNAME]
   --keyfile, -k 	Path to auth key [$OMNI_KEYFILE]
   --group, -g 		Machine group to perform task on [$OMNI_MACHINE_GROUP]
   --help, -h		show help
   --version, -v	print the version
```

## Machine Lists

A machine list looks like this:

```
[vagrants]
127.0.0.1:2222
127.0.0.1:2200
```

You have to make it by hand for now, but once that's in place, you can use Omnitool to run commands on machine groups in parallel by supplying the group name with the `-g` argument.

```
$ omnitool -g vagrants run "ls -l"
Hostname: 127.0.0.1:2222
Result:
total 8
-rw-r--r--  1 vagrant  vagrant    0 Aug 29 22:28 bbx1
drwxr-xr-x  7 vagrant  vagrant  512 Aug 29 22:29 boombox

Hostname: 127.0.0.1:2200
Result:
total 8
-rw-r--r--  1 vagrant  vagrant    0 Aug 29 22:28 bbx2
drwxr-xr-x  7 vagrant  vagrant  512 Aug 29 22:29 boombox

CMD:  [ls -l]
```

## Using with Vagrant

It's easy to use omnitool with Vagrant. Create a machine list called `vagrants` and then use the following command line arguments.

| Flag   | Default value                         |
| ------ | ------------------------------------- |
| `-u`   | vagrant                               |
| `-k`   | $HOME/.vagrant.d/insecure_private_key |
| `-g`   | vagrants                              |
