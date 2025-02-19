package clusterentities

import (
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/net"
	"github.com/stackrox/rox/pkg/networkgraph"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

// ContainerMetadata is the container metadata that is stored per instance
type ContainerMetadata struct {
	DeploymentID  string
	DeploymentTS  int64
	PodID         string
	PodUID        string
	ContainerName string
	ContainerID   string
	Namespace     string
	StartTime     *types.Timestamp
	ImageID       string
}

// PublicIPsListener is an interface for listeners on changes to the set of public IP addresses.
// Note: Implementors of this interface must ensure the methods complete in a very short time/do not block, as they
// get invoked synchronously in a critical section.
type PublicIPsListener interface {
	OnAdded(ip net.IPAddress)
	OnRemoved(ip net.IPAddress)
}

// Store is a store for managing cluster entities (currently deployments only) and allows looking them up by
// endpoint.
type Store struct {
	// ipMap maps ip addresses to sets of deployment ids this IP is associated with.
	ipMap map[net.IPAddress]map[string]struct{}
	// endpointMap maps endpoints to a (deployment id -> endpoint target info) mapping.
	endpointMap map[net.NumericEndpoint]map[string]map[EndpointTargetInfo]struct{}
	// containerIDMap maps container IDs to container metadata
	containerIDMap map[string]ContainerMetadata

	// reverseIpMap maps deployment ids to sets of IP addresses associated with this deployment.
	reverseIPMap map[string]map[net.IPAddress]struct{}
	// reverseEndpointMap maps deployment ids to sets of endpoints associated with this deployment.
	reverseEndpointMap map[string]map[net.NumericEndpoint]struct{}
	// reverseContainerIDMap maps deployment ids to sets of container IDs associated with this deployment.
	reverseContainerIDMap map[string]map[string]struct{}
	// callbackChannel is a channel to send container metadata upon resolution
	callbackChannel chan<- ContainerMetadata

	publicIPRefCounts  map[net.IPAddress]*int
	publicIPsListeners map[PublicIPsListener]struct{}

	mutex sync.RWMutex
}

// NewStore creates and returns a new store instance.
// Note: Generally, you probably do not want to call this function, but use the singleton instance returned by
// `StoreInstance()`.
func NewStore() *Store {
	store := &Store{}
	store.initMaps()
	return store
}

func (e *Store) initMaps() {
	e.ipMap = make(map[net.IPAddress]map[string]struct{})
	e.endpointMap = make(map[net.NumericEndpoint]map[string]map[EndpointTargetInfo]struct{})
	e.containerIDMap = make(map[string]ContainerMetadata)
	e.reverseIPMap = make(map[string]map[net.IPAddress]struct{})
	e.reverseEndpointMap = make(map[string]map[net.NumericEndpoint]struct{})
	e.reverseContainerIDMap = make(map[string]map[string]struct{})
	e.publicIPRefCounts = make(map[net.IPAddress]*int)
	e.publicIPsListeners = make(map[PublicIPsListener]struct{})
}

// EndpointTargetInfo is the target port for an endpoint (container port, service port etc.).
type EndpointTargetInfo struct {
	ContainerPort uint16
	PortName      string
}

// EntityData is a data structure representing the updates to be applied to the store for a given deployment.
type EntityData struct {
	ips          map[net.IPAddress]struct{}
	endpoints    map[net.NumericEndpoint][]EndpointTargetInfo
	containerIDs map[string]ContainerMetadata
}

// AddIP adds an IP address to the set of IP addresses of the respective deployment.
func (ed *EntityData) AddIP(ip net.IPAddress) {
	if ed.ips == nil {
		ed.ips = make(map[net.IPAddress]struct{})
	}
	ed.ips[ip] = struct{}{}
}

// AddEndpoint adds an endpoint along with a target info to the endpoints of the respective deployment.
func (ed *EntityData) AddEndpoint(ep net.NumericEndpoint, info EndpointTargetInfo) {
	if ed.endpoints == nil {
		ed.endpoints = make(map[net.NumericEndpoint][]EndpointTargetInfo)
	}
	ed.endpoints[ep] = append(ed.endpoints[ep], info)
}

// AddContainerID adds a container ID to the container IDs of the respective deployment.
func (ed *EntityData) AddContainerID(containerID string, container ContainerMetadata) {
	if ed.containerIDs == nil {
		ed.containerIDs = make(map[string]ContainerMetadata)
	}
	ed.containerIDs[containerID] = container
}

// Cleanup deletes all entries from store
func (e *Store) Cleanup() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.initMaps()
}

