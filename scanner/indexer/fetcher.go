package indexer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	gosync "sync"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/quay/claircore"
	"github.com/quay/claircore/indexer"
	"github.com/quay/zlog"
	"github.com/stackrox/scanner/v4/internal/sync"
	"golang.org/x/sync/errgroup"
)

var (
	_ indexer.FetchArena = (*localFetchArena)(nil)
	_ indexer.Realizer   = (*localFetcher)(nil)
)

// localFetchArena implements indexer.FetchArena.
//
// It is designed to minimize layer downloads and attempt
// to only download a layer's contents once.
//
// A localFetchArena must not be copied after first use,
// and it should be initialized via newLocalFetchArena.
type localFetchArena struct {
	// root is the directory to which layers are downloaded.
	root string
	// ko is used to ensure each layer is only downloaded once.
	ko sync.KeyedOnce[string]

	// mu protects rc.
	mu gosync.Mutex
	// rc is a map of digest to refcount.
	rc map[string]int
}

// newLocalFetchArena initializes a new localFetchArena.
func newLocalFetchArena(root string) *localFetchArena {
	return &localFetchArena{
		root: root,
		rc:   make(map[string]int),
	}
}

// Realizer returns an indexer.Realizer.
func (f *localFetchArena) Realizer(_ context.Context) indexer.Realizer {
	return &localFetcher{
		f: f,
	}
}

// Get downloads the image's manifest and returns the related claircore.Manifest.
//
// Get also downloads each previously unseen layer of the image into the arena's root directory.
func (f *localFetchArena) Get(ctx context.Context, image string, opts ...Option) (*claircore.Manifest, error) {
	// Parse the image name before doing anything else,
	// as there is no reason to do anything if the image is not properly referenced.
	ref, err := name.ParseReference(image)
	if err != nil {
		return nil, err
	}

	o := makeOptions(opts...)
	// Fetch the image's manifest from the registry.
	desc, err := remote.Get(ref, remote.WithContext(ctx), remote.WithAuth(o.auth), remote.WithPlatform(o.platform))
	if err != nil {
		return nil, err
	}

	img, err := desc.Image()
	if err != nil {
		return nil, err
	}
	d, err := img.Digest()
	if err != nil {
		return nil, err
	}
	// Convert the image manifest's digest to a claircore.Digest.
	ccd, err := claircore.ParseDigest(d.String())
	if err != nil {
		return nil, fmt.Errorf("parsing manifest digest %s: %w", d.String(), err)
	}

	manifest := &claircore.Manifest{
		Hash: ccd,
	}

	layers, err := img.Layers()
	if err != nil {
		return nil, err
	}
	manifest.Layers = make([]*claircore.Layer, len(layers))
	for i, layer := range layers {
		d, err := layer.Digest()
		if err != nil {
			return nil, err
		}
		// Convert the layer's digest to a claircore.Digest.
		ccd, err := claircore.ParseDigest(d.String())
		if err != nil {
			return nil, fmt.Errorf("parsing layer digest %s: %w", d.String(), err)
		}
		manifest.Layers[i] = &claircore.Layer{
			Hash: ccd,
		}
	}

	g, ctx := errgroup.WithContext(ctx)
	// Asynchronously download the layers, if needed.
	// This is done in a separate loop from the previous one for simpler error handling.
	for i := range layers {
		// This variable is set like this to prevent errors when reusing a loop variable.
		layer := layers[i]
		ccLayer := manifest.Layers[i]
		g.Go(f.realizeLayer(ctx, ccLayer, layer))
	}
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("could not realize layer(s) for image manifest %s: %w", manifest.Hash.String(), err)
	}

	return manifest, nil
}

