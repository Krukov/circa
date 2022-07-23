package manage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type Storage struct {
	Name    string `json:"name"`
	Setting string `json:"setting"`
}

var storagePath = regexp.MustCompile("^/api/storage/(.+)$")

func (cm *configManage) Storages(w http.ResponseWriter, r *http.Request) {
	m := storagePath.FindStringSubmatch(r.URL.Path)
	w.Header().Set("content-type", "application/json")
	if m == nil {
		if r.Method == http.MethodGet {
			cm.GetStorages(w, r)
			return
		}
		if r.Method == http.MethodPost {
			cm.CreateStorage(w, r)
			return
		}
		writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
		return
	}
	name := m[1]
	if r.Method == http.MethodGet {
		cm.GetStorage(w, r, name)
		return
	}
	if r.Method == http.MethodPatch {
		cm.ChangeStorage(w, r, name)
		return
	}
	if r.Method == http.MethodDelete {
		cm.DelStorage(w, r, name)
		return
	}
	writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
}

func (cm *configManage) GetStorages(w http.ResponseWriter, r *http.Request) {
	jsonStorage := []Storage{}
	storages, _ := cm.c.GetStorages()
	for n, s := range storages {
		jsonStorage = append(jsonStorage, Storage{Name: n, Setting: s})
	}
	_ = json.NewEncoder(w).Encode(&jsonStorage)
}

func (cm *configManage) CreateStorage(w http.ResponseWriter, r *http.Request) {
	var storage Storage
	err := json.NewDecoder(r.Body).Decode(&storage)
	if err != nil {
		writeError(w, &errorResponse{"FORMAT_ERROR", fmt.Sprintf("Can't decode body %e", err)})
		return
	}
	// TODO: Validation settings
	err = cm.c.AddStorage(storage.Name, storage.Setting)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&errorResponse{"ERROR_STORAGE_ADDING", "Can't add storage"})
		return
	}
	if r.URL.Query().Get("sync") != "" {
		cm.c.Sync()
	}
	_ = json.NewEncoder(w).Encode(&storage)
}

func (cm *configManage) GetStorage(w http.ResponseWriter, r *http.Request, name string) {
	storages, _ := cm.c.GetStorages()
	for n, s := range storages {
		if n == name {
			_ = json.NewEncoder(w).Encode(&Storage{Name: name, Setting: s})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(&errorResponse{"NOT_FOUND", "Storage not found"})
}

func (cm *configManage) ChangeStorage(w http.ResponseWriter, r *http.Request, name string) {
	var storage Storage
	err := json.NewDecoder(r.Body).Decode(&storage)
	if err != nil {
		writeError(w, &errorResponse{"FORMAT_ERROR", fmt.Sprintf("Can't decode body %e", err)})
		return
	}
	// TODO: Validation settings
	err = cm.c.AddStorage(name, storage.Setting)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&errorResponse{"ERROR_STORAGE_ADDING", "Can't add storage"})
		return
	}
	if r.URL.Query().Get("sync") != "" {
		cm.c.Sync()
	}
	cm.GetStorage(w, r, name)
}

func (cm *configManage) DelStorage(w http.ResponseWriter, r *http.Request, name string) {
	storages, _ := cm.c.GetStorages()
	for n := range storages {
		if n == name {
			err := cm.c.DelStorage(name)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(&errorResponse{"ERROR_STORAGE_DELETE", "Can't delete storage"})
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(&errorResponse{"NOT_FOUND", "Storage not found"})
}
