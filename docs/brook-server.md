## brook server

```
SRC --TCP--> brook client/relayoverbrook/dns/tproxy/GUI Client --TCP(Brook Protocol)--> brook server --TCP--> DST
SRC --UDP--> brook client/relayoverbrook/dns/tproxy/GUI Client --UDP(Brook Protocol)--> brook server --UDP--> DST
```

## Run brook server

-   Assume with port `9999` and password `hello`.
-   If there is a firewall, remember to open **TCP and UDP on port 9999**.

```
brook server --listen :9999 --password hello
```

> You can stop it with CTRL+C<br/>

#### How to connect on the client side

**assume your server IP is `1.2.3.4`**

-   brook server: `1.2.3.4:9999`
-   password: `hello`

> Connect with CLI: `brook client --server 1.2.3.4:9999 --password hello --socks5 127.0.0.1:1080`. More parameters: `brook client -h`<br/>
> Connect with GUI: add info as above

**or get brook link**

```
brook link --server 1.2.3.4:9999 --password hello
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link

**or get brook link with `name`**

```
brook link --server 1.2.3.4:9999 --password hello --name 'my brook server'
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link

#### Block domain and IP in server-side

Check these parameters

-   --blockDomainList
-   --blockCIDR4List
-   --blockCIDR6List
-   --updateListInterval

> More parameters: brook server -h

---

## Run brook server as daemon via [`joker`](https://github.com/txthinking/joker) 🔥

> We recommend running the command directly to make sure there are no errors before running it with joker

```
joker brook server --listen :9999 --password hello
```

Get the last command ID

```
joker last
```

View output and error of a command run via joker

```
joker log <ID>
```

View running commmands via joker

```
joker list
```

Stop a running command via joker

```
joker stop <ID>
```

---

## Auto start at boot via [`jinbe`](https://github.com/txthinking/jinbe)

> We recommend running the command directly to make sure there are no errors before running it via jinbe

```
jinbe joker brook server --listen :9999 --password hello
```

View added commmands via jinbe

```
jinbe list
```

Remove a added command via jinbe

```
jinbe remove <ID>
```
