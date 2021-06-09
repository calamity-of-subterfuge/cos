package client

// UnitTypeLookup allows looking up the UIDs of smart objects by their
// unit type
type UnitTypeLookup map[string]map[string]struct{}

// Add the given smart object to this index
func (l UnitTypeLookup) Add(obj *SmartObject) {
	l.AddUID(obj.UnitType, obj.GameObject.UID)
}

// AddUID adds the smart object with the given uid to this index
func (l UnitTypeLookup) AddUID(unitType string, uid string) {
	set, found := l[unitType]
	if !found {
		set = make(map[string]struct{})
		l[unitType] = set
	}
	set[uid] = struct{}{}
}

// Remove the given smart object from this index
func (l UnitTypeLookup) Remove(obj *SmartObject) {
	l.RemoveUID(obj.UnitType, obj.GameObject.UID)
}

// RemoveUID removes the smart object with the given type and uid from
// this index
func (l UnitTypeLookup) RemoveUID(unitType string, uid string) {
	set, found := l[unitType]
	if !found {
		panic("smart object not in lookup")
	}
	if _, ok := set[uid]; !ok {
		panic("smart object not in lookup")
	}

	delete(set, uid)
	if len(set) == 0 {
		delete(l, unitType)
	}
}

// Each calls the operator on every unit of the given type
func (l UnitTypeLookup) Each(unitType string, operator func(string)) {
	set, found := l[unitType]
	if !found {
		return
	}

	for uid := range set {
		operator(uid)
	}
}
