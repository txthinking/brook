## Run brook wssserver

Make sure your domain name has been resolved successfully. Brook will automatically issue a certificate for your server, assuming your domain is `domain.com`. If there is a firewall, remember to open **TCP on port 80 and 443**.

```
$ brook wssserver --domain domain.com --password hello
```

> You can stop it with CTRL+C<br/>
> More parameters: \$ brook wssserver -h

Then your brook wssserver should be: `wss://domain.com:443`

---

## Run it in the background via `nohup`

> We recommend running the command directly to make sure there are no errors before running it via nohup

```
$ nohup brook wssserver --domain domain.com --password hello &
```

Stop background brook

```
$ killall brook
```

---

## Run as daemon via [`joker`](https://github.com/txthinking/joker) 🔥

> We recommend running the command directly to make sure there are no errors before running it with joker

```
$ joker brook wssserver --domain domain.com --password hello
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

View logs of a command run via joker

> Your can get ID from output by \$ joker list

```
$ joker log <ID>
```

---

## Auto start at boot via [`jinbe`](https://github.com/txthinking/jinbe)

> We recommend running the command directly to make sure there are no errors before running it via jinbe

```
$ jinbe brook wssserver --domain domain.com --password hello
```

Or with joker

```
$ jinbe joker brook wssserver --domain domain.com --password hello
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
