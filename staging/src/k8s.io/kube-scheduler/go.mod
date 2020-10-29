// This is a generated file. Do not edit directly.

module k8s.io/kube-scheduler

go 1.15

require (
	github.com/google/go-cmp v0.5.2
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/component-base v0.0.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/google/go-cmp => github.com/google/go-cmp v0.4.0
	github.com/stretchr/testify => github.com/stretchr/testify v1.4.0
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
	k8s.io/api => ../api
	k8s.io/apimachinery => ../apimachinery
	k8s.io/client-go => ../client-go
	k8s.io/component-base => ../component-base
	k8s.io/kube-scheduler => ../kube-scheduler
)
