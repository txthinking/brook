# `brook server` Protocol

## Terminology

- **`DST Address`**: The address that the application actually wants to request, address contains IP/domain and port

    ```
    ATYP + IP/Domain + PORT
    ```
    - `ATYP`: 1 byte
        - 0x01: IPv4
        - 0x03: Domain
        - 0x04: IPv6
    - `IP/Domain`: 4/n/16 bytes
        - If ATYP is 0x01, then this is IPv4, 4 bytes
        - If ATYP is 0x03, then this is domain, n bytes, and the first byte is the domain length
        - If ATYP is 0x04, then this is IPv6, 16 bytes
    - `Port`: 2 bytes
        - Big Endian 16-bit unsigned integer

- **`KEY`**: AES key, 32 bytes
    - `KEY`: HKDF_SHA256(Password, Nonce, Info)
        - `Password`: User-defined password
        - `Nonce`: 12 bytes
        - `Info`: [0x62, 0x72, 0x6f, 0x6f, 0x6b]
- **`HKDF`**: Defined in RFC 5869
- **`SHA256`**: Defined in FIPS 180-4
- **`AES`**: Defined in U.S. Federal Information Processing Standards Publication 197
- **`AES-GCM`**: Defined in RFC 5246, 5869

## Client --TCP--> Server

```
Client Nonce + [AES_GCM(Fragment Length) + AES_GCM(Fragment)]...
```

> The maximum length of `AES_GCM(Fragment Length) + AES_GCM(Fragment)` is 2048 bytes

- `Client Nonce`: 12 bytes, randomly generated
    - The nonce should be recalculated when it is not used for the first time, the calculation method: add `1` to the first 8 bytes according to the Little Endian 64-bit unsigned integer
- `Fragment Length`: Big Endian 16-bit unsigned integer
- `Fragment`: Actual data being proxied
    - The first Fragment should be:
        ```
        Unix Timestamp + DST Address
        ```
        - [`Unix Timestamp`](https://en.wikipedia.org/wiki/Unix_time): If it is not even, it should be increased by 1. Big Endian 32-bit unsigned integer

## Server --TCP--> Client

```
Server Nonce + [AES_GCM(Fragment Length) + AES_GCM(Fragment)]...
```

> The maximum length of `AES_GCM(Fragment Length) + AES_GCM(Fragment)` is 2048 bytes

- Server Nonce: 12 bytes, randomly generated
    - The nonce should be recalculated when it is not used for the first time, the calculation method: add `1` to the first 8 bytes according to the Little Endian 64-bit unsigned integer
- `Fragment Length`: Big Endian 16-bit unsigned integer
- `Fragment`: Actual data being proxied

## Client --UDP--> Server

```
Client Nonce + AES_GCM(Fragment)
```

> The maximum length of `Client Nonce + AES_GCM(Fragment)` is 65507 bytes

- `Client Nonce`: 12 bytes, randomly generated each time
- `Fragment`:
    ```
    Unix Timestamp + DST Address + Data
    ```
    - [`Unix Timestamp`](https://en.wikipedia.org/wiki/Unix_time): Big Endian 32-bit unsigned integer
    - `Data`: Actual data being proxied


## Server --UDP--> Client

```
Server Nonce + AES_GCM(Fragment)
```

> The maximum length of `Server Nonce + AES_GCM(Fragment)` is 65507 bytes

- `Server Nonce`: 12 bytes, randomly generated each time
- `Fragment`:
    ```
    DST Address + Data
    ```
    - `Data`: Actual data being proxied
