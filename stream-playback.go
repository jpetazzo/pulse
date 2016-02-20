package pulse

// #include "client.h"
// #cgo pkg-config: libpulse
import "C"

import (
    "io"
    // "unsafe"
    // "log"
)

type PlaybackStream struct {
    *Stream
}

func NewPlaybackStream(client *Client, name string, source io.Reader) *PlaybackStream {
    rv := &PlaybackStream{
        Stream:     NewStream(client, name),
    }

    rv.Source = source

    return rv
}

func (self *PlaybackStream) Initialize() error {
    if err := self.Stream.Initialize(); err != nil {
        return err
    }

    C.pa_stream_set_state_callback(self.Stream.toNative(), (C.pa_stream_notify_cb_t)(C.pulse_stream_state_callback), self.Stream.ToUserdata())
    C.pa_stream_set_write_callback(self.Stream.toNative(), (C.pa_stream_request_cb_t)(C.pulse_stream_write_callback), self.Stream.ToUserdata())

    go func(){
        C.pa_stream_connect_playback(self.Stream.toNative(), nil, nil, (C.pa_stream_flags_t)(0), nil, nil)
    }()

//  block until a terminal stream state is reached; successful or otherwise
    select {
    case err := <-self.Stream.state:
        return err
    }
}