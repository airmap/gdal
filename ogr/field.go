package ogr

/*
#include "go_ogr_wkb.h"
#include "gdal_version.h"
*/
import "C"

import (
	"unsafe"
)

/* -------------------------------------------------------------------- */
/*      Field definition functions                                      */
/* -------------------------------------------------------------------- */

// List of well known binary geometry types
type FieldType int

const (
	FT_Integer       = FieldType(C.OFTInteger)
	FT_IntegerList   = FieldType(C.OFTIntegerList)
	FT_Real          = FieldType(C.OFTReal)
	FT_RealList      = FieldType(C.OFTRealList)
	FT_String        = FieldType(C.OFTString)
	FT_StringList    = FieldType(C.OFTStringList)
	FT_Binary        = FieldType(C.OFTBinary)
	FT_Date          = FieldType(C.OFTDate)
	FT_Time          = FieldType(C.OFTTime)
	FT_DateTime      = FieldType(C.OFTDateTime)
	FT_Integer64     = FieldType(C.OFTInteger64)
	FT_Integer64List = FieldType(C.OFTInteger64List)
)

type Justification int

const (
	J_Undefined = Justification(C.OJUndefined)
	J_Left      = Justification(C.OJLeft)
	J_Right     = Justification(C.OJRight)
)

type FieldDefinition struct {
	cval C.OGRFieldDefnH
}

type Field struct {
	cval *C.OGRField
}

// Create a new field definition
func CreateFieldDefinition(name string, fieldType FieldType) FieldDefinition {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	fieldDef := C.OGR_Fld_Create(cName, C.OGRFieldType(fieldType))
	return FieldDefinition{fieldDef}
}

// Destroy the field definition
func (fd FieldDefinition) Destroy() {
	C.OGR_Fld_Destroy(fd.cval)
}

// Fetch the name of the field
func (fd FieldDefinition) Name() string {
	name := C.OGR_Fld_GetNameRef(fd.cval)
	return C.GoString(name)
}

// Set the name of the field
func (fd FieldDefinition) SetName(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.OGR_Fld_SetName(fd.cval, cName)
}

// Fetch the type of this field
func (fd FieldDefinition) Type() FieldType {
	fType := C.OGR_Fld_GetType(fd.cval)
	return FieldType(fType)
}

// Set the type of this field
func (fd FieldDefinition) SetType(fType FieldType) {
	C.OGR_Fld_SetType(fd.cval, C.OGRFieldType(fType))
}

// Fetch the justification for this field
func (fd FieldDefinition) Justification() Justification {
	justify := C.OGR_Fld_GetJustify(fd.cval)
	return Justification(justify)
}

// Set the justification for this field
func (fd FieldDefinition) SetJustification(justify Justification) {
	C.OGR_Fld_SetJustify(fd.cval, C.OGRJustification(justify))
}

// Fetch the formatting width for this field
func (fd FieldDefinition) Width() int {
	width := C.OGR_Fld_GetWidth(fd.cval)
	return int(width)
}

// Set the formatting width for this field
func (fd FieldDefinition) SetWidth(width int) {
	C.OGR_Fld_SetWidth(fd.cval, C.int(width))
}

// Fetch the precision for this field
func (fd FieldDefinition) Precision() int {
	precision := C.OGR_Fld_GetPrecision(fd.cval)
	return int(precision)
}

// Set the precision for this field
func (fd FieldDefinition) SetPrecision(precision int) {
	C.OGR_Fld_SetPrecision(fd.cval, C.int(precision))
}

// Set defining parameters of field in a single call
func (fd FieldDefinition) Set(
	name string,
	fType FieldType,
	width, precision int,
	justify Justification,
) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	C.OGR_Fld_Set(
		fd.cval,
		cName,
		C.OGRFieldType(fType),
		C.int(width),
		C.int(precision),
		C.OGRJustification(justify),
	)
}

// Fetch whether this field should be ignored when fetching features
func (fd FieldDefinition) IsIgnored() bool {
	ignore := C.OGR_Fld_IsIgnored(fd.cval)
	return ignore != 0
}

