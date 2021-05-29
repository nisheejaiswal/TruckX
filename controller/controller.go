package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/TruckX/constants"
	"github.com/TruckX/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func RunController() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/vehicle", vehicleHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/vehicle/all", vehicleListHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/vehicle/{imei}", vehicleGetHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/vehicle/power/{imei}", vehiclePowerHandler).Methods(http.MethodPatch)
	r.HandleFunc("/api/v1/vehicle/{imei}/event/{alarm_type}", vehicleEventHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/vehicle/login/{imei}", vehicleLoginHandler).Methods(http.MethodPost)

	r.HandleFunc("/api/v1/vehicle/video/{imei}", videoHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/vehicle/video/{imei}", videoGetHandler).Methods(http.MethodGet)

	r.HandleFunc("/api/v1/location/{name}", locationHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/location/{name}", locationGetHandler).Methods(http.MethodGet)

	r.HandleFunc("/api/v1/dashcam/{imei}/{command}", dashcamPostHandler).Methods(http.MethodPost)

	corsObj := handlers.AllowedOrigins([]string{"*"})
	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}),
		handlers.AllowedMethods([]string{"Access-Control-Allow-Methods", "POST", "OPTIONS", "GET", "PUT"}),
		handlers.AllowCredentials(),
		corsObj,
	)

	err := http.ListenAndServe(":8081", corsHandler(contentTypeMiddleware(r)))
	if err != nil {
		log.Fatalln(err)
	}
}

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func errorParser(w http.ResponseWriter, errStr string) {
	for _, errorResponse := range constants.ErrorCategory {
		if strings.Contains(errStr, errorResponse.Error) {

			resp := strings.Split(errStr, errorResponse.Error)
			if len(resp) == 2 {
				responseWithStaus(w, errorResponse.StatusCode, models.ErrorResponse{
					Error:      strings.TrimSpace(resp[1]),
					StatusCode: errorResponse.StatusCode,
				})
				return
			}
		}
	}

	responseWithStaus(w, http.StatusInternalServerError, models.ErrorResponse{
		Error:      errStr,
		StatusCode: http.StatusInternalServerError,
	})
}

func responseWithStaus(w http.ResponseWriter, statusCode int, response interface{}) {
	message, err := json.Marshal(response)
	if err != nil {
		errorParser(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(message)
}
