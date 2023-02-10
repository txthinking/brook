# brook wsserver --withoutBrookProtocol protocol

<!--THEME:github-->
<!--G-R3M673HK5V-->

## Terminology

-   **`DST Address`**: The address that the application actually wants to request, address contains IP/domain and port

    ```
    ATYP + IP/Domain + PORT
    ```

    -   `ATYP`: 1 byte
        -   0x01: IPv4
        -   0x03: Domain
        -   0x04: IPv6
    -   `IP/Domain`: 4/n/16 bytes
        -   If ATYP is 0x01, then this is IPv4, 4 bytes
        -   If ATYP is 0x03, then this is domain, n bytes, and the first byte is the domain length
        -   If ATYP is 0x04, then this is IPv6, 16 bytes
    -   `Port`: 2 bytes
        -   Big Endian 16-bit unsigned integer

-   **`Password`**: User-defined password
-   **`SHA256`**: Defined in FIPS 180-4

## Client --TCP--> Server

```
[Standard WebSocket Protocol Header] + [SHA256(Password) + (DST Address Length+4) + Unix Timestamp + DST Address] + [DATA]...
```

> The maximum length of `[SHA256(Password) + (DST Address Length+4) + Unix Timestamp + DST Address]` is 2048 bytes

-   `DST Address Length+4`: Big Endian 16-bit unsigned integer
-   [`Unix Timestamp`](https://en.wikipedia.org/wiki/Unix_time): If it is not even, it should be increased by 1. Big Endian 32-bit unsigned integer
-   `DATA`: Actual data being proxied

## Server --TCP--> Client

```
[Standard WebSocket Protocol Header] + [DATA]...
```

## Client --UDP(UDP over TCP)--> Server

```
[Standard WebSocket Protocol Header] + [SHA256(Password) + (DST Address Length+4) + Unix Timestamp + DST Address] + [Fragment Length + Fragment]...
```

> The maximum length of `[SHA256(Password) + (DST Address Length+4) + Unix Timestamp + DST Address]` is 2048 bytes<br/>
> The maximum length of `[Fragment Length + Fragment]` is 65507 bytes<br/>

-   `DST Address Length+4`: Big Endian 16-bit unsigned integer
-   `Fragment Length`: Big Endian 16-bit unsigned integer
-   `Fragment`: Actual data being proxied
-   [`Unix Timestamp`](https://en.wikipedia.org/wiki/Unix_time): If it is not odd, it should be increased by 1. Big Endian 32-bit unsigned integer

## Server --UDP(UDP over TCP)--> Client

```
[Standard WebSocket Protocol Header] + [Fragment Length + Fragment]...
```

> The maximum length of `[Fragment Length + Fragment]` is 65507 bytes<br/>

-   `Fragment Length`: Big Endian 16-bit unsigned integer
-   `Fragment`: Actual data being proxied
