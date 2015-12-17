// Autogenerated by Frugal Compiler (0.0.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package event

import (
	"bytes"
	"fmt"
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Workiva/frugal-go"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

type FFoo interface {
	// Ping the server.
	Ping(frugal.Context) (err error)
	// Blah the server.
	Blah(frugal.Context, int32, string, *Event) (r int64, err error)
}

type FFooClient struct {
	FTransport       frugal.FTransport
	FProtocolFactory *frugal.FProtocolFactory
	InputProtocol    *frugal.FProtocol
	OutputProtocol   *frugal.FProtocol
	mu               sync.Mutex
}

func NewFFooClient(t frugal.FTransport, f *frugal.FProtocolFactory) *FFooClient {
	t.SetRegistry(frugal.NewClientRegistry())
	return &FFooClient{
		FTransport:       t,
		FProtocolFactory: f,
		InputProtocol:    f.GetProtocol(t),
		OutputProtocol:   f.GetProtocol(t),
	}
}

// Ping the server.
func (f *FFooClient) Ping(ctx frugal.Context) (err error) {
	oprot := f.OutputProtocol
	if oprot == nil {
		oprot = f.FProtocolFactory.GetProtocol(f.FTransport)
		f.OutputProtocol = oprot
	}
	errorC := make(chan error, 1)
	resultC := make(chan struct{}, 1)
	if err = f.FTransport.Register(ctx, f.recvPingHandler(ctx, resultC, errorC)); err != nil {
		return
	}
	f.mu.Lock()
	if err = oprot.WriteRequestHeader(ctx); err != nil {
		f.mu.Unlock()
		f.FTransport.Unregister(ctx)
		return
	}
	if err = oprot.WriteMessageBegin("ping", thrift.CALL, 0); err != nil {
		f.mu.Unlock()
		f.FTransport.Unregister(ctx)
		return
	}
	args := FooPingArgs{}
	if err = args.Write(oprot); err != nil {
		f.mu.Unlock()
		f.FTransport.Unregister(ctx)
		return
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		f.mu.Unlock()
		f.FTransport.Unregister(ctx)
		return
	}
	if err = oprot.Flush(); err != nil {
		f.mu.Unlock()
		f.FTransport.Unregister(ctx)
		return
	}
	f.mu.Unlock()

	select {
	case err = <-errorC:
		return
	case <-resultC:
		f.FTransport.Unregister(ctx)
		return
	}
}

func (f *FFooClient) recvPingHandler(ctx frugal.Context, resultC chan<- struct{}, errorC chan<- error) frugal.AsyncCallback {
	return func(tr thrift.TTransport) error {
		iprot := f.FProtocolFactory.GetProtocol(tr)
		if err := iprot.ReadResponseHeader(ctx); err != nil {
			errorC <- err
			return err
		}
		method, mTypeId, _, err := iprot.ReadMessageBegin()
		if err != nil {
			errorC <- err
			return err
		}
		if method != "ping" {
			err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "ping failed: wrong method name")
			errorC <- err
			return err
		}
		if mTypeId == thrift.EXCEPTION {
			error0 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
			var error1 error
			error1, err = error0.Read(iprot)
			if err != nil {
				errorC <- err
				return err
			}
			if err = iprot.ReadMessageEnd(); err != nil {
				errorC <- err
				return err
			}
			err = error1
			errorC <- err
			return err
		}
		if mTypeId != thrift.REPLY {
			err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "ping failed: invalid message type")
			errorC <- err
			return err
		}
		result := FooPingResult{}
		if err = result.Read(iprot); err != nil {
			errorC <- err
			return err
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			errorC <- err
			return err
		}
		resultC <- struct{}{}
		return nil
	}
}

// Blah the server.
func (f *FFooClient) Blah(ctx frugal.Context, num int32, str string, event *Event) (r int64, err error) {
	oprot := f.OutputProtocol
	if oprot == nil {
		oprot = f.FProtocolFactory.GetProtocol(f.FTransport)
		f.OutputProtocol = oprot
	}
	errorC := make(chan error, 1)
	resultC := make(chan int64, 1)
	if err = f.FTransport.Register(ctx, f.recvBlahHandler(ctx, resultC, errorC)); err != nil {
		return
	}
	f.mu.Lock()
	if err = oprot.WriteRequestHeader(ctx); err != nil {
		f.mu.Unlock()
		f.FTransport.Unregister(ctx)
		return
	}
	if err = oprot.WriteMessageBegin("blah", thrift.CALL, 0); err != nil {
		f.mu.Unlock()
		f.FTransport.Unregister(ctx)
		return
	}
	args := FooBlahArgs{
		Num:   num,
		Str:   str,
		Event: event,
	}
	if err = args.Write(oprot); err != nil {
		f.mu.Unlock()
		f.FTransport.Unregister(ctx)
		return
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		f.mu.Unlock()
		f.FTransport.Unregister(ctx)
		return
	}
	if err = oprot.Flush(); err != nil {
		f.mu.Unlock()
		f.FTransport.Unregister(ctx)
		return
	}
	f.mu.Unlock()

	select {
	case err = <-errorC:
		return
	case r = <-resultC:
		f.FTransport.Unregister(ctx)
		return
	}
}

