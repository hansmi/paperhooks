package client

type ObjectPermissions struct {
	View   ObjectPermissionPrincipals `json:"view"`
	Change ObjectPermissionPrincipals `json:"change"`
}

type ObjectPermissionPrincipals struct {
	Users  []int64 `json:"users"`
	Groups []int64 `json:"groups"`
}
