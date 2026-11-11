package security

import "fmt"

type AuthManager struct {
    Mode string
}

func NewAuthManager(mode string) *AuthManager {
    return &AuthManager{Mode: mode}
}

func (a *AuthManager) Validate(user string) error {
    if a == nil {
        return nil
    }
    if a.Mode == "local_only" && user == "" {
        return fmt.Errorf("local-only mode requires an authenticated session")
    }
    return nil
}
