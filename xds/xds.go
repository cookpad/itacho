package xds

const (
	// CDS flag
	CDS ResponseType = iota
	// RDS flag
	RDS
	// EDS flag
	EDS
)

// ResponseType to switch xDS type in itacho
type ResponseType = int

const (
	_ APIVersion = iota + 1
	// v2 xDS API
	V2
	// v3 xDS API
	V3
)

// APIVersion to switch API version in itacho
type APIVersion = int
