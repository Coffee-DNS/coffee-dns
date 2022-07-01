package record

type Record struct {
	// The hostname of the DNS record
	Hostname string

	// The value (address, cname, etc)
	Value string

	// The record type( A, AAAA, CNAME, etc)
	Type string
}
