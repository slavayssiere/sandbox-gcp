module github.com/slavayssiere/sandbox-gcp/app-grpc/consumer/aggregator

require (
	cloud.google.com/go v0.34.0
	github.com/go-redis/redis v6.15.1+incompatible
	github.com/googleapis/gax-go v2.0.2+incompatible // indirect
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.6.2
	github.com/prometheus/client_golang v0.9.2
	github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier v0.0.0-20190115125349-e0017caf23f7
	golang.org/x/oauth2 v0.0.0-20181203162652-d668ce993890
	google.golang.org/api v0.1.0
	google.golang.org/genproto v0.0.0-20181202183823-bd91e49a0898
	google.golang.org/grpc v1.17.0
)
