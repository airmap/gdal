package main

import (
	"fmt"
	"log"

	"github.com/dewberry/gdal"
)

func main() {

	err := gdal.VSIMkdir("/vsimem/a_dir")
	if err != nil {
		log.Fatal(err)
	}

	dir := "/vsimem/a_directory/asdf"
	err = gdal.VSIMkdirRecursive(dir)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		log.Fatal(err)
	}

	ds := driver.Create("/vsimem/a_directory/tmp.tif", 256, 256, 1, gdal.Byte, []string{})
	ds.Close()

	fmt.Println("Globbing all vsimem files:")
	fmt.Println(gdal.VSIReadDirRecursive("/vsimem/"), "\n")

	err = gdal.VSIRmdir("/vsimem/a_dir")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Globbing all vsimem files after rmdir:")
	fmt.Println(gdal.VSIReadDirRecursive("/vsimem/"), "\n")

	fmt.Println("Adding a second vsimem file...")
	ds2 := driver.Create("/vsimem/a_directory/tmp2.tif", 256, 256, 1, gdal.Byte, []string{})
	ds2.Close()

	fmt.Println("Globbing vsimem files in a_directory/:")
	fmt.Println(gdal.VSIReadDirRecursive("/vsimem/a_directory"), "\n")

	err = gdal.VSIMkdir("/vsimem/a_dir2")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Globbing all vsimem files:")
	fmt.Println(gdal.VSIReadDirRecursive("/vsimem/"), "\n")

	fmt.Println("Removing a_directory/*:")
	err = gdal.VSIRmdirRecursive("/vsimem/a_directory")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Globbing all vsimem files:")
	fmt.Println(gdal.VSIReadDirRecursive("/vsimem/"))

}
