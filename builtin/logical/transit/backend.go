package transit

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf)
	be, err := b.Backend.Setup(conf)
	if err != nil {
		return nil, err
	}

	return be, nil
}

func Backend(conf *logical.BackendConfig) *backend {
	var b backend
	b.Backend = &framework.Backend{
		Paths: []*framework.Path{
			// Rotate/Config needs to come before Keys
			// as the handler is greedy
			b.pathConfig(),
			b.pathRotate(),
			b.pathRewrap(),
			b.pathKeys(),
			b.pathEncrypt(),
			b.pathDecrypt(),
			b.pathDatakey(),
		},

		Secrets: []*framework.Secret{},
	}

	b.lm = newLockManager(conf.System.CachingDisabled())

	return &b
}

type backend struct {
	*framework.Backend
	lm *lockManager
}
