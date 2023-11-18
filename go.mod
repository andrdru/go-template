module github.com/andrdru/go-template

go 1.21.1

replace github.com/andrdru/go-template/graceful v0.0.0 => ./graceful

replace github.com/andrdru/go-template/configs v0.0.0 => ./configs

replace github.com/andrdru/go-template/tx v0.0.0 => ./tx

require (
	github.com/andrdru/go-template/configs v0.0.0
	github.com/andrdru/go-template/graceful v0.0.0
	github.com/andrdru/go-template/tx v0.0.0
	github.com/google/uuid v1.3.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/mailru/easyjson v0.7.7
	github.com/prometheus/client_golang v1.17.0
	golang.org/x/crypto v0.14.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gomodule/redigo v1.8.9 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_model v0.4.1-0.20230718164431-9a2bf3000d16 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)
