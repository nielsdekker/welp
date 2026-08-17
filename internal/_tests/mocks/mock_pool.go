package mocks

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path"

	"github.com/nielsdekker/welp/internal/requests"
)

//go:embed resources
var resources embed.FS
var _pool *mockPool

type mockPool struct {
	// Each mock resource will be mapped to it's own hostname. For example the
	// `/spa` folder will map to the `spa.test` hostname
	hosts map[string]*http.ServeMux
}

func GetPool() requests.Pool {
	if _pool == nil {
		p := mockPool{
			hosts: make(map[string]*http.ServeMux),
		}
		entries, err := fs.ReadDir(resources, "resources")
		if err != nil {
			fmt.Printf("Setting up the mocks failed: %v", err)
			os.Exit(1)
		}

		// Create a mux per directory in the mock resources
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			name := entry.Name()
			sub, err := fs.Sub(resources, path.Join("resources", name))
			if err != nil {
				fmt.Printf("Unable to create sub filesystem for %s: %v\n", name, err)
				os.Exit(1)
			}

			p.hosts[name+".test"] = toMux(name, sub)
		}

		_pool = &p
	}

	return _pool
}

func (p *mockPool) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	mux, ok := p.hosts[req.Host]
	if !ok {
		return nil, fmt.Errorf("No mock found for hostname: %s", req.Host)
	}

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	res := rec.Result()
	res.Request = req
	return res, nil
}

// Switch on the name, for certain folders special logic is used and
// implemented.
func toMux(name string, sub fs.FS) *http.ServeMux {
	switch name {
	case "spa":
		return toSpaMux(sub)
	default:
		return toStaticMux(sub)
	}
}

// Creates a simple mux that returns only known files in the mock
func toStaticMux(sub fs.FS) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServerFS(sub))
	return mux
}

// Creates a mux that returns the index.html for each unknown path
func toSpaMux(sub fs.FS) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServerFS(sub)
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if we can stat the file, make sure to remove the leading `/`
		// from the path because the sub fs will always throw an error if the
		// path starts with `/`.
		if _, err := fs.Stat(sub, r.URL.Path[1:]); err != nil {
			f, err := sub.Open("index.html")
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte{})
				return
			}

			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(200)
			body, _ := io.ReadAll(f)
			w.Write(body)
		} else {
			fileServer.ServeHTTP(w, r)
		}
	}))
	return mux
}
