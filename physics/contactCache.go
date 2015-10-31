package physics

type ContactCache interface {
	Add(pair PhysicsPair)
	Remove(pair PhysicsPair)
	Contains(pair PhysicsPair) bool
	MarkContactsAsOld()
	CleanOldContacts()
}

type ContactCacheImpl struct {
	contacts map[int64]map[int64]bool
}

func NewContactCache() ContactCache {
	return &ContactCacheImpl{make(map[int64]map[int64]bool)}
}

func (cc *ContactCacheImpl) Add(pair PhysicsPair) {
	cache1, ok := cc.contacts[pair.object1.id]
	if !ok {
		cache1 = make(map[int64]bool)
		cc.contacts[pair.object1.id] = cache1
	}
	//TODO:
}

func (cc *ContactCacheImpl) Remove(pair PhysicsPair) {

}

func (cc *ContactCacheImpl) Contains(pair PhysicsPair) bool {
	return false
}

func (cc *ContactCacheImpl) MarkContactsAsOld() {

}

func (cc *ContactCacheImpl) CleanOldContacts() {

}
