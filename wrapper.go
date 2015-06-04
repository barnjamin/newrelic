package newrelic

/*
#cgo LDFLAGS: -L./newrelic_sdk/lib
#cgo LDFLAGS: -lnewrelic-collector-client
#cgo LDFLAGS: -lnewrelic-common
#cgo LDFLAGS: -lnewrelic-transaction

#cgo CFLAGS: -Inewrelic_sdk/include

#include <stdlib.h>
#include "newrelic_collector_client.h"
#include "newrelic_common.h"
#include "newrelic_transaction.h"

void register_default_handler() {
    newrelic_register_message_handler(newrelic_message_handler);
}

long begin_datastore_segment(
	long transaction_id,
	long parent_segment_id,
	const char *table,
	const char *operation,
	const char *sql,
	const char *sql_trace_rollup_name) {
    return newrelic_segment_datastore_begin(transaction_id, parent_segment_id,
        table, operation, sql, sql_trace_rollup_name,
        newrelic_basic_literal_replacement_obfuscator);

}


*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"unsafe"
)

const (
	ROOT_SEGMENT = 0 //TODO:: can we import C constants?

	SELECT = "select"
	INSERT = "insert"
	UPDATE = "update"
	DELETE = "delete"
)

var (
	ErrOther                 = errors.New("Newrelic Error: Other")
	ErrDisabled              = errors.New("Newrelic Error: Disabled")
	ErrInvalidParam          = errors.New("Newrelic Error: Invalid Param")
	ErrInvalidID             = errors.New("Newrelic Error: Invalid ID")
	ErrTransactionNotStarted = errors.New("Newrelic Error: Transaction Not Started")
	ErrTransactionInProgress = errors.New("Newrelic Error: Transaction In Progress")
	ErrTransactionNotNamed   = errors.New("Newrelic Error: Transaction Not Named")
)

// Call this to get the specific error
func DecodeError(res int) error {

	if res >= 0 {
		return nil
	}

	switch res {
	case -0x10001:
		return ErrOther
	case -0x20001:
		return ErrDisabled
	case -0x30001:
		return ErrInvalidParam
	case -0x30002:
		return ErrInvalidID
	case -0x40001:
		return ErrTransactionNotStarted
	case -0x40002:
		return ErrTransactionInProgress
	case -0x40003:
		return ErrTransactionNotNamed
	}

	return fmt.Errorf("Newrelic Error: Undefined (%d)", res)
}

func Initialize(name string) {
	C.register_default_handler()

	ckey := C.CString(os.Getenv("NEWRELIC_API_KEY"))
	defer C.free(unsafe.Pointer(ckey))

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	clang := C.CString("Go")
	defer C.free(unsafe.Pointer(clang))

	cver := C.CString(runtime.Version())
	defer C.free(unsafe.Pointer(cver))

	C.newrelic_init(ckey, cname, clang, cver)
}

func Stop(reason string) {
	creason := C.CString(reason)
	defer C.free(unsafe.Pointer(creason))

	C.newrelic_request_shutdown(creason)
}

// Automatically enabled, only use this to turn off
// or back on after you turned it off (0, 1)
func EnableInsturmentation(set int) {
	C.newrelic_enable_instrumentation(C.int(set))
}

func RecordMetric(metric string, value float64) (int, error) {
	cmetric := C.CString(metric)
	defer C.free(unsafe.Pointer(cmetric))

	res := C.newrelic_record_metric(cmetric, C.double(value))
	return int(res), DecodeError(int(res))
}

func RecordCPUUsage(userTimeSec, usagePerc float64) (int, error) {
	res := C.newrelic_record_cpu_usage(C.double(userTimeSec), C.double(usagePerc))
	return int(res), DecodeError(int(res))
}

func RecordMemoryUsage(mb float64) (int, error) {
	res := C.newrelic_record_memory_usage(C.double(mb))
	return int(res), DecodeError(int(res))
}

func StartTransaction() (int, error) {
	id := C.newrelic_transaction_begin()
	return int(id), DecodeError(int(id))
}

func SetTransactionName(id int, name string) (int, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	res := C.newrelic_transaction_set_name(C.long(id), cname)
	return int(res), DecodeError(int(res))
}

func SetTransactionTypeWeb(tid int) (int, error) {
	res := C.newrelic_transaction_set_type_web(C.long(tid))
	return int(res), DecodeError(int(res))
}

func SetTransactionTypeOther(tid int) (int, error) {
	res := C.newrelic_transaction_set_type_other(C.long(tid))
	return int(res), DecodeError(int(res))
}

func SetTransactionCategory(tid int, category string) (int, error) {
	ccategory := C.CString(category)
	defer C.free(unsafe.Pointer(ccategory))

	res := C.newrelic_transaction_set_category(C.long(tid), ccategory)
	return int(res), DecodeError(int(res))
}

func SetTransactionError(tid int, errType, message, trace, traceDelimiter string) (int, error) {
	cerrType := C.CString(errType)
	defer C.free(unsafe.Pointer(cerrType))

	cmessage := C.CString(message)
	defer C.free(unsafe.Pointer(cmessage))

	ctrace := C.CString(trace)
	defer C.free(unsafe.Pointer(ctrace))

	ctraceDelimiter := C.CString(traceDelimiter)
	defer C.free(unsafe.Pointer(ctraceDelimiter))

	res := C.newrelic_transaction_notice_error(C.long(tid), cerrType, cmessage, ctrace, ctraceDelimiter)
	return int(res), DecodeError(int(res))
}

func SetTransactionRequestURL(tid int, url string) (int, error) {
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))

	res := C.newrelic_transaction_set_request_url(C.long(tid), curl)
	return int(res), DecodeError(int(res))
}

func SetTransactionMaxSegments(tid, max int) (int, error) {
	res := C.newrelic_transaction_set_max_trace_segments(C.long(tid), C.int(max))
	return int(res), DecodeError(int(res))
}

func AddTransactionAttribute(tid int, key, val string) (int, error) {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	cval := C.CString(val)
	defer C.free(unsafe.Pointer(cval))

	res := C.newrelic_transaction_add_attribute(C.long(tid), ckey, cval)
	return int(res), DecodeError(int(res))
}

func EndTransaction(tid int) (int, error) {
	res := C.newrelic_transaction_end(C.long(tid))
	return int(res), DecodeError(int(res))
}

func StartGenericSegment(tid, parentId int, name string) (int, error) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))

	id := C.newrelic_segment_generic_begin(C.long(tid), C.long(parentId), cs)

	return int(id), DecodeError(int(id))
}

func StartDatastoreSegment(tid, parentId int, table, operation, sql, rollup_name string) (int, error) {
	ctable := C.CString(table)
	defer C.free(unsafe.Pointer(ctable))

	coperation := C.CString(operation)
	defer C.free(unsafe.Pointer(coperation))

	csql := C.CString(sql)
	defer C.free(unsafe.Pointer(csql))

	crollup_name := C.CString(rollup_name)
	defer C.free(unsafe.Pointer(crollup_name))

	id := C.begin_datastore_segment(C.long(tid), C.long(parentId), ctable, coperation, csql, crollup_name)

	return int(id), DecodeError(int(id))
}

func StartExternalSegment(tid, parentId int, host, name string) (int, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	chost := C.CString(host)
	defer C.free(unsafe.Pointer(chost))

	id := C.newrelic_segment_external_begin(C.long(tid), C.long(parentId), cname, chost)

	return int(id), DecodeError(int(id))
}

func EndSegment(tid, sid int) (int, error) {
	res := C.newrelic_segment_end(C.long(tid), C.long(sid))
	return int(res), DecodeError(int(res))
}
