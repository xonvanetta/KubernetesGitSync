package sops

import (
	"fmt"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
	"github.com/getsops/sops/v3/age"
	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/cmd/sops/formats"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/version"
)

//TODO remove sops/v3/cmd imports/prettify this for kubernetes use case instead of cli

func Encrypt(yamlDocument []byte, ageRecipients string) ([]byte, error) {
	inputStore := common.StoreForFormat(formats.Yaml, config.NewStoresConfig())
	branches, err := inputStore.LoadPlainFile(yamlDocument)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error unmarshalling file: %s", err), codes.CouldNotReadInputFile)
	}

	keyGroups, err := keyGroups(ageRecipients)
	if err != nil {
		return nil, err
	}

	tree := sops.Tree{
		Branches: branches,
		Metadata: sops.Metadata{
			Version:        version.Version,
			EncryptedRegex: "^(data|stringData)$",
			KeyGroups:      keyGroups,
		},
		//FilePath: path,
	}
	dataKey, errs := tree.GenerateDataKey()
	if len(errs) > 0 {
		err = fmt.Errorf("Could not generate data key: %s", errs)
		return nil, err
	}

	err = common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey,
		Tree:    &tree,
		Cipher:  aes.NewCipher(),
	})
	if err != nil {
		return nil, err
	}

	//Todo formatter is not really the same as kube printer
	encryptedFile, err := common.StoreForFormat(formats.Yaml, config.NewStoresConfig()).EmitEncryptedFile(tree)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
	}

	return encryptedFile, nil
}

func keyGroups(ageRecipients string) ([]sops.KeyGroup, error) {
	var ageMasterKeys []keys.MasterKey
	ageKeys, err := age.MasterKeysFromRecipients(ageRecipients)
	if err != nil {
		return nil, err
	}
	for _, k := range ageKeys {
		ageMasterKeys = append(ageMasterKeys, k)
	}

	var group sops.KeyGroup
	group = append(group, ageMasterKeys...)
	return []sops.KeyGroup{group}, nil
}
