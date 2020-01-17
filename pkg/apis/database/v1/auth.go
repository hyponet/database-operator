package v1

type Secret struct {
	Name string `json:"name"`
}

type DatabaseAuth struct {
	RootPasswordSecret *Secret `json:"root_password_secret,omitempty"`
}
