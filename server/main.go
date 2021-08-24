package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"example.com/m/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/kr/pretty"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type key int

var DB *gorm.DB

const (
	requestIDKey key = 0
)

var (
	listenAddr string
	healthy    int32
)

var printroutes = flag.Bool("routes", false, "Generate router documentation")

func main() {

	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Route("/api/signup", func(r chi.Router) {
		r.Post("/", routes.ChiSignupHandler)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
		})
	})
	http.ListenAndServe(":5000", r)
	pretty.Println("Server is ready to handle requests at 5000")
	/*
		flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
		flag.Parse()
		fileServer := http.FileServer(http.Dir("./build"))

		logger := log.New(os.Stdout, "http: ", log.LstdFlags)
		logger.Println("Server is starting...")
		initDB()

		// Insert the middleware
		router := http.NewServeMux()
		router.Handle("/", index())
		router.Handle("/healthz", healthz())
		router.Handle("/build", fileServer)
		router.Handle("/api/signup/", routes.HandleSignup(DB))
		nextRequestID := func() string {
			return fmt.Sprintf("%d", time.Now().UnixNano())
		}

		handler := cors.Default().Handler(router)
		server := &http.Server{
			Addr:         listenAddr,
			Handler:      tracing(nextRequestID)(logging(logger)(handler)),
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		}

		done := make(chan bool)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		go func() {
			<-quit
			logger.Println("Server is shutting down...")
			atomic.StoreInt32(&healthy, 0)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			server.SetKeepAlivesEnabled(false)
			if err := server.Shutdown(ctx); err != nil {
				logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
			}
			close(done)
		}()

		logger.Println("Server is ready to handle requests at", listenAddr)
		atomic.StoreInt32(&healthy, 1)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
		}

		<-done
		logger.Println("Server stopped")
	*/
}

func initDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file \n", err)
	}

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		os.Getenv("PSQL_USER"), os.Getenv("PSQL_PASS"), os.Getenv("PSQL_DBNAME"), os.Getenv("PSQL_PORT"))

	log.Print("Connecting to PostgreSQL DB...", dsn)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	log.Println("Db info: ", DB)
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)

	}
	log.Println("connected")
}

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, World!")
	})
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path)
				buf, bodyErr := ioutil.ReadAll(r.Body)
				if bodyErr != nil {
					log.Print("bodyErr ", bodyErr.Error())
					http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
					return
				}

				rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
				rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
				log.Printf("BODY: %q", rdr1)
				r.Body = rdr2
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
