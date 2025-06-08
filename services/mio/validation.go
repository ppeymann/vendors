package mio

import "errors"

// ErrValidationFailed is returned request json body object is not stratified schema rules
var ErrValidationFailed = errors.New("request rejected, validation error")
var ErrBadUploadRequest = errors.New("upload request is illegal")
var ErrFileTypeNotSupported = errors.New("upload file type is not supported")
var ErrInvalidAccessToken = errors.New("invalid photo access token")
