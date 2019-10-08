package expreduce

import (
	"github.com/asaskevich/govalidator"
	"github.com/corywalker/expreduce/expreduce/atoms"
	"github.com/corywalker/expreduce/pkg/expreduceapi"
	"net"
)

func isIPAddress(s string) bool {
	return net.ParseIP(s) != nil
}

func isHost(s string) bool {
	return govalidator.IsDNSName(s)
}

func getNetworkDefinitions() (defs []Definition) {
	defs = append(defs, Definition{
		Name:    "IPAddress",
		Details: `Symbolic representation of a IPv4 or IPv6 address`,
	})
	defs = append(defs, Definition{
		Name:    "URL",
		Details: `Symbolic representation of a URL, potentially including a port number (default: 80)`,
	})
	defs = append(defs, Definition{
		Name:    "HostLookup",
		Details: `Look up an IP address or host name (DNS and reverse DNS)`,
		legacyEvalFn: func(this expreduceapi.ExpressionInterface, es expreduceapi.EvalStateInterface) expreduceapi.Ex {
			if len(this.GetParts()) != 2 {
				return this
			}
			arg1 := this.GetPart(1)
			var argAsString string
			switch arg1.(type) {
			case *atoms.Symbol:
				argAsString = arg1.(*atoms.Symbol).String()
			case *atoms.String:
				argAsString = arg1.(*atoms.String).Val
			}
			toReturn := atoms.NewExpression([]expreduceapi.Ex{})
			if isIPAddress(argAsString) {
				ips, _ := net.LookupAddr(argAsString)
				toReturn.AppendEx(atoms.NewSymbol("Network`Host"))
				toReturn.AppendEx(atoms.NewString(ips[0]))
			} else if isHost(argAsString) {
				hosts, _ := net.LookupHost(argAsString)
				toReturn.AppendEx(atoms.NewSymbol("Network`IPAddress"))
				toReturn.AppendEx(atoms.NewString(hosts[0]))
			} else {
				return this
			}
			return toReturn
		},
	})
	return defs
}
