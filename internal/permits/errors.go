package permits

import "errors"

var InvalidTokenLifetime = errors.New("invalid token lifetime")
var InvalidSessionLifetime = errors.New("invalid session lifetime")
