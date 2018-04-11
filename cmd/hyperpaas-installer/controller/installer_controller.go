// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/docker/cli/cli/compose/convert"
	composetypes "github.com/docker/cli/cli/compose/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	dockerclient "github.com/docker/docker/client"
	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-sse"
	"github.com/hyperscale/hyperpaas/cmd/hyperpaas-installer/assets"
	"github.com/hyperscale/hyperpaas/docker"
	"github.com/hyperscale/hyperpaas/docker/compose"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// Resolve image constants
const (
	defaultNetworkDriver = "overlay"
	ResolveImageAlways   = "always"
	ResolveImageChanged  = "changed"
	ResolveImageNever    = "never"
)

// Event struct
type Event struct {
	Type   string `json:"type"`
	Action string `json:"action"`
	Error  string `json:"error,omitempty"`
}

// InstallerController struct
type InstallerController struct {
	dockerClient *docker.Client
	events       chan interface{}
}

// NewInstallerController func
func NewInstallerController(dockerClient *docker.Client) (*InstallerController, error) {
	return &InstallerController{
		dockerClient: dockerClient,
		events:       make(chan interface{}, 10),
	}, nil
}

// Mount endpoints
func (c InstallerController) Mount(r *server.Router) {
	events := sse.NewServer(c.getEventsHandler)
	events.SetRetry(time.Second * 5)

	r.AddRoute("/installer/events", events).Methods(http.MethodGet)
	r.AddRouteFunc("/installer", c.postInstallerHandler).Methods(http.MethodPost)
}

func (c *InstallerController) postInstallerHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	ctx := context.Background()

	if err := r.ParseMultipartForm(2048); err != nil {
		server.FailureFromError(w, http.StatusRequestEntityTooLarge, err)

		return
	}

	log.Info().Msgf("Data: %#v", r.PostForm)

	c.events <- Event{
		Type:   "install",
		Action: "started",
	}

	content, err := assets.Asset("static/config/docker-compose.yml")
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	c.events <- Event{
		Type:   "compose",
		Action: "loading",
	}

	composeConfig, err := compose.Loader(content, map[string]string{
		"EMAIL":  r.PostForm.Get("email"),
		"DOMAIN": r.PostForm.Get("domain"),
	})
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		c.events <- Event{
			Type:   "compose",
			Action: "failed",
			Error:  err.Error(),
		}

		return
	}

	c.events <- Event{
		Type:   "compose",
		Action: "loaded",
	}

	go func() {
		namespace := convert.NewNamespace("hyperpaas")

		{
			c.events <- Event{
				Type:   "services",
				Action: "cleaning",
			}

			// prune services
			services := map[string]struct{}{}
			for _, service := range composeConfig.Services {
				services[service.Name] = struct{}{}
			}
			pruneServices(ctx, c.dockerClient, namespace, services)

			c.events <- Event{
				Type:   "services",
				Action: "cleaned",
			}
		}

		c.events <- Event{
			Type:   "networks",
			Action: "loading",
		}

		serviceNetworks := getServicesDeclaredNetworks(composeConfig.Services)
		networks, externalNetworks := convert.Networks(namespace, composeConfig.Networks, serviceNetworks)
		if err := validateExternalNetworks(ctx, c.dockerClient, externalNetworks); err != nil {
			log.Error().Err(err).Msg("validateExternalNetworks")

			c.events <- Event{
				Type:   "networks",
				Action: "failed",
				Error:  err.Error(),
			}

			return
		}

		c.events <- Event{
			Type:   "networks",
			Action: "loaded",
		}

		c.events <- Event{
			Type:   "networks",
			Action: "creating",
		}

		if err := createNetworks(ctx, c.dockerClient, namespace, networks); err != nil {
			log.Error().Err(err).Msg("createNetworks")

			c.events <- Event{
				Type:   "networks",
				Action: "failed",
				Error:  err.Error(),
			}

			return
		}

		c.events <- Event{
			Type:   "networks",
			Action: "created",
		}

		registryPwd, err := hashBcrypt(r.PostForm.Get("registry_password"))
		if err != nil {
			log.Error().Err(err).Msg("hashBcrypt")

			c.events <- Event{
				Type:   "secrets",
				Action: "failed",
				Error:  err.Error(),
			}

			return
		}

		file, err := os.OpenFile("./var/lib/hyperpaas/secrets/registry.htpasswd", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Error().Err(err).Msg("OpenFile")

			c.events <- Event{
				Type:   "secrets",
				Action: "failed",
				Error:  err.Error(),
			}

			return
		}

		defer func() {
			if err := file.Close(); err != nil {
				log.Error().Err(err).Msg("OpenFile Close")
			}
		}()

		if _, err := file.Write([]byte(fmt.Sprintf("%s:%s", r.PostForm.Get("registry_user"), registryPwd))); err != nil {
			log.Error().Err(err).Msg("OpenFile")

			c.events <- Event{
				Type:   "secrets",
				Action: "failed",
				Error:  err.Error(),
			}

			return
		}

		c.events <- Event{
			Type:   "secrets",
			Action: "loading",
		}

		secrets, err := convert.Secrets(namespace, composeConfig.Secrets)
		if err != nil {
			log.Error().Err(err).Msg("Secrets")

			return
		}

		c.events <- Event{
			Type:   "secrets",
			Action: "loaded",
		}

		c.events <- Event{
			Type:   "secrets",
			Action: "creating",
		}

		if err := createSecrets(ctx, c.dockerClient, secrets); err != nil {
			log.Error().Err(err).Msg("createSecrets")

			c.events <- Event{
				Type:   "secrets",
				Action: "failed",
				Error:  err.Error(),
			}

			return
		}

		c.events <- Event{
			Type:   "secrets",
			Action: "created",
		}

		c.events <- Event{
			Type:   "configs",
			Action: "loading",
		}

		configs, err := convert.Configs(namespace, composeConfig.Configs)
		if err != nil {
			log.Error().Err(err).Msg("Configs")

			c.events <- Event{
				Type:   "configs",
				Action: "failed",
				Error:  err.Error(),
			}

			return
		}

		c.events <- Event{
			Type:   "confings",
			Action: "loaded",
		}

		c.events <- Event{
			Type:   "configs",
			Action: "creating",
		}

		if err := createConfigs(ctx, c.dockerClient, configs); err != nil {
			log.Error().Err(err).Msg("createConfigs")

			c.events <- Event{
				Type:   "configs",
				Action: "failed",
				Error:  err.Error(),
			}

			return
		}

		c.events <- Event{
			Type:   "configs",
			Action: "created",
		}

		c.events <- Event{
			Type:   "services",
			Action: "loading",
		}

		services, err := convert.Services(namespace, composeConfig, c.dockerClient.Client)
		if err != nil {
			log.Error().Err(err).Msg("Services")

			c.events <- Event{
				Type:   "services",
				Action: "failed",
				Error:  err.Error(),
			}

			return
		}

		c.events <- Event{
			Type:   "services",
			Action: "loaded",
		}

		c.events <- Event{
			Type:   "services",
			Action: "deploying",
		}

		if err := deployServices(ctx, c.dockerClient, services, namespace, true, ResolveImageAlways); err != nil {
			log.Error().Err(err).Msg("deployServices")

			c.events <- Event{
				Type:   "services",
				Action: "failed",
				Error:  err.Error(),
			}

			//@TODO send Error Event
			return
		}

		c.events <- Event{
			Type:   "services",
			Action: "deployed",
		}

		c.events <- Event{
			Type:   "install",
			Action: "finished",
		}
	}()

	server.JSON(w, http.StatusOK, composeConfig)
}

