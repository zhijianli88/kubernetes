// This is a generated file. Do not edit directly.

module k8s.io/csi-translation-lib

go 1.15

require (
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/klog/v2 v2.2.0
)

replace (
	github.com/google/go-cmp => github.com/google/go-cmp v0.4.0
	github.com/stretchr/testify => github.com/stretchr/testify v1.4.0
	k8s.io/api => ../api
	k8s.io/apimachinery => ../apimachinery
	k8s.io/csi-translation-lib => ../csi-translation-lib
)
