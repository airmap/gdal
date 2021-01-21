package ogr

/*
#include "go_ogr_wkb.h"
#include "gdal_version.h"
*/
import "C"
import (
	"unsafe"
)

type OGRDriver struct {
	cval C.OGRSFDriverH
}

// Fetch name of driver (file format)
func (driver OGRDriver) Name() string {
	name := C.OGR_Dr_GetName(driver.cval)
	return C.GoString(name)
}

// Attempt to open file with this driver
func (driver OGRDriver) Open(filename string, update int) (newDS DataSource, ok bool) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	ds := C.OGR_Dr_Open(driver.cval, cFilename, C.int(update))
	return DataSource{ds}, ds != nil
}

// Test if this driver supports the named capability
func (driver OGRDriver) TestCapability(capability string) bool {
	cString := C.CString(capability)
	defer C.free(unsafe.Pointer(cString))
	val := C.OGR_Dr_TestCapability(driver.cval, cString)
	return val != 0
}

// Create a new data source based on this driver
func (driver OGRDriver) Create(name string, options []string) (newDS DataSource, ok bool) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	ds := C.OGR_Dr_CreateDataSource(driver.cval, cName, (**C.char)(unsafe.Pointer(&opts[0])))
	return DataSource{ds}, ds != nil
}

// Create a new datasource with this driver by copying all layers of the existing datasource
func (driver OGRDriver) Copy(source DataSource, name string, options []string) (newDS DataSource, ok bool) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	ds := C.OGR_Dr_CopyDataSource(driver.cval, source.cval, cName, (**C.char)(unsafe.Pointer(&opts[0])))
	return DataSource{ds}, ds != nil
}

// Delete a data source
func (driver OGRDriver) Delete(filename string) error {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	return C.OGR_Dr_DeleteDataSource(driver.cval, cFilename).Err()
}

// Add a driver to the list of registered drivers
func (driver OGRDriver) Register() {
	C.OGRRegisterDriver(driver.cval)
}

// Remove a driver from the list of registered drivers
func (driver OGRDriver) Deregister() {
	C.OGRDeregisterDriver(driver.cval)
}

// Fetch the number of registered drivers
func OGRDriverCount() int {
	count := C.OGRGetDriverCount()
	return int(count)
}

// Fetch the indicated driver by index
func OGRDriverByIndex(index int) OGRDriver {
	driver := C.OGRGetDriver(C.int(index))
	return OGRDriver{driver}
}

// Fetch the indicated driver by name
func OGRDriverByName(name string) OGRDriver {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	driver := C.OGRGetDriverByName(cName)
	return OGRDriver{driver}
}
