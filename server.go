package chromesvc

import "net/http"
import "encoding/json"

func StartServer() {
	http.HandleFunc("/render", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		req := &PageRenderTask{
			URL: q["url"][0],
		}
		resp, err := RenderPage(req)
		if err != nil {
			http.Error(w, "No can do!", http.StatusInternalServerError)
		}
		encode := json.NewEncoder(w)
		encode.SetIndent("", "  ")
		encode.Encode(resp)
	})

	http.ListenAndServe("0.0.0.0:1234", nil)
}
