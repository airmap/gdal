package ogr

/*
#include "go_ogr_wkb.h"
#include "gdal_version.h"
*/
import "C"
import (
	"errors"
)

func init() {
	C.OGRRegisterAll()
}

// Convert a go bool to a C int
func BoolToCInt(in bool) (out C.int) {
	if in {
		out = 1
	} else {
		out = 0
	}
	return
}

func IntSliceToCInt(data []int) []C.int {
	sliceSz := len(data)
	result := make([]C.int, sliceSz)
	for i := 0; i < sliceSz; i++ {
		result[i] = C.int(data[i])
	}
	return result
}

/* -------------------------------------------------------------------- */
/*      Significant constants.                                          */
/* -------------------------------------------------------------------- */

// List of well known binary geometry types
type GeometryType uint32

const (
	GT_Null                  = 4294967295
	GT_Unknown               = GeometryType(C.wkbUnknown)
	GT_Point                 = GeometryType(C.wkbPoint)
	GT_LineString            = GeometryType(C.wkbLineString)
	GT_Polygon               = GeometryType(C.wkbPolygon)
	GT_MultiPoint            = GeometryType(C.wkbMultiPoint)
	GT_MultiLineString       = GeometryType(C.wkbMultiLineString)
	GT_MultiPolygon          = GeometryType(C.wkbMultiPolygon)
	GT_GeometryCollection    = GeometryType(C.wkbGeometryCollection)
	GT_None                  = GeometryType(C.wkbNone)
	GT_LinearRing            = GeometryType(C.wkbLinearRing)
	GT_Point25D              = GeometryType(C.wkbPoint25D)
	GT_LineString25D         = GeometryType(C.wkbLineString25D)
	GT_Polygon25D            = GeometryType(C.wkbPolygon25D)
	GT_MultiPoint25D         = GeometryType(C.wkbMultiPoint25D)
	GT_MultiLineString25D    = GeometryType(C.wkbMultiLineString25D)
	GT_MultiPolygon25D       = GeometryType(C.wkbMultiPolygon25D)
	GT_GeometryCollection25D = GeometryType(C.wkbGeometryCollection25D)
)

var (
	ErrDebug                   = errors.New("Debug Error")
	ErrNotEnoughData           = errors.New("Not Enough Data Error")
	ErrNotEnoughMemory         = errors.New("Not Enough Memory Error")
	ErrUnsupportedGeometryType = errors.New("Unsupported Geometry Type Error")
	ErrUnsupportedOperation    = errors.New("Unsupported Operation Error")
	ErrCorruptData             = errors.New("Corrupt Data Error")
	ErrFailure                 = errors.New("OGR Failure Error")
	ErrUnsupportedSRS          = errors.New("Unsupported SRS Error")
	ErrInvalidHandle           = errors.New("Invalid Handle Error")
	ErrNonExistingFeature      = errors.New("Non Existing Feature Error")
	ErrUndefined               = errors.New("Undefined Error")
)

func (err C.OGRErr) Err() error {
	switch err {
	case 0:
		return nil
	case 1:
		return ErrNotEnoughData
	case 2:
		return ErrNotEnoughMemory
	case 3:
		return ErrUnsupportedGeometryType
	case 4:
		return ErrUnsupportedOperation
	case 5:
		return ErrCorruptData
	case 6:
		return ErrFailure
	case 7:
		return ErrUnsupportedSRS
	case 8:
		return ErrInvalidHandle
	case 9:
		return ErrNonExistingFeature
	}
	return ErrUndefined
}

/* -------------------------------------------------------------------- */
/*      Envelope functions                                              */
/* -------------------------------------------------------------------- */

type Envelope struct {
	cval C.OGREnvelope
}

func (env Envelope) MinX() float64 {
	return float64(env.cval.MinX)
}

func (env Envelope) MaxX() float64 {
	return float64(env.cval.MaxX)
}

func (env Envelope) MinY() float64 {
	return float64(env.cval.MinY)
}

func (env Envelope) MaxY() float64 {
	return float64(env.cval.MaxY)
}

func (env *Envelope) SetMinX(val float64) {
	env.cval.MinX = C.double(val)
}

func (env *Envelope) SetMaxX(val float64) {
	env.cval.MaxX = C.double(val)
}

func (env *Envelope) SetMinY(val float64) {
	env.cval.MinY = C.double(val)
}

