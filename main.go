package main

import (
	"errors"
	"flag"
	repositoryModel "github.com/iakunin/git-version-manager/models/repository"
	tagModel "github.com/iakunin/git-version-manager/models/tag"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Panicf("Unable to get currenc workDir: '%s'", err)
	}

	repoDir := flag.String("repoDir", workDir, "git-repo directory (current workDir by default)")
	prefix := flag.String("prefix", "", "tag prefix")
	suffix := flag.String("suffix", "", "tag suffix")
	isVerbose := flag.Bool("verbose", false, "enable verbose logging")
	bumpStrategy := flag.String(
		"bumpStrategy",
		string(tagModel.Patch),
		"available values: `patch`, `minor`, `major`",
	)
	successfulMessage := flag.String(
		"successfulMessage",
		"",
		"message which will be written to console in the case of successful tagging",
	)
	flag.Parse()

	if *isVerbose {
		log.SetLevel(log.TraceLevel)
	}

	repository, err := repositoryModel.Open(*repoDir)
	if err != nil {
		log.Panicf("Unable to open the repo: '%s'", err)
	}

	rawTags, err := repository.GetRawTags()
	if err != nil {
		log.Panicf("Unable to getRawTags: '%s'", err)
	}

	for _, t := range rawTags {
		log.Debugf("This is a rawTag: '%s'", *t)
	}

	tags, err := createTags(rawTags, *prefix, *suffix)
	if err != nil {
		log.Panicf("Unable to createTags: '%s'", err)
	}
	for _, t := range tags {
		log.Debugf("This is a tag (before bump): '%s'\n", t.String())
	}

	maxTag, _ := findMaxTag(tags)
	log.Infof("This is a MAX tag: '%s'", maxTag.String())

	maxTag.Bump(tagModel.BumpStrategy(*bumpStrategy))

	log.Infof("This is a BUMPED tag: '%s'\n", maxTag.String())

	for _, t := range tags {
		log.Debugf("This is a tag (after bump): '%s'\n", t.String())
	}

	err = repository.SetTag(maxTag.String())
	if err != nil {
		log.Panicf("Unable to set tag: '%s'", err)
	}

	// @TODO: push tags using `repository.PushTags()`
	// @TODO: fetch tags using `repository.FetchTags()`

	if *successfulMessage != "" {
		log.Infoln(*successfulMessage)
	}
}

func createTags(rawTags []*string, prefix string, suffix string) ([]*tagModel.Tag, error) {
	var result []*tagModel.Tag

	for i := range rawTags {
		t, err := tagModel.New(*rawTags[i], prefix, suffix)
		if err != nil {
			log.Debugf("skipping a tag '%s'", *rawTags[i])
			continue
		}

		result = append(result, t)
	}

	if len(result) == 0 {
		result = append(result, tagModel.Empty(prefix, suffix))
	}

	return result, nil
}

func findMaxTag(tags []*tagModel.Tag) (*tagModel.Tag, error) {
	length := len(tags)

	if length == 0 {
		return nil, errors.New("got empty `tags` array")
	} else if length == 1 {
		return tags[0], nil
	}

	max := tags[0]

	for _, t := range tags[1:] {
		if max.LessThan(*t) {
			max = t
		}
	}

	return max, nil
}
