## brook wsserver

```
SRC --TCP--> brook wsclient/relayoverbrook/dns/tproxy/GUI Client --TCP(Brook Protocol)--> brook wsserver --TCP--> DST
SRC --UDP--> brook wsclient/relayoverbrook/dns/tproxy/GUI Client --TCP(Brook Protocol)--> brook wsserver --UDP--> DST
```

## Run brook wsserver

-   Assume the port is `9999` and the password is `hello`
-   If there is a firewall, remember to open **TCP on port 9999**.

```
brook wsserver --listen :9999 --password hello
```

> You can stop it with CTRL+C<br/>

#### How to connect on the client side

**assuming your server IP is `1.2.3.4`**

-   brook wsserver: `ws://1.2.3.4:9999`
-   password: `hello`

> Connect with CLI: `brook wsclient --wsserver ws://1.2.3.4:9999 --password hello --socks5 127.0.0.1:1080`. More parameters: `brook wsclient -h`<br/>
> Connect with GUI: add info as above

**or get brook link**

```
brook link --server ws://1.2.3.4:9999 --password hello
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link result

**or get brook link with custom domain, the domain can be any domain, even if it's not your domain name**

```
brook link --server ws://hello.com:9999 --password hello --address 1.2.3.4:9999
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link

**or get brook link with `name`**

```
brook link --server ws://hello.com:9999 --password hello --address 1.2.3.4:9999 --name 'my brook wsserver'
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link

#### Block domain and IP in server-side

Check these parameters

-   --blockDomainList
-   --blockCIDR4List
-   --blockCIDR6List
-   --updateListInterval

> More parameters: brook wsserver -h

---

## Run brook wsserver as daemon via [`joker`](https://github.com/txthinking/joker) 🔥

> We recommend running the command directly to make sure there are no errors

```
joker brook wsserver --listen :9999 --password hello
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

> We recommend running the command directly to make sure there are no errors

```
jinbe joker brook wsserver --listen :9999 --password hello
```

View added commmands via jinbe

```
jinbe list
```

Remove a added command via jinbe

```
jinbe remove <ID>
```
