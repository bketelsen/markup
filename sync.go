package markup

// Sync synchronizes a component.
// It check all the elements associated with the component and performs changes if required.
// Returns the changed elements.
func Sync(c Componer) (changed []*Element, err error) {
	currentElem, err := ComponentRoot(c)
	if err != nil {
		return
	}

	rendered, err := render(c.Render(), c)
	if err != nil {
		return
	}

	newElem, err := Decode(rendered)
	if err != nil {
		return
	}

	_, changed, err = sync(currentElem, newElem)
	return
}

func sync(current *Element, new *Element) (parentChanged bool, changed []*Element, err error) {
	if current.Name != new.Name || !current.Attributes.equals(new.Attributes) || len(current.Children) != len(new.Children) {
		return syncElements(current, new)
	}

	currentChanged := false
	requireChanged := false

	var childChanged []*Element

	for i, child := range current.Children {
		if requireChanged, childChanged, err = sync(child, new.Children[i]); err != nil {
			return
		}

		if currentChanged {
			continue
		}

		currentChanged = requireChanged
		changed = append(changed, childChanged...)
	}

	if currentChanged {
		changed = []*Element{current}
	}
	return
}

func syncElements(current *Element, new *Element) (parentChanged bool, changed []*Element, err error) {
	switch {
	case current.Type == HTML && new.Type != HTML:
		return syncHTMLWithComponentOrText(current, new)

	case current.Type == Component && new.Type == Component:
		return syncComponentWithComponent(current, new)

	case current.Type == Component && new.Type != Component:
		return syncComponentWithTextOrHTML(current, new)

	case current.Type == Text && new.Type == Text:
		return syncTextWithText(current, new)

	case current.Type == Text && new.Type != Text:
		return syncTextWithHTMLOrComponent(current, new)

	default:
		return syncHTMLWithHTML(current, new)
	}
}

func syncHTMLWithHTML(current *Element, new *Element) (parentChanged bool, changed []*Element, err error) {
	for _, c := range current.Children {
		dismount(c)
	}

	for _, c := range new.Children {
		c.Parent = current

		if err = mount(c, current.Component, current.ContextID); err != nil {
			return
		}
	}

	current.Name = new.Name
	current.Attributes = new.Attributes
	current.Children = new.Children

	changed = []*Element{current}
	return
}

func syncHTMLWithComponentOrText(current *Element, new *Element) (parentChanged bool, changed []*Element, err error) {
	dismount(current)

	current.Name = new.Name
	current.Attributes = new.Attributes
	current.Type = new.Type
	current.Children = nil

	parentChanged = true
	err = mount(current, current.Component, current.ContextID)
	return
}

func syncComponentWithComponent(current *Element, new *Element) (parentChanged bool, changed []*Element, err error) {
	current.Attributes = new.Attributes

	if current.Name == new.Name {
		if err = updateComponentFields(current.Component, new.Attributes); err != nil {
			return
		}

		changed, err = Sync(current.Component)
		return
	}

	dismount(current)
	current.Name = new.Name

	parentChanged = true
	err = mountComponent(current, current.ContextID)
	return
}

func syncComponentWithTextOrHTML(current *Element, new *Element) (parentChanged bool, changed []*Element, err error) {
	dismount(current)

	current.Name = new.Name
	current.Attributes = new.Attributes
	current.Type = new.Type

	parentChanged = true
	err = mount(current, current.Parent.Component, current.Parent.ContextID)
	return
}

func syncTextWithText(current *Element, new *Element) (parentChanged bool, changed []*Element, err error) {
	current.Attributes = new.Attributes
	parentChanged = true
	return
}

func syncTextWithHTMLOrComponent(current *Element, new *Element) (parentChanged bool, changed []*Element, err error) {
	current.Name = new.Name
	current.Attributes = new.Attributes
	current.Type = new.Type
	current.Children = new.Children

	parentChanged = true
	err = mount(current, current.Parent.Component, current.Parent.ContextID)
	return
}
