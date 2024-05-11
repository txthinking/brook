# User System

This content introduces how to develop a user system with Brook. Your system only needs to focus on two concepts: **Token** and **User API**. To support user system, you **must use brook server/wsserver/wssserver/quicserver with the brook protocol**.

<img src="https://brook.app/images/user-system.png" width="500">

## Token

The concept of a token is similar to the user authentication systems of many other systems, such as cookie, session, it represents an identifier that is visible to the user, so **it shoud be unpredictable**. For example, suppose your user system has the concept of an auto-incrementing user ID. You might need to encrypt this user ID to generate a token to be placed on the client side. Alternatively, if your user ID is a UUID, or if you have an auto-incrementing user ID along with an additional UUID, you can also use the UUID directly as a token. This entirely depends on your own decision. If you have developed any user authentication system, this is basic knowledge. No matter what method you use to generate the token, you need to **encode it as hexadecimal**. The length should be controlled at around a few dozen bytes, do not make it too long.

For example, encrypt user id or make session, and encode in hexadecimal:

```
hex_encode(your_encrypt_or_session_function(user id))
// 3ae6afc9fad94abd8985d8ecc77afb273ae6afc9fad94abd8985d8ecc77afb273ae6afc9fad94abd8985d8ecc77afb27
```

For example, UUID:

```javascript
crypto.randomUUID().replaceAll('-', '')
// 3ae6afc9fad94abd8985d8ecc77afb27
```

## User API

Your system must provide an API for Brook Server to validate token. For example: `https://your-api-server.com/a_unpredictable_path`, yes, it is recommended to add an unpredictable path to your https API, of course, you can also use the http api for internal network communication. Brook Server will send GET request to your User API to check if token is valid, the request format is `https://your-api-server.com/a_unpredictable_path?token=xxx`. When the response is 200, the body should be the user's unique identifier, such as user ID; all other status codes are considered to represent an invalid user, and in these cases, the body should be a string describing the error.

For example, your User API is:

```
https://your-api-server.com/a_unpredictable_path
```

Brook Server will send GET request with token to your User API:

```
GET https://your-api-server.com/a_unpredictable_path?token=xxx
```

If the token is valid, your User API should response status code 200 and body should be user's unique identifier, such as user ID 9:

```
HTTP/1.1 200 OK
Content-Length: 1
Content-Type: text/plain; charset=utf-8

9
```

If the token is invalid, or because of any other reasons, the service cannot be provided to this user, such as the user has expired, your User API should response status code non-200 and body should be the short reason:

```
HTTP/1.1 400 BAD REQUEST
Content-Length: 22
Content-Type: text/plain; charset=utf-8

The user 9 has expired
```

## Run Brook Server with your User API

```
brook --serverLog /path/to/log.txt --userAPI https://your-api-server.com/a_unpredictable_path server --listen :9999 --password hello
```

You can count the traffic of each user from serverLog

```
{"bytes":"2190","dst":"8.8.8.8:53","from":"34.105.110.232:49514","network":"tcp","time":"2024-02-26T09:56:12Z","user":"9"}
{"bytes":"2237","dst":"8.8.8.8:53","from":"34.105.110.232:49331","network":"udp","time":"2024-02-26T09:57:12Z","user":"9"}
```

## Generate brook link with token

```
brook link --server 1.2.3.4:9999 --password hello --token xxx
```

## Basic reference implementation

https://github.com/TxThinkingInc/brook-user-system