func (c *InstallerController) getEventsHandler(rw sse.ResponseWriter, r *http.Request) {
	for {
		select {
		case event := <-c.events:
			data, err := json.Marshal(event)
			if err != nil {
				log.Error().Err(err).Msg("Marshal Event")

				continue
			}

			rw.Send(&sse.MessageEvent{
				Data: data,
			})

		case <-rw.CloseNotify:

			return
		}
	}
}

// pruneServices removes services that are no longer referenced in the source
func pruneServices(ctx context.Context, client *docker.Client, namespace convert.Namespace, services map[string]struct{}) {
	oldServices, err := getServices(ctx, client, namespace.Name())
	if err != nil {
		log.Error().Err(err).Msg("Failed to list services")
	}

	pruneServices := []swarm.Service{}

	for _, service := range oldServices {
		if _, exists := services[namespace.Descope(service.Spec.Name)]; !exists {
			pruneServices = append(pruneServices, service)
		}
	}

	removeServices(ctx, client, pruneServices)
}

func removeServices(ctx context.Context, client *docker.Client, services []swarm.Service) bool {
	var hasError bool

	sort.Slice(services, sortServiceByName(services))

	for _, service := range services {
		log.Info().Msgf("Removing service %s", service.Spec.Name)

		if err := client.ServiceRemove(ctx, service.ID); err != nil {
			hasError = true

			log.Error().Err(err).Msgf("Failed to remove service %s", service.ID)
		}
	}

	return hasError
}

