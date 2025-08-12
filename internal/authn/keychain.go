package authn

// This file previously contained a complex go-containerregistry keychain implementation.
// With the migration to containers/image, authentication is now handled directly
// through types.SystemContext.AuthFilePath, so the complex keychain is no longer needed.
//
// If any authentication utilities are needed in the future for containers/image,
// they can be added here.