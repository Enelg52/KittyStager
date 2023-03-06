package priv

type Privilege struct {
	Name        string
	Description string
	Enabled     bool
}

func NewPrivilege(name, description string, enable bool) *Privilege {
	return &Privilege{name, description, enable}
}

func (privilege *Privilege) GetName() string {
	return privilege.Name
}

func (privilege *Privilege) GetDescription() string {
	return privilege.Description
}

func (privilege *Privilege) GetEnable() bool {
	return privilege.Enabled
}
