package signatures

func (ms *MapSignatures) FindAll() []string {
	keys := make([]string, 0, len(ms.hm))
	for key := range ms.hm {
		keys = append(keys, key)
	}

	return keys
}

