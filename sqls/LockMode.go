package sqls

type LockMode string

const (
	FOR_UPDATE = "FOR UPDATE"
	IN_SHARE_MODE = "LOCK IN SHARE MODE"
)
