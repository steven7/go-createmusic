package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/steven7/go-createmusic/go/config"
	"github.com/steven7/go-createmusic/go/controllers"
	"github.com/steven7/go-createmusic/go/middleware"
	"github.com/steven7/go-createmusic/go/models"
	"net/http"
)


func main() {

	//boolPtr := flag.Bool("prod", true, "Provide this flag "+
	//	"in production. This ensures that a .config file is "+
	//	"provided before the application starts.")
	boolPtr := flag.Bool("prod", false, "Provide this flag "+
		"in production. This ensures that a .config file is "+
		"provided before the application starts.")
	flag.Parse()


	cfg := config.LoadConfig(*boolPtr)
	dbCfg := cfg.Database
	fmt.Println("trying with host ", dbCfg.Host)
	services, err := models.NewServices(

		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		// only log when not in prod
		models.WithLogMode(!cfg.IsProd()),
		// We want each of these services, but if we didn't need
		// one of them we could possibly skip that config func
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithTrack(),
		//models.WithImage(),
		//models.WithMusicFile(),
		models.WithFile(),
		models.WithOauth(),
	)

	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	// not set up
	//mgCfg := cfg.Mailgun
	//emailer := email.NewClient(
	//	email.WithSender("ImageCloud.com Support", "support@"+mgCfg.Domain),
	//	email.WithMailgun(mgCfg.Domain, mgCfg.APIKey, mgCfg.PublicAPIKey),
	//)

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	// emailer is nil because we arent using it now
	usersC := controllers.NewUsers(services.User, nil)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)
	//tracksC := controllers.NewTracksController(services.Track, services.Image, services.MusicFile, r)
	tracksC := controllers.NewTracksController(services.Track, services.File, r)

	//
	/*
	configs := make(map[string]*oauth2.Config)
	configs[models.OAuthDropbox] = &oauth2.Config{
		ClientID: 	  cfg.Dropbox.ID,
		ClientSecret: cfg.Dropbox.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL  :  cfg.Dropbox.AuthURL,
			TokenURL :  cfg.Dropbox.TokenURL,
		},
		RedirectURL: "http://localhost:3000/oauth/dropbox/callback",
	}
	oauthsC := controllers.NewOAuths(services.OAuth ,configs)
	 */

	//
	//
	//  services.DestructiveReset()
	//
	// be careful with this ^^^^^

	//
	// Middleware
	// csrf middleware
	//

	//b, err := rand.Bytes(32)
	//if err != nil {
	//	panic(err)
	//}
	// csrfMw := csrf.Protect(b) //, csrf.Secure(true)) // , csrf.Secure(cfg.IsProd()))
	// csrfMw := csrf.Protect(b, csrf.Secure(false)) //cfg.IsProd()))

	// this is for server side
	// csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

	corsMw := cors.New(cors.Options{
		AllowedHeaders: []string{"accept", "authorization", "content-type"},
		AllowedOrigins: []string{"http://localhost", "http://localhost:3000", "*"}, // * is for testing only not production
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.Handle("/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	// This will assign the page to the nor found handler
	var h http.Handler = http.Handler(staticC.NotFound)
	r.NotFoundHandler = h


	// Gallery routes

	//r.Handle("/galleries",
	//	requireUserMw.ApplyFn(galleriesC.Index)).
	//	Methods("GET").
	//	Name(controllers.IndexGalleries)
	//r.Handle("/galleries/new",
	//	requireUserMw.Apply(galleriesC.New)).
	//	Methods("GET")
	//r.Handle("/galleries",
	//	requireUserMw.ApplyFn(galleriesC.Create)).
	//	Methods("POST")
	//r.HandleFunc("/galleries/{id:[0-9]+}",
	//	galleriesC.Show).
	//	Methods("GET").
	//	Name(controllers.ShowGallery)
	//r.HandleFunc("/galleries/{id:[0-9]+}/edit",
	//	galleriesC.Edit).
	//	Methods("GET").
	//	Name(controllers.EditGallery)
	//r.HandleFunc("/galleries/{id:[0-9]+}/update",
	//	requireUserMw.ApplyFn(galleriesC.Update)).
	//	Methods("POST")
	//r.HandleFunc("/galleries/{id:[0-9]+}/delete",
	//	requireUserMw.ApplyFn(galleriesC.Delete)).
	//	Methods("POST")
	//r.HandleFunc("/galleries/{id:[0-9]+}/images",
	//	requireUserMw.ApplyFn(galleriesC.ImageUpload)).
	//	Methods("POST")
	//r.HandleFunc("/galleries/{id:[0-9]+}/images/link",
	//	requireUserMw.ApplyFn(galleriesC.ImageViaLink)).
	//	Methods("POST")


	// tracks

	// view tracks
	r.Handle("/tracks",
		// eventually take off this middleware to let the user preview the website
		requireUserMw.ApplyFn(tracksC.Index)).
		Methods("GET").
		Name(controllers.IndexTracks)

	// create new track
	r.Handle("/tracks/new",
		// eventually take off this middleware to let the user preview the website
		requireUserMw.Apply(tracksC.ChooseTypeView)).
		Methods("GET")


	r.HandleFunc("/tracks/{id:[0-9]+}/play",
		tracksC.Play).
		Methods("GET").
		Name(controllers.PlayTrack)

	// create new // phase one
	r.Handle("/tracks/createlocal",
		requireUserMw.ApplyFn(tracksC.CreateLocal)).
		Methods("POST")
	// create with python for singlepage js front end
	r.Handle("/tracks/createWithComposeAI",
		requireUserMw.ApplyFn(tracksC.CreateWithComposeAI)).
		Methods("POST")


	r.Handle("/tracks/createWithDJ",
		requireUserMw.ApplyFn(tracksC.ChooseDJOptions)).
		Methods("POST")

	//
	r.Handle("/tracks/createWithDJWorking",
		requireUserMw.ApplyFn(tracksC.CreateWithDJWorking)).
		Methods("POST")
	r.Handle("/tracks/createWithDJComplete",
		requireUserMw.ApplyFn(tracksC.CreateWithDJComplete)).
		Methods("POST")

	// create new looks like edit
	// when create song is pressed
	r.Handle("/tracks/createlocalcomplete",
		requireUserMw.ApplyFn(tracksC.CreateLocalComplete)).
		Methods("POST")


	// edit existing
	r.Handle("/tracks/{id:[0-9]+}/editLocalTrack",
		// eventually take off this middleware to let the user preview the website
		requireUserMw.ApplyFn(tracksC.EditLocal)).
		Methods("GET").//
		Name(controllers.EditTrack)
	r.Handle("/tracks/{id:[0-9]+}/editDJCreatedTrack",
		// eventually take off this middleware to let the user preview the website
		requireUserMw.ApplyFn(tracksC.EditDJ)).
		Methods("GET")

	//r.HandleFunc("/tracks/{id:[0-9]+}/images/link",
	//	requireUserMw.ApplyFn(galleriesC.ImageViaLink)).
	//	Methods("POST")

	//
	// create track
	//
	// upload
	r.HandleFunc("/tracks/{id:[0-9]+}/music",
		requireUserMw.ApplyFn(tracksC.MusicUpload)).
		Methods("POST")
	r.HandleFunc("/tracks/{id:[0-9]+}/images",
		requireUserMw.ApplyFn(tracksC.ImageUpload)).
		Methods("POST")

	// add to db -- when edit song pressed
	r.HandleFunc("/tracks/{id:[0-9]+}/create",
		requireUserMw.ApplyFn(tracksC.CreateLocalSongWithDB)).
		Methods("POST")
	// edit existing pressed
	r.HandleFunc("/tracks/{id:[0-9]+}/update",
		requireUserMw.ApplyFn(tracksC.EditLocalSongComplete)).
		Methods("POST")


	//r.Handle("/tracks",
	//	requireUserMw.ApplyFn(tracksC.Create)).
	//	Methods("POST")


	// Image routes
	//imageHandler := http.FileServer(http.Dir("./images/"))
	//r.PathPrefix("/images/").Handler(http.StripPrefix("/images/",imageHandler))
	imageHandler := http.FileServer(http.Dir("./userfiles/tracks/"))
	r.PathPrefix("/userfiles/tracks/").Handler(http.StripPrefix("/userfiles/tracks/",imageHandler))

	// file routes
	//imageHandler := http.FileServer(http.Dir("./images/"))
	//r.PathPrefix("/images/").Handler(http.StripPrefix("/images/",imageHandler))


	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete",
		requireUserMw.ApplyFn(galleriesC.ImageDelete)).
		Methods("POST")

	// Assets
	assetHandler := http.FileServer(http.Dir("./assets"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)



	//
	// API routes
	//
	usersCAPI := controllers.NewUsersAPI(services.User, nil)
	tracksCAPI := controllers.NewTracksAPI(services.Track, services.File, r)

	//
	//
	//  API routes
	//
	//

	//
	// Auth API
	r.HandleFunc("/api/auth/login", usersCAPI.AuthenticateWithAPI).Methods("POST")
	r.HandleFunc("/api/auth/signup", usersCAPI.CreateWithAPI).Methods("POST")

	//
	// Tracks API
	r.HandleFunc("/api/tracks/createlocal", tracksCAPI.CreateLocalWithAPI).Methods("POST")
	r.HandleFunc("/api/tracks/createWithComposeAI", tracksCAPI.CreateWithComposeAI).Methods("POST")
	//
	r.HandleFunc("/api/tracks/index", tracksCAPI.IndexWithAPI).Methods("POST")
	r.HandleFunc("/api/tracks/one", tracksCAPI.GetTrackWithAPI).Methods("POST")
	r.HandleFunc("/api/tracks/one/coverimage", tracksCAPI.GetTrackCoverFileWithAPI).Methods("POST")
	r.HandleFunc("/api/tracks/one/musicfile", tracksCAPI.GetTrackMusicFileWithAPI).Methods("POST")

	fmt.Printf("lol new line")
	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), corsMw.Handler(userMw.Apply(r)))
	//http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.Apply(r)))

}