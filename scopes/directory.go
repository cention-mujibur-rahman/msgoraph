package scopes

var (
	// ApplicationDirectoryReadAll Read directory data
	ApplicationDirectoryReadAll = Scope{
		AdminConsentRequired: true,
		DisplayString:        "Read directory data",
		Description:          "Allows the app to read data in your organization's directory, such as users, groups and apps, without a signed-in user.",
		Permission:           "Directory.Read.All",
		Type:                 PermissionTypeApplication,
	}
	// ApplicationDirectoryReadWriteAll Read and write directory data
	ApplicationDirectoryReadWriteAll = Scope{
		AdminConsentRequired: true,
		DisplayString:        "Read and write directory data",
		Description:          "Allows the app to read and write data in your organization's directory, such as users, and groups, without a signed-in user. Does not allow user or group deletion.",
		Permission:           "Directory.ReadWrite.All",
		Type:                 PermissionTypeApplication,
	}
	// DelegatedDirectoryReadAll Read directory data
	DelegatedDirectoryReadAll = Scope{
		AdminConsentRequired: true,
		DisplayString:        "Read directory data",
		Description:          "Allows the app to read data in your organization's directory, such as users, groups and apps. Note: Users may consent to applications that require this permission if the application is registered in their own organization’s tenant.",
		Permission:           "Directory.Read.All",
		Type:                 PermissionTypeDelegated,
	}
	// DelegatedDirectoryReadWriteAll Read and write directory data
	DelegatedDirectoryReadWriteAll = Scope{
		AdminConsentRequired: true,
		DisplayString:        "Read and write directory data",
		Description:          "Allows the app to read and write data in your organization's directory, such as users, and groups. It does not allow the app to delete users or groups, or reset user passwords.",
		Permission:           "Directory.ReadWrite.All",
		Type:                 PermissionTypeDelegated,
	}
	// DelegatedDirectoryAccessAsUser Access directory as the signed-in user
	DelegatedDirectoryAccessAsUser = Scope{
		AdminConsentRequired: true,
		DisplayString:        "Access directory as the signed-in user",
		Description:          "Allows the app to have the same access to information in the directory as the signed-in user.",
		Permission:           "Directory.AccessAsUser.All",
		Type:                 PermissionTypeDelegated,
	}
)
