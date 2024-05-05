package core

type ContainerListenerDef struct {
	description  string
	displayName  string
	listenerName string

	parameters map[string]string
}

func (cld *ContainerListenerDef) GetDescription() string {
	return cld.description
}

func (cld *ContainerListenerDef) SetDescription(description string) {
	cld.description = description
}

func (cld *ContainerListenerDef) GetDisplayName() string {
	return cld.displayName
}

func (cld *ContainerListenerDef) SetDisplayName(displayName string) {
	cld.displayName = displayName
}

func (cld *ContainerListenerDef) GetListenerName() string {
	return cld.listenerName
}

func (cld *ContainerListenerDef) SetListenerName(listenerName string) {
	cld.listenerName = listenerName
}

func (cld *ContainerListenerDef) GetParameterMap() map[string]string {
	return cld.parameters
}
