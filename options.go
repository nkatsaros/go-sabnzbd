package sabnzbd

type Option func(*Sabnzbd) error

func UseHttp() Option {
	return func(s *Sabnzbd) error {
		s.useHttp()
		return nil
	}
}

func UseHttps() Option {
	return func(s *Sabnzbd) error {
		s.useHttps()
		return nil
	}
}

func Addr(addr string) Option {
	return func(s *Sabnzbd) error {
		return s.setAddr(addr)
	}
}

func Path(path string) Option {
	return func(s *Sabnzbd) error {
		return s.setPath(path)
	}
}

func LoginAuth(username, password string) Option {
	return func(s *Sabnzbd) error {
		return s.setAuth(loginAuth{username, password})
	}
}

func ApikeyAuth(apikey string) Option {
	return func(s *Sabnzbd) error {
		return s.setAuth(apikeyAuth{apikey})
	}
}

func NoneAuth() Option {
	return func(s *Sabnzbd) error {
		return s.setAuth(noneAuth{})
	}
}
