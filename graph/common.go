package graph

func getValueIfNotNull[K comparable](pointer *K) (value K) {
	if pointer != nil {
		return *pointer
	}

	return
}
