package zmq

import (
	"log"
	"strings"
	"syscall"

	"github.com/iotaledger/giota"
	"github.com/pebbe/zmq4"
)

func (c *Client) getMessage() (string, error) {
	for {
		msg, err := c.socket.Recv(0)
		if err != nil {
			if zmq4.AsErrno(err) == zmq4.AsErrno(syscall.EAGAIN) {
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

func (c *Client) handleMessages() {
	msgChan := make(chan string)
	errChan := make(chan error)
	for {

		go func() {
			msg, err := c.getMessage()
			if err != nil {
				errChan <- err
				return
			}
			msgChan <- msg
		}()
		select {
		case <-c.stopChan:
			log.Println("disconnecting")
			return
		case err := <-errChan:
			log.Println("Error:", err)
			continue
		case msg := <-msgChan:

			parts := strings.Fields(msg)
			switch parts[0] {

			// Transaction
			case "tx":
				txn := Transaction{
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

				if ch, ok := c.subscriptions[TransactionMsg]; ok {
					ch <- txn
				}

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

				if ch, ok := c.subscriptions[ConfirmationMsg]; ok {
					ch <- msg
				}

			// Tip Requester Statistics
			case "rstat":

				msg := ReqStat{
					ReceiveQueueSize:   parts[1],
					BroadcastQueueSize: parts[2],
					TxnToRequest:       parts[3],
					ReplyQueueSize:     parts[4],
					NumberOfStoredTxns: parts[5],
				}

				if ch, ok := c.subscriptions[ReqStatMsg]; ok {
					ch <- msg
				}
				// Latest milestone has changed
			case "lmi":
				msg := MilestoneChange{
					Previous: giota.Trytes(parts[1]),
					Latest:   giota.Trytes(parts[2]),
				}

				if ch, ok := c.subscriptions[MilestoneChangeMsg]; ok {
					ch <- msg
				}

			// Latest SOLID SUBTANGLE milestone has changed
			case "lmsi":
				msg := MilestoneChange{
					Previous: giota.Trytes(parts[1]),
					Latest:   giota.Trytes(parts[2]),
				}

				if ch, ok := c.subscriptions[SolidSubtangleMilestoneChangeMsg]; ok {
					ch <- msg
				}

				// Latest SOLID SUBTANGLE milestone hash
			case "lmhs":
				msg := MilestoneHash{
					Milestone: giota.Trytes(parts[1]),
				}

				if ch, ok := c.subscriptions[SolidSubtangleMilestoneHashMsg]; ok {
					ch <- msg
				}
				// DNS checker validating address
			case "dnscv":

				msg := DNSCheckerChecking{
					Hostname: parts[1],
					IP:       parts[2],
				}
				if ch, ok := c.subscriptions[DNSCheckerCheckingMsg]; ok {
					ch <- msg
				}

			// DNS Check good
			case "dnscc":
				msg := DNSCheckerOK{
					Hostname: parts[1],
				}
				if ch, ok := c.subscriptions[DNSCheckerOKMsg]; ok {
					ch <- msg
				}

			// IP addressed changed
			case "dnscu":
				msg := DNSCheckerIPChanged{
					Hostname: parts[1],
					IP:       parts[2],
				}
				if ch, ok := c.subscriptions[DNSCheckerIPChangedMsg]; ok {
					ch <- msg
				}

				/*
					//RecentSeenBytes cache hit/miss ratio:
					case "hmr":
						type sn struct {
							Hit  string
							Miss string
						}

						stat := sn{
							Hit:  parts[1],
							Miss: parts[2],
						}
						pp.Print(stat)

					//Adding non-tethered neighbor:
					case "antn":
						type sn struct {
							URI string
						}

						stat := sn{
							URI: parts[1],
						}
						pp.Print(stat)

					//Refused non-tethered neighbor:
					case "rntn":
						type sn struct {
							URI             string
							MaxPeersAllowed string
						}

						stat := sn{
							URI:             parts[1],
							MaxPeersAllowed: parts[2],
						}
						pp.Print(stat)

					//Removed existing tx from request list:
					case "rtl":
						type sn struct {
							Transaction string
						}

						stat := sn{
							Transaction: parts[1],
						}
						pp.Print(stat)

					//Reason to stop: transactionViewModel == null
					case "rtsn":
						type sn struct {
							Transaction string
						}

						stat := sn{
							Transaction: parts[1],
						}
						pp.Print(stat)

					//Reason to stop: !checkSolidity
					case "rtss":
						type sn struct {
							Transaction string
						}

						stat := sn{
							Transaction: parts[1],
						}
						pp.Print(stat)

					//Reason to stop: !LedgerValidator
					case "rtsv":
						type sn struct {
							Transaction string
						}

						stat := sn{
							Transaction: parts[1],
						}
						pp.Print(stat)

					//Reason to stop: transactionViewModel==extraTip
					case "rtsd":
						type sn struct {
							Transaction string
						}

						stat := sn{
							Transaction: parts[1],
						}
						pp.Print(stat)

					//Reason to stop: TransactionViewModel is a tip
					case "rtst":
						type sn struct {
							Transaction string
						}

						stat := sn{
							Transaction: parts[1],
						}
						pp.Print(stat)

					//Reason to stop: transactionViewModel==itself
					case "rtsl":
						type sn struct {
							Transaction string
						}

						stat := sn{
							Transaction: parts[1],
						}
						pp.Print(stat)

					//Tx traversed to find tip:
					case "mctn":
						type sn struct {
							Transaction string
						}

						stat := sn{
							Transaction: parts[1],
						}
						pp.Print(stat)



				*/
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
