package zmq

import "github.com/iotaledger/giota"

type MessageType int

const (
	TransactionMsg MessageType = iota + 1
	ConfirmationMsg
	ReqStatMsg
	MilestoneChangeMsg
	SolidSubtangleMilestoneChangeMsg
	SolidSubtangleMilestoneHashMsg
	DNSCheckerCheckingMsg
	DNSCheckerOKMsg
	DNSCheckerIPChangedMsg
)

var (
	msgTypes = map[MessageType]string{
		TransactionMsg:                   "tx",
		ConfirmationMsg:                  "sn",
		ReqStatMsg:                       "rstat",
		MilestoneChangeMsg:               "lmi",
		SolidSubtangleMilestoneChangeMsg: "lmsi",
		SolidSubtangleMilestoneHashMsg:   "lmhs",
		DNSCheckerCheckingMsg:            "dnscv",
		DNSCheckerOKMsg:                  "dnscc",
		DNSCheckerIPChangedMsg:           "dnscu",
	}
)

type Message interface{}

// Transaction represents a new transaction on the network.
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

// Confirmation messages arrive when the node considers a transaction confirmed.
type Confirmation struct {
	Index       string
	Hash        string
	AddressHash string
	Trunk       string
	Branch      string
	Bundle      string
}

// ReqStat messages contain information on the state of the transaction requestor
// of the node.
type ReqStat struct {
	ReceiveQueueSize   string
	BroadcastQueueSize string
	TxnToRequest       string
	ReplyQueueSize     string
	NumberOfStoredTxns string
}

type MilestoneChange struct {
	Previous giota.Trytes
	Latest   giota.Trytes
}

type MilestoneHash struct {
	Milestone giota.Trytes
}

type DNSCheckerChecking struct {
	Hostname string
	IP       string
}

type DNSCheckerOK struct {
	Hostname string
}

type DNSCheckerIPChanged struct {
	Hostname string
	IP       string
}
