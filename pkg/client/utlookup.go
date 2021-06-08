package client

// UnitTypeLookup allows looking up the UIDs of smart objects by their
// unit type
type UnitTypeLookup map[string]map[string]struct{}

// Add the given smart object to this index
func (l UnitTypeLookup) Add(obj *SmartObject) {
	set, found := l[obj.UnitType]
	if !found {
		set = make(map[string]struct{})
		l[obj.UnitType] = set
	}
	set[obj.GameObject.UID] = struct{}{}
}

// Remove the given smart object from this index
func (l UnitTypeLookup) Remove(obj *SmartObject) {
	set, found := l[obj.UnitType]
	if !found {
		panic("smart object not in lookup")
	}
	if _, ok := set[obj.GameObject.UID]; !ok {
		panic("smart object not in lookup")
	}

	delete(set, obj.GameObject.UID)
	if len(set) == 0 {
		delete(l, obj.UnitType)
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
