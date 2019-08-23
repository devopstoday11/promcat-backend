package vendor

import (
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

type Kind string

const (
	VENDOR Kind = "Vendor"
)

type Resource struct {
	ID          string `json:"id,omitempty" yaml:"id,omitempty"`
	Kind        Kind   `json:"kind" yaml:"kind"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Icon        string `json:"icon" yaml:"icon"`
	Website     string `json:"website" yaml:"website"`
}

type resourceAlias Resource // Avoid stack overflow while marshalling / unmarshalling

func (r *Resource) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	res := resourceAlias{}
	err = unmarshal(&res)
	if err != nil {
		return
	}
	*r = Resource(res)
	r.ID = r.Hash()
	return
}

func (r *Resource) MarshalYAML() (interface{}, error) {
	x := resourceAlias(*r)
	x.ID = r.Hash()
	return yaml.Marshal(x)
}

func (r *Resource) UnmarshalJSON(data []byte) (err error) {
	res := resourceAlias{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	*r = Resource(res)
	r.ID = r.Hash()
	return
}

func (r *Resource) MarshalJSON() ([]byte, error) {
	x := resourceAlias(*r)
	x.ID = r.Hash()
	return json.Marshal(x)
}

func (r *Resource) Validate() error {
	var errors []string

	if r.Kind == "" {
		errors = append(errors, "the resource must have a defined Kind")
	}

	if r.Icon == "" {
		errors = append(errors, "the resource must have a valid icon")
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ","))
	}

	return nil
}

func (r *Resource) Hash() string {
	sum := sha1.Sum([]byte(string(r.Kind) + r.Name + r.Description))
	b32 := base32.StdEncoding.EncodeToString(sum[:])
	return b32[:20]
}