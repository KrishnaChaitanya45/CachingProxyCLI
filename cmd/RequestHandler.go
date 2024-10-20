package cmd

import (
	"fmt"
	"io"
	"log"
	"time"

	"net/http"
	"os"
	"sync"
)

type RequestHandler struct {
	cache  map[string]string
	mu     sync.RWMutex
	origin string
}

func (r *RequestHandler) ReadFromCache(w http.ResponseWriter, fileName string) {
	response, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("FAILED TO READ FROM THE CACHE %s", err.Error())
		http.Error(w, "FAILED TO READ FROM THE CACHE", http.StatusInternalServerError)
		return
	}
	fmt.Printf("\n\n👽 RETURNED THE CACHED RESPONSE 👽\n\n")
	w.Header().Add("X-Cache", "HIT")
	w.Write(response)
}

func (r *RequestHandler) SetCache(w http.ResponseWriter, urlPath string) {

	//? Doesn't exists in the cache
	//? Send Request to the origin and fetch the data
	originUrl := fmt.Sprintf("%s%s", r.origin, urlPath)
	newReq, err := http.NewRequest(http.MethodGet, originUrl, nil)
	if err != nil {
		log.Printf("ERROR WHILE CREATING A NEW REQUEST %s", err.Error())
		http.Error(w, "Unable to create a new request", http.StatusInternalServerError)
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	os.Mkdir("cache", 0755)
	resp, err := client.Do(newReq)
	if err != nil {
		log.Printf("ERROR WHILE SENDING THE REQUEST %s", err.Error())
		http.Error(w, "Unable to fetch data from origin", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error response from origin", resp.StatusCode)
		return
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR WHILE READING FROM THE RESPONSE BODY %s", err.Error())
		http.Error(w, "Unable to read from the response body", http.StatusInternalServerError)
		return
	}
	go r.CreateAndWriteToFile(urlPath, response)
	r.mu.Lock()
	r.cache[urlPath] = fmt.Sprintf("cache%s.json", urlPath)
	r.mu.Unlock()
	fmt.Printf("\n\n🦆 CACHED THE RESPONSE 🦆\n\n")
	w.Header().Add("X-Cache", "MISS")
	w.Write(response)
}

func (r *RequestHandler) CreateAndWriteToFile(fileName string, content []byte) {
	file, err := os.Create(fmt.Sprintf("cache%s.json", fileName))
	if err != nil {
		log.Printf("ERROR WHILE CREATING A NEW FILE %s", err.Error())
		return
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		log.Printf("ERROR WHILE WRITING TO THE FILE %s", err.Error())

	}
}

func (r *RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete {
		go r.ClearCache()
		w.Header().Add("X-Cache", "CLEARED")
		w.Write([]byte("CLEARING CACHE"))
		return
	}
	urlPath := req.URL.Path
	r.mu.RLock()
	fileName, ok := r.cache[urlPath]
	r.mu.RUnlock()
	if ok {
		//? exists in the cache
		r.ReadFromCache(w, fileName)

	} else {
		r.SetCache(w, urlPath)

	}
}

func (r *RequestHandler) ClearCache() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache = make(map[string]string)
	err := os.RemoveAll("cache")
	if err != nil {
		log.Printf("FAILED TO CLEAR CACHE %s", err.Error())
		return
	}
	// Recreate the cache directory
	err = os.Mkdir("cache", 0755)
	if err != nil {
		log.Printf("FAILED TO RECREATE CACHE DIRECTORY %s", err.Error())
		return
	}
	log.Println("Cache cleared successfully")
}
