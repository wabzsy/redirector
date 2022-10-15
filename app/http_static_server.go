package app

import (
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"time"
)

func (a *App) HTTPStaticServer() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "http",
		Short:   "http static server",
		Aliases: []string{"static"},
		Args:    cobra.NoArgs,
		Example: "./red http \n./red http -l 0.0.0.0:7777 \n ./red -d /dev/shm",
	}

	cmd.Flags().StringP("listen", "l", "0.0.0.0:8000", "http server listen address")
	cmd.Flags().StringP("dir", "d", "./", "target directory")

	cmd.RunE = a.HTTPStaticServerHandler
	return cmd
}

func (a *App) HTTPStaticServerHandler(cmd *cobra.Command, _ []string) error {
	addr, err := cmd.Flags().GetString("listen")
	if err != nil {
		return err
	}

	path, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}

	log.Println("Serving HTTP on " + addr + " (http://" + addr + ") ...")

	if err := http.ListenAndServe(addr, LoggerMiddleware(http.FileServer(http.Dir(path)))); err != nil {
		log.Println(err)
	}
	return nil
}

func LoggerMiddleware(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		targetMux.ServeHTTP(w, r)
		log.Printf("[%s] -- \"%s %s %s\" %s\n", r.RemoteAddr, r.Method, r.RequestURI, r.Proto, time.Since(now))
	})
}
