package physics

type ContactCache interface {
	Add(index1 int, index2 int)
	Remove(index1 int, index2 int)
	Contains(index1 int, index2 int) bool
	MarkContactsAsOld()
	CleanOldContacts()
}

func getHash(index1, index2 int) int {
	return index1 | (index2 << 16)
}

type ContactCacheImpl struct {
	contacts map[int64]map[int64]bool
}

func NewContactCache() ContactCache {
	return &ContactCacheImpl{make(map[int64]map[int64]bool)}
}

func (cc *ContactCacheImpl) Add(index1 int, index2 int) {
	cache1, ok := cc.contacts[pair.object1.id]
	if !ok {
		cache1 = make(map[int64]bool)
		cc.contacts[pair.object1.id] = cache1
	}
	//TODO:
}

func (cc *ContactCacheImpl) Remove(index1 int, index2 int) {

}

func (cc *ContactCacheImpl) Contains(index1 int, index2 int) bool {
	return false
}

func (cc *ContactCacheImpl) MarkContactsAsOld() {

}

func (cc *ContactCacheImpl) CleanOldContacts() {

}
