package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"
)

func Test_Unix_CreateService(t *testing.T) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", DefaultAddr)
			},
		},
	}
	serviceSpec := `{
  "Name": "web",
  "TaskTemplate": {
    "ContainerSpec": {
      "Image": "nginx:latest",
      "Mounts": [
        {
          "ReadOnly": true,
          "Source": "web-data",
          "Target": "/usr/share/nginx/html",
          "Type": "volume",
          "VolumeOptions": {
            "DriverConfig": {},
            "Labels": {
              "com.example.something": "something-value"
            }
          }
        }
      ],
      "Hosts": [
        "10.10.10.10 host1",
        "ABCD:EF01:2345:6789:ABCD:EF01:2345:6789 host2"
      ],
      "User": "33",
      "DNSConfig": {
        "Nameservers": [
          "8.8.8.8"
        ],
        "Search": [
          "example.org"
        ],
        "Options": [
          "timeout:3"
        ]
      },
      "Secrets": [
      ]
    },
    "LogDriver": {
      "Name": "json-file",
      "Options": {
        "max-file": "3",
        "max-size": "10M"
      }
    },
    "Placement": {},
    "Resources": {
      "Limits": {
        "MemoryBytes": 104857600
      },
      "Reservations": {}
    },
    "RestartPolicy": {
      "Condition": "on-failure",
      "Delay": 10000000000,
      "MaxAttempts": 10
    }
  },
  "Mode": {
    "Replicated": {
      "Replicas": 4
    }
  },
  "UpdateConfig": {
    "Parallelism": 2,
    "Delay": 1000000000,
    "FailureAction": "pause",
    "Monitor": 15000000000,
    "MaxFailureRatio": 0.15
  },
  "RollbackConfig": {
    "Parallelism": 1,
    "Delay": 1000000000,
    "FailureAction": "pause",
    "Monitor": 15000000000,
    "MaxFailureRatio": 0.15
  },
  "EndpointSpec": {
    "Ports": [
      {
        "Protocol": "tcp",
        "PublishedPort": 8080,
        "TargetPort": 80
      }
    ]
  },
  "Labels": {
    "foo": "bar"
  }
}`
	res, err := httpc.Post("http://unix/CreateService",
		"application/x-www-form-urlencoded",
		strings.NewReader(serviceSpec))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(res.Status)
	for k, v := range res.Header {
		fmt.Println(k, ": ", v)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	strBody := string(body)

	fmt.Println(strBody)
}
func Test_Unix_UpdateService(t *testing.T) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", DefaultAddr)
			},
		},
	}
	serviceSpec := `{
  "Name": "web",
  "TaskTemplate": {
    "ContainerSpec": {
      "Image": "nginx:latest",
      "Mounts": [
        {
          "ReadOnly": true,
          "Source": "web-data",
          "Target": "/usr/share/nginx/html",
          "Type": "volume",
          "VolumeOptions": {
            "DriverConfig": {},
            "Labels": {
              "com.example.something": "something-value"
            }
          }
        }
      ],
      "Hosts": [
        "10.10.10.10 host1",
        "ABCD:EF01:2345:6789:ABCD:EF01:2345:6789 host2"
      ],
      "User": "33",
      "DNSConfig": {
        "Nameservers": [
          "8.8.8.8"
        ],
        "Search": [
          "example.org"
        ],
        "Options": [
          "timeout:3"
        ]
      },
      "Secrets": [
      ]
    },
    "LogDriver": {
      "Name": "json-file",
      "Options": {
        "max-file": "3",
        "max-size": "10M"
      }
    },
    "Placement": {},
    "Resources": {
      "Limits": {
        "MemoryBytes": 104857600
      },
      "Reservations": {}
    },
    "RestartPolicy": {
      "Condition": "on-failure",
      "Delay": 10000000000,
      "MaxAttempts": 10
    }
  },
  "Mode": {
    "Replicated": {
      "Replicas": 2
    }
  },
  "UpdateConfig": {
    "Parallelism": 2,
    "Delay": 1000000000,
    "FailureAction": "pause",
    "Monitor": 15000000000,
    "MaxFailureRatio": 0.15
  },
  "RollbackConfig": {
    "Parallelism": 1,
    "Delay": 1000000000,
    "FailureAction": "pause",
    "Monitor": 15000000000,
    "MaxFailureRatio": 0.15
  },
  "EndpointSpec": {
    "Ports": [
      {
        "Protocol": "tcp",
        "PublishedPort": 8080,
        "TargetPort": 80
      }
    ]
  },
  "Labels": {
    "foo": "bar"
  }
}`
	serviceID := "toqjci8q0jsh"
	res, err := httpc.Post("http://unix/UpdateService/"+serviceID,
		"application/x-www-form-urlencoded",
		strings.NewReader(serviceSpec))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(res.Status)
	for k, v := range res.Header {
		fmt.Println(k, ": ", v)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	strBody := string(body)

	fmt.Println(strBody)
}
