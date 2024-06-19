package types

// GenesisState defines the module's genesis state.
type GenesisState struct {
	Param1 string `json:"param1"`
	Param2 int    `json:"param2"`
}

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() GenesisState {
	return GenesisState{
		Param1: "default",
		Param2: 0,
	}
}

// Validate performs basic genesis state validation.
func (gs GenesisState) Validate() error {
	// Add validation logic here.
	return nil
}