func sortServiceByName(services []swarm.Service) func(i, j int) bool {
	return func(i, j int) bool {
		return services[i].Spec.Name < services[j].Spec.Name
	}
}

func removeNetworks(ctx context.Context, client *docker.Client, networks []types.NetworkResource) bool {
	var hasError bool

	for _, network := range networks {
		log.Info().Msgf("Removing network %s", network.Name)

		if err := client.NetworkRemove(ctx, network.ID); err != nil {
			hasError = true

			log.Error().Err(err).Msgf("Failed to remove network %s", network.ID)
		}
	}

	return hasError
}

func removeSecrets(ctx context.Context, client *docker.Client, secrets []swarm.Secret) bool {
	var hasError bool

	for _, secret := range secrets {
		log.Info().Msgf("Removing secret %s", secret.Spec.Name)

		if err := client.SecretRemove(ctx, secret.ID); err != nil {
			hasError = true

			log.Error().Err(err).Msgf("Failed to remove secret %s", secret.ID)
		}
	}

	return hasError
}

func removeConfigs(ctx context.Context, client *docker.Client, configs []swarm.Config) bool {
	var hasError bool

	for _, config := range configs {
		log.Info().Msgf("Removing config %s", config.Spec.Name)

		if err := client.ConfigRemove(ctx, config.ID); err != nil {
			hasError = true

			log.Error().Err(err).Msgf("Failed to remove config %s", config.ID)
		}
	}

	return hasError
}

func getServicesDeclaredNetworks(serviceConfigs []composetypes.ServiceConfig) map[string]struct{} {
	serviceNetworks := map[string]struct{}{}

	for _, serviceConfig := range serviceConfigs {
		if len(serviceConfig.Networks) == 0 {
			serviceNetworks["default"] = struct{}{}
			continue
		}

		for network := range serviceConfig.Networks {
			serviceNetworks[network] = struct{}{}
		}
	}

	return serviceNetworks
}

func validateExternalNetworks(ctx context.Context, client dockerclient.NetworkAPIClient, externalNetworks []string) error {
	for _, networkName := range externalNetworks {
		if !container.NetworkMode(networkName).IsUserDefined() {
			// Networks that are not user defined always exist on all nodes as
			// local-scoped networks, so there's no need to inspect them.
			continue
		}

		network, err := client.NetworkInspect(ctx, networkName, types.NetworkInspectOptions{})

		switch {
		case dockerclient.IsErrNotFound(err):
			return errors.Errorf("network %q is declared as external, but could not be found. You need to create a swarm-scoped network before the stack is deployed", networkName)
		case err != nil:
			return err
		case network.Scope != "swarm":
			return errors.Errorf("network %q is declared as external, but it is not in the right scope: %q instead of \"swarm\"", networkName, network.Scope)
		}
	}

	return nil
}

func deleteSecrets(ctx context.Context, client *docker.Client, secrets []swarm.SecretReference) error {
	for _, secretRef := range secrets {
		secret, _, err := client.SecretInspectWithRaw(ctx, secretRef.SecretName)
		if err == nil {
			if err := client.SecretRemove(ctx, secret.ID); err != nil {
				return errors.Wrapf(err, "failed to remove secret %s", secretRef.SecretName)
			}
		}
	}

	return nil
}

func createSecrets(ctx context.Context, client *docker.Client, secrets []swarm.SecretSpec) error {
	for _, secretSpec := range secrets {
		secret, _, err := client.SecretInspectWithRaw(ctx, secretSpec.Name)

		switch {
		case err == nil:
			// secret already exists, then we update that
			if err := client.SecretUpdate(ctx, secret.ID, secret.Meta.Version, secretSpec); err != nil {
				return errors.Wrapf(err, "failed to update secret %s", secretSpec.Name)
			}
		case dockerclient.IsErrNotFound(err):
			// secret does not exist, then we create a new one.

			log.Info().Msgf("Creating secret %s", secretSpec.Name)

			if _, err := client.SecretCreate(ctx, secretSpec); err != nil {
				return errors.Wrapf(err, "failed to create secret %s", secretSpec.Name)
			}
		default:
			return err
		}
	}

	return nil
}

func createConfigs(ctx context.Context, client *docker.Client, configs []swarm.ConfigSpec) error {
	for _, configSpec := range configs {
		config, _, err := client.ConfigInspectWithRaw(ctx, configSpec.Name)

		switch {
		case err == nil:
			// config already exists, then we update that
			if err := client.ConfigUpdate(ctx, config.ID, config.Meta.Version, configSpec); err != nil {
				return errors.Wrapf(err, "failed to update config %s", configSpec.Name)
			}
		case dockerclient.IsErrNotFound(err):
			// config does not exist, then we create a new one.

			log.Info().Msgf("Creating config %s", configSpec.Name)
			if _, err := client.ConfigCreate(ctx, configSpec); err != nil {
				return errors.Wrapf(err, "failed to create config %s", configSpec.Name)
			}
		default:
			return err
		}
	}
	return nil
}