// Apply applies an update to the store. If incremental is true, data will be added; otherwise, data for each deployment
// that is a key in the map will be replaced (or deleted).
func (e *Store) Apply(updates map[string]*EntityData, incremental bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.applyNoLock(updates, incremental)
}

func (e *Store) purgeNoLock(deploymentID string) {
	for ip := range e.reverseIPMap[deploymentID] {
		set := e.ipMap[ip]
		delete(set, deploymentID)
		if len(set) == 0 {
			delete(e.ipMap, ip)
			if ip.IsPublic() {
				e.decPublicIPRefNoLock(ip)
			}
		}
	}
	for ep := range e.reverseEndpointMap[deploymentID] {
		set := e.endpointMap[ep]
		delete(set, deploymentID)
		if len(set) == 0 {
			delete(e.endpointMap, ep)
			if ipAddr := ep.IPAndPort.Address; ipAddr.IsPublic() {
				e.decPublicIPRefNoLock(ipAddr)
			}
		}
	}
	for containerID := range e.reverseContainerIDMap[deploymentID] {
		delete(e.containerIDMap, containerID)
	}

	delete(e.reverseIPMap, deploymentID)
	delete(e.reverseEndpointMap, deploymentID)
	delete(e.reverseContainerIDMap, deploymentID)
}

func (e *Store) applyNoLock(updates map[string]*EntityData, incremental bool) {
	if !incremental {
		for deploymentID := range updates {
			e.purgeNoLock(deploymentID)
		}
	}

	for deploymentID, data := range updates {
		if data == nil {
			continue
		}
		e.applySingleNoLock(deploymentID, *data)
	}
}

func (e *Store) applySingleNoLock(deploymentID string, data EntityData) {
	reverseEPs := e.reverseEndpointMap[deploymentID]
	reverseIPs := e.reverseIPMap[deploymentID]
	reverseContainerIDs := e.reverseContainerIDMap[deploymentID]

	for ep, targetInfos := range data.endpoints {
		if reverseEPs == nil {
			reverseEPs = make(map[net.NumericEndpoint]struct{})
			e.reverseEndpointMap[deploymentID] = reverseEPs
		}
		reverseEPs[ep] = struct{}{}

		epMap := e.endpointMap[ep]
		if epMap == nil {
			epMap = make(map[string]map[EndpointTargetInfo]struct{})
			e.endpointMap[ep] = epMap
			if ipAddr := ep.IPAndPort.Address; ipAddr.IsPublic() {
				e.incPublicIPRefNoLock(ipAddr)
			}
		}
		targetSet := epMap[deploymentID]
		if targetSet == nil {
			targetSet = make(map[EndpointTargetInfo]struct{})
			epMap[deploymentID] = targetSet
		}
		for _, tgtInfo := range targetInfos {
			targetSet[tgtInfo] = struct{}{}
		}
	}

	for ip := range data.ips {
		if reverseIPs == nil {
			reverseIPs = make(map[net.IPAddress]struct{})
			e.reverseIPMap[deploymentID] = reverseIPs
		}
		reverseIPs[ip] = struct{}{}

		ipMap := e.ipMap[ip]
		if ipMap == nil {
			ipMap = make(map[string]struct{})
			e.ipMap[ip] = ipMap
			if ip.IsPublic() {
				e.incPublicIPRefNoLock(ip)
			}
		}
		ipMap[deploymentID] = struct{}{}
	}

	mdsForCallback := make([]ContainerMetadata, 0, len(data.containerIDs))
	for containerID, metadata := range data.containerIDs {
		if reverseContainerIDs == nil {
			reverseContainerIDs = make(map[string]struct{})
			e.reverseContainerIDMap[deploymentID] = reverseContainerIDs
		}
		reverseContainerIDs[containerID] = struct{}{}
		e.containerIDMap[containerID] = metadata
		mdsForCallback = append(mdsForCallback, metadata)
	}

	if e.callbackChannel != nil && len(mdsForCallback) > 0 {
		go sendMetadataCallbacks(e.callbackChannel, mdsForCallback)
	}
}

func sendMetadataCallbacks(callbackC chan<- ContainerMetadata, mds []ContainerMetadata) {
	for _, md := range mds {
		callbackC <- md
	}
}

