## Run brook wsserver

Assume with port `9999` and password `hello`. If there is a firewall, remember to allow TCP on this port.

```
$ brook wsserver --listen :9999 --password hello
```

Assume your server public IP is `1.2.3.4`, then your brook wsserver is: `ws://1.2.3.4:9999`

> More parameters: $ brook wsserver -h

---

## Run in background via `nohup`

> We recommend running the command directly to make sure there are no errors before running it via nohup

```
$ nohup brook wsserver --listen :9999 --password hello &
```

Stop background brook

```
$ killall brook
```

---

## Run as daemon via [`joker`](https://github.com/txthinking/joker) 🔥

> We recommend running the command directly to make sure there are no errors before running it with joker

```
$ joker brook wsserver --listen :9999 --password hello
```

View running commmands via joker

```
$ joker list
```

Stop a running command via joker

> Your can get ID from output by $ joker list

```
$ joker stop <ID>
```

View log of a command run via joker

> Your can get ID from output by $ joker list

```
$ joker log <ID>
```

---

## Auto start at boot via [`boa`](https://github.com/brook-community/boa)

> We recommend running the command directly to make sure there are no errors before running it via boa

```
$ boa brook wsserver --listen :9999 --password hello
```

Or with joker

```
$ boa joker brook wsserver --listen :9999 --password hello
```
