package dynamics

type ContactCache interface {
	Add(index1 int, index2 int)
	Remove(index1 int, index2 int)
	Contains(index1 int, index2 int) bool
	MarkContactsAsOld()
	CleanOldContacts()
	Clear()
}

func getHash(index1, index2 int) int {
	if index1 > index2 {
		return index1 | (index2 << 16)
	}
	return index2 | (index1 << 16)
}

type ContactCacheImpl struct {
	contacts map[int]bool
}

func NewContactCache() ContactCache {
	return &ContactCacheImpl{make(map[int]bool)}
}

func (cc *ContactCacheImpl) Add(index1 int, index2 int) {
	cc.contacts[getHash(index1, index2)] = false
}

func (cc *ContactCacheImpl) Remove(index1 int, index2 int) {
	delete(cc.contacts, getHash(index1, index2))
}

func (cc *ContactCacheImpl) Contains(index1 int, index2 int) bool {
	_, ok := cc.contacts[getHash(index1, index2)]
	return ok
}

func (cc *ContactCacheImpl) MarkContactsAsOld() {
	for key, _ := range cc.contacts {
		cc.contacts[key] = true
	}
}

func (cc *ContactCacheImpl) CleanOldContacts() {
	for key, value := range cc.contacts {
		if value {
			delete(cc.contacts, key)
		}
	}
}

func (cc *ContactCacheImpl) Clear() {
	for key, _ := range cc.contacts {
		delete(cc.contacts, key)
	}
}