// RegisterContainerMetadataCallbackChannel registers the given channel as the callback channel for container metadata.
// Any previously registered callback channel will get overwritten by repeatedly calling this method. The previous
// callback channel (if any) is returned by this function.
func (e *Store) RegisterContainerMetadataCallbackChannel(callbackChan chan<- ContainerMetadata) chan<- ContainerMetadata {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	oldChan := e.callbackChannel
	e.callbackChannel = callbackChan
	return oldChan
}

// LookupResult contains the result of a lookup operation.
type LookupResult struct {
	Entity         networkgraph.Entity
	ContainerPorts []uint16
	PortNames      []string
}

// LookupByEndpoint returns possible target deployments by endpoint (if any).
func (e *Store) LookupByEndpoint(endpoint net.NumericEndpoint) []LookupResult {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.lookupNoLock(endpoint)
}

// LookupByContainerID retrieves the deployment ID by a container ID.
func (e *Store) LookupByContainerID(containerID string) (ContainerMetadata, bool) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	metadata, ok := e.containerIDMap[containerID]
	return metadata, ok
}

func (e *Store) lookupNoLock(endpoint net.NumericEndpoint) (results []LookupResult) {
	for deploymentID, targetInfoSet := range e.endpointMap[endpoint] {
		result := LookupResult{
			Entity:         networkgraph.EntityForDeployment(deploymentID),
			ContainerPorts: make([]uint16, 0, len(targetInfoSet)),
		}
		for tgtInfo := range targetInfoSet {
			result.ContainerPorts = append(result.ContainerPorts, tgtInfo.ContainerPort)
			if tgtInfo.PortName != "" {
				result.PortNames = append(result.PortNames, tgtInfo.PortName)
			}
		}
		results = append(results, result)
	}

	if len(results) > 0 {
		return
	}

	for deploymentID := range e.ipMap[endpoint.IPAndPort.Address] {
		result := LookupResult{
			Entity:         networkgraph.EntityForDeployment(deploymentID),
			ContainerPorts: []uint16{endpoint.IPAndPort.Port},
		}
		results = append(results, result)
	}

	return
}

// RegisterPublicIPsListener registers a listener that listens on changes to the set of public IP addresses.
// It returns a boolean indicating whether the listener was actually unregistered (i.e., a return value of false
// indicates that the listener was already registered).
func (e *Store) RegisterPublicIPsListener(listener PublicIPsListener) bool {
	// This mutex is pretty broad in scope, but since registering listeners occurs rarely, it's better than adding
	// another mutex that would need to get locked separately.
	e.mutex.Lock()
	defer e.mutex.Unlock()

	oldLen := len(e.publicIPsListeners)
	e.publicIPsListeners[listener] = struct{}{}

	return len(e.publicIPsListeners) > oldLen
}

// UnregisterPublicIPsListener unregisters a previously registered listener for public IP events. It returns a boolean
// indicating whether the listener was actually unregistered (i.e., a return value of false indicates that the listener
// was not registered in the first place).
func (e *Store) UnregisterPublicIPsListener(listener PublicIPsListener) bool {
	e.mutex.Lock()
	defer e.mutex.Lock()

	oldLen := len(e.publicIPsListeners)
	delete(e.publicIPsListeners, listener)

	return len(e.publicIPsListeners) < oldLen
}

func (e *Store) incPublicIPRefNoLock(addr net.IPAddress) {
	refCnt := e.publicIPRefCounts[addr]
	if refCnt == nil {
		refCnt = new(int)
		e.publicIPRefCounts[addr] = refCnt
		e.notifyPublicIPsListenersNoLock(PublicIPsListener.OnAdded, addr)
	}
	*refCnt++
}

func (e *Store) decPublicIPRefNoLock(addr net.IPAddress) {
	refCnt := e.publicIPRefCounts[addr]
	if refCnt == nil {
		utils.Should(errors.New("public IP has zero refcount already"))
		return
	}

	*refCnt--
	if *refCnt == 0 {
		delete(e.publicIPRefCounts, addr)
		e.notifyPublicIPsListenersNoLock(PublicIPsListener.OnRemoved, addr)
	}
}

func (e *Store) notifyPublicIPsListenersNoLock(notifyFunc func(PublicIPsListener, net.IPAddress), ip net.IPAddress) {
	for listener := range e.publicIPsListeners {
		notifyFunc(listener, ip)
	}
}
