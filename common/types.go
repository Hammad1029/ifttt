package common

type IntIota int

type Manipulatable interface {
	Manipulate(dependencies map[IntIota]any) error
}

type Validatable interface {
	Validate() error
}