func (env *Envelope) SetMaxY(val float64) {
	env.cval.MaxY = C.double(val)
}

func (env Envelope) IsInit() bool {
	return env.cval.MinX != 0 || env.cval.MinY != 0 || env.cval.MaxX != 0 || env.cval.MaxY != 0
}

func min(a, b C.double) C.double {
	if a < b {
		return a
	}
	return b
}

func max(a, b C.double) C.double {
	if a > b {
		return a
	}
	return b
}

// Return the union of this envelope with another one
func (env Envelope) Union(other Envelope) Envelope {
	if env.IsInit() {
		env.cval.MinX = min(env.cval.MinX, other.cval.MinX)
		env.cval.MinY = min(env.cval.MinY, other.cval.MinY)
		env.cval.MaxX = max(env.cval.MaxX, other.cval.MaxX)
		env.cval.MaxY = max(env.cval.MaxY, other.cval.MaxY)
	} else {
		env.cval.MinX = other.cval.MinX
		env.cval.MinY = other.cval.MinY
		env.cval.MaxX = other.cval.MaxX
		env.cval.MaxY = other.cval.MaxY
	}
	return env
}

// Return the intersection of this envelope with another
func (env Envelope) Intersect(other Envelope) Envelope {
	if env.Intersects(other) {
		if env.IsInit() {
			env.cval.MinX = max(env.cval.MinX, other.cval.MinX)
			env.cval.MinY = max(env.cval.MinY, other.cval.MinY)
			env.cval.MaxX = min(env.cval.MaxX, other.cval.MaxX)
			env.cval.MaxY = min(env.cval.MaxY, other.cval.MaxY)
		} else {
			env.cval.MinX = other.cval.MinX
			env.cval.MinY = other.cval.MinY
			env.cval.MaxX = other.cval.MaxX
			env.cval.MaxY = other.cval.MaxY
		}
	} else {
		env.cval.MinX = 0
		env.cval.MinY = 0
		env.cval.MaxX = 0
		env.cval.MaxY = 0
	}
	return env
}

// Test if one envelope intersects another
func (env Envelope) Intersects(other Envelope) bool {
	return env.cval.MinX <= other.cval.MaxX &&
		env.cval.MaxX >= other.cval.MinX &&
		env.cval.MinY <= other.cval.MaxY &&
		env.cval.MaxY >= other.cval.MinY
}

// Test if one envelope completely contains another
func (env Envelope) Contains(other Envelope) bool {
	return env.cval.MinX <= other.cval.MinX &&
		env.cval.MaxX >= other.cval.MaxX &&
		env.cval.MinY <= other.cval.MinY &&
		env.cval.MaxY >= other.cval.MaxY
}

/* -------------------------------------------------------------------- */
/*      Misc functions                                                  */
/* -------------------------------------------------------------------- */

// Clean up all OGR related resources
func CleanupOGR() {
	C.OGRCleanupAll()
}

/* -------------------------------------------------------------------- */
/*      Driver functions                                                */
/* -------------------------------------------------------------------- */

/* -------------------------------------------------------------------- */
/*      Style manager functions                                         */
/* -------------------------------------------------------------------- */

type StyleMgr struct {
	cval C.OGRStyleMgrH
}

type StyleTool struct {
	cval C.OGRStyleToolH
}

type StyleTable struct {
	cval C.OGRStyleTableH
}

// Unimplemented: CreateStyleManager

// Unimplemented: Destroy

// Unimplemented: InitFromFeature

// Unimplemented: InitStyleString

// Unimplemented: PartCount

// Unimplemented: PartCount

// Unimplemented: AddPart

// Unimplemented: AddStyle

// Unimplemented: CreateStyleTool

// Unimplemented: Destroy

// Unimplemented: Type

// Unimplemented: Unit

// Unimplemented: SetUnit

// Unimplemented: ParamStr

// Unimplemented: ParamNum

// Unimplemented: ParamDbl

// Unimplemented: SetParamStr

// Unimplemented: SetParamNum

// Unimplemented: SetParamDbl

// Unimplemented: StyleString

// Unimplemented: RGBFromString

// Unimplemented: CreateStyleTable

// Unimplemented: Destroy

// Unimplemented: Save

// Unimplemented: Load

// Unimplemented: Find

// Unimplemented: ResetStyleStringReading

// Unimplemented: NextStyle

// Unimplemented: LastStyleName
