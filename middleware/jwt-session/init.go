package jwt

// Manager Manager
var Manager *SessionManager

// Init Init
func Init(opts ...Options) {
	Manager = NewSessionManger(opts...)

}
