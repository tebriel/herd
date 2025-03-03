package herd

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/spf13/viper"
)

type fakeProvider struct {
}

func (p *fakeProvider) Name() string {
	return "fake"
}

func (p *fakeProvider) Prefix() string {
	return "fake:"
}

func (p *fakeProvider) Equivalent(o HostProvider) bool {
	return false
}

func (p *fakeProvider) Load(ctx context.Context, lm LoadingMessage) (Hosts, error) {
	time.Sleep(2 * time.Millisecond)
	dl, ok := ctx.Deadline()
	if ok && time.Until(dl) < 0 {
		return nil, errors.New("context deadline exceeded")
	}
	h := NewHost("test-host", "", HostAttributes{"foo": "bar"})
	return Hosts{h}, nil
}

func (p *fakeProvider) ParseViper(v *viper.Viper) error {
	return nil
}

func TestNewRegistry(t *testing.T) {
	r := NewRegistry("/tmp", "/tmp")
	if len(r.providers) > 0 {
		t.Errorf("got %d providers, expected none", len(r.providers))
	}
}

func TestGetHosts(t *testing.T) {
	r := Registry{providers: []HostProvider{&fakeProvider{}, &fakeProvider{}}}
	err := r.LoadHosts(context.Background(), func(string, bool, error) {})
	if err != nil {
		t.Errorf("%t %v", err, err)
		t.Errorf("Could not load hosts: %s", err.Error())
	}
	if len(r.hosts) != 1 {
		t.Errorf("Hosts returned by multiple providers are not merged, got %d hosts instead of 1", len(r.hosts))
	}
}

func TestGetHostsTimeout(t *testing.T) {
	r := Registry{providers: []HostProvider{&fakeProvider{}, &fakeProvider{}}}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	err := r.LoadHosts(ctx, func(string, bool, error) {})
	if err == nil {
		t.Error("Expected a timeout")
	}
}
