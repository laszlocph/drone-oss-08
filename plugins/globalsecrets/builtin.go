package globalsecrets

import (
	"github.com/woodpecker-ci/woodpecker/model"
)

type builtin struct {
	store model.GlobalSecretStore
}

// New returns a new global secret service.
func New(store model.GlobalSecretStore) model.GlobalSecretService {
	return &builtin{store}
}

func (b *builtin) GlobalSecretFind(name string) (*model.GlobalSecret, error) {
	return b.store.GlobalSecretFind(name)
}

func (b *builtin) GlobalSecretList() ([]*model.GlobalSecret, error) {
	return b.store.GlobalSecretList()
}

func (b *builtin) GlobalSecretListBuild(build *model.Build) ([]*model.GlobalSecret, error) {
	return b.store.GlobalSecretList()
}

func (b *builtin) GlobalSecretCreate(in *model.GlobalSecret) error {
	return b.store.GlobalSecretCreate(in)
}

func (b *builtin) GlobalSecretUpdate(in *model.GlobalSecret) error {
	return b.store.GlobalSecretUpdate(in)
}

func (b *builtin) GlobalSecretDelete(name string) error {
	secret, err := b.store.GlobalSecretFind(name)
	if err != nil {
		return err
	}
	return b.store.GlobalSecretDelete(secret)
}