// realizeLayer returns a function which downloads the layer once.
//
// The function attempts to increment the digest's ref count for each call.
func (f *localFetchArena) realizeLayer(ctx context.Context, ccLayer *claircore.Layer, layer v1.Layer) func() error {
	d := ccLayer.Hash.String()
	return func() error {
		path := filepath.Join(f.root, d)
		var tmp string

		select {
		case <-ctx.Done():
			return ctx.Err()
		case res := <-f.ko.DoChan(d, func() (any, error) {
			// Only the first call to DoChan will make it here.
			return f.downloadOnce(ctx, d, layer)
		}):
			if err := res.Err; err != nil {
				return fmt.Errorf("could not download layer %s: %w", d, err)
			}
			tmp = res.V.(string)
		}

		f.mu.Lock()
		defer f.mu.Unlock()

		ct, ok := f.rc[d]
		// Is this the first time we reference the layer's file?
		if !ok {
			// Did the file get removed while we were waiting on the lock?
			if _, err := os.Stat(tmp); errors.Is(err, os.ErrNotExist) {
				return err
			}
			// Move the file to its final path.
			if err := os.Rename(tmp, path); err != nil {
				return fmt.Errorf("moving layer from temporary to final path: %w", err)
			}
		}

		ct++
		f.rc[d] = ct

		// Set the layer's URI to the local filepath.
		ccLayer.URI = path

		return nil
	}
}

// downloadOnce downloads the contents of the layer into
// the arena's root directory at a temporary path.
func (f *localFetchArena) downloadOnce(ctx context.Context, digest string, layer v1.Layer) (string, error) {
	// Write the uncompressed layer, as ClairCore's indexer assumes the layer is uncompressed.
	uncompressed, err := layer.Uncompressed()
	if err != nil {
		return "", fmt.Errorf("fetching layer %s: %w", digest, err)
	}
	defer func() {
		// TODO: consider logging failures as a warning
		// and/or tracking metrics.
		_ = uncompressed.Close()
	}()

	rm := true
	file, err := os.CreateTemp(f.root, "fetch.*")
	if err != nil {
		return "", fmt.Errorf("creating temp file for layer %s: %w", digest, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			zlog.Warn(ctx).Err(err).Msg("unable to close layer file")
		}
		if rm {
			if err := os.Remove(file.Name()); err != nil {
				zlog.Warn(ctx).Err(err).Msg("unable to remove unsuccessful layer fetch")
			}
		}
	}()

	_, err = io.Copy(file, uncompressed)
	if err != nil {
		return "", fmt.Errorf("writing contents of layer %s into temp path: %w", digest, err)
	}

	rm = false
	return file.Name(), nil
}

// forget decrements the layer's refcount and "forgets" the layer
// once the refcount reaches zero.
func (f *localFetchArena) forget(d string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	ct, ok := f.rc[d]
	if !ok {
		return nil
	}

	ct--
	if ct == 0 {
		delete(f.rc, d)
		defer f.ko.Forget(d)
		return os.Remove(filepath.Join(f.root, d))
	}

	f.rc[d] = ct

	return nil
}

// Close removes all files left in the arena.
//
// It's not an error to have active fetchers, but may cause errors to have files
// unlinked underneath their users.
func (f *localFetchArena) Close(ctx context.Context) error {
	ctx = zlog.ContextWithValues(ctx,
		"component", "indexer/fetchArena.Close",
		"arena", f.root)

	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.rc) != 0 {
		zlog.Warn(ctx).
			Int("count", len(f.rc)).
			Msg("seem to have active fetchers")
		zlog.Info(ctx).
			Msg("clearing arena")
	}

	var errs []error
	for d := range f.rc {
		delete(f.rc, d)
		f.ko.Forget(d)
		if err := os.Remove(filepath.Join(f.root, d)); err != nil {
			errs = append(errs, err)
		}
	}
	if err := errors.Join(errs...); err != nil {
		return err
	}

	return nil
}

type localFetcher struct {
	f *localFetchArena
	// clean lists the layer hashes to clean up once no longer needed.
	clean []string
}

// Realize populates the local filepath for each layer.
//
// It is assumed the layer's URI is the local filesystem path to the layer.
func (f *localFetcher) Realize(_ context.Context, ls []*claircore.Layer) error {
	f.clean = make([]string, len(ls))
	for _, l := range ls {
		f.clean = append(f.clean, l.Hash.String())
		if err := l.SetLocal(l.URI); err != nil {
			return err
		}
	}
	return nil
}

// Close marks all the layers' backing files as unused.
//
// This method may actually delete the backing files.
func (f *localFetcher) Close() error {
	var errs []error
	for _, d := range f.clean {
		if err := f.f.forget(d); err != nil {
			errs = append(errs, err)
		}
	}
	if err := errors.Join(errs...); err != nil {
		return err
	}

	return nil
}
