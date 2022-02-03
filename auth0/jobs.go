package auth0

import "gopkg.in/auth0.v5/management"

type JobsAPI interface {
	Read(id string, opts ...management.RequestOption) (j *management.Job, err error)
	ImportUsers(j *management.Job, opts ...management.RequestOption) (err error)
}
