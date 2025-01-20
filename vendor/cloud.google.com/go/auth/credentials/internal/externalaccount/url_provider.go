// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package externalaccount

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
<<<<<<< HEAD
	"log/slog"
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	"net/http"

	"cloud.google.com/go/auth/internal"
	"cloud.google.com/go/auth/internal/credsfile"
<<<<<<< HEAD
	"github.com/googleapis/gax-go/v2/internallog"
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
)

const (
	fileTypeText             = "text"
	fileTypeJSON             = "json"
	urlProviderType          = "url"
	programmaticProviderType = "programmatic"
	x509ProviderType         = "x509"
)

type urlSubjectProvider struct {
	URL     string
	Headers map[string]string
	Format  *credsfile.Format
	Client  *http.Client
<<<<<<< HEAD
	Logger  *slog.Logger
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
}

func (sp *urlSubjectProvider) subjectToken(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", sp.URL, nil)
	if err != nil {
		return "", fmt.Errorf("credentials: HTTP request for URL-sourced credential failed: %w", err)
	}

	for key, val := range sp.Headers {
		req.Header.Add(key, val)
	}
<<<<<<< HEAD
	sp.Logger.DebugContext(ctx, "url subject token request", "request", internallog.HTTPRequest(req, nil))
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	resp, body, err := internal.DoRequest(sp.Client, req)
	if err != nil {
		return "", fmt.Errorf("credentials: invalid response when retrieving subject token: %w", err)
	}
<<<<<<< HEAD
	sp.Logger.DebugContext(ctx, "url subject token response", "response", internallog.HTTPResponse(resp, body))
=======
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	if c := resp.StatusCode; c < http.StatusOK || c >= http.StatusMultipleChoices {
		return "", fmt.Errorf("credentials: status code %d: %s", c, body)
	}

	if sp.Format == nil {
		return string(body), nil
	}
	switch sp.Format.Type {
	case "json":
		jsonData := make(map[string]interface{})
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			return "", fmt.Errorf("credentials: failed to unmarshal subject token file: %w", err)
		}
		val, ok := jsonData[sp.Format.SubjectTokenFieldName]
		if !ok {
			return "", errors.New("credentials: provided subject_token_field_name not found in credentials")
		}
		token, ok := val.(string)
		if !ok {
			return "", errors.New("credentials: improperly formatted subject token")
		}
		return token, nil
	case fileTypeText:
		return string(body), nil
	default:
		return "", errors.New("credentials: invalid credential_source file format type: " + sp.Format.Type)
	}
}

func (sp *urlSubjectProvider) providerType() string {
	return urlProviderType
}
