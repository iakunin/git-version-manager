package repository

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

const tagPrefix = "refs/tags/"

type Repository struct {
	dir  string
	repo *git.Repository
}

func Open(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return &Repository{dir: path, repo: repo}, nil
}

func (r *Repository) GetRawTags() ([]*string, error) {
	tagRefs, err := r.repo.Tags()

	if err != nil {
		log.Errorf("unable to get rawTags from repository: %s", err)
		return nil, err
	}

	var rawTags []*string
	err = tagRefs.ForEach(func(t *plumbing.Reference) error {
		name := t.Name().String()

		if !strings.HasPrefix(name, tagPrefix) {
			errMessage := fmt.Sprintf("tag must starts with predefined preffix (%s): %s", tagPrefix, name)
			log.Error(errMessage)
			return errors.New(errMessage)
		}

		replaced := strings.Replace(name, tagPrefix, "", 1)
		rawTags = append(rawTags, &replaced)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return rawTags, nil
}

func (r *Repository) tagExists(tag string) bool {
	tagFoundErr := "tag was found"
	tags, err := r.repo.Tags()
	if err != nil {
		log.Errorf("get tags error: %s", err)
		return false
	}
	res := false
	err = tags.ForEach(func(t *plumbing.Reference) error {
		if t.Name().String() == (tagPrefix + tag) {
			res = true
			return fmt.Errorf(tagFoundErr)
		}
		return nil
	})

	if err != nil && err.Error() != tagFoundErr {
		log.Errorf("iterate tags error: %s", err)
		return false
	}
	return res
}

func (r *Repository) SetTag(tag string) error {
	if r.tagExists(tag) {
		err := fmt.Sprintf("tag %s already exists", tag)
		log.Infof(err)
		return errors.New(err)
	}
	log.Infof("Set tag %s", tag)
	h, err := r.repo.Head()
	if err != nil {
		log.Errorf("get HEAD error: %s", err)
		return err
	}

	_, err = r.repo.CreateTag(tag, h.Hash(), nil)

	if err != nil {
		log.Errorf("create tag error: %s", err)
		return err
	}

	return nil
}

func publicKey() (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys
	sshPath := os.Getenv("HOME") + "/.ssh/github_rsa"
	sshKey, _ := ioutil.ReadFile(sshPath)
	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKey), "")
	if err != nil {
		return nil, err
	}
	return publicKey, err
}

func (r *Repository) PushTags() error {
	auth, _ := publicKey()

	options := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
		Auth:       auth,
	}

	err := r.repo.Push(options)

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Info("origin remote was up to date, no push done")
			return nil
		}
		log.Errorf("push to remote origin error: %s", err)
		return err
	}

	return nil
}

func (r *Repository) FetchTags() error {
	// @TODO: implement me
	return nil
}
