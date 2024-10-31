package user

func NewLDAPCreateUserRequest(username, password, descriptoin string) *CreateUserRequest {
	return &CreateUserRequest{
		Provider:    PROVIDER_LDAP,
		CreateType:  CREATE_TYPE_REGISTRY,
		UserName:    username,
		Password:    password,
		Description: descriptoin,
	}
}
