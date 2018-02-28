package zmq

import (
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/k0kubun/pp"
	"github.com/pebbe/zmq4"
)

func (c *Client) handleMessages() {
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
					log.Println("Error reconnecting:", err)
				}
				continue
			}

			log.Println("Unable to read message:", err)
			continue
		}

		parts := strings.Fields(msg)
		switch parts[0] {

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

				// DNS checker validating address
				case "dnscv":
					type sn struct {
						Hostname string
						IP       string
					}

					stat := sn{
						Hostname: parts[1],
						IP:       parts[2],
					}
					pp.Print(stat)

				// DNS Check good
				case "dnscc":

					type sn struct {
						Hostname string
					}

					stat := sn{
						Hostname: parts[1],
					}
					pp.Print(stat)

				// IP addressed changed
				case "dnscu":
					type sn struct {
						Hostname string
						IP       string
					}

					stat := sn{
						Hostname: parts[1],
						IP:       parts[2],
					}
					pp.Print(stat)

				// Latest milestone has changed
				case "lmi":
					type sn struct {
						Previous string
						Latest   string
					}

					stat := sn{
						Previous: parts[1],
						Latest:   parts[2],
					}
					pp.Print(stat)

				// Latest SOLID SUBTANGLE milestone has changed
				case "lmsi":
					type sn struct {
						Previous string
						Latest   string
					}

					stat := sn{
						Previous: parts[1],
						Latest:   parts[2],
					}
					pp.Print(stat)

				// Latest SOLID SUBTANGLE milestone hash
				case "lmhs":
					type sn struct {
						Milestone string
					}

					stat := sn{
						Milestone: parts[1],
					}
					pp.Print(stat)


					//pp.Print(stat)
			*/
		default:
			if parts[len(parts)-1] == "sn" {
				continue
			}
			pp.Print(msg)
			fmt.Println()

		}

	}

}
