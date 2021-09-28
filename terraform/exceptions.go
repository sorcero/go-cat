package terraform

import "errors"

var InvalidConfigError = errors.New("could not access git repository, make sure the git_url, git_username and git_password are correct")
