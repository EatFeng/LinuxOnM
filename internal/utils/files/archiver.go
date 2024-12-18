package files

type ShellArchiver interface {
	Compress(sourcePaths []string, dstFile string, secret string) error
}
