package dig2doggo

import (
	"strings"

	"github.com/kluzzebass/reflag/preprocessor"
	_ "github.com/kluzzebass/reflag/preprocessor/dig" // Register dig preprocessor
	"github.com/kluzzebass/reflag/translator"
)

func init() {
	translator.Register(&Translator{})
}

type Translator struct{}

func (t *Translator) Name() string        { return "dig2doggo" }
func (t *Translator) SourceTool() string  { return "dig" }
func (t *Translator) TargetTool() string  { return "doggo" }
func (t *Translator) IncludeInInit() bool { return false }

func (t *Translator) Translate(args []string, mode string) []string {
	// Apply preprocessing if available
	preprocessedArgs := args
	if p := preprocessor.Get("dig"); p != nil {
		preprocessedArgs = p.Preprocess(args)
	}
	return translateFlags(preprocessedArgs)
}

func translateFlags(args []string) []string {
	var result []string
	var queryName string
	var queryType string
	var nameserver string
	var queryClass string
	skipNext := false

	for i := 0; i < len(args); i++ {
		if skipNext {
			skipNext = false
			continue
		}

		arg := args[i]

		if arg == "--" {
			break
		}

		if strings.HasPrefix(arg, "@") {
			nameserver = arg[1:]
			continue
		}

		if strings.HasPrefix(arg, "+") {
			handlePlusOption(arg[1:], &result)
			continue
		}

		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			if arg[1] == '-' {
				continue
			}

			flags := arg[1:]
			for j := 0; j < len(flags); j++ {
				c := flags[j]
				switch c {
				case '4':
					result = append(result, "-4")
				case '6':
					result = append(result, "-6")
				case 'b':
					if j+1 < len(flags) {
						skipNext = false
					} else if i+1 < len(args) {
						skipNext = true
					}
				case 'c':
					var val string
					if j+1 < len(flags) {
						val = flags[j+1:]
						j = len(flags)
					} else if i+1 < len(args) {
						val = args[i+1]
						skipNext = true
					}
					if val != "" {
						queryClass = val
					}
				case 'f':
					if j+1 < len(flags) {
						j = len(flags)
					} else if i+1 < len(args) {
						skipNext = true
					}
				case 'k':
					if j+1 < len(flags) {
						j = len(flags)
					} else if i+1 < len(args) {
						skipNext = true
					}
				case 'p':
					if j+1 < len(flags) {
						j = len(flags)
					} else if i+1 < len(args) {
						skipNext = true
					}
				case 'q':
					var val string
					if j+1 < len(flags) {
						val = flags[j+1:]
						j = len(flags)
					} else if i+1 < len(args) {
						val = args[i+1]
						skipNext = true
					}
					if val != "" {
						queryName = val
					}
				case 't':
					var val string
					if j+1 < len(flags) {
						val = flags[j+1:]
						j = len(flags)
					} else if i+1 < len(args) {
						val = args[i+1]
						skipNext = true
					}
					if val != "" {
						queryType = strings.ToUpper(val)
					}
				case 'x':
					var val string
					if j+1 < len(flags) {
						val = flags[j+1:]
						j = len(flags)
					} else if i+1 < len(args) {
						val = args[i+1]
						skipNext = true
					}
					if val != "" {
						result = append(result, "-x")
						queryName = val
					}
				case 'm':
					result = append(result, "--debug")
				case 'u':
				case 'i', 'h', 'v':
				}
			}
			continue
		}

		if queryName == "" {
			queryName = arg
		} else if queryType == "" && isValidQueryType(arg) {
			queryType = strings.ToUpper(arg)
		} else if queryClass == "" && isValidQueryClass(arg) {
			queryClass = strings.ToUpper(arg)
		}
	}

	if queryName != "" {
		result = append(result, "-q", queryName)
	}

	if queryType != "" {
		result = append(result, "-t", queryType)
	}

	if queryClass != "" {
		result = append(result, "-c", queryClass)
	}

	if nameserver != "" {
		if strings.Contains(nameserver, "://") {
			result = append(result, "-n", "@"+nameserver)
		} else {
			result = append(result, "-n", nameserver)
		}
	}

	return result
}

func handlePlusOption(opt string, result *[]string) {
	isNegated := strings.HasPrefix(opt, "no")
	if isNegated {
		opt = opt[2:]
	}

	if strings.Contains(opt, "=") {
		parts := strings.SplitN(opt, "=", 2)
		opt = parts[0]
		val := parts[1]

		switch opt {
		case "timeout", "time":
			*result = append(*result, "--timeout", val+"s")
		case "ndots":
			*result = append(*result, "--ndots", val)
		case "bufsize":
		case "edns":
		case "subnet":
			*result = append(*result, "--ecs", val)
		}
		return
	}

	switch opt {
	case "short":
		if !isNegated {
			*result = append(*result, "--short")
		}
	case "tcp", "vc":
		if !isNegated {
			*result = append(*result, "-n", "@tcp://")
		}
	case "trace":
	case "recurse":
		if !isNegated {
			*result = append(*result, "--rd")
		}
	case "dnssec":
		if !isNegated {
			*result = append(*result, "--do")
		}
	case "aa", "aaonly", "aaflag":
		if !isNegated {
			*result = append(*result, "--aa")
		}
	case "ad", "adflag":
		if !isNegated {
			*result = append(*result, "--ad")
		}
	case "cd", "cdflag":
		if !isNegated {
			*result = append(*result, "--cd")
		}
	case "nsid":
		if !isNegated {
			*result = append(*result, "--nsid")
		}
	case "cookie":
		if !isNegated {
			*result = append(*result, "--cookie")
		}
	case "padding":
		if !isNegated {
			*result = append(*result, "--padding")
		}
	case "ede":
		if !isNegated {
			*result = append(*result, "--ede")
		}
	case "search":
		if !isNegated {
			*result = append(*result, "--search")
		}
	case "stats", "cmd", "question", "answer", "authority", "additional", "comments", "rrcomments":
	case "ttlid", "cl", "qr", "split", "identify", "multiline", "onesoa", "nssearch":
	case "fail", "besteffort", "keepopen", "ignore", "crypto", "defname", "expire":
	case "idnout", "ednsnegotiation", "ednsflags", "ednsopt":
	}
}

func isValidQueryType(s string) bool {
	s = strings.ToUpper(s)
	validTypes := map[string]bool{
		"A": true, "AAAA": true, "AFSDB": true, "APL": true, "CAA": true,
		"CDNSKEY": true, "CDS": true, "CERT": true, "CNAME": true, "CSYNC": true,
		"DHCID": true, "DLV": true, "DNAME": true, "DNSKEY": true, "DS": true,
		"EUI48": true, "EUI64": true, "HINFO": true, "HIP": true, "HTTPS": true,
		"IPSECKEY": true, "KEY": true, "KX": true, "LOC": true, "MX": true,
		"NAPTR": true, "NS": true, "NSEC": true, "NSEC3": true, "NSEC3PARAM": true,
		"OPENPGPKEY": true, "PTR": true, "RP": true, "RRSIG": true, "SIG": true,
		"SMIMEA": true, "SOA": true, "SRV": true, "SSHFP": true, "SVCB": true,
		"TA": true, "TKEY": true, "TLSA": true, "TSIG": true, "TXT": true,
		"URI": true, "ZONEMD": true, "ANY": true, "AXFR": true, "IXFR": true,
	}
	return validTypes[s]
}

func isValidQueryClass(s string) bool {
	s = strings.ToUpper(s)
	return s == "IN" || s == "CH" || s == "HS" || s == "CHAOS" || s == "HESIOD"
}
