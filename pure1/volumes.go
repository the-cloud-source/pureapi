/*
   Copyright 2018 David Evans

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package pure1

import "encoding/json"

// VolumeService type creates a service to expose volume endpoints
type VolumeService struct {
	client *Client
}

// GetVolumes returns a list of volume objects
func (v *VolumeService) GetVolumes(params map[string]string) ([]Volume, error) {
	req, err := v.client.NewRequest("GET", "volumes", params, nil)
	if err != nil {
		return nil, err
	}

	r := &pure1Response{}

	_, err = v.client.Do(req, r, false)
	if err != nil {
		return nil, err
	}

	m := []Volume{}

	for len(m) < r.TotalItems {
		for _, item := range r.Items {
			i := Volume{}
			str, _ := json.Marshal(item)
			json.Unmarshal([]byte(str), &i)
			m = append(m, i)
		}

		if len(m) < r.TotalItems {
			if r.ContinuationToken != nil {
				if params == nil {
					params = map[string]string{"continuation_token": r.ContinuationToken.(string)}
				} else {
					params["continuation_token"] = r.ContinuationToken.(string)
				}
				req, err := v.client.NewRequest("GET", "volumes", params, nil)
				if err != nil {
					return nil, err
				}

				_, err = v.client.Do(req, r, false)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return m, err
}
