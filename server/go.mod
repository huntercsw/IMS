module server

go 1.13

require (
	google.golang.org/grpc v1.30.0
	imsPb v0.0.0
	imsRegister v0.0.0
	middleWare v0.0.0
)

replace (
	github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	imsLb => ../imsLb
	imsPb => ../imsPb
	imsRegister => ../imsRegister
	middleWare => ../middleWare
)
