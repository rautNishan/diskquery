package types

type TransactionId uint32
type Oid uint32

type ListCell struct {
	PtrValue interface{}
	IntValue int
	OidValue Oid
	XidValue TransactionId
}

type List struct {
	Type            NodeTag
	Length          int
	MaxLength       int
	Elements        *ListCell
	InitialElements []ListCell
}
