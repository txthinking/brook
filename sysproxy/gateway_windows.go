package sysproxy

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/txthinking/brook/sysproxy/internal/winhelper"
)

const bufSize = 32 * 1024 // 32KB buffer for GetAdaptersAddresses result, definitely enough

// readInterfaceForGateway parses one AF_INET gateway addr from _IP_ADAPTER_ADDRESSES_LH struct and returns the next adapter in the linked list
func readInterfaceForGateway(reader winhelper.RawPointerReader) (next winhelper.RawPointerReader, gateway string) {

	/*
		typedef struct _IP_ADAPTER_ADDRESSES_LH {
			union {
				ULONGLONG Alignment;
				struct {
				ULONG    Length;
				IF_INDEX IfIndex;
				};
			};
			struct _IP_ADAPTER_ADDRESSES_LH    *Next;
			PCHAR                              AdapterName;
			PIP_ADAPTER_UNICAST_ADDRESS_LH     FirstUnicastAddress;
			PIP_ADAPTER_ANYCAST_ADDRESS_XP     FirstAnycastAddress;
			PIP_ADAPTER_MULTICAST_ADDRESS_XP   FirstMulticastAddress;
			PIP_ADAPTER_DNS_SERVER_ADDRESS_XP  FirstDnsServerAddress;
			PWCHAR                             DnsSuffix;
			PWCHAR                             Description;
			PWCHAR                             FriendlyName;
			BYTE                               PhysicalAddress[MAX_ADAPTER_ADDRESS_LENGTH];
			ULONG                              PhysicalAddressLength;
			union {
				ULONG Flags;
				struct {
				ULONG DdnsEnabled : 1;
				ULONG RegisterAdapterSuffix : 1;
				ULONG Dhcpv4Enabled : 1;
				ULONG ReceiveOnly : 1;
				ULONG NoMulticast : 1;
				ULONG Ipv6OtherStatefulConfig : 1;
				ULONG NetbiosOverTcpipEnabled : 1;
				ULONG Ipv4Enabled : 1;
				ULONG Ipv6Enabled : 1;
				ULONG Ipv6ManagedAddressConfigurationSupported : 1;
				};
			};
			ULONG                              Mtu;
			IFTYPE                             IfType;
			IF_OPER_STATUS                     OperStatus;
			IF_INDEX                           Ipv6IfIndex;
			ULONG                              ZoneIndices[16];
			PIP_ADAPTER_PREFIX_XP              FirstPrefix;
			ULONG64                            TransmitLinkSpeed;
			ULONG64                            ReceiveLinkSpeed;
			PIP_ADAPTER_WINS_SERVER_ADDRESS_LH FirstWinsServerAddress;
			PIP_ADAPTER_GATEWAY_ADDRESS_LH     FirstGatewayAddress;

			// The rest part is useless, whatever:)

			ULONG                              Ipv4Metric;
			ULONG                              Ipv6Metric;
			IF_LUID                            Luid;
			SOCKET_ADDRESS                     Dhcpv4Server;
			NET_IF_COMPARTMENT_ID              CompartmentId;
			NET_IF_NETWORK_GUID                NetworkGuid;
			NET_IF_CONNECTION_TYPE             ConnectionType;
			TUNNEL_TYPE                        TunnelType;
			SOCKET_ADDRESS                     Dhcpv6Server;
			BYTE                               Dhcpv6ClientDuid[MAX_DHCPV6_DUID_LENGTH];
			ULONG                              Dhcpv6ClientDuidLength;
			ULONG                              Dhcpv6Iaid;
			PIP_ADAPTER_DNS_SUFFIX             FirstDnsSuffix;
		} IP_ADAPTER_ADDRESSES_LH, *PIP_ADAPTER_ADDRESSES_LH;
	*/

	_, reader = reader.ULongLong() // Alignment
	next, reader = reader.Deref()  // Next
	_, reader = reader.PChar()     // AdapterName
	_, reader = reader.Ptr()       // FirstUnicastAddress
	_, reader = reader.Ptr()       // FirstAnycastAddress
	_, reader = reader.Ptr()       // FirstMulticastAddress
	_, reader = reader.Ptr()       // FirstDnsServerAddress
	_, reader = reader.PWChar()    // DnsSuffix
	_, reader = reader.PWChar()    // Description
	_, reader = reader.PWChar()    // FriendlyName
	_, reader = reader.Bytes(8)    // PhysicalAddress[MAX_ADAPTER_ADDRESS_LENGTH=8]
	_, reader = reader.ULong()     // PhysicalAddressLength
	_, reader = reader.ULong()     // Flags
	_, reader = reader.ULong()     // Mtu
	_, reader = reader.DWord()     // IfType
	var status int32
	status, reader = reader.Int() // OperStatus
	if status != 1 {
		// Link is down.
		// return
	}
	_, reader = reader.ULong() // Ipv6IfIndex
	for i := 0; i < 16; i++ {
		_, reader = reader.ULong() // ZoneIndices
	}
	_, reader = reader.Ptr()     // FirstPrefix
	_, reader = reader.ULong64() // TransmitLinkSpeed
	_, reader = reader.ULong64() // ReceiveLinkSpeed
	_, reader = reader.Ptr()     // FirstWinsServerAddress
	if !reader.IsNilPointer() {
		var gatewayReader winhelper.RawPointerReader // Type is _IP_ADAPTER_GATEWAY_ADDRESS_LH
		gatewayReader, _ = reader.Deref()

		for {
			/*
				typedef struct _IP_ADAPTER_GATEWAY_ADDRESS_LH {
				  union {
				    ULONGLONG Alignment;
				    struct {
				      ULONG Length;
				      DWORD Reserved;
				    };
				  };
				  struct _IP_ADAPTER_GATEWAY_ADDRESS_LH *Next;
				  SOCKET_ADDRESS                        Address;
				} IP_ADAPTER_GATEWAY_ADDRESS_LH, *PIP_ADAPTER_GATEWAY_ADDRESS_LH;

			*/
			_, gatewayReader = gatewayReader.ULongLong() // Alignment
			next, gatewayReader = gatewayReader.Deref()

			var addr winhelper.RawPointerReader // Type is sockaddr_in
			addr, gatewayReader = gatewayReader.Deref()

			/*
				typedef struct sockaddr_in {
					#if(_WIN32_WINNT < 0x0600)
					    short   sin_family;
					#else //(_WIN32_WINNT < 0x0600)
					    ADDRESS_FAMILY sin_family;
					#endif //(_WIN32_WINNT < 0x0600)
					    USHORT sin_port;
					    IN_ADDR sin_addr;
					    CHAR sin_zero[8];
				} SOCKADDR_IN, *PSOCKADDR_IN;
			*/
			_, gatewayReader = gatewayReader.Int() // AddrLength
			var family uint16
			family, addr = addr.UShort() // sin_family
			_, addr = addr.UShort()      // sin_port
			if family == syscall.AF_INET {
				gtBytes, _ := addr.Bytes(4)
				gateway = fmt.Sprintf("%d.%d.%d.%d", gtBytes[0], gtBytes[1], gtBytes[2], gtBytes[3])
				return
			}

			gatewayReader = next
		}
	}

	return
}