func getStackNetworks(
	ctx context.Context,
	client dockerclient.APIClient,
	namespace string,
) ([]types.NetworkResource, error) {
	return client.NetworkList(ctx, types.NetworkListOptions{Filters: getStackFilter(namespace)})
}

func createNetworks(ctx context.Context, client *docker.Client, namespace convert.Namespace, networks map[string]types.NetworkCreate) error {
	existingNetworks, err := getStackNetworks(ctx, client, namespace.Name())
	if err != nil {
		return err
	}

	existingNetworkMap := make(map[string]types.NetworkResource)
	for _, network := range existingNetworks {
		existingNetworkMap[network.Name] = network
	}

	for internalName, createOpts := range networks {
		name := namespace.Scope(internalName)
		if _, exists := existingNetworkMap[name]; exists {
			continue
		}

		if createOpts.Driver == "" {
			createOpts.Driver = defaultNetworkDriver
		}

		log.Info().Msgf("Creating network %s", name)
		if _, err := client.NetworkCreate(ctx, name, createOpts); err != nil {
			return errors.Wrapf(err, "failed to create network %s", internalName)
		}
	}

	return nil
}

func deployServices(
	ctx context.Context,
	client *docker.Client,
	services map[string]swarm.ServiceSpec,
	namespace convert.Namespace,
	sendAuth bool,
	resolveImage string,
) error {
	existingServices, err := getServices(ctx, client, namespace.Name())
	if err != nil {
		return err
	}

	existingServiceMap := make(map[string]swarm.Service)
	for _, service := range existingServices {
		existingServiceMap[service.Spec.Name] = service
	}

	for internalName, serviceSpec := range services {
		name := namespace.Scope(internalName)

		encodedAuth := ""
		image := serviceSpec.TaskTemplate.ContainerSpec.Image
		if sendAuth {
			// Retrieve encoded auth token from the image reference
			/*encodedAuth, err = command.RetrieveAuthTokenFromImage(ctx, dockerCli, image)
			if err != nil {
				return err
			}*/
		}

		if service, exists := existingServiceMap[name]; exists {
			log.Info().Msgf("Updating service %s (id: %s)", name, service.ID)

			updateOpts := types.ServiceUpdateOptions{EncodedRegistryAuth: encodedAuth}

			switch {
			case resolveImage == ResolveImageAlways || (resolveImage == ResolveImageChanged && image != service.Spec.Labels[convert.LabelImage]):
				// image should be updated by the server using QueryRegistry
				updateOpts.QueryRegistry = true
			case image == service.Spec.Labels[convert.LabelImage]:
				// image has not changed; update the serviceSpec with the
				// existing information that was set by QueryRegistry on the
				// previous deploy. Otherwise this will trigger an incorrect
				// service update.
				serviceSpec.TaskTemplate.ContainerSpec.Image = service.Spec.TaskTemplate.ContainerSpec.Image
			}
			response, err := client.ServiceUpdate(
				ctx,
				service.ID,
				service.Version,
				serviceSpec,
				updateOpts,
			)
			if err != nil {
				return errors.Wrapf(err, "failed to update service %s", name)
			}

			for _, warning := range response.Warnings {
				log.Error().Msg(warning)
			}
		} else {
			log.Info().Msgf("Creating service %s", name)

			createOpts := types.ServiceCreateOptions{EncodedRegistryAuth: encodedAuth}

			// query registry if flag disabling it was not set
			if resolveImage == ResolveImageAlways || resolveImage == ResolveImageChanged {
				createOpts.QueryRegistry = true
			}

			if _, err := client.ServiceCreate(ctx, serviceSpec, createOpts); err != nil {
				return errors.Wrapf(err, "failed to create service %s", name)
			}
		}
	}

	return nil
}

func getStackFilter(namespace string) filters.Args {
	filter := filters.NewArgs()
	filter.Add("label", convert.LabelNamespace+"="+namespace)

	return filter
}

func getServiceFilter(namespace string) filters.Args {
	return getStackFilter(namespace)
}

func getServices(ctx context.Context, client *docker.Client, namespace string) ([]swarm.Service, error) {
	return client.ServiceList(ctx, types.ServiceListOptions{Filters: getServiceFilter(namespace)})
}

func hashBcrypt(password string) (hash string, err error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	return string(passwordBytes), nil
}
