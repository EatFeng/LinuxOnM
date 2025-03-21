package constant

type DBContext string

const (
	DateTimeLayout     = "2006-01-02 15:04:05"
	DateTimeSlimLayout = "20060102150405"
)

const (
	TypeSystem = "system"

	DB DBContext = "db"
)

const (
	SelfSigned = "selfSigned"
	SSLReady   = "ready"
)
