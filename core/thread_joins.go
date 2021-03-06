package core

import (
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// joinInitial creates an outgoing join block for an emtpy thread
func (t *Thread) joinInitial() (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	msg, err := t.buildJoin(t.node().Identity.Pretty())
	if err != nil {
		return nil, err
	}

	res, err := t.commitBlock(msg, pb.ThreadBlock_JOIN, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, repo.JoinBlock, "", ""); err != nil {
		return nil, err
	}

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	log.Debugf("added JOIN to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// join creates an outgoing join block
func (t *Thread) join(inviterId peer.ID) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	msg, err := t.buildJoin(inviterId.Pretty())
	if err != nil {
		return nil, err
	}

	res, err := t.commitBlock(msg, pb.ThreadBlock_JOIN, nil)
	if err != nil {
		return nil, err
	}

	if err := t.indexBlock(res, repo.JoinBlock, "", ""); err != nil {
		return nil, err
	}

	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	log.Debugf("added JOIN to %s: %s", t.Id, res.hash.B58String())

	return res.hash, nil
}

// handleJoinBlock handles an incoming join block
func (t *Thread) handleJoinBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadJoin, error) {
	msg := new(pb.ThreadJoin)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	if err := t.indexBlock(&commitResult{
		hash:   hash,
		header: block.Header,
	}, repo.JoinBlock, "", ""); err != nil {
		return nil, err
	}

	// collect author as an unwelcomed peer
	if msg.Contact != nil {
		pid, err := peer.IDB58Decode(block.Header.Author)
		if err != nil {
			return nil, err
		}
		if err := t.addOrUpdatePeer(pid, protoContactToModel(msg.Contact)); err != nil {
			return nil, err
		}
	}

	return msg, nil
}

// buildJoin builds up a join block
func (t *Thread) buildJoin(inviterId string) (*pb.ThreadJoin, error) {
	msg := &pb.ThreadJoin{
		Inviter: inviterId,
	}
	contact, err := t.getContact(t.node().Identity.Pretty())
	if err != nil {
		return nil, err
	}
	if contact != nil {
		msg.Contact = repoContactToProto(contact)
	}
	return msg, nil
}
