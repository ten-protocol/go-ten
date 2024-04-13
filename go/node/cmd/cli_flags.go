package main

// Flags unique to Node deployment further service flags in `./go/config/flag_usage.go
const (
	NodeNameFlag     = "nodeName"
	IsSGXEnabledFlag = "isSGXEnabled"
	PccsAddrFlag     = "pccsAddr"
	HostImage        = "hostImage"
	EnclaveImage     = "enclaveImage"
	EdgelessDBImage  = "edgelessDBImage"
)

func nodeFlagUsageMap() map[string]string {
	return map[string]string{
		NodeNameFlag:     "Common name for containers and reference",
		IsSGXEnabledFlag: "Use SGX or simulation",
		PccsAddrFlag:     "SGX attestation address",
		HostImage:        "Docker image for host service",
		EnclaveImage:     "Docker image for enclave service",
		EdgelessDBImage:  "Docker image for edgeless DB (enclave persistence)",
	}
}
