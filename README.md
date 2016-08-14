# Omnitool

Omnitool is a tool for using SSH on multiple machines in parallel. It is particularly useful for times when you need a bunch of machines for a short period of time but don't want to invest the time in building up a more robust management system.

Omnitool's goal is to let you think in terms of one machine while working with N machines.

## Installing

```shell
$ go get github.com/jmsdnns/omnitool
$ go install github.com/jmsdnns/omnitool
```

## Using It

Omnitool has help built-in.

```shell
Usage:
  omnitool [command]

Available Commands:
  copy        Copies file to host group
  run         Runs a command on host group

Flags:
  -g, --group string       host group for task
  -h, --help               help for omnitool
      --hostsfile string   path to hosts file (default "hosts.list")
  -k, --keyfile string     path to ssh key
  -u, --username string    username for ssh

Use "omnitool [command] --help" for more information about a command.
```

Running a command on a host group looks like this:

```shell
$ ./omnitool -u ... -k ... -g ... run "ls -l"
CMD: ls -l

Host: 33.33.33.12:22
Result:
total 0
-rw-rw-r-- 1 vagrant vagrant 0 Aug 14 02:53 machine_2

Host: 33.33.33.11:22
Result:
total 0
-rw-rw-r-- 1 vagrant vagrant 0 Aug 14 02:53 machine_1
```

Copying a file to a host group looks like this:

```shell
$ ./omnitool -u ... -k ... -g ... copy hosts.list hosts.list
Host: 33.33.33.12:22
Result: ok

Host: 33.33.33.11:22
Result: ok
```

## Hosts File

A hosts file looks like this:

```
[vagrants]
127.0.0.1:2222
127.0.0.1:2200

[jms labs]
192.168.0.5:22
192.168.0.6:22
```

## Using with Vagrant

Using Omnitool with Vagrant is easy, but you will need to tell Vagrant not to automatically generate SSH keys for each host. This lets us use a single key for each host, which more accurately reflects what you'd have in AWS anyway.

_For more info, read about the `config.ssh.insert_key` flag in [Vagrant's SSH documentation](https://www.vagrantup.com/docs/vagrantfile/ssh_settings.html)._

| Flag   | Default value                         |
| ------ | ------------------------------------- |
| `-u`   | vagrant                               |
| `-k`   | $HOME/.vagrant.d/insecure_private_key |
