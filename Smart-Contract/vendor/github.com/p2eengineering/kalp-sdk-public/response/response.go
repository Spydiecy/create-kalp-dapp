package response

import "github.com/hyperledger/fabric-protos-go/peer"

// A response with a representation similar to an HTTP response that can
// be used within another message.
type Response struct {
	peer.Response
}
