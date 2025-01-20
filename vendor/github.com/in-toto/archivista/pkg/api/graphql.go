<<<<<<< HEAD
// Copyright 2023 The Witness Contributors
=======
// Copyright 2023-2024 The Witness Contributors
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
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

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

<<<<<<< HEAD
type graphQLError struct {
	Message string `json:"message"`
}

type graphQLResponse[T any] struct {
	Data   T              `json:"data,omitempty"`
	Errors []graphQLError `json:"errors,omitempty"`
}

type graphQLRequestBody[TVars any] struct {
	Query     string `json:"query"`
	Variables TVars  `json:"variables,omitempty"`
}

func GraphQlQuery[TRes any, TVars any](ctx context.Context, baseUrl, query string, vars TVars) (TRes, error) {
=======
const RetrieveSubjectsQuery = `query($gitoid: String!) {
	subjects(
		where: {
			hasStatementWith:{
        hasDsseWith:{
          gitoidSha256: $gitoid
        }
      }
		}
	) {
		edges {
      node{
        name
        subjectDigests{
          algorithm
          value
        }
      }
    }
  }
}`

const SearchQuery = `query($algo: String!, $digest: String!) {
  dsses(
    where: {
      hasStatementWith: {
        hasSubjectsWith: {
          hasSubjectDigestsWith: {
            value: $digest,
            algorithm: $algo
          }
        }
      }
    }
  ) {
    edges {
      node {
        gitoidSha256
        statement {
          attestationCollections {
            name
            attestations {
              type
            }
          }
        }
      }
    }
  }
}`

func GraphQlQuery[TRes any, TVars any](ctx context.Context, baseUrl, query string, vars TVars) (TRes, error) {
	return GraphQlQueryWithHeaders[TRes, TVars](ctx, baseUrl, query, vars, nil)
}

func GraphQlQueryWithHeaders[TRes any, TVars any](ctx context.Context, baseUrl, query string, vars TVars, headers map[string]string) (TRes, error) {
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	var response TRes
	queryUrl, err := url.JoinPath(baseUrl, "query")
	if err != nil {
		return response, err
	}

<<<<<<< HEAD
	requestBody := graphQLRequestBody[TVars]{
=======
	requestBody := GraphQLRequestBodyGeneric[TVars]{
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
		Query:     query,
		Variables: vars,
	}

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		return response, err
	}

<<<<<<< HEAD
	req, err := http.NewRequestWithContext(ctx, "POST", queryUrl, bytes.NewReader(reqBody))
=======
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, queryUrl, bytes.NewReader(reqBody))
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	if err != nil {
		return response, err
	}

<<<<<<< HEAD
=======
	for k, v := range headers {
		req.Header.Set(k, v)
	}

>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	req.Header.Set("Content-Type", "application/json")
	hc := &http.Client{}
	res, err := hc.Do(req)
	if err != nil {
		return response, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		errMsg, err := io.ReadAll(res.Body)
		if err != nil {
			return response, err
		}

		return response, errors.New(string(errMsg))
	}

	dec := json.NewDecoder(res.Body)
<<<<<<< HEAD
	gqlRes := graphQLResponse[TRes]{}
=======
	gqlRes := GraphQLResponseGeneric[TRes]{}
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	if err := dec.Decode(&gqlRes); err != nil {
		return response, err
	}

	if len(gqlRes.Errors) > 0 {
		return response, fmt.Errorf("graph ql query failed: %v", gqlRes.Errors)
	}

	return gqlRes.Data, nil
}
