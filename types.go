package zmq

import "github.com/iotaledger/giota"

type MessageType int

const (
	AllMessages MessageType = iota + 1
	TransactionMsg
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
		AllMessages:                      "",
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

// The Message interface is used to pass ZeroMQ messages.
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

// MilestoneChange contains information for milestone changes.
type MilestoneChange struct {
	Previous giota.Trytes
	Latest   giota.Trytes
}

// Milestone hash contains the hash of the latest milestone.
type MilestoneHash struct {
	Milestone giota.Trytes
}

// DNSCheckerChecking contains information on the hostname and IP address of a
// neighbor that is being checked for changes.
type DNSCheckerChecking struct {
	Hostname string
	IP       string
}

// DNSCheckerOK informs you that the IP address for Hostname has not changed.
type DNSCheckerOK struct {
	Hostname string
}

// DNSCheckerIPChanged informs you that the IP address for Hostname has changed.
type DNSCheckerIPChanged struct {
	Hostname string
	IP       string
}
