## Run brook wsserver

Assume the port is `9999` and the password is `hello`. If there is a firewall, remember to open **TCP on port 9999**.

```
$ brook wsserver --listen :9999 --password hello
```

Assuming your server public IP is `1.2.3.4`, your brook wsserver is `ws://1.2.3.4:9999`

> You can stop it with CTRL+C<br/>
> More parameters: \$ brook wsserver -h

---

## Run in the background via `nohup`

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

> Your can get ID from output by \$ joker list

```
$ joker stop <ID>
```

View log of a command run via joker

> Your can get ID from output by \$ joker list

```
$ joker log <ID>
```

---

## Auto start at boot via [`jinbe`](https://github.com/txthinking/jinbe)

> We recommend running the command directly to make sure there are no errors before running it via jinbe

```
$ jinbe brook wsserver --listen :9999 --password hello
```

Or with joker

```
$ jinbe joker brook wsserver --listen :9999 --password hello
```

View added commmands via jinbe

```
$ jinbe list
```

Remove a added command via jinbe

> Your can get ID from output by \$ jinbe list

```
$ jinbe remove <ID>
```
