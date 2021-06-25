package handler

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/emacampolo/link-tracker/internal/link"
	"github.com/emacampolo/link-tracker/internal/platform/shell"
)

type Link struct {
	linkService link.Service
}

func NewLink(l link.Service) *Link {
	return &Link{
		linkService: l,
	}
}

func (lnk *Link) Create() shell.Handler {
	type request struct {
		Link     string `json:"link"`
		Password string `json:"password"`
	}

	type response struct {
		ID int `json:"id"`
	}

	return func(req *shell.Request) (*shell.Response, error) {
		var r request
		if err := shell.Decode(req, &r); err != nil {
			return nil, err
		}

		if r.Link == "" {
			return nil, errors.New("link is missing")
		}

		if r.Password == "" {
			return nil, errors.New("password is missing")
		}

		l, err := lnk.linkService.Create(req.Context(), r.Link, r.Password)
		if err != nil {
			return nil, err
		}

		resp := response{
			ID: l.ID,
		}

		return shell.Respond(req.Context(), resp), nil
	}
}

func (lnk *Link) Redirect() shell.Handler {
	type request struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}

	type response struct {
		Msg string `json:"msg"`
	}

	return func(req *shell.Request) (*shell.Response, error) {
		var r request
		if err := shell.Decode(req, &r); err != nil {
			return nil, err
		}

		idParam := r.ID
		if idParam == "" {
			return nil, errors.New("is is missing")
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			return nil, err
		}

		if r.Password == "" {
			return nil, errors.New("password is missing")
		}

		ll, err := lnk.linkService.Redirect(req.Context(), id, r.Password)
		if err != nil {
			return nil, err
		}

		resp := response{
			Msg: fmt.Sprintf("Redirecting to:%s", ll.URL),
		}

		return shell.Respond(req.Context(), resp), nil
	}
}

func (lnk *Link) Metrics() shell.Handler {
	type request struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}

	type response struct {
		ID       int    `json:"id"`
		URL      string `json:"url"`
		Count    int    `json:"count"`
		Inactive bool   `json:"inactive"`
	}

	return func(req *shell.Request) (*shell.Response, error) {
		var r request
		if err := shell.Decode(req, &r); err != nil {
			return nil, err
		}

		idParam := r.ID
		if idParam == "" {
			return nil, errors.New("is is missing")
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			return nil, err
		}

		l, err := lnk.linkService.FindByID(req.Context(), id)
		if err != nil {
			return nil, err
		}

		resp := response{
			ID:       l.ID,
			URL:      l.URL,
			Count:    l.Count,
			Inactive: l.Inactive,
		}

		return shell.Respond(req.Context(), resp), nil
	}
}

func (lnk *Link) Inactivate() shell.Handler {
	type request struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}

	type response struct {
		Msg string `json:"msg"`
	}

	return func(req *shell.Request) (*shell.Response, error) {
		var r request
		if err := shell.Decode(req, &r); err != nil {
			return nil, err
		}

		idParam := r.ID
		if idParam == "" {
			return nil, errors.New("is is missing")
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			return nil, err
		}

		if err := lnk.linkService.Inactivate(req.Context(), id); err != nil {
			return nil, err
		}

		resp := response{
			Msg: fmt.Sprintf("Link: %d deleted", id),
		}

		return shell.Respond(req.Context(), resp), nil
	}
}
