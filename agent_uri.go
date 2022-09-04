package fastreq

func (a *Agent) Host(host string) *Agent {
	a.req.URI().SetHost(host)

	return a
}

func (a *Agent) HostBytes(host []byte) *Agent {
	a.req.URI().SetHostBytes(host)

	return a
}

func (a *Agent) QueryString(queryString string) *Agent {
	a.req.URI().SetQueryString(queryString)

	return a
}

func (a *Agent) QueryStringBytes(queryString []byte) *Agent {
	a.req.URI().SetQueryStringBytes(queryString)

	return a
}

func (a *Agent) BasicAuth(username, password string) *Agent {
	a.req.URI().SetUsername(username)
	a.req.URI().SetPassword(password)

	return a
}

func (a *Agent) BasicAuthBytes(username, password []byte) *Agent {
	a.req.URI().SetUsernameBytes(username)
	a.req.URI().SetPasswordBytes(password)

	return a
}
