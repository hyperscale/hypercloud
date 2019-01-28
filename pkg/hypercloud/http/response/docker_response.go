// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package response

import (
	"time"

	"github.com/docker/docker/api/types/swarm"
)

// Version represents the internal object version.
type Version struct {
	Index uint64 `json:"index,omitempty"`
}

// Meta is a base object inherited by most of the other once.
type Meta struct {
	Version   Version   `json:"version,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Annotations represents how to describe an object.
type Annotations struct {
	Name   string            `json:"name,omitempty"`
	Labels map[string]string `json:"labels"`
}

// Driver represents a driver (network, logging, secrets backend).
type Driver struct {
	Name    string            `json:"name,omitempty"`
	Options map[string]string `json:"options,omitempty"`
}

// TLSInfo represents the TLS information about what CA certificate is trusted,
// and who the issuer for a TLS certificate is
type TLSInfo struct {
	// TrustRoot is the trusted CA root certificate in PEM format
	TrustRoot string `json:"trust_root,omitempty"`

	// CertIssuer is the raw subject bytes of the issuer
	CertIssuerSubject []byte `json:"cert_issuer_subject,omitempty"`

	// CertIssuerPublicKey is the raw public key bytes of the issuer
	CertIssuerPublicKey []byte `json:"cert_issuer_public_key,omitempty"`
}

// Node represents a node.
type Node struct {
	ID string `json:"id"`

	Meta
	// Spec defines the desired state of the node as specified by the user.
	// The system will honor this and will *never* modify it.
	Spec NodeSpec `json:"spec,omitempty"`
	// Description encapsulates the properties of the Node as reported by the
	// agent.
	Description NodeDescription `json:"description,omitempty"`
	// Status provides the current status of the node, as seen by the manager.
	Status NodeStatus `json:"status,omitempty"`
	// ManagerStatus provides the current status of the node's manager
	// component, if the node is a manager.
	ManagerStatus *ManagerStatus `json:"manager_status,omitempty"`
}

// NodeSpec represents the spec of a node.
type NodeSpec struct {
	Annotations
	Role         swarm.NodeRole         `json:"role,omitempty"`
	Availability swarm.NodeAvailability `json:"availability,omitempty"`
}

// NodeDescription represents the description of a node.
type NodeDescription struct {
	Hostname  string            `json:"hostname,omitempty"`
	Platform  Platform          `json:"platform,omitempty"`
	Resources Resources         `json:"resources,omitempty"`
	Engine    EngineDescription `json:"engine,omitempty"`
	TLSInfo   TLSInfo           `json:"tls_info,omitempty"`
}

// Platform represents the platform (Arch/OS).
type Platform struct {
	Architecture string `json:"architecture,omitempty"`
	OS           string `json:"os,omitempty"`
}

// EngineDescription represents the description of an engine.
type EngineDescription struct {
	EngineVersion string              `json:"engine_version,omitempty"`
	Labels        map[string]string   `json:"labels,omitempty"`
	Plugins       []PluginDescription `json:"plugins,omitempty"`
}

// PluginDescription represents the description of an engine plugin.
type PluginDescription struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

// NodeStatus represents the status of a node.
type NodeStatus struct {
	State   swarm.NodeState `json:"state,omitempty"`
	Message string          `json:"message,omitempty"`
	Addr    string          `json:"addr,omitempty"`
}

// ManagerStatus represents the status of a manager.
type ManagerStatus struct {
	Leader       bool               `json:"leader,omitempty"`
	Reachability swarm.Reachability `json:"reachability,omitempty"`
	Addr         string             `json:"addr,omitempty"`
}

// Resources represents resources (CPU/Memory).
type Resources struct {
	NanoCPUs         int64             `json:"nano_cpus,omitempty"`
	MemoryBytes      int64             `json:"memory_bytes,omitempty"`
	GenericResources []GenericResource `json:"generic_resources,omitempty"`
}

// GenericResource represents a "user defined" resource which can
// be either an integer (e.g: SSD=3) or a string (e.g: SSD=sda1)
type GenericResource struct {
	NamedResourceSpec    *NamedGenericResource    `json:"named_resource_spec,omitempty"`
	DiscreteResourceSpec *DiscreteGenericResource `json:"discrete_resource_spec,omitempty"`
}

// NamedGenericResource represents a "user defined" resource which is defined
// as a string.
// "Kind" is used to describe the Kind of a resource (e.g: "GPU", "FPGA", "SSD", ...)
// Value is used to identify the resource (GPU="UUID-1", FPGA="/dev/sdb5", ...)
type NamedGenericResource struct {
	Kind  string `json:"kind,omitempty"`
	Value string `json:"value,omitempty"`
}

// DiscreteGenericResource represents a "user defined" resource which is defined
// as an integer
// "Kind" is used to describe the Kind of a resource (e.g: "GPU", "FPGA", "SSD", ...)
// Value is used to count the resource (SSD=5, HDD=3, ...)
type DiscreteGenericResource struct {
	Kind  string `json:"kind,omitempty"`
	Value int64  `json:"value,omitempty"`
}
