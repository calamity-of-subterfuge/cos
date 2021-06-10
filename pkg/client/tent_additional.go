package client

import (
	"github.com/calamity-of-subterfuge/cos/pkg/srvpkts"
	"github.com/calamity-of-subterfuge/cos/pkg/unitdets"
	"github.com/calamity-of-subterfuge/cos/pkg/utils"
)

// TentOffer describes a single offer in a tent. The offer is visible on both
// the initiating and target teams tents. You can only see offers on your own
// tent.
type TentOffer struct {
	// InitiatingTeam is the team initiating the offer.
	InitiatingTeam int

	// TargetTeam is the team recieving the offer
	TargetTeam int

	// Offer is the offer of resources being made, which will be taken from the
	// initiating team and given to the target team if the offer is accepted.
	Offer map[string]int

	// Request is the request of resources, which will be taken from the target
	// team and given to the initiating team if the offer is accepted.
	Request map[string]int
}

// TentAdditional is the additional information for a tent
type TentAdditional struct {
	// IncomingOffers are the offers from other teams which have been made
	// to the team owning this tent. The keys are uids of the offers.
	IncomingOffers map[string]TentOffer

	// OutgoingOffers are the offers from this team which have been made to
	// other teams. The keys are the uids of the offers.
	OutgoingOffers map[string]TentOffer
}

func (a *TentAdditional) Update(packet *srvpkts.SmartObjectUpdatePacket) error {
	var updateDetails unitdets.TentUpdateDetails
	_, err := utils.DecodeWithType(packet.Additional.(map[string]interface{}), &updateDetails)
	if err != nil {
		return err
	}

	for _, uid := range updateDetails.RemovedOutgoingOffers {
		delete(a.OutgoingOffers, uid)
	}
	for _, uid := range updateDetails.RemovedIncomingOffers {
		delete(a.IncomingOffers, uid)
	}
	for uid, offer := range updateDetails.AddedOutgoingOffers {
		a.OutgoingOffers[uid] = TentOffer(offer)
	}
	for uid, offer := range updateDetails.AddedIncomingOffers {
		a.IncomingOffers[uid] = TentOffer(offer)
	}
	return nil
}

func init() {
	registerSmartObjectUnitAdditional("tent", func(sync *srvpkts.SmartObjectSync) (SmartObjectAdditional, error) {
		var syncDetails unitdets.TentSyncDetails
		_, err := utils.DecodeWithType(sync.Additional.(map[string]interface{}), &syncDetails)
		if err != nil {
			return nil, err
		}

		incomingOffers := make(map[string]TentOffer, len(syncDetails.IncomingOffers))
		for uid, offer := range syncDetails.IncomingOffers {
			incomingOffers[uid] = TentOffer(offer)
		}
		outgoingOffers := make(map[string]TentOffer, len(syncDetails.OutgoingOffers))
		for uid, offer := range syncDetails.OutgoingOffers {
			outgoingOffers[uid] = TentOffer(offer)
		}

		return &TentAdditional{
			IncomingOffers: incomingOffers,
			OutgoingOffers: outgoingOffers,
		}, nil
	})
}