// Set whether this field should be ignored when fetching features
func (fd FieldDefinition) SetIgnored(ignore bool) {
	C.OGR_Fld_SetIgnored(fd.cval, BoolToCInt(ignore))
}

// Fetch human readable name for the field type
func (ft FieldType) Name() string {
	name := C.OGR_GetFieldTypeName(C.OGRFieldType(ft))
	return C.GoString(name)
}

/* -------------------------------------------------------------------- */
/*      Feature definition functions                                    */
/* -------------------------------------------------------------------- */

type FeatureDefinition struct {
	cval C.OGRFeatureDefnH
}

// Create a new feature definition object
func CreateFeatureDefinition(name string) FeatureDefinition {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	fd := C.OGR_FD_Create(cName)
	return FeatureDefinition{fd}
}

// Destroy a feature definition object
func (fd FeatureDefinition) Destroy() {
	C.OGR_FD_Destroy(fd.cval)
}

// Drop a reference, and delete object if no references remain
func (fd FeatureDefinition) Release() {
	C.OGR_FD_Release(fd.cval)
}

// Fetch the name of this feature definition
func (fd FeatureDefinition) Name() string {
	name := C.OGR_FD_GetName(fd.cval)
	return C.GoString(name)
}

// Fetch the number of fields in the feature definition
func (fd FeatureDefinition) FieldCount() int {
	count := C.OGR_FD_GetFieldCount(fd.cval)
	return int(count)
}

// Fetch the definition of the indicated field
func (fd FeatureDefinition) FieldDefinition(index int) FieldDefinition {
	fieldDefn := C.OGR_FD_GetFieldDefn(fd.cval, C.int(index))
	return FieldDefinition{fieldDefn}
}

// Fetch the index of the named field
func (fd FeatureDefinition) FieldIndex(name string) int {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	index := C.OGR_FD_GetFieldIndex(fd.cval, cName)
	return int(index)
}

// Add a new field definition to this feature definition
func (fd FeatureDefinition) AddFieldDefinition(fieldDefn FieldDefinition) {
	C.OGR_FD_AddFieldDefn(fd.cval, fieldDefn.cval)
}

// Delete a field definition from this feature definition
func (fd FeatureDefinition) DeleteFieldDefinition(index int) error {
	return C.OGR_FD_DeleteFieldDefn(fd.cval, C.int(index)).Err()
}

// Fetch the geometry base type of this feature definition
func (fd FeatureDefinition) GeometryType() GeometryType {
	gt := C.OGR_FD_GetGeomType(fd.cval)
	return GeometryType(gt)
}

// Set the geometry base type for this feature definition
func (fd FeatureDefinition) SetGeometryType(geomType GeometryType) {
	C.OGR_FD_SetGeomType(fd.cval, C.OGRwkbGeometryType(geomType))
}

// Fetch if the geometry can be ignored when fetching features
func (fd FeatureDefinition) IsGeometryIgnored() bool {
	isIgnored := C.OGR_FD_IsGeometryIgnored(fd.cval)
	return isIgnored != 0
}

// Set whether the geometry can be ignored when fetching features
func (fd FeatureDefinition) SetGeometryIgnored(val bool) {
	C.OGR_FD_SetGeometryIgnored(fd.cval, BoolToCInt(val))
}

// Fetch if the style can be ignored when fetching features
func (fd FeatureDefinition) IsStyleIgnored() bool {
	isIgnored := C.OGR_FD_IsStyleIgnored(fd.cval)
	return isIgnored != 0
}

// Set whether the style can be ignored when fetching features
func (fd FeatureDefinition) SetStyleIgnored(val bool) {
	C.OGR_FD_SetStyleIgnored(fd.cval, BoolToCInt(val))
}

// Increment the reference count by one
func (fd FeatureDefinition) Reference() int {
	count := C.OGR_FD_Reference(fd.cval)
	return int(count)
}

// Decrement the reference count by one
func (fd FeatureDefinition) Dereference() int {
	count := C.OGR_FD_Dereference(fd.cval)
	return int(count)
}

// Fetch the current reference count
func (fd FeatureDefinition) ReferenceCount() int {
	count := C.OGR_FD_GetReferenceCount(fd.cval)
	return int(count)
}
