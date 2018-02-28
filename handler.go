package zmq

import (
	"log"
	"strconv"
	"strings"
	"syscall"

	"github.com/iotaledger/giota"
	"github.com/pebbe/zmq4"
)

// Loop until we receive a valid message or an error
func (c *Client) getMessage() (string, error) {
	for {
		msg, err := c.socket.Recv(0)
		if err != nil {
			if zmq4.AsErrno(err) == zmq4.AsErrno(syscall.EAGAIN) {
				// This happens after a few seconds and from
				// what I've seen is perfectly normal.
				continue
			}

			if err == zmq4.ETIMEDOUT {
				log.Println("zmq timeout")
				err := c.connect()
				if err != nil {
					return msg, err
				}
				continue
			}

			log.Println("Unable to read message:", err)
			continue
		}
		return msg, err
	}
}
func (c *Client) sendMessage(msg Message, msgType MessageType) {
	if ch, ok := c.subscriptions[msgType]; ok {
		ch <- msg
	}
	if ch, ok := c.subscriptions[AllMessages]; ok {
		ch <- msg
	}
}
func (c *Client) handleMessages() {
	msgChan := make(chan string)
	errChan := make(chan error)
	for {

		// Launch a go routine to check for new messages
		go func() {
			msg, err := c.getMessage()
			if err != nil {
				errChan <- err
				return
			}
			msgChan <- msg
		}()

		select {

		// Check to see if we should disconnect and stop processing
		case <-c.stopChan:
			log.Println("disconnecting")
			return

		// Check to see if the getMessage go routine returned an error
		case err := <-errChan:
			log.Println("Error:", err)
			continue

		// Handle the message
		case msg := <-msgChan:

			parts := strings.Fields(msg)
			switch parts[0] {

			// Transaction
			case "tx":
				msg := Transaction{
					Hash:         parts[1],
					Address:      parts[2],
					Value:        parts[3],
					Tag:          parts[4],
					Timestamp:    parts[5],
					CurrentIndex: parts[6],
					LastIndex:    parts[7],
					Bundle:       parts[8],
					Trunk:        parts[9],
					Branch:       parts[10],
					ArrivalDate:  parts[11],
				}

				c.sendMessage(msg, TransactionMsg)

			// Transaction confirmed
			case "sn":
				msg := Confirmation{
					Index:       parts[1],
					Hash:        parts[2],
					AddressHash: parts[3],
					Trunk:       parts[4],
					Branch:      parts[5],
					Bundle:      parts[6],
				}

				c.sendMessage(msg, ConfirmationMsg)

			// Tip Requester Statistics
			case "rstat":
				msg := ReqStat{
					ReceiveQueueSize:   parts[1],
					BroadcastQueueSize: parts[2],
					TxnToRequest:       parts[3],
					ReplyQueueSize:     parts[4],
					NumberOfStoredTxns: parts[5],
				}

				c.sendMessage(msg, ReqStatMsg)

				// Latest milestone has changed
			case "lmi":
				msg := MilestoneChange{
					Previous: giota.Trytes(parts[1]),
					Latest:   giota.Trytes(parts[2]),
				}

				c.sendMessage(msg, MilestoneChangeMsg)

			// Latest SOLID SUBTANGLE milestone has changed
			case "lmsi":
				msg := MilestoneChange{
					Previous: giota.Trytes(parts[1]),
					Latest:   giota.Trytes(parts[2]),
				}

				c.sendMessage(msg, SolidSubtangleMilestoneChangeMsg)

				// Latest SOLID SUBTANGLE milestone hash
			case "lmhs":
				msg := TransactionHash{
					Hash: giota.Trytes(parts[1]),
				}

				c.sendMessage(msg, SolidSubtangleMilestoneHashMsg)

				// DNS checker validating address
			case "dnscv":

				msg := DNSCheckerChecking{
					Hostname: parts[1],
					IP:       parts[2],
				}
				c.sendMessage(msg, DNSCheckerCheckingMsg)

				// DNS Check good
			case "dnscc":
				msg := DNSCheckerOK{
					Hostname: parts[1],
				}
				c.sendMessage(msg, DNSCheckerOKMsg)

			// IP addressed changed
			case "dnscu":
				msg := DNSCheckerIPChanged{
					Hostname: parts[1],
					IP:       parts[2],
				}
				c.sendMessage(msg, DNSCheckerIPChangedMsg)

			//Tx traversed to find tip:
			case "mctn":
				count, err := strconv.Atoi(parts[1])
				if err != nil {
					continue
				}
				msg := TransactionTraversalCount{
					Count: count,
				}

				c.sendMessage(msg, TipTraversalCountMsg)

			//Removed existing tx from request list:
			case "rtl":
				msg := TransactionHash{
					Hash: giota.Trytes(parts[1]),
				}

				c.sendMessage(msg, TransactionRequestRemovedMsg)
			//RecentSeenBytes cache hit/miss ratio:
			case "hmr":
				hit, err := strconv.Atoi(parts[1])
				if err != nil {
					continue
				}
				miss, err := strconv.Atoi(parts[2])
				if err != nil {
					continue
				}
				msg := RecentSeenBytesHitMiss{
					Hit:  hit,
					Miss: miss,
				}
				c.sendMessage(msg, RecentSeenBytesHitMissMsg)

			//Adding non-tethered neighbor:
			case "antn":
				msg := AddedNonTetheredNeighbor{
					URI: parts[1],
				}
				c.sendMessage(msg, AddedNonTetheredNeighborMsg)

			//Refused non-tethered neighbor:
			case "rntn":
				max, err := strconv.Atoi(parts[1])
				if err != nil {
					continue
				}

				msg := RefusedNonTetheredNeighbor{
					URI:             parts[1],
					MaxPeersAllowed: max,
				}
				c.sendMessage(msg, RefusedNonTetheredNeighborMsg)

			//Reason to stop: transactionViewModel == null
			case "rtsn":
				msg := TransactionHash{
					Hash: giota.Trytes(parts[1]),
				}

				c.sendMessage(msg, TipSelectionStoppedNullMsg)
			//Reason to stop: !checkSolidity
			case "rtss":
				msg := TransactionHash{
					Hash: giota.Trytes(parts[1]),
				}

				c.sendMessage(msg, TipSelectionStoppedSolidityCheckMsg)
			//Reason to stop: !LedgerValidator
			case "rtsv":
				msg := TransactionHash{
					Hash: giota.Trytes(parts[1]),
				}

				c.sendMessage(msg, TipSelectionStoppedLedgerValidatorMsg)
			//Reason to stop: transactionViewModel==extraTip
			case "rtsd":
				msg := TransactionHash{
					Hash: giota.Trytes(parts[1]),
				}

				c.sendMessage(msg, TipSelectionStoppedExtraTipMsg)
			//Reason to stop: TransactionViewModel is a tip
			case "rtst":
				msg := TransactionHash{
					Hash: giota.Trytes(parts[1]),
				}

				c.sendMessage(msg, TipSelectionStoppedIsTipMsg)
			//Reason to stop: transactionViewModel==itself
			case "rtsl":
				msg := TransactionHash{
					Hash: giota.Trytes(parts[1]),
				}

				c.sendMessage(msg, TipSelectionStoppedSelfMsg)
			default:
				// Filter out this funky sn message
				if parts[len(parts)-1] == "sn" {
					continue
				}
				// Unknown message type here
			}
		}

	}
}
