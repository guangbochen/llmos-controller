/*
Copyright 2024 llmos.ai.

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
// Code generated by main. DO NOT EDIT.

package v1

import (
	v1 "github.com/llmos-ai/llmos-operator/pkg/apis/management.llmos.ai/v1"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/v3/pkg/generic"
	"github.com/rancher/wrangler/v3/pkg/schemes"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	schemes.Register(v1.AddToScheme)
}

type Interface interface {
	ManagedAddon() ManagedAddonController
	Setting() SettingController
	Upgrade() UpgradeController
	User() UserController
}

func New(controllerFactory controller.SharedControllerFactory) Interface {
	return &version{
		controllerFactory: controllerFactory,
	}
}

type version struct {
	controllerFactory controller.SharedControllerFactory
}

func (v *version) ManagedAddon() ManagedAddonController {
	return generic.NewController[*v1.ManagedAddon, *v1.ManagedAddonList](schema.GroupVersionKind{Group: "management.llmos.ai", Version: "v1", Kind: "ManagedAddon"}, "managedaddons", true, v.controllerFactory)
}

func (v *version) Setting() SettingController {
	return generic.NewNonNamespacedController[*v1.Setting, *v1.SettingList](schema.GroupVersionKind{Group: "management.llmos.ai", Version: "v1", Kind: "Setting"}, "settings", v.controllerFactory)
}

func (v *version) Upgrade() UpgradeController {
	return generic.NewController[*v1.Upgrade, *v1.UpgradeList](schema.GroupVersionKind{Group: "management.llmos.ai", Version: "v1", Kind: "Upgrade"}, "upgrades", true, v.controllerFactory)
}

func (v *version) User() UserController {
	return generic.NewNonNamespacedController[*v1.User, *v1.UserList](schema.GroupVersionKind{Group: "management.llmos.ai", Version: "v1", Kind: "User"}, "users", v.controllerFactory)
}
