package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/spl3g/tcpinger/pinger"
	"github.com/spl3g/tcpinger/utils"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var serverError = ErrorResponse{"internal server error"}

func handleCheck(validate *validator.Validate) http.HandlerFunc {
	type request struct {
		IP      string         `json:"ip" validate:"required,ip|hostname"`
		Port    uint           `json:"port" validate:"required,gte=0,lte=65535"`
		Timeout utils.Duration `json:"timeout" validate:"required"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := utils.ReadJSON[request](r.Body)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{err.Error()})
			return
		}

		err = validate.StructCtx(r.Context(), &req)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), req.Timeout.Duration)
		defer cancel()

		pgr := pinger.NewTCPPinger(req.IP, req.Port)
		checkResp, err := pgr.Check(ctx)
		if errors.Is(err, &pinger.CheckDNSError{}) {
			utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		} else if err != nil {
			log.Printf("error: %s", err)
			utils.WriteJSON(w, http.StatusInternalServerError, serverError)
		} else {
			checkResp.FormatJSON(w)
		}
	})
}
