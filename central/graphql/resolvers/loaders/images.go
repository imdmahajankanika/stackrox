package loaders

import (
	"context"
	"errors"
	"reflect"

	"github.com/stackrox/rox/central/image/datastore"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/sync"
)

var imageLoadertype = reflect.TypeOf(storage.Image{})

func init() {
	RegisterTypeFactory(reflect.TypeOf(storage.Image{}), func() interface{} {
		return NewImageLoader(datastore.Singleton())
	})
}

// NewImageLoader creates a new loader for image data.
func NewImageLoader(ds datastore.DataStore) ImageLoader {
	return &imageLoaderImpl{
		loaded: make(map[string]*storage.Image),
		ds:     ds,
	}
}

// GetImageLoader returns the ImageLoader from the context if it exists.
func GetImageLoader(ctx context.Context) (ImageLoader, error) {
	loader, err := GetLoader(ctx, imageLoadertype)
	if err != nil {
		return nil, err
	}
	return loader.(ImageLoader), nil
}

// ImageLoader loads image data, and stores already loaded images for other ops in the same context to use.
type ImageLoader interface {
	FromIDs(ctx context.Context, ids []string) ([]*storage.Image, error)
	FromID(ctx context.Context, id string) (*storage.Image, error)
	FromQuery(ctx context.Context, query *v1.Query) ([]*storage.Image, error)

	CountFromQuery(ctx context.Context, query *v1.Query) (int32, error)
	CountAll(ctx context.Context) (int32, error)
}

// imageLoaderImpl implements the ImageDataLoader interface.
type imageLoaderImpl struct {
	lock   sync.RWMutex
	loaded map[string]*storage.Image

	ds datastore.DataStore
}

// FromIDs loads a set of images from a set of ids.
func (idl *imageLoaderImpl) FromIDs(ctx context.Context, ids []string) ([]*storage.Image, error) {
	images, err := idl.load(ctx, ids)
	if err != nil {
		return nil, err
	}
	return images, nil
}

// FromID loads an image from an ID.
func (idl *imageLoaderImpl) FromID(ctx context.Context, id string) (*storage.Image, error) {
	images, err := idl.load(ctx, []string{id})
	if err != nil {
		return nil, err
	}
	return images[0], nil
}

// FromQuery loads a set of images that match a query.
func (idl *imageLoaderImpl) FromQuery(ctx context.Context, query *v1.Query) ([]*storage.Image, error) {
	results, err := idl.ds.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	return idl.FromIDs(ctx, search.ResultsToIDs(results))
}

func (idl *imageLoaderImpl) CountFromQuery(ctx context.Context, query *v1.Query) (int32, error) {
	results, err := idl.ds.Search(ctx, query)
	if err != nil {
		return 0, err
	}
	return int32(len(results)), nil
}

func (idl *imageLoaderImpl) CountAll(ctx context.Context) (int32, error) {
	count, err := idl.ds.CountImages(ctx)
	return int32(count), err
}

func (idl *imageLoaderImpl) load(ctx context.Context, ids []string) ([]*storage.Image, error) {
	images, missing := idl.readAll(ids)
	if len(missing) > 0 {
		var err error
		images, err = idl.ds.GetImagesBatch(ctx, collectMissing(ids, missing))
		if err != nil {
			return nil, err
		}
		idl.setAll(images)
		images, missing = idl.readAll(ids)
	}
	if len(missing) > 0 {
		return nil, errors.New("not all images could be found")
	}
	return images, nil
}

func (idl *imageLoaderImpl) setAll(images []*storage.Image) {
	idl.lock.Lock()
	defer idl.lock.Unlock()

	for _, image := range images {
		idl.loaded[image.GetId()] = image
	}
}

func (idl *imageLoaderImpl) readAll(ids []string) (images []*storage.Image, missing []int) {
	idl.lock.RLock()
	defer idl.lock.RUnlock()

	for idx, id := range ids {
		image, isLoaded := idl.loaded[id]
		if !isLoaded {
			missing = append(missing, idx)
		} else {
			images = append(images, image)
		}
	}
	return
}

func collectMissing(ids []string, missing []int) []string {
	missingIds := make([]string, 0, len(missing))
	for _, missingIdx := range missing {
		missingIds = append(missingIds, ids[missingIdx])
	}
	return missingIds
}
