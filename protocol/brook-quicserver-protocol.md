# brook quicserver protocol

<!--THEME:github-->
<!--G-R3M673HK5V-->

## Client --TCP over QUIC Stream--> Server

This is the same as the brook server protocol, except that the **QUIC Stream** is used instead of **TCP**

## Server --TCP over QUIC Stream--> Client

This is the same as the brook server protocol, except that the **QUIC Stream** is used instead of **TCP**

## Client --UDP over QUIC Datagram--> Server

This is the same as the brook server protocol, except that the **QUIC Datagram** is used instead of **UDP**

> The maximum length of datagram is [1197](https://github.com/quic-go/quic-go/blob/a81365ece88ce9d4601ef140073abadc7657fec8/internal/protocol/params.go#L137) now, and may change in the [future](https://datatracker.ietf.org/doc/html/rfc9221#section-3)

## Server --UDP over QUIC Datagram--> Client

This is the same as the brook server protocol, except that the **QUIC Datagram** is used instead of **UDP**

> The maximum length of datagram is [1197](https://github.com/quic-go/quic-go/blob/a81365ece88ce9d4601ef140073abadc7657fec8/internal/protocol/params.go#L137) now, and may change in the [future](https://datatracker.ietf.org/doc/html/rfc9221#section-3)
