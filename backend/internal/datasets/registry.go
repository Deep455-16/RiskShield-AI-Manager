package datasets

type Registry struct {
	Adapters map[string]DatasetAdapter
}

func NewRegistry(dataDir string) *Registry {
	return &Registry{
		Adapters: map[string]DatasetAdapter{
			"sfindset": &SFinDSetAdapter{BaseDir: dataDir + "/sfindset"},
			"global_bank": &GlobalBankAdapter{BaseDir: dataDir + "/global_bank"},
			"bank_marketing": &BankMarketingAdapter{BaseDir: dataDir + "/bank_marketing"},
		},
	}
}

func (r *Registry) ListDatasets() []DatasetMetadata {
	var meta []DatasetMetadata
	for _, a := range r.Adapters {
		meta = append(meta, a.GetMetadata())
	}
	return meta
}

func (r *Registry) GetAdapter(id string) DatasetAdapter {
	return r.Adapters[id]
}
