// Copyright 2016 Mirantis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Mirantis/k8s-AppController/interfaces"
)

// Base is a base struct that contains data common for all resources
type Base struct {
	meta map[string]string
}

// Meta returns metadata parameter with given name, or empty string,
// if no metadata were provided or such parameter does not exist.
func (b Base) Meta(paramName string) string {
	if b.meta == nil {
		return ""
	}

	val, ok := b.meta[paramName]

	if !ok {
		return ""
	}

	return val
}

// KindToResourceTemplate is a map mapping kind strings to empty structs representing proper resources
// structs implement interfaces.ResourceTemplate
var KindToResourceTemplate = map[string]interfaces.ResourceTemplate{
	"daemonset":             DaemonSet{},
	"job":                   Job{},
	"petset":                PetSet{},
	"pod":                   Pod{},
	"replicaset":            ReplicaSet{},
	"service":               Service{},
	"configmap":             ConfigMap{},
	"secret":                Secret{},
	"deployment":            Deployment{},
	"persistentvolumeclaim": PersistentVolumeClaim{},
}

// Kinds is slice of keys from KindToResourceTemplate
var Kinds = getKeys(KindToResourceTemplate)

func getKeys(m map[string]interfaces.ResourceTemplate) (keys []string) {
	for key := range m {
		keys = append(keys, key)
	}

	return keys
}

func resourceListReady(resources []interfaces.BaseResource) (string, error) {
	for _, r := range resources {
		log.Printf("Checking status for resource %s", r.Key())
		status, err := r.Status(nil)
		if err != nil {
			return "error", err
		}
		if status != "ready" {
			return "not ready", fmt.Errorf("Resource %s is not ready", r.Key())
		}
	}
	return "ready", nil
}

func getPercentage(factorName string, meta map[string]string) (int32, error) {
	var factor string
	var ok bool
	if meta == nil {
		factor = "100"
	} else if factor, ok = meta[factorName]; !ok {
		factor = "100"
	}

	f, err := strconv.ParseInt(factor, 10, 32)
	if (f < 0 || f > 100) && err == nil {
		err = fmt.Errorf("%s factor not between 0 and 100", factorName)
	}
	return int32(f), err
}

func checkExistence(r interfaces.BaseResource) error {
	log.Println("Looking for ", r.Key())
	status, err := r.Status(nil)

	if err == nil {
		log.Printf("Found %s, status: %s ", r.Key(), status)
		return nil
	}

	return err
}

func createExistingResource(r interfaces.BaseResource) error {
	if err := checkExistence(r); err != nil {
		log.Fatalf("Resource %s not found", r.Key())
		return errors.New("Resource not found")
	}
	return nil
}
