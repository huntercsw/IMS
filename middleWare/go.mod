module middleWare

go 1.13

require (
	github.com/golang/protobuf v1.4.2 // indirect
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	google.golang.org/grpc v1.23.0
	imsPb v0.0.0
)

replace (
	github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	imsPb => ../imsPb
)