func (f *FFooClient) recvBlahHandler(ctx frugal.Context, resultC chan<- int64, errorC chan<- error) frugal.AsyncCallback {
	return func(tr thrift.TTransport) error {
		iprot := f.FProtocolFactory.GetProtocol(tr)
		if err := iprot.ReadResponseHeader(ctx); err != nil {
			errorC <- err
			return err
		}
		method, mTypeId, _, err := iprot.ReadMessageBegin()
		if err != nil {
			errorC <- err
			return err
		}
		if method != "blah" {
			err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "blah failed: wrong method name")
			errorC <- err
			return err
		}
		if mTypeId == thrift.EXCEPTION {
			error0 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
			var error1 error
			error1, err = error0.Read(iprot)
			if err != nil {
				errorC <- err
				return err
			}
			if err = iprot.ReadMessageEnd(); err != nil {
				errorC <- err
				return err
			}
			err = error1
			errorC <- err
			return err
		}
		if mTypeId != thrift.REPLY {
			err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "blah failed: invalid message type")
			errorC <- err
			return err
		}
		result := FooBlahResult{}
		if err = result.Read(iprot); err != nil {
			errorC <- err
			return err
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			errorC <- err
			return err
		}
		if result.Awe != nil {
			errorC <- result.Awe
			return result.Awe
		}
		resultC <- result.GetSuccess()
		return nil
	}
}

type FFooProcessor struct {
	processorMap map[string]frugal.FProcessorFunction
	handler      FFoo
	writeMu      *sync.Mutex
	errors       chan error
}

func (p *FFooProcessor) GetProcessorFunction(key string) (processor frugal.FProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return
}

func NewFFooProcessor(handler FFoo) *FFooProcessor {
	writeMu := &sync.Mutex{}
	errors := make(chan error, 1)
	p := &FFooProcessor{
		handler:      handler,
		processorMap: make(map[string]frugal.FProcessorFunction),
		writeMu:      writeMu,
		errors:       errors,
	}
	p.processorMap["ping"] = &fooFPing{
		handler: handler,
		writeMu: writeMu,
		errors:  errors,
	}
	p.processorMap["blah"] = &fooFBlah{
		handler: handler,
		writeMu: writeMu,
		errors:  errors,
	}
	return p
}

func (p *FFooProcessor) Errors() <-chan error {
	return p.errors
}

func (p *FFooProcessor) Process(iprot, oprot *frugal.FProtocol) {
	ctx, err := iprot.ReadRequestHeader()
	if err != nil {
		p.errors <- err
		return
	}
	name, _, _, err := iprot.ReadMessageBegin()
	if err != nil {
		p.errors <- err
		return
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		processor.Process(ctx, iprot, oprot)
		return
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x3 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	p.writeMu.Lock()
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, 0)
	x3.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	p.writeMu.Unlock()
	p.errors <- x3
}

type fooFPing struct {
	handler FFoo
	writeMu *sync.Mutex
	errors  chan<- error
}

func (p *fooFPing) Process(ctx frugal.Context, iprot, oprot *frugal.FProtocol) {
	args := FooPingArgs{}
	var err error
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		p.writeMu.Lock()
		oprot.WriteMessageBegin("ping", thrift.EXCEPTION, 0)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		p.writeMu.Unlock()
		p.errors <- err
		return
	}

	iprot.ReadMessageEnd()
	result := FooPingResult{}
	var err2 error
	if err2 = p.handler.Ping(ctx); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing ping: "+err2.Error())
		p.writeMu.Lock()
		oprot.WriteMessageBegin("ping", thrift.EXCEPTION, 0)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		p.writeMu.Unlock()
		p.errors <- err2
		return
	}
	p.writeMu.Lock()
	if err2 = oprot.WriteResponseHeader(ctx); err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageBegin("ping", thrift.REPLY, 0); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	p.writeMu.Unlock()
	if err != nil {
		p.errors <- err
	}
}

type fooFBlah struct {
	handler FFoo
	writeMu *sync.Mutex
	errors  chan<- error
}

func (p *fooFBlah) Process(ctx frugal.Context, iprot, oprot *frugal.FProtocol) {
	args := FooBlahArgs{}
	var err error
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		p.writeMu.Lock()
		oprot.WriteMessageBegin("blah", thrift.EXCEPTION, 0)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		p.writeMu.Unlock()
		p.errors <- err
		return
	}

	iprot.ReadMessageEnd()
	result := FooBlahResult{}
	var err2 error
	var retval int64
	if retval, err2 = p.handler.Blah(ctx, args.Num, args.Str, args.Event); err2 != nil {
		switch v := err2.(type) {
		case *AwesomeException:
			result.Awe = v
		default:
			x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing blah: "+err2.Error())
			p.writeMu.Lock()
			oprot.WriteMessageBegin("blah", thrift.EXCEPTION, 0)
			x.Write(oprot)
			oprot.WriteMessageEnd()
			oprot.Flush()
			p.writeMu.Unlock()
			p.errors <- err2
			return
		}
	} else {
		result.Success = &retval
	}
	p.writeMu.Lock()
	if err2 = oprot.WriteResponseHeader(ctx); err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageBegin("blah", thrift.REPLY, 0); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	p.writeMu.Unlock()
	if err != nil {
		p.errors <- err
	}
}
