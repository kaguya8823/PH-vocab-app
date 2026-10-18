package router

router.HandleFunc(
	".api/vocabularies",
	handler.GetVocabularies,
).Methods("GET")