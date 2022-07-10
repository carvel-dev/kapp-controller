package step

type Step interface {
	PreInteract() error
	Interact() error
	PostInteract() error
}

func Run(step Step) error {
	err := step.PreInteract()
	if err != nil {
		return err
	}
	err = step.Interact()
	if err != nil {
		return err
	}
	return step.PostInteract()
}
