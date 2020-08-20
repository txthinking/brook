## Run brook server

Assume with port `9999` and password `hello`

```
$ brook server -l :9999 -p hello
```

Assume your server public IP is `1.2.3.4`, then your brook server is: `1.2.3.4:9999`

> You can stop it with CTRL+C<br/>
> More parameters: $ brook server -h

## Run in background

```
$ nohup brook server -l :9999 -p hello &
```

## Stop background brook

```
$ killall brook
```
