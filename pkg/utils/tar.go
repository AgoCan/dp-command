package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func Untargz(tarball, target string) error {
	file, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := target + "/" + header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			dir := path[:len(path)-len(header.FileInfo().Name())]
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			writer, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(writer, tarReader); err != nil {
				writer.Close()
				return err
			}
			writer.Close()
		default:
			return fmt.Errorf("unkown type: %v in %s", header.Typeflag, header.Name)
		}
	}
	return nil
}
