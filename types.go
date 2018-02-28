package zmq

import "github.com/iotaledger/giota"

type MessageType int

const (
	// Subscribe to all messages sent by IRI
	AllMessages MessageType = iota + 1

	// Transaction being seen by the node for the first time
	TransactionMsg

	// Transactions newly confirmed
	ConfirmationMsg

	// Stats about some stuff
	ReqStatMsg

	// Milestone changed
	MilestoneChangeMsg

	// Subtangle milestone changed
	SolidSubtangleMilestoneChangeMsg

	// Subtangle milestone hash
	SolidSubtangleMilestoneHashMsg

	// Checking the domain for a neighbor for an updated IP address
	DNSCheckerCheckingMsg

	// Domain check returned the same IP we already have
	DNSCheckerOKMsg

	// Domain check returned a new IP for neighbor
	DNSCheckerIPChangedMsg

	// Count of the number of transactions traversed to find a tip
	TipTraversalCountMsg

	// Removed an existing transaction from the request queue
	TransactionRequestRemovedMsg

	// RecentSeenBytes cache hit/miss ratio
	RecentSeenBytesHitMissMsg

	// Adding a non-tethered neighbor
	AddedNonTetheredNeighborMsg

	// Refusing connection from non-tethered neighbor
	RefusedNonTetheredNeighborMsg

	// Tip selection stopped due to transactionViewModel == null
	TipSelectionStoppedNullMsg

	// Tip selection stopped due to !checkSolidity
	TipSelectionStoppedSolidityCheckMsg

	// Tip selection stopped due to !LedgerValidator
	TipSelectionStoppedLedgerValidatorMsg

	// Tip selection stopped due to transactionViewModel==extraTip
	TipSelectionStoppedExtraTipMsg

	// Tip selection stopped due to TransactionViewModel is a tip
	TipSelectionStoppedIsTipMsg

	// Tip selection stopped due to transactionViewModel==itself
	TipSelectionStoppedSelfMsg
)

var (
	msgTypes = map[MessageType]string{
		AllMessages:     "",
		TransactionMsg:  "tx",
		ConfirmationMsg: "sn",
		ReqStatMsg:      "rstat",

		MilestoneChangeMsg:               "lmi",
		SolidSubtangleMilestoneChangeMsg: "lmsi",
		SolidSubtangleMilestoneHashMsg:   "lmhs",

		DNSCheckerCheckingMsg:  "dnscv",
		DNSCheckerOKMsg:        "dnscc",
		DNSCheckerIPChangedMsg: "dnscu",

		TipTraversalCountMsg:         "mctn",
		TransactionRequestRemovedMsg: "rtl",

		RecentSeenBytesHitMissMsg: "hmr",

		AddedNonTetheredNeighborMsg:   "antn",
		RefusedNonTetheredNeighborMsg: "rntn",

		TipSelectionStoppedNullMsg:            "rtsn",
		TipSelectionStoppedSolidityCheckMsg:   "rtss",
		TipSelectionStoppedLedgerValidatorMsg: "rtsv",
		TipSelectionStoppedExtraTipMsg:        "rtsd",
		TipSelectionStoppedIsTipMsg:           "rtst",
		TipSelectionStoppedSelfMsg:            "rtsl",
	}
)

// The Message interface is used to pass ZeroMQ messages.
type Message interface{}

// Transaction represents a first time seen transaction on the current node.
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

// TransactionHash contains the hash of a transaction.
type TransactionHash struct {
	Hash giota.Trytes
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

// TransactionTraversalCount contains the number of transactions traveresed to
// find a requested tip.
type TransactionTraversalCount struct {
	Count int
}

// RecentSeenBytes cache hit/miss ratio
type RecentSeenBytesHitMiss struct {
	Hit  int
	Miss int
}

// AddedNonTetheredNeighbor contains the URI of a non-tethered neighbor that was
// added.
type AddedNonTetheredNeighbor struct {
	URI string
}

// RefusedNonTetheredNeighbor contains the URI of a non-tethered neighbor that
// was refused a connection along with the MaxPeersAllowed count.
type RefusedNonTetheredNeighbor struct {
	URI             string
	MaxPeersAllowed int
}
