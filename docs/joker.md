## Run brook server as daemon with &nbsp; [joker](https://github.com/txthinking/joker)

Install joker

```
$ nami install github.com/txthinking/joker
```

You may want more information on [joker github page](https://github.com/txthinking/joker)

> We recommend running the command directly to make sure there are no errors before running it with joker

```
$ joker brook server -l :9999 -p hello
```

## View running commmands with joker

```
$ joker list
```

## Stop a running command with joker

> Your can get ID from output by $ joker list

```
$ joker stop <ID>
```

## View log of a running or stopped command with joker

> Your can get ID from output by $ joker list

```
$ joker log <ID>
```

