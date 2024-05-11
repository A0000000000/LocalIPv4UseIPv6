package main

type Local struct {
	IP   string `json:"ip"`
	Mask string `json:"mask"`
}

type Remote []struct {
	Ipv4 string `json:"ipv4"`
	Ipv6 string `json:"ipv6"`
}

type Config struct {
	Local  `json:"local"`
	Remote `json:"remote"`
}
