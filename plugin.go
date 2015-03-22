package timeshift

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/ugorji/go/codec"
	"github.com/yosisa/fluxion/buffer"
	"github.com/yosisa/fluxion/message"
	"github.com/yosisa/fluxion/plugin"
)

type Config struct {
	Listen string
}

type Plugin struct {
	env  *plugin.Env
	conf Config
	ts   *TimeSlice
}

func (p *Plugin) Init(env *plugin.Env) error {
	p.env = env
	p.ts = NewTimeSlice(500*500, func() interface{} {
		return NewTimeSlice(500, func() interface{} {
			return NewTimeSlice(1, func() interface{} {
				return NewMsgpackBuffer()
			})
		})
	})
	return env.ReadConfig(&p.conf)
}

func (p *Plugin) Start() error {
	go http.ListenAndServe(p.conf.Listen, p)
	return nil
}

func (p *Plugin) Encode(ev *message.Event) (buffer.Sizer, error) {
	if err := p.ts.Add(ev.Time.Unix(), ev); err != nil {
		p.env.Log.Error(err)
	}
	return nil, nil
}

func (p *Plugin) Write(l []buffer.Sizer) (int, error) {
	return len(l), nil
}

func (p *Plugin) Close() error {
	return nil
}

func (p *Plugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var since, until int64
	if s := r.URL.Query().Get("until"); s != "" {
		until, _ = strconv.ParseInt(s, 10, 64)
	}
	if s := r.URL.Query().Get("since"); s != "" {
		since, _ = strconv.ParseInt(s, 10, 64)
	}
	if until == 0 {
		until = time.Now().Unix()
	}
	if since == 0 {
		since = until - 60
	}

	var tag *regexp.Regexp
	if s := r.URL.Query().Get("tag"); s != "" {
		if tag, err = regexp.Compile(s); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	vals := p.ts.Range(since, until)
	for _, buf := range vals {
		dec := codec.NewDecoderBytes(buf.(*MsgpackBuffer).Bytes(), mh)
		for {
			var ev message.Event
			if err := dec.Decode(&ev); err != nil {
				break
			}
			if tag == nil || tag.MatchString(ev.Tag) {
				enc.Encode(&ev)
			}
		}
	}
}

func Factory() plugin.Plugin {
	return &Plugin{}
}
