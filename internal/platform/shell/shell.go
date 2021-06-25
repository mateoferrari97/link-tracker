package shell

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

const (
	MinActionLength = 1
	MinParamLength  = 2
)

type Application struct {
	in         io.Reader
	out        io.Writer
	handlerFns map[string]Handler
}

func New(in io.Reader, out io.Writer) *Application {
	return &Application{
		in:         in,
		out:        out,
		handlerFns: make(map[string]Handler),
	}
}

func (app *Application) Run() error {
	serverErrors := make(chan error, 1)
	serverInputs := make(chan string, 1)

	defer close(serverErrors)
	defer close(serverInputs)

	go func() {
		log.SetOutput(app.out)
		log.Println("API listening for instructions")

		scanner := bufio.NewScanner(app.in)
		for scanner.Scan() {
			serverInputs <- fmt.Sprintf(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			serverErrors <- err
		}
	}()

	for {
		select {
		case in := <-serverInputs:
			if err := app.handleInput(in); err != nil {
				_, _ = fmt.Fprintln(app.out, err.Error())
			}
		case err := <-serverErrors:
			return err
		}
	}
}

func (app *Application) handleInput(input string) error {
	in := strings.Trim(input, "\n")
	fields := strings.Fields(in)
	if len(fields) < MinActionLength {
		return errors.New("could not handle empty input")
	}

	action := fields[0]
	fields = fields[1:]

	params := make(map[string]string)
	for _, field := range fields {
		param := strings.Split(field, ":")
		if len(param) != MinParamLength {
			return errors.New("could not handle param")
		}

		paramName, paramValue := param[0], param[1]
		params[paramName] = paramValue
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req := Request{
		ctx:  context.Background(),
		body: body,
	}

	h, exist := app.handlerFns[action]
	if !exist {
		return errors.New("could not handle input: handler not found")
	}

	resp, err := h(&req)
	if err != nil {
		return err
	}

	return app.respond(resp)
}

func (app *Application) AddHandlerFunc(pattern string, h Handler) {
	app.handlerFns[pattern] = h
}

type Handler func(req *Request) (*Response, error)

type Request struct {
	ctx  context.Context
	body []byte
}

func (req *Request) Context() context.Context {
	if req.ctx != nil {
		return req.ctx
	}

	return context.Background()
}

func Decode(req *Request, v interface{}) error {
	if err := json.Unmarshal(req.body, v); err != nil {
		return fmt.Errorf("could not decode request: %v", err)
	}

	return nil
}

type Response struct {
	ctx  context.Context
	body interface{}
}

func Respond(ctx context.Context, resp interface{}) *Response {
	return &Response{
		ctx:  ctx,
		body: resp,
	}
}

func (app *Application) respond(resp *Response) error {
	body, err := json.Marshal(resp.body)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintln(app.out, string(body))

	return nil
}
