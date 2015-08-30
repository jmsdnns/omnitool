# Omnitool

Omnitool is a tool for using SSH on multiple machines in parallel.

Omnitool's goal is to make managing networks of computers easier. Think in terms of one machine while working on N machines.

## Using It

First, install it.

```
$ go get github.com/jmsdnns/omnitool
$ go get github.com/codegangsta/cli
$ go install github.com/jmsdnns/omnitool
```

Omnitool has help built-in and is explicit about what default values are used.

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
   --list, -l "machines.list"                                       Path to machine list file [$OMNI_MACHINE_LIST]
   --username, -u "vagrant"						                    Username for machine group [$OMNI_USERNAME]
   --keyfile, -k "/Users/jmsdnns/.vagrant.d/insecure_private_key"	Path to auth key [$OMNI_KEYFILE]
   --group, -g "vagrants"                                           Machine group to perform task on [$OMNI_MACHINE_GROUP]
   --help, -h                                                       show help
   --version, -v							                        print the version
```

The help is generated from [Jeremy Saenz's `cli`](https://github.com/codegangsta/) library, which omnitool uses to parse command line input.

## It's Early

The tool is new and the ideas are early. Some steps are manual. You have to create the machine list by hand, for example.

## Machine Lists

I use [Boombox](https://github.com/jmsdnns/boombox) to instantiate multiple FreeBSD VM's, and then I put their IP's in a file called `machines.list`.

_This list will eventually be generated_

A machine list looks like this:

```
[vagrants]
127.0.0.1:2222
127.0.0.1:2200
```

Once that's in place, you can use Omnitool to run commands on machine groups in parallel by supplying the group name with the `-g` argument.

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

## With Vagrant

Omnitool's default values assume you want to use Vagrant.

## Not With Vagrant

To override the authentication details, pass the `-u` flag and supply a username and pass a `-k` flag and supply the path to your SSH key.

```
$ omnitool -u jmsdnns -k ~/.ssh/id_rsa -g apiservers run "ls"
```

## Next Up

* SFTP
* Machine list generation
* Provisioning
