package main

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"log"
	"net/http"
)

func createHttpServer(config Config) {
	s, err := rest.NewServer(rest.RestConf{
		Host: "0.0.0.0",
		Port: 8080,
	})
	if err != nil {
		log.Printf("web server error: %s\n", err)
	}

	s.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(writer http.ResponseWriter, request *http.Request) {
			httpx.OkJson(writer, "App is running.")
		},
	})

	s.AddRoutes([]rest.Route{
		{
			Method: http.MethodGet,
			Path:   "/config",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				httpx.OkJson(writer, config)
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/config/:ipv4",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				type Params struct {
					Ipv4 string `path:"ipv4"`
				}
				var params Params
				httpx.Parse(request, &params)
				for _, item := range config.Remote {
					if params.Ipv4 == item.Ipv4 {
						httpx.OkJson(writer, item)
						return
					}
				}
				httpx.OkJson(writer, "not found ipv4")
			},
		},
		{
			Method: http.MethodPost,
			Path:   "/config",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				type Params struct {
					Ipv4 string `json:"ipv4"`
					Ipv6 string `json:"ipv6"`
				}
				var params Params
				httpx.Parse(request, &params)
				config.Remote = append(config.Remote, params)
				httpx.OkJson(writer, config)
			},
		},
		{
			Method: http.MethodPut,
			Path:   "/config",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				type Params struct {
					Ipv4 string `json:"ipv4"`
					Ipv6 string `json:"ipv6"`
				}
				var params Params
				httpx.Parse(request, &params)
				for index, item := range config.Remote {
					if item.Ipv4 == params.Ipv4 {
						config.Remote[index].Ipv6 = params.Ipv6
					}
				}
				httpx.OkJson(writer, config)
			},
		},
		{
			Method: http.MethodDelete,
			Path:   "/config/:ipv4",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				type Params struct {
					Ipv4 string `path:"ipv4"`
				}
				var params Params
				httpx.Parse(request, &params)
				for index := range config.Remote {
					if params.Ipv4 == config.Remote[index].Ipv4 {
						if index == len(config.Remote)-1 {
							config.Remote = config.Remote[:index]
						} else {
							config.Remote = append(config.Remote[:index], config.Remote[index+1:]...)
						}
						break
					}
				}
				httpx.OkJson(writer, config)
			},
		},
	})

	defer s.Stop()
	s.Start()
}