// GetDefaultGateway returns default gateway
func GetDefaultGateway() (string, error) {
	h, err := syscall.LoadLibrary("iphlpapi.dll")
	if err != nil {
		return "", err
	}
	f, err := syscall.GetProcAddress(h, "GetAdaptersAddresses")
	if err != nil {
		return "", err
	}

	var buf [bufSize]byte

	var bufSizeVar uint32 = bufSize

	ret, _, _ := syscall.Syscall6(
		uintptr(f),
		5,
		syscall.AF_INET, // Family AF_INET=2
		0x0080|0x0001|0x0002|0x0004|0x0008, // Flags GAA_FLAG_INCLUDE_GATEWAYS | GAA_FLAG_SKIP_UNICAST | GAA_FLAG_SKIP_ANYCAST | GAA_FLAG_SKIP_MULTICAST | GAA_FLAG_SKIP_DNS_SERVER
		0,                                    // Reserved
		uintptr(unsafe.Pointer(&buf)),        // buf
		uintptr(unsafe.Pointer(&bufSizeVar)), // bufSize
		0,
	)
	switch ret {
	case 0:
	case 8:
		return "", errors.New("not enough memory")
	case 87:
		return "", errors.New("invalid parameter")
	case 111:
		return "", errors.New("buffer Overflow")
	case 232:
		return "", errors.New("no data")
	default:
		return "", errors.New("unknown error")
	}

	var reader winhelper.RawPointerReader
	reader = winhelper.New(uintptr(unsafe.Pointer(&buf)), bufSize)
	for reader.IsValid() {
		next, gt := readInterfaceForGateway(reader)
		if gt != "" {
			return gt, nil
		}
		reader = next
	}

	_ = buf // Prevent GC of buf
	return "", errors.New("All interfaces are consumed")
}
