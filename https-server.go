// test with: curl --cacert rootCa.crt https://localhost:443
// test with modified hosts file: curl --cacert rootCa.crt https://www.mydomain.com
// curl --tlsv1.1 --cacert rootCa.crt https://www.mydomain.com

package main

import (
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func main() {
	http.HandleFunc("/", HelloServer)

	err := http.ListenAndServeTLS(":443", "mydomain.com.crt", "mydomain.com.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
