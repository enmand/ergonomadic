package irc

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

func IPString(addr net.Addr) Name {
	addrStr := addr.String()
	ipaddr, _, err := net.SplitHostPort(addrStr)
	if err != nil {
		return Name(addrStr)
	}
	return Name(ipaddr)
}

func AddrLookupHostname(addr net.Addr) Name {
	return LookupHostname(IPString(addr))
}

func LookupHostname(addr Name) Name {
	names, err := net.LookupAddr(addr.String())
	if err != nil {
		return Name(addr)
	}

	hostname := strings.TrimSuffix(names[0], ".")
	return Name(hostname)
}

func MaskHostname(addr Name, mask []byte) Name {
	hostname := string(addr)

	host := strings.SplitN(hostname, ".", 2)
	if len(host) == 1 {
		// "simple" hostname, likely a local address
		return NewName(hostname)
	}

	isIP := net.ParseIP(hostname) != nil
	if isIP {
		octets := strings.Split(hostname, ".")
		octetMask := make([]string, 4)

		for i, o := range octets {
			masked := createMask(o, mask)
			octetMask[i] = masked
		}

		hostname = strings.Join(octetMask, ".")
	} else {
		masked := createMask(host[0], mask)
		hostname = fmt.Sprintf("%s.%s", masked, host[1])
	}

	return NewName(hostname)

}

func createMask(s string, mask []byte) string {
	hasher := md5.New()
	hasher.Write(mask)
	hasher.Write([]byte(s))
	sum := hex.EncodeToString(hasher.Sum(nil))
	masked := string(sum)

	return masked[0:8]
}
