package routes

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/kr/pretty"
)

//--
// Request and Response payloads for the REST api.
//
// The payloads embed the data model objects an
//
// In a real-world project, it would make sense to put these payloads
// in another file, or another sub-package.
//--

// User fixture data
var users = []*User{
	{ID: "100", Email: "Peter", Password: "Peter"},
	{ID: "200", Email: "Julia", Password: "Peter"},
}

// Projects fixture data
var projects = []*Project{
	{
		Owner:            "1",
		Id:               "1",
		ProjectName:      "My first Project",
		HasuraConsoleURL: "http://3.22.214.239:8080/console",
		PostgresURL:      "postgresql://market_data_access:1234@3.22.214.239:5432/market_data",
	},
}

type UserPayload struct {
	*User
	// Role string `json:"role"`
}

// ArticleRequest is the request payload for Article data model.
//
// NOTE: It's good practice to have well defined request and response payloads
// so you can manage the specific inputs and outputs for clients, and also gives
// you the opportunity to transform data on input or output, for example
// on request, we'd like to protect certain fields and on output perhaps
// we'd like to include a computed field based on other values that aren't
// in the data model. Also, check out this awesome blog post on struct composition:
// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
type SignupRequest struct {
	// User *User `json:"user,omitempty"`
	*User
	// ProtectedID string `json:"id"` // override 'id' json to have more control
}

// Create a struct that models the structure of a user, both in the request body, and in the DB
type User struct {
	ID       string
	Email    string
	Password string
}

type Project struct {
	Id               string
	Owner            string
	ProjectName      string
	HasuraConsoleURL string
	PostgresURL      string
}

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Username string
	Password string
}

// ArticleResponse is the response payload for the Article data model.
// See NOTE above in ArticleRequest as well.
//
// In the ArticleResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type SignupResponse struct {
	*User
	Token   string
	Elapsed int
}

type ProjectResponse struct {
	*Project
	Elapsed int
}

func NewUserCreatedResponse(user *User, token string) *SignupResponse {
	resp := &SignupResponse{
		User:  user,
		Token: token,
	}

	return resp
}

func GetProjectResponse(project *Project) *ProjectResponse {
	resp := &ProjectResponse{
		Project: project,
	}
	return resp
}

func ProjectListResponse(projects []*Project) []render.Renderer {
	list := []render.Renderer{}
	for _, project := range projects {
		list = append(list, GetProjectResponse(project))
	}
	return list
}

func (u *SignupRequest) Bind(r *http.Request) error {
	// u.User is nil if no User fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if u.User == nil {
		return errors.New("missing required Article fields.")
	}
	// a.User is nil if no Userpayload fields are sent in the request. In this app
	// this won't cause a panic, but checks in this Bind method may be required if
	// a.User or futher nested fields like a.User.Name are accessed elsewhere.

	// just a post-process after a decode..
	u.User.ID = "" // unset the protected ID
	// a.Article.Title = strings.ToLower(a.Article.Title) // as an example, we down-case
	return nil
}

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func dbNewUser(user *User) (string, error) {
	user.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
	users = append(users, user)
	return user.ID, nil
}

func dbGetProject(projectId string) (*Project, error) {
	return projects[0], nil
}

func ProjectCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// _, claims, _ := jwtauth.FromContext(r.Context())
		projectId := chi.URLParam(r, "projectId")
		project, err := dbGetProject(projectId)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "project", project)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rd *SignupResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}
func (rd *ProjectResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, ProjectListResponse(projects)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// GetProject returns the specific Project. You'll notice it just
// fetches the Project right off the context, as its understood that
// if we made it this far, the Project must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func GetProject(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the project
	// context because this handler is a child of the ProjectCtx
	// middleware. The worst case, the recoverer middleware will save us.
	project := r.Context().Value("project").(*Project)

	if err := render.Render(w, r, GetProjectResponse(project)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func ChiSignupHandler(w http.ResponseWriter, r *http.Request) {
	data := &SignupRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user := &User{
		Email:    data.Email,
		Password: data.Password,
	}
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": data.Email})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
	token := tokenString
	dbNewUser(user)
	pretty.Println(users)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewUserCreatedResponse(user, token))
}
