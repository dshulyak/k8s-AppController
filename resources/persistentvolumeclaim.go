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
	"log"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/unversioned"

	"github.com/Mirantis/k8s-AppController/client"
	"github.com/Mirantis/k8s-AppController/interfaces"
	"github.com/Mirantis/k8s-AppController/report"
)

type PersistentVolumeClaim struct {
	Base
	PersistentVolumeClaim *api.PersistentVolumeClaim
	Client                unversioned.PersistentVolumeClaimInterface
}

func persistentVolumeClaimKey(name string) string {
	return "persistentvolumeclaim/" + name
}

func (p PersistentVolumeClaim) Key() string {
	return persistentVolumeClaimKey(p.PersistentVolumeClaim.Name)
}

func persistentVolumeClaimStatus(p unversioned.PersistentVolumeClaimInterface, name string) (string, error) {
	persistentVolumeClaim, err := p.Get(name)
	if err != nil {
		return "error", err
	}

	if persistentVolumeClaim.Status.Phase == api.ClaimBound {
		return "ready", nil
	}

	return "not ready", nil
}

func (p PersistentVolumeClaim) Create() error {
	if err := checkExistence(p); err != nil {
		log.Println("Creating ", p.Key())
		p.PersistentVolumeClaim, err = p.Client.Create(p.PersistentVolumeClaim)
		return err
	}
	return nil
}

// Delete deletes persistentVolumeClaim from the cluster
func (p PersistentVolumeClaim) Delete() error {
	return p.Client.Delete(p.PersistentVolumeClaim.Name)
}

func (p PersistentVolumeClaim) Status(meta map[string]string) (string, error) {
	return persistentVolumeClaimStatus(p.Client, p.PersistentVolumeClaim.Name)
}

// NameMatches gets resource definition and a name and checks if
// the PersistentVolumeClaim part of resource definition has matching name.
func (p PersistentVolumeClaim) NameMatches(def client.ResourceDefinition, name string) bool {
	return def.PersistentVolumeClaim != nil && def.PersistentVolumeClaim.Name == name
}

// New returns new PersistentVolumeClaim based on resource definition
func (p PersistentVolumeClaim) New(def client.ResourceDefinition, c client.Interface) interfaces.Resource {
	return NewPersistentVolumeClaim(def.PersistentVolumeClaim, c.PersistentVolumeClaims(), def.Meta)
}

// NewExisting returns new ExistingPersistentVolumeClaim based on resource definition
func (p PersistentVolumeClaim) NewExisting(name string, c client.Interface) interfaces.Resource {
	return NewExistingPersistentVolumeClaim(name, c.PersistentVolumeClaims())
}

func NewPersistentVolumeClaim(persistentVolumeClaim *api.PersistentVolumeClaim, client unversioned.PersistentVolumeClaimInterface, meta map[string]string) interfaces.Resource {
	return report.SimpleReporter{BaseResource: PersistentVolumeClaim{Base: Base{meta}, PersistentVolumeClaim: persistentVolumeClaim, Client: client}}
}

type ExistingPersistentVolumeClaim struct {
	Base
	Name   string
	Client unversioned.PersistentVolumeClaimInterface
}

func (p ExistingPersistentVolumeClaim) Key() string {
	return persistentVolumeClaimKey(p.Name)
}

func (p ExistingPersistentVolumeClaim) Create() error {
	return createExistingResource(p)
}

func (p ExistingPersistentVolumeClaim) Status(meta map[string]string) (string, error) {
	return persistentVolumeClaimStatus(p.Client, p.Name)
}

// Delete deletes persistentVolumeClaim from the cluster
func (p ExistingPersistentVolumeClaim) Delete() error {
	return p.Client.Delete(p.Name)
}

func NewExistingPersistentVolumeClaim(name string, client unversioned.PersistentVolumeClaimInterface) interfaces.Resource {
	return report.SimpleReporter{BaseResource: ExistingPersistentVolumeClaim{Name: name, Client: client}}
}
