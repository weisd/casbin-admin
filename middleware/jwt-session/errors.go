package jwt

// ValidError ValidError
type ValidError struct {
	Massage string
}

func (p *ValidError) Error() string {
	return p.Massage
}

func New() {

}
