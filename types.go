package zmq

type MessageType int

const (
	TransactionMsg MessageType = iota + 1
	ConfirmationMsg
	ReqStatMsg
)

var (
	msgTypes = map[MessageType]string{
		TransactionMsg:  "tx",
		ConfirmationMsg: "sn",
		ReqStatMsg:      "rstat",
	}
)

type Message interface {
	Type() MessageType
}

type Transaction struct {
	Hash         string
	Address      string
	Value        string
	Tag          string
	Timestamp    string
	CurrentIndex string
	LastIndex    string
	Bundle       string
	Trunk        string
	Branch       string
	ArrivalDate  string
}

func (t Transaction) Type() MessageType { return TransactionMsg }

type Confirmation struct {
	Index       string
	Hash        string
	AddressHash string
	Trunk       string
	Branch      string
	Bundle      string
}

func (t Confirmation) Type() MessageType { return ConfirmationMsg }

type ReqStat struct {
	ReceiveQueueSize   string
	BroadcastQueueSize string
	TxnToRequest       string
	ReplyQueueSize     string
	NumberOfStoredTxns string
}

func (t ReqStat) Type() MessageType { return ReqStatMsg }
