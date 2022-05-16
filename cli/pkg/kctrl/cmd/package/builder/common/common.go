package common

type Step interface {
	PreInteract() error
	PostInteract() error
	Interact() error
}

func ValidateInputYesOrNo(input string) (bool, bool) {
	if input == "y" || input == "Y" {
		return true, true
	} else if input == "n" || input == "N" {
		return false, true
	}
	return false, false
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
	err = step.PostInteract()
	if err != nil {
		return err
	}
	return nil
}
