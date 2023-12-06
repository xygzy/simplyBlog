
#export GIN_MODE=release

# Windows x86 64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build blog.go
mv blog.exe publish/blog_64.exe

# Windows 386
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build blog.go
mv blog.exe publish/blog_32.exe

# Linux arm 32
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build blog.go
mv blog publish/blog_arm_32

# Linux mac 32
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build blog.go
mv blog publish/blog_mac_64

# Linux x86 64
go build blog.go
mv blog publish/blog_64

# commit the update
svn commit -m "BLOG pulished"

