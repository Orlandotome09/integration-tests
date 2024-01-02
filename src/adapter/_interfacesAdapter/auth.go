package interfacesAdapter

type AuthRepository interface {
	GetAccessToken() (string, error)
}
