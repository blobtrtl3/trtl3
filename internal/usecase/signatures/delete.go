package signatures

func (ms *MapSignatures) Delete(key string) {
	delete(ms.hm, key)
}
