package ml

import "fmt"

func Render(c Componer) ([]*Element, error) {
	currentElem, ok := compoElements[c]
	if !ok {
		return nil, fmt.Errorf("can't render not mounted component: %T %+v", c, c)
	}

	rendered, err := parseTemplate(c.Render(), c)
	if err != nil {
		return nil, err
	}

	newElem, err := decodeString(rendered)
	if err != nil {
		return nil, err
	}

	_, dirtyElems, err := mergeElements(currentElem, newElem)
	return dirtyElems, err
}

func mergeElements(current *Element, new *Element) (bool, []*Element, error) {
	if current.Name != new.Name || !current.Attributes.equals(new.Attributes) || len(current.Children) != len(new.Children) {
		switch {
		case current.tagType == htmlTag && new.tagType == htmlTag:
			if err := mergeHTMLHTML(current, new); err != nil {
				return false, nil, err
			}

			return false, []*Element{current}, nil

		case current.tagType == htmlTag && new.tagType == componentTag:
		case current.tagType == htmlTag && new.tagType == textTag:
		case current.tagType == componentTag && new.tagType == componentTag:
		case current.tagType == componentTag && new.tagType == htmlTag:
		case current.tagType == componentTag && new.tagType == textTag:
		case current.tagType == textTag && new.tagType == textTag:
		case current.tagType == textTag && new.tagType == htmlTag:
		case current.tagType == textTag && new.tagType == componentTag:
		}

	}

	var dirtyChildElems []*Element
	isDirty := false

	for i, c := range current.Children {
		isParentDirty, dirtyElems, err := mergeElements(c, new.Children[i])

		if err != nil {
			return false, nil, err
		}

		if !isDirty {
			isDirty = isParentDirty
		}

		dirtyChildElems = append(dirtyChildElems, dirtyElems...)
	}

	if isDirty {
		return false, []*Element{current}, nil
	}

	return false, dirtyChildElems, nil
}

func mergeHTMLHTML(current *Element, new *Element) error {
	current.Name = new.Name
	current.Attributes = new.Attributes
	current.tagType = new.tagType

	for _, c := range current.Children {
		if err := dismount(c); err != nil {
			return err
		}
	}

	for _, c := range new.Children {
		c.Parent = current

		if err := mount(c, current.Component, current.Context); err != nil {
			return err
		}
	}

	current.Children = new.Children
	return nil
}
