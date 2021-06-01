package srvpkts

// SmartObjectSync provides all the information required to know about a
// SmartObject for a client who has never seen it before.
type SmartObjectSync struct {
	GameObjectSync

	// UnitType is the identifier of the Unit this smart object is. Generally
	// AIs will use this to determine what orders are available for the unit,
	// e.g, a "market" unit might have a "buy" order.
	UnitType string `mapstructure:"unit_type" json:"unit_type"`

	// CurrentHealth is the current amount of health this object has. This
	// value may be a rounded representation of the objects true health.
	CurrentHealth int `mapstructure:"current_health" json:"current_health"`

	// MaxHealth is the maximum amount of health this object can have.
	MaxHealth int `mapstructure:"max_health" json:"max_health"`

	// ControllingTeam is the team which controls this object.
	ControllingTeam int `mapstructure:"controlling_team" json:"controlling_team"`

	// ControllingRole is the role required to control this object. See core.Role
	// and core.RoleToName.
	ControllingRole string `mapstructure:"controlling_role" json:"controlling_role"`

	// Additional contains the SyncInfo on the Unit controlling this smart object,
	// and depends on the type of unit.
	Additional interface{} `mapstructure:"additional" json:"additional"`
}
