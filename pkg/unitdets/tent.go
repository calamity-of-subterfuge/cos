package unitdets

// TentOffer describes a single offer in a tent. The offer is visible
// on both the initiating and target teams tents. You can only see
// offers on your own tent.
type TentOffer struct {
	// InitiatingTeam is the team initiating the offer.
	InitiatingTeam int `json:"initiating_team" mapstructure:"initiating_team"`

	// TargetTeam is the team recieving the offer
	TargetTeam int `json:"target_team" mapstructure:"target_team"`

	// Offer is the offer of resources being made, which will be taken from the
	// initiating team and given to the target team if the offer is accepted.
	Offer map[string]int `json:"offer" mapstructure:"offer"`

	// Request is the request of resources, which will be taken from the target
	// team and given to the initiating team if the offer is accepted.
	Request map[string]int `json:"request" mapstructure:"request"`
}

// TentSyncDetails is the information provided under Additional for a tent unit
// during a sync event, typically a GameSync or SmartObjectAdded. It provides
// all the custom information about the tent. If the Tent is not owned by the
// player then the offers will always be empty maps.
type TentSyncDetails struct {
	// IncomingOffers are the offers from other teams which have been made
	// to the team owning this tent. The keys are uids of the offers.
	IncomingOffers map[string]TentOffer `json:"incoming_offers" mapstructure:"incoming_offers"`

	// OutgoingOffers are the offers from this team which have been made to
	// other teams. The keys are the uids of the offers.
	OutgoingOffers map[string]TentOffer `json:"outgoing_offers" mapstructure:"outgoing_offers"`
}

// TentUpdateDetails contains the information under Additional for a tent unit
// during an update event, typically SmartObjectUpdate.
type TentUpdateDetails struct {
	// RemovedIncomingOffers is the slice of incoming offer uids which are no
	// longer in this tent, either because they were responded to or withdrawn
	RemovedIncomingOffers []string

	// RemovedOutgoings is the slice of outgoing offer uids which are no longer
	// in this tent, either because they were responded to or withdrawn
	RemovedOutgoingOffers []string

	// AddedIncomingOffers contains all the new incoming offers on this tent,
	// where the keys are uids
	AddedIncomingOffers map[string]TentOffer

	// AddedOutgoingOffers contains all the new outgoing offers on this tent,
	// where the keys are uids
	AddedOutgoingOffers map[string]TentOffer
}
