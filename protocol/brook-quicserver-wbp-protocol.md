# brook quicserver --withoutBrookProtocol protocol

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

## Client --TCP over QUIC Stream--> Server

```
[SHA256(Password) + (DST Address Length+4) + Unix Timestamp + DST Address] + [DATA]...
```

> The maximum length of `[SHA256(Password) + (DST Address Length+4) + Unix Timestamp + DST Address]` is 2048 bytes

-   `DST Address Length+4`: Big Endian 16-bit unsigned integer
-   [`Unix Timestamp`](https://en.wikipedia.org/wiki/Unix_time): If it is not even, it should be increased by 1. Big Endian 32-bit unsigned integer
-   `DATA`: Actual data being proxied

## Server --TCP over QUIC Stream--> Client

```
[DATA]...
```

## Client --UDP over QUIC Datagram--> Server

```
SHA256(Password) + Unix Timestamp + DST Address + Data
```

> The maximum length of datagram is [1197](https://github.com/quic-go/quic-go/blob/a81365ece88ce9d4601ef140073abadc7657fec8/internal/protocol/params.go#L137) now, and may change in the [future](https://datatracker.ietf.org/doc/html/rfc9221#section-3)

-   [`Unix Timestamp`](https://en.wikipedia.org/wiki/Unix_time): Big Endian 32-bit unsigned integer
-   `Data`: Actual data being proxied

## Server --UDP over QUIC Datagram--> Client

```
DST Address + Data
```

> The maximum length of datagram is [1197](https://github.com/quic-go/quic-go/blob/a81365ece88ce9d4601ef140073abadc7657fec8/internal/protocol/params.go#L137) now, and may change in the [future](https://datatracker.ietf.org/doc/html/rfc9221#section-3)
